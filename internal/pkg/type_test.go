package pkg

import "testing"

func TestDataType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		d    DataType
		want bool
	}{
		{"Valid DataTypeCredential", DataTypeCredential, true},
		{"Valid DataTypeCard", DataTypeCard, true},
		{"Valid DataTypeText", DataTypeText, true},
		{"Valid DataTypeDocument", DataTypeDocument, true},
		{"Invalid DataTypeUnknown", DataType("DataTypeUnknown"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
