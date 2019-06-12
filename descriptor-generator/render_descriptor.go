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
	buildServiceFromMethods(fileDescriptorProto, renderer)
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

			fieldDescr := descriptor.FieldDescriptorProto{}
			fieldDescr.Name = &f.Name
			fieldDescr.Number = &(ctr)

			label := descriptor.FieldDescriptorProto_LABEL_OPTIONAL
			fieldDescr.Label = &label

			// TODO: Use switch
			if f.Kind == surface_v1.FieldKind_SCALAR {
				protoType := getProtoTypeForField(f)
				fieldDescr.Type = &protoType
			} else if f.Kind == surface_v1.FieldKind_REFERENCE {
				// TODO: Could this also be enum?
				protoType := descriptor.FieldDescriptorProto_TYPE_MESSAGE
				fieldDescr.Type = &protoType
				fieldDescr.TypeName = &f.Type
			} else if f.Kind == surface_v1.FieldKind_ARRAY {
				label = descriptor.FieldDescriptorProto_LABEL_REPEATED
				protoType := getProtoTypeForField(f)
				fieldDescr.Type = &protoType
			} else if f.Kind == surface_v1.FieldKind_MAP {
				// TODO
			}
			message.Field = append(message.Field, &fieldDescr)
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

func buildServiceFromMethods(descr *descriptor.FileDescriptorProto, renderer *Renderer) {
	methods := renderer.Model.Methods

	service := &descriptor.ServiceDescriptorProto{}
	descr.Service = []*descriptor.ServiceDescriptorProto{service}

	serviceName := strings.Title(renderer.Package)
	service.Name = &serviceName

	for _, method := range methods {
		// TODO: How to transfer information about http transcoding annotations? Currently inside UninterpretedOption
		// TODO: ClientStreaming
		// TODO: ServerStreaming

		mUninterpretedOptions := descriptor.UninterpretedOption{}
		identIfierValue := method.Path + ";" + method.Method
		mUninterpretedOptions.IdentifierValue = &identIfierValue

		mOptionsDescr := descriptor.MethodOptions{}
		mOptionsDescr.UninterpretedOption = []*descriptor.UninterpretedOption{&mUninterpretedOptions}

		mDescr := descriptor.MethodDescriptorProto{}
		mDescr.Name = &method.Name
		mDescr.InputType = &method.ParametersTypeName
		mDescr.OutputType = &method.ResponsesTypeName
		mDescr.Options = &mOptionsDescr

		service.Method = append(service.Method, &mDescr)
	}
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
