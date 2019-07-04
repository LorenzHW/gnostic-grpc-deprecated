package generator

import (
	"github.com/golang/protobuf/proto"
	openapiv3 "github.com/googleapis/gnostic/OpenAPIv3"
	surface "github.com/googleapis/gnostic/surface"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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
	input := "testfiles/parameters/test.pb"

	protoData, err := runGeneratorWithoutEnvironment(input)
	if err != nil {
		t.Errorf("Error while executing the descriptor generator")
		t.Errorf(err.Error())
	}

	checkContents(t, string(protoData), "goldstandard/parameter.proto")

	erroneousInput := []string{"testfiles/parameters/errors/invalid_path_param.pb", "testfiles/parameters/errors/invalid_query_param.pb"}

	for _, errorInput := range erroneousInput {
		errorMessages := map[string]bool{
			"The path parameter with the Name param1 is invalid. The path template may refer to one or more fields in the gRPC request message, as long as each field is a non-repeated field with a primitive (non-message) type": true,
			"The query parameter with the Name param1 is invalid. Note that fields which are mapped to URL query parameters must have a primitive type or a repeated primitive type or a non-repeated message type.":               true,
		}
		protoData, err = runGeneratorWithoutEnvironment(errorInput)
		if _, ok := errorMessages[err.Error()]; !ok {
			// If we don't get an error from the generator the test fails!
			t.Errorf("Error while executing the descriptor generator")
			t.Errorf(err.Error())
		}
	}

}

func TestFileDescriptorGeneratorRequestBodies(t *testing.T) {
	input := "testfiles/requestBodies/test.pb"

	protoData, err := runGeneratorWithoutEnvironment(input)
	if err != nil {
		t.Errorf("Error while executing the protoc-generator")
		t.Errorf(err.Error())
	}

	checkContents(t, string(protoData), "goldstandard/requestbodies.proto")

}

func TestFileDescriptorGeneratorResponses(t *testing.T) {
	input := "testfiles/responses/test.pb"

	protoData, err := runGeneratorWithoutEnvironment(input)
	if err != nil {
		t.Errorf("Error while executing the protoc-generator")
		t.Errorf(err.Error())
	}
	checkContents(t, string(protoData), "goldstandard/responses.proto")
}

func runGeneratorWithoutEnvironment(input string) ([]byte, error) {
	surfaceModel := buildSurfaceModel(input)
	r, _ := NewRenderer(surfaceModel)
	r.Package = "testPackage"

	fdSet, err := r.RunFileDescriptorSetGenerator()
	r.FdSet = fdSet
	if err != nil {
		return nil, err
	}
	f, err := r.RenderProto("")
	return f.Data, err
}

func buildSurfaceModel(input string) *surface.Model {
	apiData, _ := ioutil.ReadFile(input)
	documentv3 := &openapiv3.Document{}
	proto.Unmarshal(apiData, documentv3)
	surfaceModel, _ := surface.NewModelFromOpenAPI3(documentv3)
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
