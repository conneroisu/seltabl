package analysis

import (
	"github.com/conneroisu/seltabl/tools/pkg/parsers"
)

func (s *State) getUrlsAndIgnores(src string) ([]string, []string, error) {
	urls, err := parsers.ExtractUrls(src)
	if err != nil {
		s.Logger.Printf("failed to extract urls: %s\n", err)
		return nil, nil, err
	}
	ignores, _ := parsers.ExtractIgnores(src)
	return urls, ignores, nil
}
