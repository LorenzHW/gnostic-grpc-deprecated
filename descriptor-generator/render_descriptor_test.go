package descriptor_generator

import (
	"github.com/golang/protobuf/proto"
	openapiv3 "github.com/googleapis/gnostic/OpenAPIv3"
	surface "github.com/googleapis/gnostic/surface"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestFileDescriptorGeneratorParameters(t *testing.T) {
	input := "test_data/parameters/input/test.pb"
	output := "test_data/parameters/output/test.descr"

	fileDescriptorData, err := runDescriptorGeneratorWithoutEnv(input)
	if err != nil {
		t.Errorf("Error while executing the protoc-generator")
		t.Errorf(err.Error())
	}
	writeFile(output, fileDescriptorData)

}

func TestFileDescriptorGeneratorRequestBodies(t *testing.T) {
	input := "test_data/requestBodies/input/test.pb"
	output := "test_data/requestBodies/output/test.descr"

	fileDescriptorData, err := runDescriptorGeneratorWithoutEnv(input)
	if err != nil {
		t.Errorf("Error while executing the protoc-generator")
		t.Errorf(err.Error())
	}
	writeFile(output, fileDescriptorData)

}

func TestFileDescriptorGeneratorResponses(t *testing.T) {
	input := "test_data/responses/input/test.pb"
	output := "test_data/responses/output/test.descr"

	fileDescriptorData, err := runDescriptorGeneratorWithoutEnv(input)
	if err != nil {
		t.Errorf("Error while executing the protoc-generator")
		t.Errorf(err.Error())
	}
	writeFile(output, fileDescriptorData)

}

func buildSurfaceModel(input string) *surface.Model {
	apiData, _ := ioutil.ReadFile(input)
	documentv3 := &openapiv3.Document{}
	proto.Unmarshal(apiData, documentv3)
	surfaceModel, _ := surface.NewModelFromOpenAPI3(documentv3)
	return surfaceModel
}

func runDescriptorGeneratorWithoutEnv(input string) ([]byte, error) {
	packageName := "testPackage"

	surfaceModel := buildSurfaceModel(input)
	descriptorRenderer, _ := NewDescriptorRenderer(surfaceModel)
	descriptorRenderer.Package = packageName
	fileDescriptorSetData, err := descriptorRenderer.RenderFileDescriptorSet()
	return fileDescriptorSetData, err
}

func writeFile(output string, protoData []byte) {
	dir := path.Dir(output)
	os.MkdirAll(dir, 0755)
	f, _ := os.Create(output)
	defer f.Close()
	f.Write(protoData)
}
