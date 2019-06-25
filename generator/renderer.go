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
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugins "github.com/googleapis/gnostic/plugins"
	surface "github.com/googleapis/gnostic/surface"
	prDesc "github.com/jhump/protoreflect/desc"
	prPrint "github.com/jhump/protoreflect/desc/protoprint"
	_ "os"
)

// Renderer generates code for a surface.Model.
type Renderer struct {
	// TODO: Maybe create a util/generic package(?),
	// TODO: because same struct is used in gnostic-go-generator
	Model   *surface.Model
	Package string // package name
}

// NewRenderer creates a renderer.
func NewRenderer(model *surface.Model) (renderer *Renderer, err error) {
	renderer = &Renderer{}
	renderer.Model = model
	return renderer, nil
}

// Generate runs the renderer to generate the named files.
func (renderer *Renderer) Render(response *plugins.Response, fileName string) (err error) {
	file := &plugins.File{Name: fileName}
	fdSet, err := renderer.RunFileDescriptorSetGenerator()

	if err != nil {
		return err
	}

	if false { //TODO: If we want to generate the descriptor file, we need an additional flag here!
		f, err := renderer.RenderDescriptor(fdSet)
		if err != nil {
			return err
		}
		response.Files = append(response.Files, f)
	}

	file.Data, err = renderer.RenderProto(fdSet)
	response.Files = append(response.Files, file)

	return
}

func (renderer *Renderer) RenderProto(fdSet *dpb.FileDescriptorSet) ([]byte, error) {

	// Creates a protoreflect FileDescriptor, which is then used for printing.
	prFd, err := prDesc.CreateFileDescriptorFromSet(fdSet)
	if err != nil {
		return nil, err
	}

	// Print the protoreflect FileDescriptor.
	p := prPrint.Printer{}
	res, err := p.PrintProtoToString(prFd)
	if err != nil {
		return nil, err
	}

	f := NewLineWriter()
	f.WriteLine(res)

	return f.Bytes(), err
}

func (renderer *Renderer) RenderDescriptor(fdSet *dpb.FileDescriptorSet) (*plugins.File, error) {
	fdSetData, err := proto.Marshal(fdSet)
	if err != nil {
		return nil, err
	}

	descriptorFile := &plugins.File{Name: renderer.Package + ".descr"}
	descriptorFile.Data = fdSetData
	return descriptorFile, nil
}
