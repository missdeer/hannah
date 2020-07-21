package provider

import (
	"testing"
)

func TestGetProvider(t *testing.T) {
	names := []string{
		"netease",
		"xiami",
		"qq",
		"kugou",
		"kuwo",
		"bilibili",
		"migu",
	}
	for _, name := range names {
		p := GetProvider(name)
		if p.Name() != name {
			t.Error("name mismatched")
		}
	}
	invalidNames := []string{
		"invalid", "names", "5ting", "baidu",
	}
	for _, name := range invalidNames {
		p := GetProvider(name)
		if p != nil {
			t.Error("provider should not exists for", name)
		}
	}
}
