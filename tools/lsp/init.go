package lsp

// InitializeRequest is a struct for the initialize request
type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

// InitializeRequestParams is a struct for the initialize request params
type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	ID         int         `json:"id"`
	Params     struct {
		ProcessID        int    `json:"processId"`
		WorkDoneToken    string `json:"workDoneToken"`
		WorkspaceFolders any    `json:"workspaceFolders"`
		Trace            string `json:"trace"`
		RootPath         any    `json:"rootPath"`
		Capabilities     struct {
			General struct {
				PositionEncodings []string `json:"positionEncodings"`
			} `json:"general"`
			Window struct {
				ShowMessage struct {
					MessageActionItem struct {
						AdditionalPropertiesSupport bool `json:"additionalPropertiesSupport"`
					} `json:"messageActionItem"`
				} `json:"showMessage"`
				ShowDocument struct {
					Support bool `json:"support"`
				} `json:"showDocument"`
				WorkDoneProgress bool `json:"workDoneProgress"`
			} `json:"window"`
			TextDocument struct {
				RangeFormatting struct {
					DynamicRegistration bool `json:"dynamicRegistration"`
				} `json:"rangeFormatting"`
				CodeAction struct {
					DataSupport              bool `json:"dataSupport"`
					CodeActionLiteralSupport struct {
						CodeActionKind struct {
							ValueSet []string `json:"valueSet"`
						} `json:"codeActionKind"`
					} `json:"codeActionLiteralSupport"`
					ResolveSupport struct {
						Properties []string `json:"properties"`
					} `json:"resolveSupport"`
					DynamicRegistration bool `json:"dynamicRegistration"`
					IsPreferredSupport  bool `json:"isPreferredSupport"`
				} `json:"codeAction"`
				Diagnostic struct {
					DynamicRegistration bool `json:"dynamicRegistration"`
				} `json:"diagnostic"`
				DocumentHighlight struct {
					DynamicRegistration bool `json:"dynamicRegistration"`
				} `json:"documentHighlight"`
				References struct {
					DynamicRegistration bool `json:"dynamicRegistration"`
				} `json:"references"`
				Synchronization struct {
					WillSaveWaitUntil   bool `json:"willSaveWaitUntil"`
					DidSave             bool `json:"didSave"`
					WillSave            bool `json:"willSave"`
					DynamicRegistration bool `json:"dynamicRegistration"`
				} `json:"synchronization"`
				Completion struct {
					ContextSupport bool `json:"contextSupport"`
					CompletionList struct {
						ItemDefaults []string `json:"itemDefaults"`
					} `json:"completionList"`
					CompletionItemKind struct {
						ValueSet []int `json:"valueSet"`
					} `json:"completionItemKind"`
					DynamicRegistration bool `json:"dynamicRegistration"`
					CompletionItem      struct {
						DocumentationFormat     []string `json:"documentationFormat"`
						DeprecatedSupport       bool     `json:"deprecatedSupport"`
						PreselectSupport        bool     `json:"preselectSupport"`
						CommitCharactersSupport bool     `json:"commitCharactersSupport"`
						SnippetSupport          bool     `json:"snippetSupport"`
					} `json:"completionItem"`
				} `json:"completion"`
				Hover struct {
					DynamicRegistration bool     `json:"dynamicRegistration"`
					ContentFormat       []string `json:"contentFormat"`
				} `json:"hover"`
				TypeDefinition struct {
					LinkSupport bool `json:"linkSupport"`
				} `json:"typeDefinition"`
				Rename struct {
					DynamicRegistration bool `json:"dynamicRegistration"`
					PrepareSupport      bool `json:"prepareSupport"`
				} `json:"rename"`
				InlayHint struct {
					DynamicRegistration bool `json:"dynamicRegistration"`
					ResolveSupport      struct {
						Properties []string `json:"properties"`
					} `json:"resolveSupport"`
				} `json:"inlayHint"`
				SemanticTokens struct {
					AugmentsSyntaxTokens    bool     `json:"augmentsSyntaxTokens"`
					ServerCancelSupport     bool     `json:"serverCancelSupport"`
					MultilineTokenSupport   bool     `json:"multilineTokenSupport"`
					OverlappingTokenSupport bool     `json:"overlappingTokenSupport"`
					TokenTypes              []string `json:"tokenTypes"`
					Requests                struct {
						Full struct {
							Delta bool `json:"delta"`
						} `json:"full"`
						Range bool `json:"range"`
					} `json:"requests"`
					Formats             []string `json:"formats"`
					TokenModifiers      []string `json:"tokenModifiers"`
					DynamicRegistration bool     `json:"dynamicRegistration"`
				} `json:"semanticTokens"`
				DocumentSymbol struct {
					HierarchicalDocumentSymbolSupport bool `json:"hierarchicalDocumentSymbolSupport"`
					DynamicRegistration               bool `json:"dynamicRegistration"`
					SymbolKind                        struct {
						ValueSet []int `json:"valueSet"`
					} `json:"symbolKind"`
				} `json:"documentSymbol"`
				SignatureHelp struct {
					DynamicRegistration  bool `json:"dynamicRegistration"`
					SignatureInformation struct {
						ActiveParameterSupport bool     `json:"activeParameterSupport"`
						DocumentationFormat    []string `json:"documentationFormat"`
						ParameterInformation   struct {
							LabelOffsetSupport bool `json:"labelOffsetSupport"`
						} `json:"parameterInformation"`
					} `json:"signatureInformation"`
				} `json:"signatureHelp"`
				Definition struct {
					LinkSupport         bool `json:"linkSupport"`
					DynamicRegistration bool `json:"dynamicRegistration"`
				} `json:"definition"`
				Implementation struct {
					LinkSupport bool `json:"linkSupport"`
				} `json:"implementation"`
				CallHierarchy struct {
					DynamicRegistration bool `json:"dynamicRegistration"`
				} `json:"callHierarchy"`
				PublishDiagnostics struct {
					DataSupport        bool `json:"dataSupport"`
					RelatedInformation bool `json:"relatedInformation"`
					TagSupport         struct {
						ValueSet []int `json:"valueSet"`
					} `json:"tagSupport"`
				} `json:"publishDiagnostics"`
				Declaration struct {
					LinkSupport bool `json:"linkSupport"`
				} `json:"declaration"`
				Formatting struct {
					DynamicRegistration bool `json:"dynamicRegistration"`
				} `json:"formatting"`
			} `json:"textDocument"`
			Workspace struct {
				WorkspaceEdit struct {
					ResourceOperations []string `json:"resourceOperations"`
				} `json:"workspaceEdit"`
				DidChangeConfiguration struct {
					DynamicRegistration bool `json:"dynamicRegistration"`
				} `json:"didChangeConfiguration"`
				Configuration bool `json:"configuration"`
				InlayHint     struct {
					RefreshSupport bool `json:"refreshSupport"`
				} `json:"inlayHint"`
				SemanticTokens struct {
					RefreshSupport bool `json:"refreshSupport"`
				} `json:"semanticTokens"`
				ApplyEdit             bool `json:"applyEdit"`
				WorkspaceFolders      bool `json:"workspaceFolders"`
				DidChangeWatchedFiles struct {
					DynamicRegistration    bool `json:"dynamicRegistration"`
					RelativePatternSupport bool `json:"relativePatternSupport"`
				} `json:"didChangeWatchedFiles"`
				Symbol struct {
					DynamicRegistration bool `json:"dynamicRegistration"`
					SymbolKind          struct {
						ValueSet []int `json:"valueSet"`
					} `json:"symbolKind"`
				} `json:"symbol"`
			} `json:"workspace"`
		} `json:"capabilities"`
		RootURI    any `json:"rootUri"`
		ClientInfo struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"clientInfo"`
	} `json:"params"`
}

// ClientInfo is a struct for the client info
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResponse is a struct for the initialize response
type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

// InitializeResult is a struct for the initialize result
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

// ServerCapabilities is a struct for the server capabilities
type ServerCapabilities struct {
	TextDocumentSync int `json:"textDocumentSync"`

	HoverProvider      bool           `json:"hoverProvider"`
	DefinitionProvider bool           `json:"definitionProvider"`
	CodeActionProvider bool           `json:"codeActionProvider"`
	CompletionProvider map[string]any `json:"completionProvider"`
}

// ServerInfo is a struct for the server info
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// NewInitializeResponse returns a new initialize response
func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
				CodeActionProvider: true,
				CompletionProvider: map[string]any{},
			},
			ServerInfo: ServerInfo{
				Name:    "seltabl-lsp",
				Version: "0.0.0.0.0.0-beta1.final",
			},
		},
	}
}
