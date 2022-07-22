package spider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetComments(t *testing.T) {
	resp, err := GetComments("ZZ3jr3ZcoD4D33Go")
	assert.Nil(t, err)
	assert.Zero(t, resp.Status)
}
