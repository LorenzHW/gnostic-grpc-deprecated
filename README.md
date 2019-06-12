[![Build Status](https://travis-ci.com/LorenzHW/gnostic-protoc-generator.svg?branch=master)](https://travis-ci.com/LorenzHW/gnostic-protoc-generator)

# gnostic Protos Generator Plugin
[GSoC 2019 project](https://summerofcode.withgoogle.com/projects/#5244822191865856)

To run this plugin run following commands inside this directory:

    go build
    
To run the descriptor generator:
    
    ./gnostic-protoc-generator -input example/input/test.pb -output example/output/

This generator takes in a binary format of an OpenAPI specification (`texample/input/test.pb`
created with gnostic) and creates a file descriptor set `example/output/output.descr`.


To run the protoc generator:
 
    ./gnostic-protoc-generator -input example/output/output.descr -output example/output/

This generator takes in a file descriptor set `example/output/output.descr` and generates a
protocol buffer definition (`example/output/output.proto`)