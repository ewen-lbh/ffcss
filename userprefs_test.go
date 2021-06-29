package main

import (
	"sort"
	"strings"
	"testing"
)

func TestToUserJS(t *testing.T) {
	value, err := ToUserJSFile(map[string]interface{}{
		"browser.tabs.tabClipWidth":              90,
		"svg.context-properties.content.enabled": true,
	})
	if err != nil {
		panic(err)
	}
	Assert(t, sortLines(value),
		`user_pref("browser.tabs.tabClipWidth", 90);
user_pref("svg.context-properties.content.enabled", true);`,
	)
}

func sortLines(s string) string {
	lines := strings.Split(s, "\n")
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}
