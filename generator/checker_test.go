package generator

import (
	"github.com/golang/protobuf/proto"
	openapiv3 "github.com/googleapis/gnostic/OpenAPIv3"
	plugins "github.com/googleapis/gnostic/plugins"
	"io/ioutil"
	"testing"
)

func TestNewFeatureCheckerParameters(t *testing.T) {
	input := "testfiles/parameters/test.pb"
	documentv3 := ReadOpenAPIBinary(input)

	checker := NewFeatureChecker(documentv3)
	messages := checker.Run()
	expectedMessageTexts := []string{
		"Fields: Explode are not supported for parameter: param2",
		"Fields: Default are not supported for the schema: Items of param2",
		"Field: Enum is not generated as enum in .proto for schema: Items of param2",
		"Fields: Default are not supported for the schema: param4",
		"Field: Enum is not generated as enum in .proto for schema: param4",
	}
	validateMessages(t, expectedMessageTexts, messages)
}

func TestFeatureCheckerRequestBodies(t *testing.T) {
	input := "testfiles/requestBodies/test.pb"
	documentv3 := ReadOpenAPIBinary(input)

	checker := NewFeatureChecker(documentv3)
	messages := checker.Run()
	expectedMessageTexts := []string{
		"Fields: Required are not supported for the schema: Person",
		"Fields: Example are not supported for the schema: name",
		"Fields: Xml are not supported for the schema: photoUrls",
		"Fields: Required are not supported for the request: RequestBody",
	}
	validateMessages(t, expectedMessageTexts, messages)
}

func TestFeatureCheckerResponses(t *testing.T) {
	input := "testfiles/responses/test.pb"
	documentv3 := ReadOpenAPIBinary(input)

	checker := NewFeatureChecker(documentv3)
	messages := checker.Run()
	expectedMessageTexts := []string{
		"Fields: Required are not supported for the schema: Error",
		"Fields: Required are not supported for the schema: Person",
		"Fields: Example are not supported for the schema: name",
		"Fields: Xml are not supported for the schema: photoUrls",
	}
	validateMessages(t, expectedMessageTexts, messages)
}

func validateMessages(t *testing.T, expectedMessageTexts []string, messages []*plugins.Message) {
	if len(expectedMessageTexts) != len(messages) {
		t.Errorf("Number of messages from FeatureChecker does not match expected number")
	}
	for i, msg := range messages {
		if msg.Text != expectedMessageTexts[i] {
			t.Errorf("Message text does not match expected message text: %s != %s", msg.Text, expectedMessageTexts[i])
		}
	}
}

func ReadOpenAPIBinary(input string) *openapiv3.Document {
	apiData, _ := ioutil.ReadFile(input)
	documentv3 := &openapiv3.Document{}
	proto.Unmarshal(apiData, documentv3)
	return documentv3
}
