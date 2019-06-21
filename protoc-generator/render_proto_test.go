package protoc_generator

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestProtoGenerator(t *testing.T) {
	input := "test_data/parameters/test.descr"
	output := "test_data/parameters/test.proto"

	protoData, err := RunProtoGeneratorWithoutEnv(input)
	if err != nil {
		t.Errorf("Error while executing the protoc-generator")
		t.Errorf(err.Error())
	}
	WriteFile(output, protoData)

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
