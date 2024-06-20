package analysis

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/data"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/data/master"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/internal/config"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/parsers"
	"github.com/yosssi/gohtml"
)

// State is the state of the document analysis
type State struct {
	// Map of file names to contents
	Documents map[string]string
	// Selectors is the map of file names to selectors
	Selectors map[string][]master.Selector
	// Database is the database for the state
	Database data.Database[master.Queries]
	// Logger is the logger for the state
	Logger *log.Logger
}

// NewState returns a new state with no documents
func NewState(config *config.Config) (state State, err error) {
	db, err := data.NewDb(
		context.Background(),
		master.New,
		&data.Config{
			Schema:   master.MasterSchema,
			URI:      "sqlite://uri.sqlite",
			FileName: path.Join(config.ConfigPath, "uri.sqlite"),
		},
	)
	if err != nil {
		return state, fmt.Errorf("failed to create database: %w", err)
	}
	logger := getLogger(path.Join(config.ConfigPath, "state.log"))
	state = State{
		Documents: make(map[string]string),
		Selectors: make(map[string][]master.Selector),
		Database:  *db,
		Logger:    logger,
	}
	return state, nil
}

// getLogger returns a logger that writes to a file
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

// getSelectors gets all the selectors from the given URL and appends them to the selectors slice
func (s State) getSelectors(
	ctx context.Context,
	url []string,
	ignores []string,
) (selectors []master.Selector, err error) {
	u, err := s.Database.Queries.GetURLByValue(
		ctx,
		master.GetURLByValueParams{Value: url[0]},
	)
	if err == nil {
		rows, err := s.Database.Queries.GetSelectorsByURL(
			ctx,
			master.GetSelectorsByURLParams{Value: u.Value},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get selectors by url: %w", err)
		}
		if rows != nil {
			var selectors []master.Selector
			for _, row := range rows {
				selectors = append(selectors, *row)
			}
			return selectors, nil
		}
	}
	doc, err := parsers.GetMinifiedDoc(url[0], ignores)
	if err != nil {
		s.Logger.Printf("failed to get minified doc: %s\n", err)
	}
	docHTML, err := doc.Html()
	if err != nil {
		s.Logger.Printf("failed to get html: %s\n", err)
	}
	HTML, err := s.Database.Queries.InsertHTML(
		ctx,
		master.InsertHTMLParams{Value: docHTML},
	)
	if err != nil {
		s.Logger.Printf("failed to insert html: %s\n", err)
	}
	URL, err := s.Database.Queries.InsertURL(
		ctx,
		master.InsertURLParams{Value: url[0], HtmlID: HTML.ID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert url: %w", err)
	}
	selectorStrs, err := parsers.GetAllSelectors(doc)
	if err != nil {
		s.Logger.Printf("failed to get selectors: %s\n", err)
	}
	for _, sel := range selectorStrs {
		context, err := doc.Find(sel).Parent().Html()
		if err != nil {
			s.Logger.Printf("failed to get html: %s\n", err)
		}
		if err != nil {
			s.Logger.Printf("failed to get urls: %s\n", err)
		}
		context = gohtml.Format(context)
		selector, err := s.Database.Queries.InsertSelector(
			ctx,
			master.InsertSelectorParams{
				Value:   sel,
				UrlID:   URL.ID,
				Context: context,
			},
		)
		if err != nil {
			s.Logger.Printf("failed to insert selector: %s\n", err)
		}
		selectors = append(selectors, *selector)
	}
	return selectors, nil
}
