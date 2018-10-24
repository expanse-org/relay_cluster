package registry

import (
	"testing"

	motan "github.com/expanse-org/motan-go/core"
)

func TestGetRegistry(t *testing.T) {
	defaultExtFactory := &motan.DefaultExtentionFactory{}
	defaultExtFactory.Initialize()
	RegistDefaultRegistry(defaultExtFactory)
	url := &motan.URL{
		Protocol:   "direct",
		Host:       "127.0.0.1",
		Port:       4072,
		Path:       "weibo.com",
		Group:      "yf",
		Parameters: make(map[string]string),
	}
	registry := defaultExtFactory.GetRegistry(url)
	if registry.GetName() != url.Protocol {
		t.Error("GetName Error")
	}
}
