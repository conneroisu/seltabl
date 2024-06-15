package analysis

import (
	"context"
	"log"

	"github.com/conneroisu/seltabl/tools/data"
	"github.com/conneroisu/seltabl/tools/data/master"
	"github.com/conneroisu/seltabl/tools/pkg/lsp"
)

// State is the state of the document analysis
type State struct {
	// Map of file names to contents
	Documents map[string]string
	// Database is the database for the state
	Database *data.Database[master.Queries]
	// Logger is the logger for the state
	Logger *log.Logger
}

// NewState returns a new state with no documents
func NewState() State {
	db, err := data.NewDb(
		context.Background(),
		master.New,
		&data.Config{
			URI:    "",
			Schema: master.MasterSchema,
		},
		"urls.sqlite",
	)
	if err != nil {
		panic(err)
	}
	logger := getLogger("./seltabl.log")
	return State{
		Documents: map[string]string{},
		Logger:    logger,
		Database:  db,
	}
}

// LineRange returns a range of a line in a document
func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
