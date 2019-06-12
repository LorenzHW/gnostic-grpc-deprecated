package protoc_generator

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugins "github.com/googleapis/gnostic/plugins"
	"io/ioutil"
	"os"
	"path"
)

// TODO: This could be merged into gnostic/plugins. Basically only adapted it for FileDescriptorSet
func NewEnvironment() (env *plugins.Environment, err error) {
	env = &plugins.Environment{
		Invocation: os.Args[0],
		Response:   &plugins.Response{},
	}

	input := flag.String("input", "", "API description (in binary protocol buffer form)")
	output := flag.String("output", "-", "Output file or directory")
	flag.Parse()

	apiData, err := ioutil.ReadFile(*input)
	if len(apiData) == 0 {
		env.RespondAndExitIfError(fmt.Errorf("no input data"))
	}
	env.Request = &plugins.Request{}
	env.Request.OutputPath = *output
	env.Request.SourceName = path.Base(*input)

	fileDescriptorSet := &descriptor.FileDescriptorSet{}
	err = proto.Unmarshal(apiData, fileDescriptorSet)

	if err == nil {
		env.Request.AddModel("descriptor.set.Model", fileDescriptorSet)
		return env, err
	}
	// If we get here, we don't know what we got
	err = errors.New("Unrecognized format for input")
	return env, err

}
