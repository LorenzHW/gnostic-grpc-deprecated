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

package main

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugins "github.com/googleapis/gnostic/plugins"
	_ "os"
)

// Renderer generates code for a FileDescriptorProto.
type Renderer struct {
	fileDescriptorSet  descriptor.FileDescriptorSet
	currFileDescriptor *descriptor.FileDescriptorProto
}

// NewServiceRenderer creates a renderer.
func NewServiceRenderer(fileDescriptorSet *descriptor.FileDescriptorSet) (renderer *Renderer, err error) {
	renderer = &Renderer{}
	renderer.fileDescriptorSet = *fileDescriptorSet
	return renderer, nil
}

// Generate runs the renderer to generate the named files.
func (renderer *Renderer) Render(response *plugins.Response) (err error) {

	for _, fileDescriptor := range renderer.fileDescriptorSet.File {
		file := &plugins.File{}
		file.Name = *fileDescriptor.Name
		file.Data, err = renderer.RenderProto(fileDescriptor)
		if err != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("ERROR %v", err))
		}
		response.Files = append(response.Files, file)
	}
	return

}
