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

package generator

import (
	"github.com/golang/protobuf/proto"
	openapiv3 "github.com/googleapis/gnostic/OpenAPIv3"
	surface "github.com/googleapis/gnostic/surface"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

const (
	// When false, test behaves normally, checking output against golden test files.
	// But when changed to true, running test will actually re-generate golden test
	// files (which assumes output is correct).
	regenerateMode = false

	testFilesDirectory = "testfiles"
)

func TestFileDescriptorGeneratorParameters(t *testing.T) {
	input := "testfiles/parameters.yaml"

	protoData, err := runGeneratorWithoutEnvironment(input)
	if err != nil {
		handleError(err, t)
	}

	checkContents(t, string(protoData), "goldstandard/parameter.proto")
}

func TestFileDescriptorGeneratorRequestBodies(t *testing.T) {
	input := "testfiles/requestBodies.yaml"

	protoData, err := runGeneratorWithoutEnvironment(input)
	if err != nil {
		handleError(err, t)
	}

	checkContents(t, string(protoData), "goldstandard/requestbodies.proto")

}

func TestFileDescriptorGeneratorResponses(t *testing.T) {
	input := "testfiles/responses.yaml"

	protoData, err := runGeneratorWithoutEnvironment(input)
	if err != nil {
		handleError(err, t)
	}
	checkContents(t, string(protoData), "goldstandard/responses.proto")
}

func TestFileDescriptorGeneratorOther(t *testing.T) {
	input := "testfiles/other.yaml"

	protoData, err := runGeneratorWithoutEnvironment(input)
	if err != nil {
		handleError(err, t)
	}
	checkContents(t, string(protoData), "goldstandard/other.proto")
}

func runGeneratorWithoutEnvironment(input string) ([]byte, error) {
	surfaceModel := buildSurfaceModel(input)
	r := NewRenderer(surfaceModel)
	r.Package = "testPackage"

	fdSet, err := r.RunFileDescriptorSetGenerator()
	r.FdSet = fdSet
	if err != nil {
		return nil, err
	}

	f, err := r.RenderProto("")
	if err != nil {
		return nil, err
	}
	return f.Data, err
}

func buildSurfaceModel(input string) *surface.Model {
	cmd := exec.Command("gnostic", "--pb-out=-", input)
	b, _ := cmd.Output()
	documentv3, _ := createOpenAPIdocFromGnosticOutput(b)
	surfaceModel, _ := surface.NewModelFromOpenAPI3(documentv3, input)
	return surfaceModel
}

func writeFile(output string, protoData []byte) {
	dir := path.Dir(output)
	os.MkdirAll(dir, 0755)
	f, _ := os.Create(output)
	defer f.Close()
	f.Write(protoData)
}

func checkContents(t *testing.T, actualContents string, goldenFileName string) {
	goldenFileName = filepath.Join(testFilesDirectory, goldenFileName)

	if regenerateMode {
		writeFile(goldenFileName, []byte(actualContents))
	}

	// verify that output matches golden test files
	b, err := ioutil.ReadFile(goldenFileName)
	if err != nil {
		t.Errorf("Error while reading goldstandard file")
		t.Errorf(err.Error())
	}
	goldstandard := string(b)
	if goldstandard != actualContents {
		t.Errorf("File contents does not match.")
	}
}

func createOpenAPIdocFromGnosticOutput(binaryInput []byte) (*openapiv3.Document, error) {
	document := &openapiv3.Document{}
	err := proto.Unmarshal(binaryInput, document)
	if err != nil {
		// If we execute gnostic with argument: '-pb-out=-' we get an EOF
		if err.Error() != "unexpected EOF" {
			return nil, err
		}
	}
	return document, nil
}

func handleError(err error, t *testing.T) {
	t.Errorf("Error while executing the protoc-generator")
	if strings.Contains(err.Error(), "included an unresolvable reference") {
		t.Errorf("This could be due to the fact that 'typeName' is set wrong on a FieldDescriptorProto." +
			"For every FieldDescriptorProto where the type == 'FieldDescriptorProto_TYPE_MESSAGE' the correct typeName needs to be set.")
	}
	t.Errorf(err.Error())
}

// Sometimes I need
//func buildFdsetFromProto() {
//	b, err := ioutil.ReadFile("temp.descr")
//	if err != nil {
//		fmt.Print(err.Error())
//	}
//	fdSet := &descriptor.FileDescriptorSet{}
//	err = proto.Unmarshal(b, fdSet)
//	if err != nil {
//		fmt.Print(err.Error())
//	}
//}
