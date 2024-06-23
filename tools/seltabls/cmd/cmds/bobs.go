package cmds

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/sync/errgroup"
)

// WriteResponse writes a message to the writer
func WriteResponse(
	ctx context.Context,
	writer *io.Writer,
	msg rpc.MethodActor,
) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		reply, err := rpc.EncodeMessage(msg)
		if err != nil {
			return fmt.Errorf(
				"failed to encode response to request (%s): %w",
				msg.Method(),
				err,
			)
		}
		res, err := (*writer).Write([]byte(reply))
		if err != nil {
			return fmt.Errorf(
				"failed to encode response to request (%s): %w",
				msg.Method(),
				err,
			)
		}
		if res != len(reply) {
			return fmt.Errorf(
				"failed to write all of response to request (%s): %w",
				msg.Method(),
				err,
			)
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to write message (%s): %w", msg.Method(), err)
	}
	return nil
}

// CreateConfigDir creates a new config directory and returns the path.
func CreateConfigDir() (string, error) {
	path, err := homedir.Expand("~/.config/seltabls/")
	if err != nil {
		return "", fmt.Errorf("failed to expand home directory: %w", err)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		if os.IsExist(err) {
			return path, nil
		}
		return "", fmt.Errorf("failed to create or find config directory: %w", err)
	}
	return path, nil
}
