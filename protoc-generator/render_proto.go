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
	"github.com/golang/protobuf/descriptor"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/ptypes/empty"
	prDesc "github.com/jhump/protoreflect/desc"
	prPrint "github.com/jhump/protoreflect/desc/protoprint"
)

func (renderer *Renderer) RenderProto(fdSet *dpb.FileDescriptorSet) ([]byte, error) {

	buildDependenciesForProtoReflect(fdSet)

	f := NewLineWriter()
	p := prPrint.Printer{}
	prFd, err := prDesc.CreateFileDescriptorFromSet(fdSet)
	//prFd.GetOptions()
	//prMethodOptions := prFd.GetServices()[0].GetMethods()[0].GetOptions()
	//prMethodOptions2 := prFd.GetServices()[0].GetMethods()[0].GetMethodOptions()
	//if prMethodOptions != nil && prMethodOptions2 != nil {
	//	eHttp, _ := proto.GetExtension(prMethodOptions, annotations.E_Http)
	//	if eHttp != nil {
	//
	//	}
	//}
	res, err := p.PrintProtoToString(prFd)

	f.WriteLine(res)

	return f.Bytes(), err
}

func buildDependenciesForProtoReflect(fdSet *dpb.FileDescriptorSet) {
	empt := empty.Empty{}
	fd, _ := descriptor.ForMessage(&empt)

	// According to the documentation of prDesc.CreateFileDescriptorFromSet
	// the file I want to print needs to be at the end of the array.
	fdSet.File = append([]*dpb.FileDescriptorProto{fd}, fdSet.File...)

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

//TODO: Flatten URL Path parameters (query params don't need to be flattened!) (This could actually be done inside the descriptor-generator): So if: RPC request param && inside path && NOT_SCALAR -- > flatten

//TODO: Take a look at this reflect package noah mentioned (to generate protos from filedescriptor input)
//TODO: Take a look at body parameter
