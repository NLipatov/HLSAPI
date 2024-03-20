package Infrastructure

import (
	"bufio"
	"io"
	"os"
)

type DiskInteractor struct {
}

func (DiskInteractor) Store(writeTo io.Writer, readFrom io.Reader) error {
	_, err := io.Copy(writeTo, readFrom)
	if err != nil {
		return err
	}

	return nil
}

func (DiskInteractor) Get(writeTo io.Writer, path string) error {
	filePtr, err := os.Open(path)

	if err != nil {
		return err
	}

	defer filePtr.Close()

	reader := bufio.NewReader(filePtr)
	fInfo, err := filePtr.Stat()
	if err != nil {
		return err
	}
	_, err = io.CopyN(writeTo, reader, fInfo.Size())

	if err != nil {
		return err
	}

	return nil
}
