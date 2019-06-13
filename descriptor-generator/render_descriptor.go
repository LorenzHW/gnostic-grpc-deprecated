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
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	surface_v1 "github.com/googleapis/gnostic/surface"
	"google.golang.org/genproto/googleapis/api/annotations"
	"strings"
)

var typeMapping = getTypeMapping()

func (renderer *Renderer) RenderFileDescriptorSet() (res []byte, err error) {
	fileDescriptorProto := &descriptor.FileDescriptorProto{
		Name:    &renderer.Package,
		Package: &renderer.Package,
	}
	fileDescrSet := descriptor.FileDescriptorSet{
		File: []*descriptor.FileDescriptorProto{fileDescriptorProto},
	}
	buildDependencies(fileDescriptorProto)
	err = buildServiceFromMethods(fileDescriptorProto, renderer)
	buildMessagesFromTypes(fileDescriptorProto, renderer)
	res, err = proto.Marshal(&fileDescrSet)
	return res, err
}

func buildMessagesFromTypes(descr *descriptor.FileDescriptorProto, renderer *Renderer) {
	types := renderer.Model.Types

	for _, t := range types {
		message := descriptor.DescriptorProto{}
		message.Name = &t.Name

		for i, f := range t.Fields {
			ctr := int32(i + 1)
			label := descriptor.FieldDescriptorProto_LABEL_OPTIONAL
			protoType := getProtoTypeForField(f)

			fieldDescr := &descriptor.FieldDescriptorProto{
				Name:   &f.Name,
				Number: &ctr,
				Label:  &label,
				Type:   &protoType,
			}

			switch f.Kind {
			case surface_v1.FieldKind_REFERENCE:
				// TODO: Could this also be enum?
				protoType = descriptor.FieldDescriptorProto_TYPE_MESSAGE
				fieldDescr.TypeName = &f.Type
			case surface_v1.FieldKind_ARRAY:
				label = descriptor.FieldDescriptorProto_LABEL_REPEATED
			case surface_v1.FieldKind_MAP:
				// TODO
			}
			message.Field = append(message.Field, fieldDescr)
		}
		descr.MessageType = append(descr.MessageType, &message)
	}
}

func buildDependencies(descr *descriptor.FileDescriptorProto) {
	dependencies := []string{"google/api/annotations.proto", "google/protobuf/empty.proto"}

	for _, dep := range dependencies {
		descr.Dependency = append(descr.Dependency, dep)
	}
}

func buildServiceFromMethods(descr *descriptor.FileDescriptorProto, renderer *Renderer) (err error) {
	methods := renderer.Model.Methods
	serviceName := strings.Title(renderer.Package)

	service := &descriptor.ServiceDescriptorProto{
		Name: &serviceName,
	}
	descr.Service = []*descriptor.ServiceDescriptorProto{service}

	for _, method := range methods {
		// TODO: ClientStreaming
		// TODO: ServerStreaming

		mOptionsDescr := &descriptor.MethodOptions{}
		httpRule := getHttpRuleForMethod(method)
		if err := proto.SetExtension(mOptionsDescr, annotations.E_Http, &httpRule); err != nil {
			return err
		}

		mDescr := &descriptor.MethodDescriptorProto{
			Name:       &method.Name,
			InputType:  &method.ParametersTypeName,
			OutputType: &method.ResponsesTypeName,
			Options:    mOptionsDescr,
		}

		service.Method = append(service.Method, mDescr)
	}
	return nil
}

func getHttpRuleForMethod(method *surface_v1.Method) annotations.HttpRule {
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
	return httpRule
}

func getTypeMapping() map[string]descriptor.FieldDescriptorProto_Type {
	typeMapping := make(map[string]descriptor.FieldDescriptorProto_Type)
	typeMapping["double"] = descriptor.FieldDescriptorProto_TYPE_DOUBLE
	typeMapping["float"] = descriptor.FieldDescriptorProto_TYPE_FLOAT
	typeMapping["int64"] = descriptor.FieldDescriptorProto_TYPE_INT64
	typeMapping["uint64"] = descriptor.FieldDescriptorProto_TYPE_UINT64
	typeMapping["int32"] = descriptor.FieldDescriptorProto_TYPE_INT32
	typeMapping["fixed64"] = descriptor.FieldDescriptorProto_TYPE_FIXED64

	typeMapping["fixed32"] = descriptor.FieldDescriptorProto_TYPE_FIXED32
	typeMapping["bool"] = descriptor.FieldDescriptorProto_TYPE_BOOL
	typeMapping["string"] = descriptor.FieldDescriptorProto_TYPE_STRING
	typeMapping["bytes"] = descriptor.FieldDescriptorProto_TYPE_BYTES
	typeMapping["uint32"] = descriptor.FieldDescriptorProto_TYPE_UINT32
	typeMapping["sfixed32"] = descriptor.FieldDescriptorProto_TYPE_SFIXED32
	typeMapping["sfixed64"] = descriptor.FieldDescriptorProto_TYPE_SFIXED64
	typeMapping["sint32"] = descriptor.FieldDescriptorProto_TYPE_SINT32
	typeMapping["sint64"] = descriptor.FieldDescriptorProto_TYPE_SINT64
	return typeMapping
}

func getProtoTypeForField(f *surface_v1.Field) descriptor.FieldDescriptorProto_Type {
	fieldType := f.Format
	if fieldType == "" {
		fieldType = f.Type
	}
	return typeMapping[fieldType]
}
