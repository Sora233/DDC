package spider

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetClips(t *testing.T) {
	resp, err := GetClips(1802011210)
	assert.Nil(t, err)
	spew.Dump(resp)
}
