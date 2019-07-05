package generator

import (
	"github.com/golang/protobuf/proto"
	openapiv3 "github.com/googleapis/gnostic/OpenAPIv3"
	"io/ioutil"
	"testing"
)

func TestFeatureChecker(t *testing.T) {
	input := "testfiles/requestBodies/test.pb"
	documentv3 := ReadOpenAPIBinary(input)

	checker := NewFeatureChecker(documentv3)
	messages := checker.Run()
	if messages != nil {

	}

}

func ReadOpenAPIBinary(input string) *openapiv3.Document {
	apiData, _ := ioutil.ReadFile(input)
	documentv3 := &openapiv3.Document{}
	proto.Unmarshal(apiData, documentv3)
	return documentv3
}
