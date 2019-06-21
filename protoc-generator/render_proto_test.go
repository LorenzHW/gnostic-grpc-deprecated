package protoc_generator

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
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

func TestProtoGeneratorParameters(t *testing.T) {
	input := "testfiles/input/parameter.descr"

	protoData, err := RunProtoGeneratorWithoutEnv(input)
	if err != nil {
		t.Errorf("Error while executing the protoc-generator")
		t.Errorf(err.Error())
	}
	checkContents(t, string(protoData), "goldstandard/parameter.proto")

}

func TestProtoGeneratorRequestBodies(t *testing.T) {
	input := "testfiles/input/requestbodies.descr"

	protoData, err := RunProtoGeneratorWithoutEnv(input)
	if err != nil {
		t.Errorf("Error while executing the protoc-generator")
		t.Errorf(err.Error())
	}
	checkContents(t, string(protoData), "goldstandard/requestbodies.proto")

}

func TestRunProtocGeneratorResponses(t *testing.T) {
	input := "testfiles/input/responses.descr"

	protoData, err := RunProtoGeneratorWithoutEnv(input)
	if err != nil {
		t.Errorf("Error while executing the protoc-generator")
		t.Errorf(err.Error())
	}
	checkContents(t, string(protoData), "goldstandard/responses.proto")

}

func RunProtoGeneratorWithoutEnv(input string) ([]byte, error) {
	fileDescriptorSetData, _ := ioutil.ReadFile(input)
	fileDescr := &descriptor.FileDescriptorSet{}
	proto.Unmarshal(fileDescriptorSetData, fileDescr)
	renderer, err := NewProtoRenderer(fileDescr)
	protoData, err := renderer.RenderProto(fileDescr)
	return protoData, err
}

func WriteFile(output string, protoData []byte) {
	dir := path.Dir(output)
	os.MkdirAll(dir, 0755)
	f, _ := os.Create(output)
	defer f.Close()
	f.Write(protoData)
}

func checkContents(t *testing.T, actualContents string, goldenFileName string) {
	goldenFileName = filepath.Join(testFilesDirectory, goldenFileName)

	//if regenerateMode {
	//	err := ioutil.WriteFile(goldenFileName, []byte(actualContents), 0666)
	//	Ok(t, err)
	//}

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
