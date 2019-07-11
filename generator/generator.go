// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"errors"
	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	openapiv3 "github.com/googleapis/gnostic/OpenAPIv3"
	surface_v1 "github.com/googleapis/gnostic/surface"
	"google.golang.org/genproto/googleapis/api/annotations"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var protoBufScalarTypes = getProtobufTypes()
var openAPITypesToProtoBuf = getOpenAPITypesToProtoBufTypes()
var openAPIScalarTypes = getOpenAPIScalarTypes()

// Gathers all external dependencies we generated in recursive calls.
var generatedDependencies = make(map[string]bool, 0)

// Gathers all messages that have been generated from external dependencies in recursive calls.
var generatedMessages = make(map[string]string, 0)

// Uses the output of gnostic to return a dpb.FileDescriptorSet (in bytes). 'renderer' contains
// the 'model' (surface model) which has all the relevant data to create the dpb.FileDescriptorSet.
// There are four main steps:
// 		1. buildDependencies to build all static FileDescriptorProto we need.
// 		2. buildExternalDependencies recursively executes this plugin to generate all FileDescriptorSet based on external dependencies.
//		2. buildMessagesFromTypes is called to create all messages which will be rendered in .proto
//		3. buildServiceFromMethods is called to create a RPC service which will be rendered in .proto
func (renderer *Renderer) runFileDescriptorSetGenerator() (fdSet *dpb.FileDescriptorSet, err error) {
	syntax := "proto3"
	n := renderer.Package + ".proto"

	// mainProto is the proto we ultimately want to render.
	mainProto := &dpb.FileDescriptorProto{
		Name:    &n,
		Package: &renderer.Package,
		Syntax:  &syntax,
	}
	fdSet = &dpb.FileDescriptorSet{
		File: []*dpb.FileDescriptorProto{mainProto},
	}

	buildDependencies(fdSet)
	err = buildExternalDependencies(fdSet, renderer)
	if err != nil {
		return nil, err
	}

	addDependencies(fdSet)

	err = buildMessagesFromTypes(mainProto, renderer)
	if err != nil {
		return nil, err
	}

	err = buildServiceFromMethods(mainProto, renderer)
	if err != nil {
		return nil, err
	}

	return fdSet, err
}

// Adds the dependencies to the FileDescriptor we want to render. This essentially makes the 'import' statements
// inside the .proto definition.
func addDependencies(fdSet *dpb.FileDescriptorSet) {
	// At last, we need to add the dependencies to the FileDescriptorProto in order to get them rendered.
	lastFdProto := getLast(fdSet.File)
	for _, fd := range fdSet.File {
		if fd != lastFdProto {
			lastFdProto.Dependency = append(lastFdProto.Dependency, *fd.Name)
		}
	}
}

// buildExternalDependencies recursively generates all .proto definitions to external OpenAPI descriptions (URLs to other
// descriptions inside the description).
func buildExternalDependencies(fdSet *dpb.FileDescriptorSet, renderer *Renderer) (err error) {
	externalDependencies := renderer.Model.Dependencies
	externalDependencies = trimAndRemoveDuplicates(externalDependencies)

	externalFileDescriptorProtos := make([]*dpb.FileDescriptorProto, 0)
	for _, dep := range externalDependencies {
		if _, alreadyGenerated := generatedDependencies[dep]; !alreadyGenerated {
			generatedDependencies[dep] = true

			// Lets get the standard gnostic output from the external dependency.
			cmd := exec.Command("gnostic", "--pb-out=-", dep)
			b, err := cmd.Output()
			if err != nil {
				return err
			}

			// Construct an OpenAPI document v3.
			document, err := createOpenAPIDocFromGnosticOutput(b)
			if err != nil {
				return err
			}

			// Create the surface model. Keep in mind that this resolves the dependencies of the external dep again!
			surfaceModel, err := surface_v1.NewModelFromOpenAPI3(document, dep)
			if err != nil {
				return err
			}

			// Recursively call the generator.
			recursiveRenderer := NewRenderer(surfaceModel)
			fileName := path.Base(dep)
			recursiveRenderer.Package = strings.TrimSuffix(fileName, filepath.Ext(fileName))
			newFdSet, err := recursiveRenderer.runFileDescriptorSetGenerator()
			if err != nil {
				return err
			}
			renderer.ExternalFdSets = append(renderer.ExternalFdSets, newFdSet)

			dependencyProto := getLast(newFdSet.File)
			externalFileDescriptorProtos = append(externalFileDescriptorProtos, dependencyProto)
		}
	}

	fdSet.File = append(externalFileDescriptorProtos, fdSet.File...)
	return nil
}

// Protoreflect needs all the dependencies that are used inside of the FileDescriptorProto (that gets rendered)
// to work properly. Those dependencies are google/protobuf/empty.proto, google/api/annotations.proto,
// "google/protobuf/descriptor.proto" and "google/protobuf/any". For all those dependencies the corresponding
// FileDescriptorProto has to be added to the FileDescriptorSet. Protoreflect won't work
// if a reference is missing.
func buildDependencies(fdSet *dpb.FileDescriptorSet) {
	// Dependency to google/api/annotations.proto for gRPC-HTTP transcoding. Here a couple of problems arise:
	// 1. Problem: 	We cannot call descriptor.ForMessage(&annotations.E_Http), which would be our
	//				required dependency. However, we can call descriptor.ForMessage(&http) and
	//				then construct the extension manually.
	// 2. Problem: 	The name is set wrong.
	// 3. Problem: 	google/api/annotations.proto has a dependency to google/protobuf/descriptor.proto.
	http := annotations.Http{}
	fd, _ := descriptor.ForMessage(&http)

	extensionName := "http"
	n := "google/api/annotations.proto"
	l := dpb.FieldDescriptorProto_LABEL_OPTIONAL
	t := dpb.FieldDescriptorProto_TYPE_MESSAGE
	tName := "google.api.HttpRule"
	extendee := ".google.protobuf.MethodOptions"

	httpExtension := &dpb.FieldDescriptorProto{
		Name:     &extensionName,
		Number:   &annotations.E_Http.Field,
		Label:    &l,
		Type:     &t,
		TypeName: &tName,
		Extendee: &extendee,
	}

	fd.Extension = append(fd.Extension, httpExtension)                        // 1. Problem
	fd.Name = &n                                                              // 2. Problem
	fd.Dependency = append(fd.Dependency, "google/protobuf/descriptor.proto") //3.rd Problem

	// Build other required dependencies
	e := empty.Empty{}
	fdp := dpb.DescriptorProto{}
	a := any.Any{}
	fd2, _ := descriptor.ForMessage(&e)
	fd3, _ := descriptor.ForMessage(&fdp)
	fd4, _ := descriptor.ForMessage(&a)
	dependencies := []*dpb.FileDescriptorProto{fd, fd2, fd3, fd4}

	// According to the documentation of protoReflect.CreateFileDescriptorFromSet the file I want to print
	// needs to be at the end of the array. All other FileDescriptorProto are dependencies.
	fdSet.File = append(dependencies, fdSet.File...)
}

// Builds protobuf messages from the surface model types. If the type is a RPC request parameter
// the fields have to follow certain rules, and therefore have to be validated.
func buildMessagesFromTypes(descr *dpb.FileDescriptorProto, renderer *Renderer) (err error) {
	types := renderer.Model.Types

	for _, t := range types {
		messageName := strings.Title(t.Name)

		message := dpb.DescriptorProto{
			Name: &messageName,
		}

		for i, f := range t.Fields {
			if isRequestParameter(t) {
				if f.Position == surface_v1.Position_PATH {
					f, err = validatePathParameter(f, types)
					if err != nil {
						return err
					}
				}

				if f.Position == surface_v1.Position_QUERY {
					f, err = validateQueryParameter(f)
					if err != nil {
						return err
					}
				}
			}
			ctr := int32(i + 1)
			fieldDescriptor := &dpb.FieldDescriptorProto{Number: &ctr}
			setFieldDescriptorLabel(fieldDescriptor, f)
			setFieldDescriptorName(fieldDescriptor, f)
			setFieldDescriptorType(fieldDescriptor, f)
			setFieldDescriptorTypeName(fieldDescriptor, f, renderer.Package)

			if strings.Contains(f.Type, "map") {
				mapDescriptorProto := buildMapDescriptorProto(f)
				fieldDescriptor.TypeName = mapDescriptorProto.Name
				message.NestedType = append(message.NestedType, mapDescriptorProto)
			}
			message.Field = append(message.Field, fieldDescriptor)
		}
		descr.MessageType = append(descr.MessageType, &message)
		generatedMessages[*message.Name] = renderer.Package + "." + *message.Name
	}
	return nil
}

// Builds a protobuf RPC service. For every method the corresponding gRPC-HTTP transcoding options (https://github.com/googleapis/googleapis/blob/master/google/api/http.proto)
// have to be set.
func buildServiceFromMethods(descr *dpb.FileDescriptorProto, renderer *Renderer) (err error) {
	methods := renderer.Model.Methods
	serviceName := strings.Title(renderer.Package)

	service := &dpb.ServiceDescriptorProto{
		Name: &serviceName,
	}
	descr.Service = []*dpb.ServiceDescriptorProto{service}

	for _, method := range methods {
		// TODO: ClientStreaming
		// TODO: ServerStreaming

		mOptionsDescr := &dpb.MethodOptions{}
		requestBody := getRequestBodyForRequestParameters(method.ParametersTypeName, renderer.Model.Types)
		httpRule := getHttpRuleForMethod(method, requestBody)
		if err := proto.SetExtension(mOptionsDescr, annotations.E_Http, &httpRule); err != nil {
			return err
		}

		if method.ParametersTypeName == "" {
			method.ParametersTypeName = "google.protobuf.Empty"
		}
		if method.ResponsesTypeName == "" {
			method.ResponsesTypeName = "google.protobuf.Empty"
		}

		mDescr := &dpb.MethodDescriptorProto{
			Name:       &method.Name,
			InputType:  &method.ParametersTypeName,
			OutputType: &method.ResponsesTypeName,
			Options:    mOptionsDescr,
		}

		service.Method = append(service.Method, mDescr)
	}
	return nil
}

// Builds the necessary descriptor to render a map. (https://developers.google.com/protocol-buffers/docs/proto3#maps)
// A map is represented as nested message with two fields: 'key', 'value' and the Options set accordingly.
func buildMapDescriptorProto(field *surface_v1.Field) *dpb.DescriptorProto {
	isMapEntry := true
	n := field.Name + "Entry"

	mapDP := &dpb.DescriptorProto{
		Name:    &n,
		Field:   buildKeyValueFields(field),
		Options: &dpb.MessageOptions{MapEntry: &isMapEntry},
	}
	return mapDP
}

// Builds the necessary 'key', 'value' fields for the map descriptor.
func buildKeyValueFields(field *surface_v1.Field) []*dpb.FieldDescriptorProto {
	k, v := "key", "value"
	var n1, n2 int32 = 1, 2
	l := dpb.FieldDescriptorProto_LABEL_OPTIONAL
	t := dpb.FieldDescriptorProto_TYPE_STRING
	keyField := &dpb.FieldDescriptorProto{
		Name:   &k,
		Number: &n1,
		Label:  &l,
		Type:   &t,
	}

	valueType := field.Type[11:] // This transforms a string like 'map[string]int32' to 'int32'. In other words: the type of the value from the map.
	valueField := &dpb.FieldDescriptorProto{
		Name:     &v,
		Number:   &n2,
		Label:    &l,
		Type:     getProtoTypeForMapValueType(valueType),
		TypeName: getTypeNameForMapValueType(valueType),
	}
	return []*dpb.FieldDescriptorProto{keyField, valueField}
}

// Validates if the path parameter has the requested structure.
// This is necessary according to: https://github.com/googleapis/googleapis/blob/a8ee1416f4c588f2ab92da72e7c1f588c784d3e6/google/api/http.proto#L62
func validatePathParameter(field *surface_v1.Field, types []*surface_v1.Type) (*surface_v1.Field, error) {

	if field.Kind != surface_v1.FieldKind_SCALAR {
		if field.Kind == surface_v1.FieldKind_REFERENCE {
			// We got a reference. Let's try to flatten!
			field, err := flattenPathParameter(field, types)
			if err == nil {
				return field, nil
			}
		}
		return nil, errors.New("The path parameter with the Name " + field.Name + " is invalid. " +
			"The path template may refer to one or more fields in the gRPC request message, as" +
			" long as each field is a non-repeated field with a primitive (non-message) type")
	}

	return field, nil
}

// Validates if the query parameter has the requested structure.
// This is necessary according to: https://github.com/googleapis/googleapis/blob/a8ee1416f4c588f2ab92da72e7c1f588c784d3e6/google/api/http.proto#L119
func validateQueryParameter(field *surface_v1.Field) (*surface_v1.Field, error) {
	if !(field.Kind == surface_v1.FieldKind_SCALAR ||
		(field.Kind == surface_v1.FieldKind_ARRAY && openAPIScalarTypes[field.Type]) ||
		(field.Kind == surface_v1.FieldKind_REFERENCE)) {
		return nil, errors.New("The query parameter with the Name " + field.Name + " is invalid. " +
			"Note that fields which are mapped to URL query parameters must have a primitive type or" +
			" a repeated primitive type or a non-repeated message type.")
	}

	return field, nil
}

// If 'field' is a reference and holds a path parameter, then it will be flattened, meaning that the
// values of the reference, will be written into 'field'.
func flattenPathParameter(field *surface_v1.Field, types []*surface_v1.Type) (*surface_v1.Field, error) {
	// We got a reference to a parameter. Let's get the actual type.
	t, err := getType(field.Type, types)
	if err != nil {
		return nil, err
	}

	if t.Fields[0].Position != surface_v1.Position_PATH {
		return field, nil
	}
	if len(t.Fields) > 1 || t.Fields[0].Kind != surface_v1.FieldKind_SCALAR {
		return nil, errors.New("Not possible to flatten multiple fields or non-scalar values. ")
	}

	// Ok, it is possible to flatten the path parameter.
	field.Type = t.Fields[0].Type
	field.Name = t.Fields[0].Name
	field.Format = t.Fields[0].Format
	field.Kind = t.Fields[0].Kind
	field.Position = surface_v1.Position_PATH
	return field, nil
}

// Checks whether 't' is a type that will be used as a request parameter for a RPC method.
func isRequestParameter(t *surface_v1.Type) bool {
	if strings.Contains(t.Description, t.GetName()+" holds parameters to") {
		return true
	}
	return false
}

// Sets the Type of 'fd' according to the information from the surface field 'f'.
func setFieldDescriptorType(fd *dpb.FieldDescriptorProto, f *surface_v1.Field) {
	var protoType dpb.FieldDescriptorProto_Type
	if t, ok := protoBufScalarTypes[f.Format]; ok { // Let's see if we can get the type from f.format
		protoType = t
	} else if t, ok := protoBufScalarTypes[f.Type]; ok { // Maybe this works.
		protoType = t
	} else if t, ok := openAPITypesToProtoBuf[f.Type]; ok { // Safety check
		protoType = t
	} else {
		// TODO: What about Enums?
		// Ok, is it either a reference or an array of non scalar-types or a map. All of those get represented as message
		// inside the descriptor.
		protoType = dpb.FieldDescriptorProto_TYPE_MESSAGE
	}
	fd.Type = &protoType

}

// Sets the Name of 'fd'. The convention inside .proto is, that all field names are
// lowercase and all messages and types are capitalized if they are not scalar types (int64, string, ...).
func setFieldDescriptorName(fd *dpb.FieldDescriptorProto, f *surface_v1.Field) {
	name := strings.ToLower(f.Name)

	if name == "200" {
		name = "ok"
	}
	if name == "400" {
		name = "badRequest"
	}
	fd.Name = &name
}

// Sets a Label for 'fd'. If it is an array we need the 'repeated' label.
func setFieldDescriptorLabel(fd *dpb.FieldDescriptorProto, f *surface_v1.Field) {
	label := dpb.FieldDescriptorProto_LABEL_OPTIONAL
	if f.Kind == surface_v1.FieldKind_ARRAY || strings.Contains(f.Type, "map") {
		label = dpb.FieldDescriptorProto_LABEL_REPEATED
	}
	fd.Label = &label
}

// Sets the TypeName of 'fd'. A TypeName has to be set if the field is a reference to another message. Otherwise it is nil.
// The convention inside .proto is, that all field names are lowercase and all messages and types are capitalized if
// they are not scalar types (int64, string, ...).
func setFieldDescriptorTypeName(fd *dpb.FieldDescriptorProto, f *surface_v1.Field, packageName string) {
	typeName := ""
	// A field with a type of Message always has a typeName associated with it (the name of the Message).
	if *fd.Type == dpb.FieldDescriptorProto_TYPE_MESSAGE {
		// Lets default to the type we get from the surface model
		typeName = packageName + "." + strings.Title(f.Type)

		// Check whether we generated this message already inside of another dependency. If so we will use that name instead.
		if n, ok := generatedMessages[f.Type]; ok {
			typeName = n
		}

	}
	if typeName != "" {
		fd.TypeName = &typeName
	}
}

// Finds the corresponding surface model type for 'name' and returns the name of the field
// that is a request body. If no such field is found it returns nil.
func getRequestBodyForRequestParameters(name string, types []*surface_v1.Type) *string {
	requestParameterType := &surface_v1.Type{}

	for _, t := range types {
		if t.Name == name {
			requestParameterType = t
		}
	}

	for _, f := range requestParameterType.Fields {
		if f.Position == surface_v1.Position_BODY {
			return &f.Name
		}
	}
	return nil
}

// Constructs a HttpRule from google/api/http.proto. Enables gRPC-HTTP transcoding on 'method'.
// If not nil, body is also set.
func getHttpRuleForMethod(method *surface_v1.Method, body *string) annotations.HttpRule {
	var httpRule annotations.HttpRule
	switch method.Method {
	case "GET":
		httpRule = annotations.HttpRule{
			Pattern: &annotations.HttpRule_Get{
				Get: method.Path,
			},
		}
	case "POST":
		httpRule = annotations.HttpRule{
			Pattern: &annotations.HttpRule_Post{
				Post: method.Path,
			},
		}
	case "PUT":
		httpRule = annotations.HttpRule{
			Pattern: &annotations.HttpRule_Put{
				Put: method.Path,
			},
		}
	case "PATCH":
		httpRule = annotations.HttpRule{
			Pattern: &annotations.HttpRule_Patch{
				Patch: method.Path,
			},
		}
	case "DELETE":
		httpRule = annotations.HttpRule{
			Pattern: &annotations.HttpRule_Delete{
				Delete: method.Path,
			},
		}
	}

	if body != nil {
		httpRule.Body = *body
	}

	return httpRule
}

// Searches all types from the surface model for a given type 'name'. Returns a type if there is
// a match, nil if there is no match, and error if there are multiple types.
func getType(name string, types []*surface_v1.Type) (*surface_v1.Type, error) {
	var result []*surface_v1.Type
	for _, t := range types {
		if name == t.Name {
			result = append(result, t)
		}
	}
	if len(result) > 1 {
		return nil, errors.New("Multiple types with the same name exist. This is due to the fact" +
			" that there are multiple components inside the OpenAPI specification with the same Name. ")
	}
	if len(result) == 1 {
		return result[0], nil
	}
	return nil, nil
}

// Returns the type name for the given 'valueType'. A type name for a field is only set if it is some kind of
// reference (non-scalar values) otherwise it is nil.
func getTypeNameForMapValueType(valueType string) *string {
	if _, ok := protoBufScalarTypes[valueType]; ok {
		// Ok it is a scalar. For scalar values we don't set the TypeName of the field.
		return nil
	}
	if strings.Contains(valueType, "[]") {
		// We got an array as value type. This can't be represented inside .proto. So let's return the 'any' type.
		anyType := ".google.protobuf.Any"
		return &anyType
	}
	return &valueType
}

// Returns the 'protoType' for the given 'valueType'. If we don't have a scalar 'protoType', we have some kind of
// reference to another object and therefore return the 'Message' type. d
func getProtoTypeForMapValueType(valueType string) *dpb.FieldDescriptorProto_Type {
	protoType := dpb.FieldDescriptorProto_TYPE_MESSAGE
	if protoType, ok := protoBufScalarTypes[valueType]; ok {
		return &protoType
	}
	return &protoType
}

// Uses the 'binaryInput' from gnostic to create a OpenAPI document.
func createOpenAPIDocFromGnosticOutput(binaryInput []byte) (*openapiv3.Document, error) {
	document := &openapiv3.Document{}
	err := proto.Unmarshal(binaryInput, document)
	if err != nil {
		// If we execute gnostic with argument: '-pb-out=-' we get an EOF. So lets only return other errors.
		if err.Error() != "unexpected EOF" {
			return nil, err
		}
	}
	return document, nil
}

// 'externalDependencies' is a list of URLs to other OpenAPI descriptions. We need
// the base of all URLs and no duplicates.
func trimAndRemoveDuplicates(externalDependencies []string) []string {
	result := make([]string, 0)
	for _, dep := range externalDependencies {
		parts := strings.Split(dep, "#")
		if !isDuplicate(result, parts[0]) {
			result = append(result, parts[0])
		}
	}
	return result
}

// Returns true if 's' is inside result.
func isDuplicate(result []string, s string) bool {
	for _, s2 := range result {
		if s == s2 {
			return true
		}
	}
	return false
}

func getLast(protos []*dpb.FileDescriptorProto) *dpb.FileDescriptorProto {
	return protos[len(protos)-1]
}

// A map for this: https://developers.google.com/protocol-buffers/docs/proto3#scalar
func getProtobufTypes() map[string]dpb.FieldDescriptorProto_Type {
	typeMapping := make(map[string]dpb.FieldDescriptorProto_Type)
	typeMapping["double"] = dpb.FieldDescriptorProto_TYPE_DOUBLE
	typeMapping["float"] = dpb.FieldDescriptorProto_TYPE_FLOAT
	typeMapping["int64"] = dpb.FieldDescriptorProto_TYPE_INT64
	typeMapping["uint64"] = dpb.FieldDescriptorProto_TYPE_UINT64
	typeMapping["int32"] = dpb.FieldDescriptorProto_TYPE_INT32
	typeMapping["fixed64"] = dpb.FieldDescriptorProto_TYPE_FIXED64

	typeMapping["fixed32"] = dpb.FieldDescriptorProto_TYPE_FIXED32
	typeMapping["bool"] = dpb.FieldDescriptorProto_TYPE_BOOL
	typeMapping["string"] = dpb.FieldDescriptorProto_TYPE_STRING
	typeMapping["bytes"] = dpb.FieldDescriptorProto_TYPE_BYTES
	typeMapping["uint32"] = dpb.FieldDescriptorProto_TYPE_UINT32
	typeMapping["sfixed32"] = dpb.FieldDescriptorProto_TYPE_SFIXED32
	typeMapping["sfixed64"] = dpb.FieldDescriptorProto_TYPE_SFIXED64
	typeMapping["sint32"] = dpb.FieldDescriptorProto_TYPE_SINT32
	typeMapping["sint64"] = dpb.FieldDescriptorProto_TYPE_SINT64
	return typeMapping
}

// Maps OpenAPI data types (https://swagger.io/docs/specification/data-models/data-types/)
// to protobuf data types.
func getOpenAPITypesToProtoBufTypes() map[string]dpb.FieldDescriptorProto_Type {
	return map[string]dpb.FieldDescriptorProto_Type{
		"string":  dpb.FieldDescriptorProto_TYPE_STRING,
		"integer": dpb.FieldDescriptorProto_TYPE_INT32,
		"number":  dpb.FieldDescriptorProto_TYPE_FLOAT,
		"boolean": dpb.FieldDescriptorProto_TYPE_BOOL,
		"object":  dpb.FieldDescriptorProto_TYPE_MESSAGE,
		// Array not set: could be either scalar or non-scalar value.
	}
}

// All scalar types from OpenAPI.
func getOpenAPIScalarTypes() map[string]bool {
	return map[string]bool{
		"string":  true,
		"integer": true,
		"number":  true,
		"boolean": true,
	}
}
