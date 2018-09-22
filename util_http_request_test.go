package pgway

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlSanitize(t *testing.T) {

	baseStr := "//fe/fae////fe"
	expectedStr := "/fe/fae/fe"
	assert.Equal(t, expectedStr, UrlSanitize(baseStr))
}
