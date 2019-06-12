// Copyright 2019 Google Inc. All Rights Reserved.
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

package protoc_generator

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"strconv"
	"unicode/utf8"
)

var descriptorTypeToProtoType = createMapping()
var labelMapping = createLabelMapping()

func (renderer *Renderer) RenderProto(fileDescrProto *descriptor.FileDescriptorProto) ([]byte, error) {
	renderer.currFileDescriptor = fileDescrProto

	f := NewLineWriter()
	//removePackageFromNames(renderer)

	// TODO: print license
	f.WriteLine("// GENERATED FILE: DO NOT EDIT!")
	f.WriteLine(``)
	f.WriteLine(`syntax = "proto3";`)
	renderDependencies(f, renderer)
	f.WriteLine(`package ` + *fileDescrProto.Package + `;`)
	f.WriteLine(``)

	renderService(f, renderer.currFileDescriptor.Service)
	renderMessages(f, renderer.currFileDescriptor.MessageType)
	renderEnums(f, renderer.currFileDescriptor.EnumType)

	return f.Bytes(), nil
}

func removePackageFromNames(renderer *Renderer) {
	numberOfCharsToRemove := utf8.RuneCountInString(*renderer.currFileDescriptor.Package) + 2

	for _, service := range renderer.currFileDescriptor.Service {
		if service != nil {
			for _, method := range service.Method {
				*method.InputType = (*method.InputType)[numberOfCharsToRemove:]
				*method.OutputType = (*method.OutputType)[numberOfCharsToRemove:]
			}
		}
	}

	for _, message := range renderer.currFileDescriptor.MessageType {
		for _, field := range message.Field {
			if field.TypeName != nil {
				*field.TypeName = (*field.TypeName)[numberOfCharsToRemove:]
			}
		}
	}

}

func renderDependencies(f *LineWriter, renderer *Renderer) {
	dependencies := renderer.currFileDescriptor.Dependency

	f.WriteLine(``)
	for _, dependency := range dependencies {
		f.WriteLine(`import "` + dependency + `";`)
		f.WriteLine(``)
	}
}

func renderService(f *LineWriter, services []*descriptor.ServiceDescriptorProto) {
	for _, service := range services {
		f.WriteLine(`service ` + *service.Name + ` {`)
		for _, method := range service.Method {
			renderRPCsignature(f, method)
			renderOptions(f, method.Options)
			f.WriteLine(`  }`) // Closing bracket of method
			f.WriteLine(``)
		}
		f.WriteLine(`}`) // Closing bracket of RPC service
		f.WriteLine(``)
	}
}

func renderRPCsignature(f *LineWriter, method *descriptor.MethodDescriptorProto) {
	if *method.InputType == "" {
		*method.InputType = "google.protobuf.Empty"
	}

	if *method.OutputType == "" {
		*method.OutputType = "google.protobuf.Empty"
	}

	f.WriteLine(`  rpc ` + *method.Name + ` (` + *method.InputType + `) ` + `returns` + ` (` + *method.OutputType + `) {`)
}

func renderOptions(f *LineWriter, options *descriptor.MethodOptions) {
	// TODO: Problem: We don't have any information about HTTP transcoding
}

func renderMessages(f *LineWriter, messages []*descriptor.DescriptorProto) {

	for _, message := range messages {
		f.WriteLine(`message ` + *message.Name + ` {`)
		renderFields(f, message.Field)
		renderEnums(f, message.EnumType)
		// Render nested messages
		renderMessages(f, message.NestedType)

		f.WriteLine(`}`)
		f.WriteLine(``)
	}

}

func renderEnums(f *LineWriter, enums []*descriptor.EnumDescriptorProto) {
	for _, enum := range enums {
		f.WriteLine(`enum ` + *enum.Name + ` {`)
		for _, value := range enum.Value {
			f.WriteLine(" " + *value.Name + " = " + strconv.Itoa(int(*value.Number)) + ";")
		}
		f.WriteLine(`}`)
		f.WriteLine(``)
	}
}

func renderFields(f *LineWriter, fields []*descriptor.FieldDescriptorProto) {
	for _, field := range fields {
		protobufType := descriptorTypeToProtoType[*field.Type]
		if protobufType == "" && field.TypeName != nil {
			protobufType = *field.TypeName
		}
		f.WriteLine(labelMapping[*field.Label] + " " + protobufType + " " + *field.Name + " = " + strconv.Itoa(int(*field.Number)) + `;`)
	}
}

func createMapping() map[descriptor.FieldDescriptorProto_Type]string {
	typeMapping := make(map[descriptor.FieldDescriptorProto_Type]string)
	typeMapping[descriptor.FieldDescriptorProto_TYPE_DOUBLE] = "double"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_FLOAT] = "float"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_INT64] = "int64"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_UINT64] = "uint64"

	typeMapping[descriptor.FieldDescriptorProto_TYPE_INT32] = "int32"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_FIXED64] = "fixed64"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_FIXED32] = "fixed32"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_BOOL] = "bool"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_STRING] = "string"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_MESSAGE] = ""
	typeMapping[descriptor.FieldDescriptorProto_TYPE_BYTES] = "bytes"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_UINT32] = "uint32"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_ENUM] = ""
	typeMapping[descriptor.FieldDescriptorProto_TYPE_SFIXED32] = "sfixed32"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_SFIXED64] = "sfixed64"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_SINT32] = "sint32"
	typeMapping[descriptor.FieldDescriptorProto_TYPE_SINT64] = "sint64"
	return typeMapping
}

func createLabelMapping() map[descriptor.FieldDescriptorProto_Label]string {
	labelMapping := make(map[descriptor.FieldDescriptorProto_Label]string)
	labelMapping[descriptor.FieldDescriptorProto_LABEL_OPTIONAL] = ""
	labelMapping[descriptor.FieldDescriptorProto_LABEL_REPEATED] = " repeated"
	labelMapping[descriptor.FieldDescriptorProto_LABEL_REQUIRED] = " required"
	return labelMapping
}

// WATCH OUT FOR:
// The path template may refer to one or more fields in the gRPC request message, as long
// as each field is a non-repeated field with a primitive (non-message) type.
// see: https://github.com/googleapis/googleapis/blob/a8ee1416f4c588f2ab92da72e7c1f588c784d3e6/google/api/http.proto#L62
// AND:
// Note that fields which are mapped to URL query parameters must have a
// primitive type or a repeated primitive type or a non-repeated message type.
// see: https://github.com/googleapis/googleapis/blob/a8ee1416f4c588f2ab92da72e7c1f588c784d3e6/google/api/http.proto#L119

//TODO: handle enum. Not sure if possible, because of
//TODO: https://github.com/googleapis/googleapis/blob/a8ee1416f4c588f2ab92da72e7c1f588c784d3e6/google/api/http.proto#L62

//TODO: Flatten URL Path parameters (query params don't need to be flattened!)

//TODO: Take a look a look at comments from Noah
