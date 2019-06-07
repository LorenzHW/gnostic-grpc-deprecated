# gnostic Protos Generator Plugin
[GSoC 2019 project](https://summerofcode.withgoogle.com/projects/#5244822191865856)

To run this plugin run following commands inside this directory:

    go build
    ./gnostic-protoc-generator -input test/helloworld.descr -output test/generated/
    
   
This plugin uses a file descriptor set as input (`test/helloworld.descr`) which is retrieved from
`test/hellworld.proto`. This file descriptor set is then used to generated the content inside `test/generated`.