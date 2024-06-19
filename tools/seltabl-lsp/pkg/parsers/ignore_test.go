package parsers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnore(t *testing.T) {
	type args struct {
		ctx    context.Context
		src    string
		expect []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Ignore Valid Case",
			args: args{
				ctx: context.Background(),
				src: `
				// @ignore-elements: div, script
				`,
				expect: []string{"div", "script"},
			},
			wantErr: false,
		},
		{
			name: "Test Ignore Invalid Case",
			args: args{ctx: context.Background(), src: `
				// @ignore-elements: 
				`, expect: []string{}},
			wantErr: true,
		},
		{
			name: "Test Ignore Partially Invalid Case 2",
			args: args{ctx: context.Background(), src: `
				// @ignore-elements: div, script 
				`, expect: []string{"div", "script"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractIgnores(tt.args.ctx, tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractIgnores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.args.expect, got)
		})
	}
}
