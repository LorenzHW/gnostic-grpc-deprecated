### End-to-end flow

This example demonstrates an end-to-end flow for generating a gRPC API with HTTP transcoding from an
OpenAPI description.


#### What we will build:
![alt text](https://camo.githubusercontent.com/e75a8b46b078a3c1df0ed9966a16c24add9ccb83/68747470733a2f2f646f63732e676f6f676c652e636f6d2f64726177696e67732f642f3132687034435071724e5046686174744c5f63496f4a707446766c41716d35774c513067677149356d6b43672f7075623f773d37343926683d333730 "gRPC with Transcoding")

This tutorial has six steps:

1. Generate a gRPC service (.proto) from an OpenAPI description.
2. Generate server-side support code for the gRPC service.
3. Implement the server logic.
4. Set up a proxy that provides HTTP transcoding.
5. Run the proxy and the server.
6. Test your API with with curl and a gRPC client.

#### Prerequisite
Install [gnostic](https://github.com/googleapis/gnostic), [gnostic-protoc-generator](https://github.com/LorenzHW/gnostic-protoc-generator),
[go plugin for protoc](https://github.com/golang/protobuf/protoc-gen-go), [gRPC gateway plugin](https://github.com/grpc-ecosystem/grpc-gateway)
and [gRPC](https://grpc.io/)

    go get -u github.com/googleapis/gnostic
    go get -u github.com/LorenzHW/gnostic-protoc-generator
    go get -u github.com/golang/protobuf/protoc-gen-go
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
    go get -u google.golang.org/grpc
    
For simplicity lets create a temporary environment variable inside your terminal:
    
    export ANNOTATIONS="third-party/googleapis"
    
In order for this tutorial to work you should work inside this directory under `GOPATH`.

#### 1. Step

Use [gnostic](https://github.com/googleapis/gnostic) to generate the Protocol buffer 
description (`bookstore.proto`) in the current directory:

    gnostic --protoc-generator-out=. bookstore.yaml

#### 2. Step
Generate the gRPC stubs:
    
    protoc --proto_path=. --proto_path=${ANNOTATIONS} --go_out=plugins=grpc:bookstore bookstore.proto
    
 This generates `bookstore/bookstore.pb.go`.

#### 3. Step
We added an example implementation of the server using the generated gRPC stubs inside `bookstore/server.go`.
    
#### 4. Step
Generate the reverse proxy with the gRPC gateway plugin:

    protoc --proto_path=. --proto_path=${ANNOTATIONS} --grpc-gateway_out=bookstore bookstore.proto
    
We provided a sample implementation on how to use the proxy inside `bookstore/proxy.go`.

#### 5. Step
Start the proxy and the server:

    go run main.go
    
    
    
#### 6. Step

##### curl
Inside of a new terminal test your API:

Let's create a shelf first:

    curl -X POST \
      http://localhost:8081/shelves \
      -H 'Content-Type: application/json' \
      -d '{
        "name": "Books I need to read",
        "theme": "Non-fiction"
    }'
    
Get all existing shelves:

    curl -X GET http://localhost:8081/shelves
    
Create a book for the shelve with the id `1`:
    
    curl -X POST \
      http://localhost:8081/shelves/1/books \
      -H 'Content-Type: application/json' \
      -d '{
        "author": "Hans Rosling",
        "name": "Factfulness",
        "title": "Factfulness: Ten Reasons We'\''re wrong about the world - and Why Things Are Better Than You Think"
    }'
    
    
List all books for the shelve with the id `1`:

    curl -X GET http://localhost:8081/shelves/1/books
    
    
##### gRPC client

A sample gRPC client is provided inside `grpc-client/client.go` that lists all themes of your shelves:

    go run grpc-client/client.go