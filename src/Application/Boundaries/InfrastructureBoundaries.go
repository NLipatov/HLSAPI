package Boundaries

import "io"

type StoreBoundary interface {
	Store(writeTo io.Writer, readFrom io.Reader) error
}

type GetBoundary interface {
	Get(writeTo io.Writer, path string) error
}

type MediaConverterBoundary interface {
	ConvertToM3U8(inputFilename string, outputFilename string) error
}
