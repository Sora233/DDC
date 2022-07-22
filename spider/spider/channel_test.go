package spider

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetChannel(t *testing.T) {
	resp, err := GetChannel()
	assert.Nil(t, err)
	spew.Dump(resp)
}
