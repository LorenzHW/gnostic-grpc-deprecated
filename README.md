[![Build Status](https://travis-ci.com/LorenzHW/gnostic-protoc-generator.svg?branch=master)](https://travis-ci.com/LorenzHW/gnostic-protoc-generator)

# gnostic Protos Generator Plugin
[GSoC 2019 project](https://summerofcode.withgoogle.com/projects/#5244822191865856)

This tool converts an OpenAPI v3.0 API description into an equivalent .proto representation.

## High level overview:
![alt text](https://drive.google.com/uc?export=view&id=1tqDvZLiXK40ISK_LgINQGsno9-MymRQP "High Level Overview")

Under the hood the generator first creates a FileDescriptorSet (`bookststore.descr`) from the input
data. Then [protoreflect](https://github.com/jhump/protoreflect/) is used to print the output file. 

## Prerequisite:
Use [gnostic](https://github.com/googleapis/gnostic) to generate `examples/bookstore/bookstore.pb`
by running following command inside this directory:
    
    gnostic --pb-out=examples/bookstore examples/bookstore/bookstore.yaml

## How to use:

To run this plugin run following commands inside this directory:

    go build
    
To run the generator:
    
    ./gnostic-protoc-generator -input examples/bookstore/bookstore.pb -output examples/bookstore

The generator takes in a binary format of an OpenAPI specification 
(`examples/bookstore/bookstore.pb` created with gnostic) and generates a protocol buffer definition
`examples/bookstore/bookstore.proto`.

## What conversions are currently supported?

Given an [OpenAPI object](https://swagger.io/specification/#oasObject) following fields will be
represented inside a .proto file

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