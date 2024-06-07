package seltabl

import (
	"fmt"
	"io"
)

// Decoder is a struct for decoding a reader into a slice of structs.
//
// It is used by the NewDecoder function.
//
// It is not intended to be used directly.
//
// Example:
//
//	type TableStruct struct {
//		A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
//		B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
//	}
//
//	func main() {
//		r := strings.NewReader(`
//		<table>
//			<tr>
//				<td>a</td>
//				<td>b</td>
//			</tr>
//			<tr>
//				<td> 1 </td>
//				<td>2</td>
//			</tr>
//			<tr>
//				<td>3 </td>
//				<td> 4</td>
//			</tr>
//			<tr>
//				<td> 5 </td>
//				<td> 6</td>
//			</tr>
//			<tr>
//				<td>7 </td>
//				<td> 8</td>
//			</tr>
//		</table>
//		`)
//		p, err := seltabl.NewDecoder[TableStruct](r)
//		if err != nil {
//			panic(err)
//		}
//		for _, pp := range p {
//			fmt.Printf("pp %+v\n", pp)
//		}
//	}
type Decoder[T any] struct {
	reader io.ReadCloser
}

// NewDecoder parses a reader into a slice of structs.
//
// It is used by the NewFromReader function.
//
// This allows for decoding a reader into a slice of structs.
//
// Similar to the json.Decoder for brevity.
func NewDecoder[T any](r io.ReadCloser) *Decoder[T] {
	return &Decoder[T]{
		reader: r,
	}
}

// Decode parses a reader into a slice of structs.
//
// It is used by the Decoder.Decode function.
//
// This allows for decoding a reader into a slice of structs.
//
// Similar to the json.Decoder for brevity.
func (d *Decoder[T]) Decode() ([]T, error) {
	defer d.reader.Close()
	var result []T
	result, err := NewFromReader[T](d.reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}
	return result, nil
}
