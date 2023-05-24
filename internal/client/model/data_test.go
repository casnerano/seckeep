package model

import (
	"testing"

	"github.com/casnerano/seckeep/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestDataCard_Type(t *testing.T) {
	card := DataCard{}
	assert.Equal(t, shared.DataTypeCard, card.Type())
}

func TestDataCredential_Type(t *testing.T) {
	credential := DataCredential{}
	assert.Equal(t, shared.DataTypeCredential, credential.Type())
}

func TestDataDocument_Type(t *testing.T) {
	document := DataDocument{}
	assert.Equal(t, shared.DataTypeDocument, document.Type())
}

func TestDataText_Type(t *testing.T) {
	text := DataText{}
	assert.Equal(t, shared.DataTypeText, text.Type())
}
