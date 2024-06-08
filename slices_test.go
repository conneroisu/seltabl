package seltabl

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/testdata"
	"github.com/stretchr/testify/assert"
)

// TestNewFromString tests the NewFromString function
// for all the different types of tables that we have in the
// testdata package
func TestNewFromString(t *testing.T) {
	t.Run("SuperNova", func(t *testing.T) {
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
				name: "Test supernova table",
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
					t.Errorf("NewFromString() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

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
				name: "TestNewFromStringWithInvalidHTML",
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

	t.Run("Fixture Tables", func(t *testing.T) {
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
					t.Errorf("NewFromString() got = %v, want %v", got, tt.want)
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

	type NoHSelector struct {
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
				name: "TestNewFromStringWithNoHeaderSelector",
				args: args{
					htmlInput: testdata.FixtureABNumTable,
					typ:       reflect.TypeOf(NoHSelector{}),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "TestNewFromStringWithInvalidHTML",
				args: args{
					htmlInput: "invalid",
					typ:       reflect.TypeOf(NoHSelector{}),
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
				got, err := New[NoHSelector](doc)
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

	type NoDataSelector struct {
		A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" cSel:"$text"`
		B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" cSel:"$text"`
	}

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
					typ:       reflect.TypeOf(NoDataSelector{}),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "TestNewFromStringWithInvalidHTML",
				args: args{
					htmlInput: "invalid",
					typ:       reflect.TypeOf(NoDataSelector{}),
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
				got, err := New[NoDataSelector](doc)
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

	type NoCellSelector struct {
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
				name: "TestNewFromStringWithNoCellSelector",
				args: args{
					htmlInput: testdata.FixtureABNumTable,
					typ:       reflect.TypeOf(NoCellSelector{}),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "TestNewFromStringWithInvalidHTML",
				args: args{
					htmlInput: "invalid",
					typ:       reflect.TypeOf(NoCellSelector{}),
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
				got, err := New[NoCellSelector](doc)
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

	type InvalidTGenericType func(a, b int) int

	t.Run("TestNewFromStringWithInvalidGenericType", func(t *testing.T) {
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
					typ:       reflect.TypeOf(func(a, b int) int { return a + b }),
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
	})

	// test a struct with no seltabl field or blank one
	type NoSeltablField struct {
		A string `json:"a" dSel:"tr td:nth-child(1)" cSel:"$text"`
		B string `json:"b" seltabl:"b" dSel:"tr td:nth-child(2)" cSel:"$text"`
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
				name: "TestNewFromStringWithNoSeltablField",
				args: args{
					htmlInput: testdata.FixtureABNumTable,
					typ:       reflect.TypeOf(NoSeltablField{}),
				},
				want:    nil,
				wantErr: true,
			},
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
		B string `json:"b" seltabl:"b" dSel:"tr td:nth-child(2)" cSel:"$text"`
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
