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

package descriptor_generator

import (
	"errors"
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	surface_v1 "github.com/googleapis/gnostic/surface"
	"google.golang.org/genproto/googleapis/api/annotations"
	"strings"
)

var protoBufTypes = getProtobufTypes()
var openAPITypesToProtoBuf = getOpenAPITypesToProtoBufTypes()
var openAPIScalarTypes = getOpenAPIScalarTypes()

func (renderer *Renderer) RenderFileDescriptorSet() (res []byte, err error) {
	syntax := "proto3"

	fileDescriptorProto := &dpb.FileDescriptorProto{
		Name:    &renderer.Package,
		Package: &renderer.Package,
		Syntax:  &syntax,
	}
	fileDescrSet := dpb.FileDescriptorSet{
		File: []*dpb.FileDescriptorProto{fileDescriptorProto},
	}

	buildDependencies(fileDescriptorProto)

	err = buildMessagesFromTypes(fileDescriptorProto, renderer)
	if err != nil {
		return nil, err
	}

	err = buildServiceFromMethods(fileDescriptorProto, renderer)
	if err != nil {
		return nil, err
	}

	res, err = proto.Marshal(&fileDescrSet)
	if err != nil {
		return nil, err
	}
	return res, err
}

func buildMessagesFromTypes(descr *dpb.FileDescriptorProto, renderer *Renderer) (err error) {
	types := renderer.Model.Types

	for _, t := range types {
		messageName := strings.Title(t.Name)

		message := dpb.DescriptorProto{
			Name: &messageName,
		}

		for i, f := range t.Fields {
			if isRequestParameter(t) {
				f, err = flattenPathParameter(f, types)
				if err != nil {
					return err
				}
			}

			ctr := int32(i + 1)
			label := getLabelForField(f)
			name := getNameForField(f)
			typeName := getTypeNameForField(f)

			protoType, err := getProtoTypeForField(f)
			if err != nil {
				return err
			}

			fieldDescr := &dpb.FieldDescriptorProto{
				Name:     name,
				Number:   &ctr,
				Label:    label,
				Type:     protoType,
				TypeName: typeName,
			}

			message.Field = append(message.Field, fieldDescr)
		}
		descr.MessageType = append(descr.MessageType, &message)
	}
	return nil
}

// Checks whether 't' is a type that will be used as a request parameter for a RPC method.
func isRequestParameter(t *surface_v1.Type) bool {
	if strings.Contains(t.Description, t.GetName()+" holds parameters to") {
		return true
	}
	return false
}

// If 'field' is a reference and holds a path parameter, then it will be flattened, meaning that the
// values of the reference, will be written into 'field'.
// This is necessary according to: https://github.com/googleapis/googleapis/blob/a8ee1416f4c588f2ab92da72e7c1f588c784d3e6/google/api/http.proto#L62
func flattenPathParameter(field *surface_v1.Field, types []*surface_v1.Type) (*surface_v1.Field, error) {
	if field.Kind == surface_v1.FieldKind_REFERENCE {
		// We got a reference to a parameter. Let's get the actual type.
		t, err := getType(field.Type, types)
		if err != nil {
			return nil, err
		}

		if len(t.Fields) > 1 {
			return nil, errors.New("Unable to flatten input parameters. ")
		}

		if t.Fields[0].Position == surface_v1.Position_PATH {
			// Ok, the referenced parameter is a path parameter. Let's flatten it.
			field.Type = t.Fields[0].Type
			field.Name = t.Fields[0].Name
			field.Format = t.Fields[0].Format
			field.Kind = t.Fields[0].Kind
			field.Position = surface_v1.Position_PATH
		}
	}
	return field, nil
}

// Searches all types from the surface model for a given type 'name'. Returns a type if there is
// a match, nil if there is no match, and and error if there are multiple types.
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

func getLabelForField(f *surface_v1.Field) *dpb.FieldDescriptorProto_Label {
	res := dpb.FieldDescriptorProto_LABEL_OPTIONAL
	if f.Kind == surface_v1.FieldKind_ARRAY {
		res = dpb.FieldDescriptorProto_LABEL_REPEATED
	}
	return &res
}

func buildDependencies(descr *dpb.FileDescriptorProto) {
	dependencies := []string{"google/api/annotations.proto", "google/protobuf/empty.proto"}

	for _, dep := range dependencies {
		descr.Dependency = append(descr.Dependency, dep)
	}
}

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

func getProtoTypeForField(f *surface_v1.Field) (*dpb.FieldDescriptorProto_Type, error) {
	if protoType, ok := protoBufTypes[f.Format]; ok {
		return &protoType, nil
	}

	if protoType, ok := protoBufTypes[f.Type]; ok {
		return &protoType, nil
	}

	if protoType, ok := openAPITypesToProtoBuf[f.Type]; ok {
		return &protoType, nil
	}

	if f.Kind == surface_v1.FieldKind_REFERENCE || (f.Kind == surface_v1.FieldKind_ARRAY && !openAPIScalarTypes[f.Type]) {
		// It is either a reference or an array of non scalar-types.
		protoType := dpb.FieldDescriptorProto_TYPE_MESSAGE // TODO: Could also be ENUM?
		return &protoType, nil
	}

	return nil, errors.New("Unable to find a protobuf type for the surface model type ")

}

// Returns the name of the protobuf field. The convention inside .proto is, that all field names are
// lowercase and all messages and types are capitalized if they are not scalar types (int64, string, ...).
func getNameForField(f *surface_v1.Field) *string {
	name := strings.ToLower(f.Name)

	if name == "200" {
		name = "ok"
	}
	return &name
}

// Returns the type of the reference. The convention inside .proto is, that all field names are
// lowercase and all messages and types are capitalized if they are not scalar types (int64, string, ...).
func getTypeNameForField(f *surface_v1.Field) *string {
	if f.Kind == surface_v1.FieldKind_REFERENCE || (f.Kind == surface_v1.FieldKind_ARRAY && !openAPIScalarTypes[f.Type]) {
		// It is either a reference or an array of non scalar-types.
		typeName := strings.Title(f.Type)
		return &typeName
	}

	return nil
}

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

func getOpenAPIScalarTypes() map[string]bool {
	return map[string]bool{
		"string":  true,
		"integer": true,
		"number":  true,
		"boolean": true,
	}
}
