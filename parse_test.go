package seltabl

import (
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
			t.Run(tt.name, func(t *testing.T) {
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
			t.Run(tt.name, func(t *testing.T) {
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
}
