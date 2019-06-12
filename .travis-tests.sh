go test ./descriptor-generator
cp ./descriptor-generator/test_data/parameters/output/test.descr ./protoc-generator/test_data/parameters/
cp ./descriptor-generator/test_data/requestBodies/output/test.descr ./protoc-generator/test_data/requestBodies/
cp ./descriptor-generator/test_data/responses/output/test.descr ./protoc-generator/test_data/responses/
go test ./protoc-generator/