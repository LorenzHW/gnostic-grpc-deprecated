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
