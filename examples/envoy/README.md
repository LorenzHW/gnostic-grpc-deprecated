### gRPC server with envoy proxy
This tutorial demonstrates how to set up a gRPC API using [envoy](https://www.envoyproxy.io/) as proxy. 

#### Prerequisite
Follow steps 1-3 from this [tutorial](https://github.com/LorenzHW/gnostic-grpc/tree/master/examples/end-to-end) in order
to get an implementation of a gRPC server. After those steps you should have a `bookstore/server.go` and a `bookstore/bookstore.pb.go`
file.

For simplicity lets create a temporary environment variable inside your terminal:
    
    export ANNOTATIONS="third-party/googleapis"

In order for this tutorial to work you should work inside this directory under `GOPATH`.

#### 1. Step
Given `bookstore.proto` generate the descriptor set.
    
    protoc --proto_path=${ANNOTATIONS} --proto_path=. --include_imports --include_source_info --descriptor_set_out=envoy-proxy/proto.pb bookstore.proto
    
This generates `envoy-proxy/proto.pb`.

#### 2. Step 
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
    
#### 3. Step
Run the gRPC server on port 50051:

    go run main.go
    
#### 4. Step

##### curl
Inside of a new terminal test your API with the envoy proxy:

Let's create a shelf first:

    curl -X POST \
      http://localhost:51051/shelves \
      -H 'Content-Type: application/json' \
      -d '{
        "name": "Books I need to read",
        "theme": "Non-fiction"
    }'
    
Get all existing shelves:

    curl -X GET http://localhost:51051/shelves
    
Create a book for the shelve with the id `1`:
    
    curl -X POST \
      http://localhost:51051/shelves/1/books \
      -H 'Content-Type: application/json' \
      -d '{
        "author": "Hans Rosling",
        "name": "Factfulness",
        "title": "Factfulness: Ten Reasons We'\''re wrong about the world - and Why Things Are Better Than You Think"
    }'
    
    
List all books for the shelve with the id `1`:

    curl -X GET http://localhost:51051/shelves/1/books
    
    
##### gRPC client

A sample gRPC client is provided inside `grpc-client/client.go` that lists all themes of your shelves:

    go run grpc-client/client.go