package generate

import (
	"context"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/sashabaranov/go-openai"
)

// GenerateSuite generates a suite for a given name
func GenerateSuite(
	ctx context.Context,
	client *openai.Client,
	name string,
	url string,
	ignoreElements []string,
	htmlBody string,
	selectors []master.Selector,
) (err error) {
	// err = NewTestFile(name, url, ignoreElements).Generate()
	// if err != nil {
	//         return fmt.Errorf("failed to generate test file: %w", err)
	// }
	// err = NewConfigFile(name, url, ignoreElements).Generate()
	// if err != nil {
	//         return fmt.Errorf("failed to generate config file: %w", err)
	// }
	// err = NewStructFile(name, url, ignoreElements).Generate()
	// if err != nil {
	//         return fmt.Errorf("failed to generate struct file: %w", err)
	// }
	return nil
}
