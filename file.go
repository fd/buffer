package buffer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/djherbis/buffer/wrapio"
)

type File interface {
	Name() string
	Stat() (fi os.FileInfo, err error)
	io.ReaderAt
	io.WriterAt
	Close() error
}

type fileBuffer struct {
	file File
	*wrapio.Wrapper
}

func NewFile(N int64, file File) Buffer {
	return &fileBuffer{
		file:    file,
		Wrapper: wrapio.NewWrapper(file, N),
	}
}

func init() {
	gob.Register(&fileBuffer{})
}

func (buf *fileBuffer) MarshalBinary() ([]byte, error) {
	fullpath, err := filepath.Abs(filepath.Dir(buf.file.Name()))
	if err != nil {
		return nil, err
	}
	base := filepath.Base(buf.file.Name())
	buf.file.Close()

	buffer := bytes.NewBuffer(nil)
	fmt.Fprintln(buffer, filepath.Join(fullpath, base))
	fmt.Fprintln(buffer, buf.Wrapper.N, buf.Wrapper.L, buf.Wrapper.O)
	return buffer.Bytes(), nil
}

func (buf *fileBuffer) UnmarshalBinary(data []byte) error {
	buffer := bytes.NewBuffer(data)
	var filename string
	var N, L, O int64
	_, err := fmt.Fscanln(buffer, &filename)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	buf.file = file

	_, err = fmt.Fscanln(buffer, &N, &L, &O)
	buf.Wrapper = wrapio.NewWrapper(file, N)
	buf.Wrapper.L = L
	buf.Wrapper.O = O
	return nil
}
