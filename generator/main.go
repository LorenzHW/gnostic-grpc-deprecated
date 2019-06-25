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

// descriptor_renderer generates a FileDescriptorSet from a gnostic output file.
package generator

import (
	"errors"
	"go/format"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	plugins "github.com/googleapis/gnostic/plugins"
	surface "github.com/googleapis/gnostic/surface"
)

// This is the main function for the code generation plugin.
func RunDescriptorGenerator() {
	env, err := plugins.NewEnvironment()
	env.RespondAndExitIfError(err)

	fileName := env.Request.SourceName
	extension := filepath.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(extension)]

	packageName, err := resolvePackageName(fileName)
	env.RespondAndExitIfError(err)

	for _, model := range env.Request.Models {
		switch model.TypeUrl {
		case "surface.v1.Model":
			surfaceModel := &surface.Model{}
			err = proto.Unmarshal(model.Value, surfaceModel)
			if err == nil {
				// Create the renderer.
				renderer, err := NewDescriptorRenderer(surfaceModel)
				renderer.Package = packageName
				env.RespondAndExitIfError(err)

				// Run the renderer to generate files and add them to the response object.
				err = renderer.Render(env.Response, packageName+".proto")
				env.RespondAndExitIfError(err)

				// Return with success.
				env.RespondAndExit()
			}
		}
	}
	err = errors.New("No generated code surface model is available.")
	env.RespondAndExitIfError(err)
}

// resolvePackageName converts a path to a valid package name or
// error if path can't be resolved or resolves to an invalid package name.
func resolvePackageName(p string) (string, error) {
	p, err := filepath.Abs(p)
	if err == nil {
		p = filepath.Base(p)
		_, err = format.Source([]byte("package " + p))
	}
	if err != nil {
		return "", errors.New("invalid package name " + p)
	}
	return p, nil
}