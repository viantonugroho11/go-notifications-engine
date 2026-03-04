package schema

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// validatePayloadSchema memvalidasi n.Data terhadap JSON schema (PayloadSchema) menggunakan gojsonschema.
// Jika schema nil atau kosong, validasi di-skip. Jika data nil, dianggap sebagai {}.
func ValidatePayloadSchema(schema map[string]any, data map[string]any) error {
	if len(schema) == 0 {
		return nil
	}
	schemaLoader := gojsonschema.NewGoLoader(schema)
	document := data
	if document == nil {
		document = make(map[string]any)
	}
	documentLoader := gojsonschema.NewGoLoader(document)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("payload schema validation: %w", err)
	}
	if result.Valid() {
		return nil
	}
	var errMsgs []string
	for _, e := range result.Errors() {
		errMsgs = append(errMsgs, e.String())
	}
	return fmt.Errorf("payload tidak sesuai schema: %s", strings.Join(errMsgs, "; "))
}

// validateTemplateSchema untuk validasi schema yang di inputkan pada template
func ValidateTemplateSchema(schema map[string]any) error {
	if len(schema) == 0 {
		return nil
	}
	schemaLoader := gojsonschema.NewGoLoader(schema)
	document := make(map[string]any)
	documentLoader := gojsonschema.NewGoLoader(document)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("template schema validation: %w", err)
	}
	if result.Valid() {
		return nil
	}
	var errMsgs []string
	for _, e := range result.Errors() {
		errMsgs = append(errMsgs, e.String())
	}
	return fmt.Errorf("template tidak sesuai schema: %s", strings.Join(errMsgs, "; "))
}