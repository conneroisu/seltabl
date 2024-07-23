package seltabl

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/testdata"
	"github.com/stretchr/testify/assert"
)

// NoDataSelectorStruct is a test struct
type NoDataSelectorStruct struct {
	A string `seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" cSel:"$text"`
	B string `seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" cSel:"$text"`
}

// TestNewFromString tests the NewFromString function
// for all the different types of tables that we have in the
// testdata package.
func TestNewFromString(t *testing.T) {
	t.Run(
		"Test that NewFromString returns the correct result when the html is valid with the SuperNova table",
		func(t *testing.T) {
			t.Parallel()
			type args struct {
				htmlInput string
				typ       interface{}
			}
			tests := []struct {
				name    string
				args    args
				want    interface{}
				wantErr bool
			}{
				{
					name: "Test that NewFromString returns the correct result when the html is valid",
					args: args{
						htmlInput: testdata.SuperNovaTable,
						typ:       reflect.TypeOf(testdata.SuperNovaStruct{}),
					},
					want:    testdata.SuperNovaTableResult,
					wantErr: false,
				},
				{
					name: "TestNewFromStringWithInvalidHTML",
					args: args{
						htmlInput: "invalid",
						typ:       reflect.TypeOf(testdata.SuperNovaStruct{}),
					},
					want:    nil,
					wantErr: true,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					doc, err := goquery.NewDocumentFromReader(
						strings.NewReader(tt.args.htmlInput),
					)
					assert.Nil(t, err)
					got, err := New[testdata.SuperNovaStruct](doc)
					if (err != nil) != tt.wantErr {
						t.Errorf(
							"NewFromString() error = %v, wantErr %v",
							err,
							tt.wantErr,
						)
						return
					}
					if tt.wantErr {
						assert.Error(t, err)
						return
					}
					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf(
							"NewFromString() got = %v, want %v",
							got,
							tt.want,
						)
					}
				})
			}
		},
	)
	t.Run("Numbered Table", func(t *testing.T) {
		t.Parallel()
		type args struct {
			htmlInput string
			typ       interface{}
		}
		tests := []struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			{
				name: "Test numbered table",
				args: args{
					htmlInput: testdata.NumberedTable,
					typ:       reflect.TypeOf(testdata.NumberedStruct{}),
				},
				want:    testdata.NumberedTableResult,
				wantErr: false,
			},
			{
				name: "Test that NewFromString returns an error when the html is invalid",
				args: args{
					htmlInput: "invalid",
					typ:       reflect.TypeOf(testdata.NumberedStruct{}),
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.args.htmlInput),
				)
				assert.Nil(t, err)
				got, err := New[testdata.NumberedStruct](doc)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewFromString() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewFromString() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
	t.Run(
		"Test correct operation with testdata fixture tables",
		func(t *testing.T) {
			t.Parallel()
			type args struct {
				htmlInput string
				typ       interface{}
			}
			tests := []struct {
				name    string
				args    args
				want    interface{}
				wantErr bool
			}{
				{
					name: "Test fixture table",
					args: args{
						htmlInput: testdata.FixtureABNumTable,
						typ:       reflect.TypeOf(testdata.FixtureStruct{}),
					},
					want:    testdata.FixtureABNumTableResult,
					wantErr: false,
				},
				{
					name: "TestNewFromStringWithInvalidHTML",
					args: args{
						htmlInput: "invalid",
						typ:       reflect.TypeOf(testdata.FixtureStruct{}),
					},
					want:    nil,
					wantErr: true,
				},
			}
			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					doc, err := goquery.NewDocumentFromReader(
						strings.NewReader(tt.args.htmlInput),
					)
					assert.Nil(t, err)
					got, err := New[testdata.FixtureStruct](doc)
					if (err != nil) != tt.wantErr {
						t.Errorf(
							"NewFromString() error = %v, wantErr %v",
							err,
							tt.wantErr,
						)
						return
					}
					if tt.wantErr {
						assert.Error(t, err)
						return
					}
					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf(
							"NewFromString() got = %v, want %v",
							got,
							tt.want,
						)
					}
				})
			}
		},
	)

	t.Run(
		"Test that NewFromString returns an error when the html is invalid",
		func(t *testing.T) {
			t.Parallel()
			type args struct {
				htmlInput string
				typ       interface{}
			}
			tests := []struct {
				name    string
				args    args
				want    interface{}
				wantErr bool
			}{
				{
					name: "TestNewFromStringWithInvalidHTML",
					args: args{
						htmlInput: "invalid",
						typ:       reflect.TypeOf(testdata.FixtureStruct{}),
					},
					want:    nil,
					wantErr: true,
				},
			}
			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					doc, err := goquery.NewDocumentFromReader(
						strings.NewReader(tt.args.htmlInput),
					)
					assert.Nil(t, err)
					_, err = New[testdata.FixtureStruct](doc)
					if (err != nil) != tt.wantErr {
						t.Errorf(
							"NewFromString() error = %v, wantErr %v",
							err,
							tt.wantErr,
						)
						return
					}
					if tt.wantErr {
						assert.Error(t, err)
						return
					}
				})
			}
		},
	)

	t.Run("TestNewFromStringWithInvalidJSON", func(t *testing.T) {
		t.Parallel()
		type args struct {
			htmlInput string
			typ       interface{}
		}
		tests := []struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			{
				name: "TestNewFromStringWithInvalidJSON",
				args: args{
					htmlInput: "invalid",
					typ:       reflect.TypeOf(testdata.FixtureStruct{}),
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.args.htmlInput),
				)
				assert.Nil(t, err)
				_, err = New[testdata.FixtureStruct](doc)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewFromString() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
			})
		}
	})

	t.Run("TestNewFromStringWithInvalidHTML", func(t *testing.T) {
		t.Parallel()
		type args struct {
			htmlInput string
			typ       interface{}
		}
		tests := []struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			{
				name: "TestNewFromStringWithInvalidHTML",
				args: args{
					htmlInput: "invalid",
					typ:       reflect.TypeOf(testdata.FixtureStruct{}),
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.args.htmlInput),
				)
				assert.Nil(t, err)
				_, err = New[testdata.FixtureStruct](doc)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewFromString() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
			})
		}
	})

	// NoHeaderSelectorStruct is a test struct
	type NoHeaderSelectorStruct struct {
		A string `json:"a" seltabl:"a" dSel:"tr td:nth-child(1)" cSel:"$text"`
		B string `json:"b" seltabl:"b" dSel:"tr td:nth-child(2)" cSel:"$text"`
	}

	t.Run("TestNewFromStringWithNoHeaderSelector", func(t *testing.T) {
		t.Parallel()
		type args struct {
			htmlInput string
			typ       interface{}
		}
		tests := []struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			{
				name: "TestNewFromStringWithInvalidHTML",
				args: args{
					htmlInput: "invalid",
					typ:       reflect.TypeOf(NoHeaderSelectorStruct{}),
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.args.htmlInput),
				)
				assert.Nil(t, err)
				got, err := New[NoHeaderSelectorStruct](doc)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewFromString() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewFromString() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
	t.Run("TestNewFromStringWithNoDataSelector", func(t *testing.T) {
		t.Parallel()
		type args struct {
			htmlInput string
			typ       interface{}
		}
		tests := []struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			{
				name: "TestNewFromStringWithNoDataSelector",
				args: args{
					htmlInput: testdata.FixtureABNumTable,
					typ:       reflect.TypeOf(NoDataSelectorStruct{}),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "TestNewFromStringWithInvalidHTML",
				args: args{
					htmlInput: "invalid",
					typ:       reflect.TypeOf(NoDataSelectorStruct{}),
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.args.htmlInput),
				)
				assert.Nil(t, err)
				got, err := New[NoDataSelectorStruct](doc)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewFromString() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewFromString() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// NoCellSelectorStruct is a test struct
	type NoCellSelectorStruct struct {
		A string `json:"a" seltabl:"a" dSel:"tr td:nth-child(1)" cSel:""`
		B string `json:"b" seltabl:"b" dSel:"tr td:nth-child(2)" cSel:""`
	}

	t.Run("TestNewFromStringWithNoCellSelector", func(t *testing.T) {
		t.Parallel()
		type args struct {
			htmlInput string
			typ       interface{}
		}
		tests := []struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			{
				name: "TestNewFromStringWithInvalidHTML",
				args: args{
					htmlInput: "invalid",
					typ:       reflect.TypeOf(NoCellSelectorStruct{}),
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.args.htmlInput),
				)
				assert.Nil(t, err)
				got, err := New[NoCellSelectorStruct](doc)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewFromString() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewFromString() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// InvalidTGenericType is a test struct that is invalid.
	type InvalidTGenericType func(a, b int) int

	t.Run(
		"Test that NewFromString returns an error when the generic type is invalid like when the test passes in a function",
		func(t *testing.T) {
			t.Parallel()
			type args struct {
				htmlInput string
				typ       interface{}
			}
			tests := []struct {
				name    string
				args    args
				want    interface{}
				wantErr bool
			}{
				{
					name: "TestNewFromStringWithInvalidGenericType",
					args: args{
						htmlInput: "invalid",
						typ: reflect.TypeOf(
							func(a, b int) int { return a + b },
						),
					},
					want:    nil,
					wantErr: true,
				},
				{
					name: "Test that NewFromString returns an error when the generic type is invalid like when the test passes in a function",
					args: args{
						htmlInput: "invalid",
						typ: reflect.TypeOf(
							func(b int) int { return b },
						),
					},
					want:    nil,
					wantErr: true,
				},
			}
			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					doc, err := goquery.NewDocumentFromReader(
						strings.NewReader(tt.args.htmlInput),
					)
					assert.Nil(t, err)
					_, err = New[InvalidTGenericType](doc)
					if (err != nil) != tt.wantErr {
						t.Errorf(
							"NewFromString() error = %v, wantErr %v",
							err,
							tt.wantErr,
						)
						return
					}
					if tt.wantErr {
						assert.Error(t, err)
						return
					}
				})
			}
		},
	)

	// test a struct with no seltabl field or blank one.
	type NoSeltablField struct {
		A string `json:"a" dSel:"tr td:nth-child(1)" cSel:"$text"`
		B string `json:"b" dSel:"tr td:nth-child(2)" cSel:"$text" seltabl:"b"`
	}

	t.Run("TestNewFromStringWithNoSeltablField", func(t *testing.T) {
		t.Parallel()
		type args struct {
			htmlInput string
			typ       interface{}
		}
		tests := []struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			{
				name: "TestNewFromStringWithInvalidHTML",
				args: args{
					htmlInput: "invalid",
					typ:       reflect.TypeOf(NoSeltablField{}),
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.args.htmlInput),
				)
				assert.Nil(t, err)
				_, err = New[NoSeltablField](doc)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewFromString() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
			})
		}
	})
}

// TestNewFromUrl tests the NewFromURL function.
func TestNewFromUrl(t *testing.T) {
	t.Run("TestNewFromUrl", func(t *testing.T) {
		t.Parallel()
		type args struct {
			url string
		}
		tests := []struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			{
				name: "TestNewFromUrlWithInvalidURL",
				args: args{
					url: "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := NewFromURL[testdata.FixtureStruct](tt.args.url)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewFromURL() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewFromURL() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
	// test a struct with no seltabl field or blank one
	type NoSeltablField struct {
		A string `json:"a" dSel:"tr td:nth-child(1)" cSel:"$text"`
		B string `json:"b" dSel:"tr td:nth-child(2)" cSel:"$text" seltabl:"b"`
	}
	t.Run("TestNewFromUrlWithNoSeltablField", func(t *testing.T) {
		t.Parallel()
		type args struct {
			url string
		}
		tests := []struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			{
				name: "TestNewFromUrlWithInvalidURL",
				args: args{
					url: "ttps://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := NewFromURL[NoSeltablField](tt.args.url)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewFromURL() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewFromURL() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// test with a undecodable body in respinse to NewFromURL
	t.Run("TestNewFromUrlWithInvalidBody", func(t *testing.T) {
		t.Parallel()
		// Mock HTTP server for testing
		server := httpTestServer(t, ":86;&")
		defer server.Close()

		type args struct {
			url string
		}
		tests := []struct {
			name    string
			args    args
			want    interface{}
			wantErr bool
		}{
			{
				name: "TestNewFromUrlWithInvalidBody",
				args: args{
					url: server.URL,
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "TestNewFromUrlWithInvalidURL",
				args: args{
					url: "http//github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
				},
				want:    nil,
				wantErr: true,
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := NewFromURL[testdata.FixtureStruct](tt.args.url)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"NewFromURL() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewFromURL() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

// httpTestServer sets up a test HTTP server
func httpTestServer(t *testing.T, body string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write([]byte(body))
			if err != nil {
				t.Fatalf("Failed to write response: %v", err)
			}
		}),
	)
}

// Helper function to create a goquery document from a string
func createDocFromString(htmlStr string) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
}

// TestieStruct is a test struct that is used for testing. It
// has a field with the tag hSel, dSel, and cSel.
type TestieStruct struct {
	A string `json:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr:not(:first-child) td:nth-child(1)" cSel:"$text"`
	B string `json:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr:not(:first-child) td:nth-child(2)" cSel:"$text"`
}

// TestNew_ValidStruct tests the New function with a valid struct.
func TestNew_ValidStruct(t *testing.T) {
	t.Parallel()
	html := `
		<table>
			<tr> <td>a</td> <td>b</td> </tr>
			<tr> <td>1</td> <td>2</td> </tr>
			<tr> <td>3</td> <td>4</td> </tr>
			<tr> <td>5</td> <td>6</td> </tr>
			<tr> <td>7</td> <td>8</td> </tr>
		</table>`
	doc, _ := createDocFromString(html)
	result, err := New[TestieStruct](doc)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, "1", result[0].A)
	assert.Equal(t, "2", result[0].B)
}

// TestNew_InvalidType tests the New function with an invalid type.
func TestNew_InvalidType(t *testing.T) {
	html := `
		<table>
			<tr> <td>a</td> <td>b</td> </tr>
			<tr> <td>1</td> <td>2</td> </tr>
		</table>`
	doc, _ := createDocFromString(html)
	type InvalidType int
	_, err := New[InvalidType](doc)
	assert.Error(t, err)
}

// TestNew_MissingSelectors tests the New function with a struct that is missing selectors.
func TestNew_MissingSelectors(t *testing.T) {
	html := `
		<table>
			<tr> <td>a</td> <td>b</td> </tr>
			<tr> <td>1</td> <td>2</td> </tr>
		</table>`
	doc, _ := createDocFromString(html)
	type MissingSelectors struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	_, err := New[MissingSelectors](doc)
	assert.Error(t, err)
}

// TestNew_MissingMustBePresent tests the New function with a struct that is missing a mustBePresent selector.
func TestNew_MissingMustBePresent(t *testing.T) {
	html := `
		<table>
			<tr> <td>a</td> <td>b</td> </tr>
			<tr> <td>1</td> <td>2</td> </tr>
		</table>`
	doc, _ := createDocFromString(html)
	type MustBePresentStruct struct {
		A string `json:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr:not(:first-child) td:nth-child(1)" cSel:"$text" mustBePresent:"c"`
	}
	_, err := New[MustBePresentStruct](doc)
	assert.Error(t, err)
}

// TestNew_NoDataFound tests the New function with a struct that does not have data.
func TestNew_NoDataFound(t *testing.T) {
	html := `
		<table>
			<tr> <td>a</td> <td>b</td> </tr>
		</table>`
	doc, _ := createDocFromString(html)
	_, err := New[TestieStruct](doc)
	assert.Error(t, err)
}

// TestNewFromString_ValidHTML tests the NewFromString function with valid HTML.
func TestNewFromString_ValidHTML(t *testing.T) {
	t.Parallel()
	html := `
		<table>
			<tr> <td>a</td> <td>b</td> </tr>
			<tr> <td>1</td> <td>2</td> </tr>
			<tr> <td>3</td> <td>4</td> </tr>
			<tr> <td>5</td> <td>6</td> </tr>
			<tr> <td>7</td> <td>8</td> </tr>
		</table>`
	result, err := NewFromString[TestieStruct](html)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, "1", result[0].A)
	assert.Equal(t, "2", result[0].B)
}

// TestNewFromString_InvalidHTML tests the NewFromString function with invalid HTML.
func TestNewFromString_InvalidHTML(t *testing.T) {
	html := `
		<table>
			<tr> <td>a</td> <td>b</td> </tr>
			<tr> <td>1</td> <td>2</td> </tr>
		`
	_, err := NewFromString[TestieStruct](html)
	assert.Error(t, err)
}

// TestNewFromReader_ValidReader tests the NewFromReader function with a
func TestNewFromReader_ValidReader(t *testing.T) {
	t.Parallel()
	html := `
		<table>
			<tr> <td>a</td> <td>b</td> </tr>
			<tr> <td>1</td> <td>2</td> </tr>
			<tr> <td>3</td> <td>4</td> </tr>
			<tr> <td>5</td> <td>6</td> </tr>
			<tr> <td>7</td> <td>8</td> </tr>
		</table>`
	reader := strings.NewReader(html)
	result, err := NewFromReader[TestieStruct](reader)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, "1", result[0].A)
	assert.Equal(t, "2", result[0].B)
}

// TestNewFromReader_InvalidReader tests the NewFromReader function with an invalid reader.
func TestNewFromReader_InvalidReader(t *testing.T) {
	t.Parallel()
	reader := strings.NewReader("<html><body><table>")
	_, err := NewFromReader[TestieStruct](reader)
	assert.Error(t, err)
}

// TestNewFromURL_ValidURL tests the NewFromURL function with a valid URL.
func TestNewFromURL_ValidURL(t *testing.T) {
	t.Parallel()
	html := `
		<table>
			<tr> <td>a</td> <td>b</td> </tr>
			<tr> <td>1</td> <td>2</td> </tr>
			<tr> <td>3</td> <td>4</td> </tr>
			<tr> <td>5</td> <td>6</td> </tr>
			<tr> <td>7</td> <td>8</td> </tr>
		</table>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, html)
	}))
	defer server.Close()

	result, err := NewFromURL[TestieStruct](server.URL)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, "1", result[0].A)
	assert.Equal(t, "2", result[0].B)
}

// TestNewFromURL_InvalidURL tests the NewFromURL function with an invalid URL.
func TestNewFromURL_InvalidURL(t *testing.T) {
	t.Parallel()
	_, err := NewFromURL[TestieStruct]("http://invalid.url")
	assert.Error(t, err)
}

// TestNewFromURL_InvalidHTMLContent tests the NewFromURL function with invalid HTML content.
func TestNewFromURL_InvalidHTMLContent(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "<html><body><table>")
	}))
	defer server.Close()

	_, err := NewFromURL[TestieStruct](server.URL)
	assert.Error(t, err)
}
