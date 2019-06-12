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

// protoc_generator converts a FileDescriptorSet into protobuf specification (.protos)
package protoc_generator

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// This is the main function for the code generation plugin.
func RunProtocGenerator() {
	env, err := NewEnvironment()
	env.RespondAndExitIfError(err)

	for _, model := range env.Request.Models {
		switch model.TypeUrl {
		case "descriptor.set.Model":
			fileDescriptorSetModel := &descriptor.FileDescriptorSet{}
			err = proto.Unmarshal(model.Value, fileDescriptorSetModel)
			if err == nil {
				// Create the renderer.
				renderer, err := NewProtoRenderer(fileDescriptorSetModel)
				// Run the renderer to generate files and add them to the response object.
				err = renderer.Render(env.Response)
				env.RespondAndExitIfError(err)

				// Return with success.
				env.RespondAndExit()

			}
		}
	}

	err = errors.New("no Model with the FileDescriptorSet data")
	env.RespondAndExitIfError(err)
}
