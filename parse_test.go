package seltabl

import (
	"reflect"
	"testing"

	"github.com/conneroisu/seltabl/testdata"
	"github.com/stretchr/testify/assert"
)

// TestFixtureTables tests the parsing of a table with headers.
func TestFixtureTables(t *testing.T) {
	p, err := NewFromString[testdata.FixtureStruct](testdata.FixtureABNumTable)
	assert.Nil(t, err)
	assert.Equal(t, "1", p[0].A)
	assert.Equal(t, "2", p[0].B)
	assert.Equal(t, "3", p[1].A)
	assert.Equal(t, "4", p[1].B)
	assert.Equal(t, "5", p[2].A)
	assert.Equal(t, "6", p[2].B)
	assert.Equal(t, "7", p[3].A)
	assert.Equal(t, "8", p[3].B)
}

// TestNumberedTable tests the parsing of a table with numbered headers.
func TestNumberedTable(t *testing.T) {
	p, err := NewFromString[testdata.NumberedStruct](testdata.NumberedTable)
	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Equal(t, "Row 1, Cell 1", p[0].Header1)
	assert.Equal(t, "Row 1, Cell 2", p[0].Header2)
	assert.Equal(t, "Row 1, Cell 3", p[0].Header3)
	assert.Equal(t, "Row 2, Cell 1", p[1].Header1)
	assert.Equal(t, "Row 2, Cell 2", p[1].Header2)
	assert.Equal(t, "Row 2, Cell 3", p[1].Header3)
	assert.Equal(t, "Row 3, Cell 1", p[2].Header1)
	assert.Equal(t, "Row 3, Cell 2", p[2].Header2)
	assert.Equal(t, "Row 3, Cell 3", p[2].Header3)
}

// / TestNewFromString tests the NewFromString function
func TestNewFromString(t *testing.T) {
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
			want:    nil, // nil is an error
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFromString[testdata.SuperNovaStruct](tt.args.htmlInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromString() error = %v, wantErr %v", err, tt.wantErr)
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
}
