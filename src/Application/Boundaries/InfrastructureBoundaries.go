package Boundaries

import "io"

type StoreBoundary interface {
	Store(writeTo io.Writer, readFrom io.Reader) error
}

type GetBoundary interface {
	Get(writeTo io.Writer, path string) error
}

type MediaConverterBoundary interface {
	ConvertToM3U8(workdirAbsolutePath string, inputVideoFilename string) (string, error)
}

type EnvironmentBoundary interface {
	GetAppRootPath() string
}
