package schema

import (
	"strings"
	"testing"
)

func TestValidatePayloadSchema(t *testing.T) {
	validSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"user_id": map[string]any{"type": "string"},
			"amount":  map[string]any{"type": "number"},
		},
		"required": []any{"user_id"},
	}

	tests := []struct {
		name    string
		schema  map[string]any
		data    map[string]any
		wantErr bool
		contains string
	}{
		{
			name:    "schema kosong di-skip",
			schema:  map[string]any{},
			data:    map[string]any{"user_id": "u1"},
			wantErr: false,
		},
		{
			name:    "schema nil di-skip (len 0)",
			schema:  nil,
			data:    map[string]any{"user_id": "u1"},
			wantErr: false,
		},
		{
			name:    "data nil dianggap {}",
			schema:  validSchema,
			data:    nil,
			wantErr: true,
			contains: "payload tidak sesuai schema",
		},
		{
			name:    "data valid",
			schema:  validSchema,
			data:    map[string]any{"user_id": "u1", "amount": 100.5},
			wantErr: false,
		},
		{
			name:    "data valid hanya required",
			schema:  validSchema,
			data:    map[string]any{"user_id": "u1"},
			wantErr: false,
		},
		{
			name:    "data invalid - missing required",
			schema:  validSchema,
			data:    map[string]any{"amount": 100},
			wantErr: true,
			contains: "payload tidak sesuai schema",
		},
		{
			name:    "data invalid - wrong type",
			schema:  validSchema,
			data:    map[string]any{"user_id": 123},
			wantErr: true,
			contains: "payload tidak sesuai schema",
		},
		{
			name: "schema invalid - bukan JSON schema valid",
			schema: map[string]any{
				"type": "invalid_type_value",
			},
			data:    map[string]any{},
			wantErr: true,
			contains: "payload schema validation",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePayloadSchema(tt.schema, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePayloadSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.contains != "" && err != nil && !strings.Contains(err.Error(), tt.contains) {
				t.Errorf("ValidatePayloadSchema() error = %v, want contains %q", err.Error(), tt.contains)
			}
		})
	}
}

func TestValidateTemplateSchema(t *testing.T) {
	validSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
	}

	tests := []struct {
		name     string
		schema   map[string]any
		wantErr  bool
		contains string
	}{
		{
			name:    "schema kosong di-skip",
			schema:  map[string]any{},
			wantErr: false,
		},
		{
			name:    "schema nil di-skip",
			schema:  nil,
			wantErr: false,
		},
		{
			name:    "schema valid",
			schema:  validSchema,
			wantErr: false,
		},
		{
			name: "schema invalid",
			schema: map[string]any{
				"type": "not_a_valid_type",
			},
			wantErr:  true,
			contains: "template",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplateSchema(tt.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplateSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.contains != "" && err != nil && !strings.Contains(err.Error(), tt.contains) {
				t.Errorf("ValidateTemplateSchema() error = %v, want contains %q", err.Error(), tt.contains)
			}
		})
	}
}
