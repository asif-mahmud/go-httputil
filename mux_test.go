package gohttputil_test

import (
	"testing"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/stretchr/testify/assert"
)

func TestMuxCreate(t *testing.T) {
	m := gohttputil.New()

	assert.NotNil(t, m)
}
