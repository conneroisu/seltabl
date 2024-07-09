package config

import (
	"fmt"
	"testing"
	"text/template"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

// TestSections generates a sections from a given url
func TestSections(t *testing.T) {
	a := assert.New(t)
	tmpl := template.New("sections_file_template")
	tmpl, err := tmpl.Parse(sectionsTmpl)
	a.NoError(err)
	out, err := NewSectionsErrorPrompt(
		fmt.Errorf(
			"failed to get the content of the url: Get https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html: dial tcp: lookup github.com on 127.0.0.53:53: read udp 127.0.0.1:53448->127.0.0.53:53: i/o timeout",
		),
	)
	a.NoError(err)
	a.NotEmpty(out)
	t.Logf("sections: %s", out)
}
