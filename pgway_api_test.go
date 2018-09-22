package pgway

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPgwayApi_IsSamePath(t *testing.T) {
	api := Api{
		Path: "fe",
	}

	assert.Equal(t, true, api.IsSamePath("fe"))

}
