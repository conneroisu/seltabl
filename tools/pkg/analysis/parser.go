package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/pkg/parsers"
)

func (s *State) getUrlsAndIgnores(src string) ([]string, []string, error) {
	var urls []string
	var err error
	var ignores []string
	ctx := context.Background()
	urls, err = parsers.ExtractUrls(src)
	if err != nil {
		s.Logger.Printf("failed to extract urls: %s\n", err)
		return nil, nil, err
	}
	ignores, err = parsers.ExtractIgnores(ctx, src)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract ignores: %w", err)
	}
	return urls, ignores, nil
}
