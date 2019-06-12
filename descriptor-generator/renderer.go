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
	"fmt"
	plugins "github.com/googleapis/gnostic/plugins"
	surface "github.com/googleapis/gnostic/surface"
	_ "os"
)

// Renderer generates code for a surface.Model.
type Renderer struct {
	// TODO: Maybe create a util/generic package(?),
	// TODO: because same struct is used in gnostic-go-generator
	Model   *surface.Model
	Package string // package name
}

// NewDescriptorRenderer creates a renderer.
func NewDescriptorRenderer(model *surface.Model) (renderer *Renderer, err error) {
	renderer = &Renderer{}
	renderer.Model = model
	return renderer, nil
}

// Generate runs the renderer to generate the named files.
func (renderer *Renderer) Render(response *plugins.Response, files []string) (err error) {
	for _, filename := range files {
		file := &plugins.File{Name: filename}
		file.Data, err = renderer.RenderFileDescriptorSet()
		if err != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("ERROR %v", err))
		}
		response.Files = append(response.Files, file)
	}
	return
}
