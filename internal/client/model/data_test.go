package model

import (
	"testing"

	"github.com/casnerano/seckeep/internal/pkg"
	"github.com/stretchr/testify/assert"
)

func TestDataCard_Type(t *testing.T) {
	card := DataCard{}
	assert.Equal(t, pkg.DataTypeCard, card.Type())
}

func TestDataCredential_Type(t *testing.T) {
	credential := DataCredential{}
	assert.Equal(t, pkg.DataTypeCredential, credential.Type())
}

func TestDataDocument_Type(t *testing.T) {
	document := DataDocument{}
	assert.Equal(t, pkg.DataTypeDocument, document.Type())
}

func TestDataText_Type(t *testing.T) {
	text := DataText{}
	assert.Equal(t, pkg.DataTypeText, text.Type())
}
