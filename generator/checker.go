package generator

import (
	openapiv3 "github.com/googleapis/gnostic/OpenAPIv3"
	plugins "github.com/googleapis/gnostic/plugins"
)

type FeatureChecker struct {
	// The document to be analyzed
	document *openapiv3.Document
	// The messages that are displayed to the user with information of what is not being processed
	messages []*plugins.Message
}

func NewFeatureChecker(document *openapiv3.Document) *FeatureChecker {
	return &FeatureChecker{document: document, messages: make([]*plugins.Message, 0)}
}

func (c *FeatureChecker) Run() []*plugins.Message {
	c.analyzeOpenAPIdocument()
	return c.messages
}

func (c *FeatureChecker) analyzeOpenAPIdocument() {

	if c.document.Servers != nil || c.document.Security != nil || c.document.Tags != nil || c.document.ExternalDocs != nil {
		msg := &plugins.Message{
			Code:  "DOCUMENTFIELDS",
			Level: plugins.Message_WARNING,
			Text:  "Fields: Servers, Security, Tags, and ExternalDocs are not supported for Document with title: " + c.document.Info.Title,
			Keys:  []string{"Document"},
		}
		c.messages = append(c.messages, msg)
	}

	c.analyzeComponents()
	c.analyzePaths()
}

func (c *FeatureChecker) analyzeComponents() {
	components := c.document.Components

	if components.Examples != nil || components.Headers != nil || components.SecuritySchemes != nil ||
		components.Links != nil || components.Callbacks != nil {
		msg := &plugins.Message{
			Code:  "COMPONENTSFIELDS",
			Level: plugins.Message_WARNING,
			Text:  "Fields: Examples, Headers, Links, and Callbacks are not supported for the component",
			Keys:  []string{"Component"},
		}
		c.messages = append(c.messages, msg)
	}

	if schemas := components.GetSchemas(); schemas != nil {
		for _, pair := range schemas.AdditionalProperties {
			c.analyzeSchema(pair.Name, pair.Value)
		}
	}

	if responses := components.GetResponses(); responses != nil {
		for _, pair := range responses.AdditionalProperties {
			c.analyzeResponse(pair)
		}
	}

	if parameters := components.GetParameters(); parameters != nil {
		for _, pair := range parameters.AdditionalProperties {
			c.analyzeParameter(pair.Value)
		}
	}

	if requestBodies := components.GetRequestBodies(); requestBodies != nil {
		for _, pair := range requestBodies.AdditionalProperties {
			c.analyzeRequestBody(pair)
		}
	}
}

func (c *FeatureChecker) analyzePaths() {
	for _, pathItem := range c.document.Paths.Path {
		c.analyzePathItem(pathItem)
	}
}

func (c *FeatureChecker) analyzePathItem(pair *openapiv3.NamedPathItem) {
	pathItem := pair.Value

	if pathItem.Head != nil || pathItem.Options != nil || pathItem.Trace != nil || pathItem.Servers != nil ||
		pathItem.Parameters != nil {
		msg := &plugins.Message{
			Code:  "PATHFIELDS",
			Level: plugins.Message_WARNING,
			Text:  "Fields: Head, Options, Trace, Servers, and Parameters are not supported for path: " + pair.Name,
			Keys:  []string{"Paths", pair.Name, "Operation"},
		}
		c.messages = append(c.messages, msg)
	}

	operations := getValidOperations(pathItem)
	for _, op := range operations {
		c.analyzeOperation(op)
	}
}

func (c *FeatureChecker) analyzeOperation(operation *openapiv3.Operation) {

	if operation.Tags != nil || operation.ExternalDocs != nil || operation.Callbacks != nil || operation.Deprecated ||
		operation.Security != nil || operation.Servers != nil {
		msg := &plugins.Message{
			Code:  "OPERATIONFIELDS",
			Level: plugins.Message_WARNING,
			Text: "Fields: Tags, ExternalDocs, Callbacks, Deprecated, Security, and Servers " +
				"are not supported for operation: " + operation.OperationId,
			Keys: []string{"Operation", operation.OperationId, "Callbacks"},
		}
		c.messages = append(c.messages, msg)
	}

	for _, param := range operation.Parameters {
		c.analyzeParameter(param)
	}
}

func (c *FeatureChecker) analyzeParameter(paramOrRef *openapiv3.ParameterOrReference) {
	if parameter := paramOrRef.GetParameter(); parameter != nil {
		if parameter.Required || parameter.Deprecated || parameter.AllowEmptyValue || parameter.Style != "" ||
			parameter.Explode || parameter.AllowReserved || parameter.Example != nil || parameter.Examples != nil ||
			parameter.Content != nil {
			msg := &plugins.Message{
				Code:  "PARAMETERFIELDS",
				Level: plugins.Message_WARNING,
				Text: "Fields: Required, Deprecated, AllowEmptyValue, Style, Explode, AllowReserved, Example, Examples" +
					" and Content are not supported for parameter: " + parameter.Name,
				Keys: []string{"Parameter", parameter.Name},
			}
			c.messages = append(c.messages, msg)
		}
		c.analyzeSchema(parameter.Name, parameter.Schema)
	}
}

func (c *FeatureChecker) analyzeSchema(identifier string, schemaOrReference *openapiv3.SchemaOrReference) {
	if schema := schemaOrReference.GetSchema(); schema != nil {
		if schema.Nullable || schema.Discriminator != nil || schema.ReadOnly || schema.WriteOnly || schema.Xml != nil ||
			schema.ExternalDocs != nil || schema.Example != nil || schema.Deprecated || schema.Title != "" ||
			schema.MultipleOf != 0 || schema.Maximum != 0 || schema.ExclusiveMaximum || schema.Minimum != 0 ||
			schema.ExclusiveMinimum || schema.MaxLength != 0 || schema.MinLength != 0 || schema.Pattern != "" ||
			schema.MaxItems != 0 || schema.MinItems != 0 || schema.UniqueItems || schema.MaxProperties != 0 ||
			schema.MinProperties != 0 || schema.Required != nil || schema.AllOf != nil || schema.OneOf != nil ||
			schema.AnyOf != nil || schema.Not != nil || schema.Default != nil {

			msg := &plugins.Message{
				Code:  "SCHEMAFIELDS",
				Level: plugins.Message_WARNING,
				Text: "Fields: Nullable, Discriminator, ReadOnly, WriteOnly, Xml, ExternalDocs, Example, Deprecated, " +
					"Title, MultipleOf, Maximum, ExclusiveMaximum, Minimum, ExclusiveMinimum, MaxLength, MinLength, " +
					"Pattern, MaxItems, MinItems, UniqueItems, MaxProperties, MinProperties, Required, AllOf, OneOf, " +
					"AnyOf, Not, Default are not supported for the schema: " + identifier,
				Keys: []string{identifier, "Schema"},
			}
			c.messages = append(c.messages, msg)
		}

		if enum := schema.Enum; enum != nil {
			msg := &plugins.Message{
				Code:  "SCHEMAFIELDS",
				Level: plugins.Message_WARNING,
				Text:  "Field: Enum is not generated as enum in .proto for schema: " + identifier,
				Keys:  []string{identifier, "Schema"},
			}
			c.messages = append(c.messages, msg)
		}

		if properties := schema.Properties; properties != nil {
			for _, pair := range properties.AdditionalProperties {
				c.analyzeSchema(pair.Name, pair.Value)
			}
		}

		if additionalProperties := schema.AdditionalProperties; additionalProperties != nil {
			c.analyzeSchema("AdditionalPropertiesSchema", additionalProperties.GetSchemaOrReference())
		}
	}
}

func (c *FeatureChecker) analyzeResponse(pair *openapiv3.NamedResponseOrReference) {
	if response := pair.Value.GetResponse(); response != nil {
		if response.Links != nil || response.Headers != nil {
			msg := &plugins.Message{
				Code:  "RESPONSEFIELDS",
				Level: plugins.Message_WARNING,
				Text:  "Fields: Links, and Headers are not supported for response: " + pair.Name,
				Keys:  []string{"Response", pair.Name},
			}
			c.messages = append(c.messages, msg)
		}

		for _, pair := range response.Content.AdditionalProperties {
			c.analyzeContent(pair)
		}
	}
}

func (c *FeatureChecker) analyzeRequestBody(pair *openapiv3.NamedRequestBodyOrReference) {
	if requestBody := pair.Value.GetRequestBody(); requestBody != nil {
		if requestBody.Required {
			msg := &plugins.Message{
				Code:  "REQUESTBODYFIELDS",
				Level: plugins.Message_WARNING,
				Text:  "Fields: Required are not supported for the request: " + pair.Name,
				Keys:  []string{"RequestBody", pair.Name},
			}
			c.messages = append(c.messages, msg)
		}
		for _, pair := range requestBody.Content.AdditionalProperties {
			c.analyzeContent(pair)
		}
	}
}

func (c *FeatureChecker) analyzeContent(pair *openapiv3.NamedMediaType) {
	mediaType := pair.Value

	if mediaType.Examples != nil || mediaType.Example != nil || mediaType.Encoding != nil {
		msg := &plugins.Message{
			Code:  "MEDIATYPEFIELDS",
			Level: plugins.Message_WARNING,
			Text:  "Fields: Examples, Example, and Encoding are not supported for the mediatype: " + pair.Name,
			Keys:  []string{"MediaType", pair.Name},
		}
		c.messages = append(c.messages, msg)
		c.analyzeSchema(pair.Name, mediaType.Schema)
	}
}

func getValidOperations(pathItem *openapiv3.PathItem) []*openapiv3.Operation {
	operations := make([]*openapiv3.Operation, 0)

	if pathItem.Get != nil {
		operations = append(operations, pathItem.Get)
	}
	if pathItem.Put != nil {
		operations = append(operations, pathItem.Put)
	}
	if pathItem.Post != nil {
		operations = append(operations, pathItem.Post)
	}
	if pathItem.Delete != nil {
		operations = append(operations, pathItem.Delete)
	}
	if pathItem.Patch != nil {
		operations = append(operations, pathItem.Patch)
	}
	return operations
}
