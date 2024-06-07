package seltabl

import (
	"io"
	"strings"
	"testing"
)

// DecodeExStruct is a test struct
type DecodeExStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
}

// TestDecoder_Decode tests the Decoder.Decode function
func TestDecoder_Decode(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []DecodeExStruct
		hasError bool
	}{
		{
			name: "Valid input",
			input: `
                <table>
                    <tr>
                        <td>a</td>
                        <td>b</td>
                    </tr>
                    <tr>
                        <td>1</td>
                        <td>2</td>
                    </tr>
                    <tr>
                        <td>3</td>
                        <td>4</td>
                    </tr>
                </table>
            `,
			expected: []DecodeExStruct{
				{A: "1", B: "2"},
				{A: "3", B: "4"},
			},
			hasError: false,
		},
		{
			name: "Invalid input",
			input: `
                <table>
                    <tr>
                        <td>a</td>
                        <td>b</td>
                    </tr>
                    <tr>
                        <td>1</td>
                    </tr>
                </table>
            `,
			expected: nil,
			hasError: true,
		},
		{
			name: "Invalid input with invalid html",
			input: `
                <table>
                    <tr>
                        <td>a</td>
                        <td>b</td>
                    </tr>
                    <tr>
                        <td>1</td>
                    </tr>
                </table>
            `,
			expected: nil,
			hasError: true,
		},
		{
			name: "Invalid input with invalid json",
			input: `
                <table>
                    <tr>
                        <td>a</td>
                        <td>b</td>
                    </tr>
                    <tr>
                        <td>1</td>
                    </tr>
                </table>
            `,
			expected: nil,
			hasError: true,
		},
		{
			name: "Invalid input with invalid json",
			input: `
                <table>
                    <tr>
                        <td>a</td>
                        <td>b</td>
                        <td>1</td>
                    </tr>
                </table>
            `,
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := io.NopCloser(strings.NewReader(tc.input))
			decoder := NewDecoder[DecodeExStruct](r)
			result, err := decoder.Decode()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(result) != len(tc.expected) {
				t.Errorf(
					"Expected %d results, but got %d",
					len(tc.expected),
					len(result),
				)
			}

			for i, expected := range tc.expected {
				if result[i].A != expected.A || result[i].B != expected.B {
					t.Errorf("Expected %+v, but got %+v", expected, result[i])
				}
			}
		})
	}
}
