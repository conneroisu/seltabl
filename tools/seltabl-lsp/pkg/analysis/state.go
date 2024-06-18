package analysis

import (
	"context"
	"log"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/data"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/parsers"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

// State is the state of the document analysis
type State struct {
	// Map of file names to contents
	Documents map[string]string
	// Selectors is the map of file names to selectors
	Selectors map[string][]data.Selector
	// Database is the database for the state
	Database *bun.DB
	// Logger is the logger for the state
	Logger *log.Logger
}

func getLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(logFile, "[seltabl-lsp#state]", log.LstdFlags)
}

// NewState returns a new state with no documents
func NewState(getenv func(string) string) State {
	db, err := data.NewDb(
		context.Background(),
		getenv,
		"urls.sqlite",
		data.SetupMasterDatabase,
		sqlitedialect.New(),
	)
	if err != nil {
		panic(err)
	}
	logger := getLogger("./state.log")
	return State{
		Documents: map[string]string{},
		Selectors: map[string][]data.Selector{},
		Database:  db,
		Logger:    logger,
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

// getSelectors gets all the selectors from the given URL and appends them to the selectors slice
func (s State) getSelectors(
	url string,
	ignores []string,
) ([]data.Selector, error) {
	ctx := context.Background()
	res := []data.Selector{}
	err := s.Database.NewSelect().Model(&res).Scan(ctx)
	if err != nil || len(res) == 0 {
		doc, err := parsers.GetMinifiedDoc(url, ignores)
		if err != nil {
			s.Logger.Printf("failed to get minified doc: %s\n", err)
		}
		got := parsers.GetAllSelectors(doc)
		var selectors []data.Selector
		for _, sel := range got {
			htm, err := doc.Find(sel).Parent().Html()
			if err != nil {
				s.Logger.Printf("failed to get html: %s\n", err)
			}
			var url []data.URL
			err = s.Database.NewSelect().Model(&url).Where("url = ?", sel).Scan(ctx, &url)
			if err != nil {
				s.Logger.Printf("failed to get urls: %s\n", err)
			}
			var u data.URL
			if len(url) == 0 {
				u = data.URL{
					URL: sel,
				}
				s.Database.NewInsert().Model(&u).Exec(ctx)
			} else {
				u = url[0]
			}
			selector := data.Selector{
				Selector: sel,
				URL:      u,
				Context:  htm,
			}
			selectors = append(selectors, selector)
			if _, err := s.Database.NewInsert().Model(selector).Exec(ctx); err != nil {
				s.Logger.Printf("failed to insert selector: %s\n", err)
			}
		}
		return selectors, nil
	}
	return res, nil
}
