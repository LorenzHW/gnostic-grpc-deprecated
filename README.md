[![Build Status](https://travis-ci.com/LorenzHW/gnostic-protoc-generator.svg?branch=master)](https://travis-ci.com/LorenzHW/gnostic-protoc-generator)

# gnostic Protos Generator Plugin
[GSoC 2019 project](https://summerofcode.withgoogle.com/projects/#5244822191865856)

This tool converts an OpenAPI v3.0 API description into an equivalent .proto representation.

## High level overview:
![alt text](https://drive.google.com/uc?export=view&id=1tqDvZLiXK40ISK_LgINQGsno9-MymRQP "High Level Overview")

## Prerequisite:
Use [gnostic](https://github.com/googleapis/gnostic) to generate `examples/bookstore/input/bookstore.pb`
by running following command inside this directory:
    
    gnostic --pb-out=examples/bookstore/input examples/bookstore/input/bookstore.yaml

## How to use:

To run this plugin run following commands inside this directory:

    go build
    
To run the descriptor generator:
    
    ./gnostic-protoc-generator -input examples/bookstore/input/bookstore.pb -output examples/bookstore/output/

This command triggers the descriptor-generator. The generator takes in a binary format of an OpenAPI
specification (`examples/bookstore/input/bookstore.pb` created with gnostic) and creates a file descriptor set
`examples/bookstore/output/bookstore.descr`.


To run the protoc generator:
 
    ./gnostic-protoc-generator -input examples/bookstore/output/bookstore.descr -output examples/bookstore/output/

This command triggers the proto-generator. This generator takes in a file descriptor set
`examples/bookstore/output/bookstore.descr` and generates a protocol buffer definition (`examples/bookstore/output/bookstore.proto`)


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