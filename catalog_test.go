package ffcss

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeTheme(name string) Theme {
	t := NewTheme()
	t.ExplicitName = name
	return t
}

func TestCatalogLookup(t *testing.T) {
	catalog := Catalog{
		"lorem": makeTheme("lorem"),
		"ipsum": makeTheme("ipsum"),
		"dolor": makeTheme("a_completely_different_name"),
		"bacon": makeTheme("bacon"),
	}

	actual, err := catalog.Lookup("lorem")
	assert.NoError(t, err)
	assert.Equal(t, makeTheme("lorem"), actual)

	actual, err = catalog.Lookup("a.completely different-name")
	assert.NoError(t, err)
	assert.Equal(t, makeTheme("a_completely_different_name"), actual)

	actual, err = catalog.Lookup("dolor")
	assert.Error(t, err, `theme "dolor" not found`)
	assert.Equal(t, Theme{}, actual)

	actual, err = catalog.Lookup("baaaa con")
	assert.Error(t, err, `theme "baaaa con" not found. did you mean [blue][bold]bacon[reset]?`)
	assert.Equal(t, Theme{}, actual)
}

func TestLookupPreprocess(t *testing.T) {
	cases := []struct{ in, out string }{
		{"same", "same"},
		{"B  O   N    K", "bonk"},
		{"you. shall. not. pass.", "youshallnotpass"},
		{"wHaTevER", "whatever"},
		{"all-thE.puncT_uAtio    n", "allthepunctuation"},
		{"with great power\u037e comes great responsibility.", "withgreatpower;comesgreatresponsibility"},
		{"𝔽𝔸ℕℂ𝕐", "fancy"},
	}

	for _, caze := range cases {
		assert.Equal(t, caze.out, lookupPreprocess(caze.in))
	}
}

func TestLoadCatalog(t *testing.T) {
	actual, err := LoadCatalog(filepath.Join(testarea, "catalogs/empty"))
	assert.NoError(t, err)
	assert.Equal(t, make(Catalog), actual)

	actual, err = LoadCatalog(filepath.Join(testarea, "catalogs/various"))
	assert.NoError(t, err)
	actualNames := make([]string, 0, len(actual))
	for name := range actual {
		actualNames = append(actualNames, name)
	}
	assert.Equal(t, []string{"alpenblue", "flyingfox", "montereyfox", "sometheme"}, actualNames)
}
