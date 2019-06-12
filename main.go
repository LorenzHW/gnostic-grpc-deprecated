package main

import (
	descriptor_generator "github.com/LorenzHW/gnostic-protoc-generator/descriptor-generator"
	protoc_generator "github.com/LorenzHW/gnostic-protoc-generator/protoc-generator"
	"os"
	"strings"
)

func main() {
	inputPath := os.Args[1:][1]

	generatorType := ""
	// This is not done very well: Trying to figure
	// out which generator to run based on the input
	// but things will change anyway probably. So let's
	// keep it dirty for the moment.
	if strings.Contains(inputPath, ".descr") {
		generatorType = "proto-generator"
	} else if strings.Contains(inputPath, ".pb") {
		generatorType = "descriptor-generator"
	}

	if generatorType == "descriptor-generator" {
		descriptor_generator.RunDescriptorGenerator()
	} else if generatorType == "proto-generator" {
		protoc_generator.RunProtocGenerator()
	}

}
