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
	input := "testfiles/parameters/test.pb"
	output := "../protoc-generator/testfiles/input/parameter.descr"

	fileDescriptorData, err := runDescriptorGeneratorWithoutEnv(input)
	if err != nil {
		t.Errorf("Error while executing the descriptor generator")
		t.Errorf(err.Error())
	}
	writeFile(output, fileDescriptorData)

	erroneousInput := []string{"testfiles/parameters/errors/invalid_path_param.pb", "testfiles/parameters/errors/invalid_query_param.pb"}

	for _, errorInput := range erroneousInput {
		errorMessages := map[string]bool{
			"The path parameter with the Name param1 is invalid. The path template may refer to one or more fields in the gRPC request message, as long as each field is a non-repeated field with a primitive (non-message) type": true,
			"The query parameter with the Name param1 is invalid. Note that fields which are mapped to URL query parameters must have a primitive type or a repeated primitive type or a non-repeated message type.":               true,
		}
		fileDescriptorData, err = runDescriptorGeneratorWithoutEnv(errorInput)
		if _, ok := errorMessages[err.Error()]; !ok {
			t.Errorf("Error while executing the descriptor generator")
			t.Errorf(err.Error())
		}
	}

}

func TestFileDescriptorGeneratorRequestBodies(t *testing.T) {
	input := "testfiles/requestBodies/test.pb"
	output := "../protoc-generator/testfiles/input/requestbodies.descr"

	fileDescriptorData, err := runDescriptorGeneratorWithoutEnv(input)
	if err != nil {
		t.Errorf("Error while executing the protoc-generator")
		t.Errorf(err.Error())
	}
	writeFile(output, fileDescriptorData)

}

func TestFileDescriptorGeneratorResponses(t *testing.T) {
	input := "testfiles/responses/test.pb"
	output := "../protoc-generator/testfiles/input/responses.descr"

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
	fileDescriptorSetData, err := descriptorRenderer.BuildFileDescriptorSet()
	return fileDescriptorSetData, err
}

func writeFile(output string, protoData []byte) {
	dir := path.Dir(output)
	os.MkdirAll(dir, 0755)
	f, _ := os.Create(output)
	defer f.Close()
	f.Write(protoData)
}
