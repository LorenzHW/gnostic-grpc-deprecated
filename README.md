[![Build Status](https://travis-ci.com/LorenzHW/gnostic-grpc.svg?branch=master)](https://travis-ci.com/LorenzHW/gnostic-grpc)

# gnostic gRPC plugin
[GSoC 2019 project](https://summerofcode.withgoogle.com/projects/#5244822191865856)

This tool converts an OpenAPI v3.0 API description into an equivalent .proto representation.

## High level overview:
![alt text](https://drive.google.com/uc?export=view&id=1tqDvZLiXK40ISK_LgINQGsno9-MymRQP "High Level Overview")

Under the hood the plugin first creates a FileDescriptorSet (`bookststore.descr`) from the input
data. Then [protoreflect](https://github.com/jhump/protoreflect/) is used to print the output file. 

## How to use:    
Install the plugin:

    go get -u github.com/LorenzHW/gnostic-grpc

Run the plugin:

    gnostic --grpc-out=examples/bookstore examples/bookstore/bookstore.yaml

This generates a protocol buffer definition `examples/bookstore/bookstore.proto`.

## End-to-end example
The directory `examples/end-to-end` contains a tutorial on how to build a gRPC API with an OpenAPI specification.

## What conversions are currently supported?

Given an [OpenAPI object](https://swagger.io/specification/#oasObject) following fields will be
represented inside a .proto file:

| Object        | Fields        | Supported  |
| ------------- |:-------------:| -----:|
| OpenAPI object|               |       |
|               | openapi       |    No |
|               | info          |    No |
|               | servers       |    No |
|               | paths         |   Yes |
|               | components    |   Yes |
|               | security      |    No |
|               | tags          |    No |
|               | externalDocs  |    No |