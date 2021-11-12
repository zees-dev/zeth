package node

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_prefixCheck(t *testing.T) {
	is := assert.New(t)

	httpURL := "http://localhost:9000"
	websocketURL := "ws://localhost:9000"

	url, err := url.Parse(httpURL)
	is.NoError(err)
	is.Equal("localhost:9000", url.Host)

	url, err = url.Parse(websocketURL)
	is.NoError(err)
	is.Equal("localhost:9000", url.Host)
}
