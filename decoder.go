package seltabl

import "io"

// Decoder is a struct for decoding a reader into a slice of structs.
type Decoder[T any] struct {
	reader io.ReadCloser
}

// NewDecoder parses a reader into a slice of structs.
func NewDecoder[T any](r io.ReadCloser) *Decoder[T] {
	return &Decoder[T]{
		reader: r,
	}
}

// Decode parses a reader into a slice of structs.
func (d *Decoder[T]) Decode(value *T) ([]T, error) {
	defer d.reader.Close()
	return NewFromReader[T](d.reader)
}
