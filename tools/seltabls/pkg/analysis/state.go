package analysis

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/conneroisu/seltabl/tools/seltabls/data"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/mitchellh/go-homedir"
	"github.com/yosssi/gohtml"
)

// State is the state of the document analysis.
type State struct {
	// Map of file names to contents.
	Documents map[string]string
	// Selectors is the map of file names to selectors.
	Selectors map[string][]master.Selector
	// URLs is the map of file names to urls.
	URLs map[string][]string
	// Database is the database for the state.
	Database data.Database[master.Queries]
	// ResponseCtxs is the map of to response contexts indexed by grpc request id.
	ResponseCtxs map[int]string
	// Logger is the logger for the state.
	Logger *log.Logger
}

// NewState returns a new state with no documents
func NewState() (state State, err error) {
	ctx := context.Background()
	configPath, err := CreateConfigDir("~/.config/seltabls/")
	if err != nil {
		return state, fmt.Errorf("failed to create config directory: %w", err)
	}
	db, err := data.NewDb(
		ctx,
		master.New,
		&data.Config{
			Schema:   master.MasterSchema,
			URI:      "sqlite://uri.sqlite",
			FileName: path.Join(configPath, "uri.sqlite"),
		},
	)
	if err != nil {
		return state, fmt.Errorf("failed to create database: %w", err)
	}
	logger := getLogger(path.Join(configPath, "state.log"))
	state = State{
		Documents: make(map[string]string),
		Selectors: make(map[string][]master.Selector),
		Database:  *db,
		Logger:    logger,
		URLs:      make(map[string][]string),
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
	return log.New(logFile, "[seltabls#state]", log.LstdFlags)
}

// GetSelectors gets all the selectors from the given URL and appends them to the selectors slice
func GetSelectors(
	ctx context.Context,
	state *State,
	url string,
	ignores []string,
) (selectors []master.Selector, err error) {
	rows, err := state.Database.Queries.GetSelectorsByURL(
		ctx,
		master.GetSelectorsByURLParams{Value: url},
	)
	if err == nil && rows != nil {
		var selectors []master.Selector
		for _, row := range rows {
			selectors = append(selectors, *row)
		}
		return selectors, nil
	}
	doc, err := parsers.GetMinifiedDoc(url, ignores)
	if err != nil {
		return nil, fmt.Errorf("failed to get minified doc: %w", err)
	}
	docHTML, err := doc.Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get html: %w", err)
	}
	HTML, err := state.Database.Queries.InsertHTML(
		ctx,
		master.InsertHTMLParams{Value: docHTML},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert html: %w", err)
	}
	URL, err := state.Database.Queries.InsertURL(
		ctx,
		master.InsertURLParams{Value: url, HtmlID: HTML.ID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert url: %w", err)
	}
	selectorStrs, err := parsers.GetAllSelectors(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to get selectors: %w", err)
	}
	for _, sel := range selectorStrs {
		context, err := doc.Find(sel).Parent().Find(":nth-child(1)").Html()
		if err != nil {
			state.Logger.Printf("failed to get html: %s\n", err)
		}
		context = gohtml.Format(context)
		selector, err := state.Database.Queries.InsertSelector(
			ctx,
			master.InsertSelectorParams{
				Value:   sel,
				UrlID:   URL.ID,
				Context: context,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert selector: %w", err)
		}
		selectors = append(selectors, *selector)
	}
	return selectors, nil
}

// CreateConfigDir creates a new config directory and returns the path.
func CreateConfigDir(dirPath string) (string, error) {
	path, err := homedir.Expand(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to expand home directory: %w", err)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		if os.IsExist(err) {
			return path, nil
		}
		return "", fmt.Errorf(
			"failed to create or find config directory: %w",
			err,
		)
	}
	return path, nil
}
