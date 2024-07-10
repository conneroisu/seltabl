package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/stretchr/testify/assert"
)

// increment increments a given integer
func increment(i int) int {
	return i + 1
}

// TestIncrement increments a given integer
func TestIncrement(t *testing.T) {
	a := assert.New(t)
	i := increment(1)
	a.Equal(2, i)
}

func TestGenerate(t *testing.T) {
	t.Run("Test Generate", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		cfg := &domain.ConfigFile{
			Name:           "TestStruct",
			URL:            "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
			IgnoreElements: []string{"script", "style", "link", "img", "footer", "header"},
			HTMLBody:       "<html><body><table><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>",
			Selectors: master.Selectors{
				{
					Value:      "a",
					Occurances: 1,
					Context:    "html head meta[name=csrf-param]",
				},
				{
					Value:      "b",
					Occurances: 1,
					Context:    "html",
				},
				{
					Value:      "c",
					Occurances: 1,
					Context:    "html body div.footer",
				},
				{
					Value:      "d",
					Occurances: 1,
					Context:    "html > head > title",
				},
			},
		}

		path := filepath.Join(".", "seltabl.yaml")
		log.Debugf("Writing config file to path: %s", path)
		defer log.Debugf("Config file written to path: %s", path)
		f, err := os.Create(path)
		a.NoError(err)
		defer f.Close()
		err = Generate(context.Background(), nil, cfg, f)
		a.NoError(err)
	})
}
