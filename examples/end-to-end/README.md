# End-to-examples
This directory contains end-to-end flows for generating a gRPC API with HTTP transcoding from an OpenAPI description.

The first example uses the  [gRPC gateway plugin](https://github.com/grpc-ecosystem/grpc-gateway) for the proxy.

The second example uses [envoy](https://www.envoyproxy.io/) for the proxy.

## End-to-end flow with gRPC gateway plugin

This example demonstrates an end-to-end flow for generating a gRPC API with HTTP transcoding from an
OpenAPI description.


#### What we will build:

![alt text](https://drive.google.com/uc?export=view&id=118eI8Tb88gJF47nclHbLxqOS1N_ygt4o "gRPC with Transcoding")

This tutorial has six steps:

1. Generate a gRPC service (.proto) from an OpenAPI description.
2. Generate server-side support code for the gRPC service.
3. Implement the server logic.
4. Set up a proxy that provides HTTP transcoding.
5. Run the proxy and the server.
6. Test your API with with curl and a gRPC client.

#### Prerequisite
Install [gnostic](https://github.com/googleapis/gnostic), [gnostic-grpc](https://github.com/LorenzHW/gnostic-grpc),
[go plugin for protoc](https://github.com/golang/protobuf/protoc-gen-go), [gRPC gateway plugin](https://github.com/grpc-ecosystem/grpc-gateway)
and [gRPC](https://grpc.io/)

    go get -u github.com/googleapis/gnostic
    go get -u github.com/LorenzHW/gnostic-grpc
    go get -u github.com/golang/protobuf/protoc-gen-go
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
    go get -u google.golang.org/grpc
    
For simplicity lets create a temporary environment variable inside your terminal:
    
    export ANNOTATIONS="third-party/googleapis"
    
In order for this tutorial to work you should work inside this directory under `GOPATH`.

#### 1. Step

Use [gnostic](https://github.com/googleapis/gnostic) to generate the Protocol buffer 
description (`bookstore.proto`) in the current directory:

    gnostic --grpc-out=. bookstore.yaml

#### 2. Step
Generate the gRPC stubs:
    
    protoc --proto_path=. --proto_path=${ANNOTATIONS} --go_out=plugins=grpc:bookstore bookstore.proto
    
 This generates `bookstore/bookstore.pb.go`.

#### 3. Step
We added an example implementation of the server using the generated gRPC stubs inside `bookstore/server.go`.
    
#### 4. Step
Generate the reverse proxy with the gRPC gateway plugin:

    protoc --proto_path=. --proto_path=${ANNOTATIONS} --grpc-gateway_out=bookstore bookstore.proto

This generates `bookstore/bookstore.pb.gw.go`.

We provided a sample implementation on how to use the proxy inside `bookstore/proxy.go`.

#### 5. Step
Start the proxy and the server:

    go run main.go
    
#### 6. Step

##### cURL
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
    

## End-to-end flow with envoy
Follow steps 1-3 from the other tutorial.

#### 4. Step
Given `bookstore.proto` generate the descriptor set.
    
    protoc --proto_path=${ANNOTATIONS} --proto_path=. --include_imports --include_source_info \
    --descriptor_set_out=envoy-proxy/proto.pb bookstore.proto
    
This generates `envoy-proxy/proto.pb`.

#### 5. Step 
The file `envoy-proxy/envoy.yaml` contains an envoy configuration with a gRPC-JSON [transcoder](https://www.envoyproxy.io/docs/envoy/latest/configuration/http_filters/grpc_json_transcoder_filter).
According to the configuration, port 51051 proxies gRPC requests to a gRPC server running on localhost:50051 and uses 
the gRPC-JSON transcoder filter to provide the RESTful JSON mapping. I.e.: you can either make gRPC or RESTful JSON 
requests to localhost:51051.
  
Get the envoy docker image:

    docker pull envoyproxy/envoy-dev:bcc66c6b74c365d1d2834cfe15b847ae13be0eb6  
  
The file `envoy-proxy/Dockerfile` uses the envoy image we just pulled as base image and copies `envoy.yaml`
and `proto.pb` to the filesystem of the docker container.  

Build a docker image:

    docker build -t envoy:v1 envoy-proxy
    
Run the docker container with the created image on port 51051:

    docker run -d --name envoy -p 9901:9901 -p 51051:51051 envoy:v1
    
#### 6. Step
Run the gRPC server on port 50051 (if you haven't done the other tutorial, this will also start the gRPC gateway proxy):

    go run main.go
    
#### 7. Step

##### cURL
Now you can test the envoy proxy with the same cURL calls as in 6. step of the other tutorial except that you have to change
the port to 51051, e.g.:

    curl -X POST \
      http://localhost:51051/shelves \
      -H 'Content-Type: application/json' \
      -d '{
        "name": "Books I have read",
        "theme": "Biography"
    }'
    
##### gRPC client:

A sample gRPC client is provided inside `grpc-client/client.go` that lists all themes of your shelves:

    go run grpc-client/client.go

This client calls the gRPC server directly (port 50051). You can also call the envoy proxy (port 51051).