package deepfindexe

import (
	"bytes"
	"github.com/golang/snappy"
	"io"
	"io/ioutil"
)

// Snappy facilitates Snappy compression.
type Snappy struct{}

// Decompress reads in, decompresses it, and writes it to out.
func (s *Snappy) Decompress(in io.Reader, out io.Writer) error {
	r := snappy.NewReader(in)
	_, err := io.Copy(out, r)
	return err
}

// WalkByMime calls walkFn for each visited item in archive.
func (gz *Snappy) WalkByMime(f File, walkFn WalkFunc) error {
	r := snappy.NewReader(bytes.NewReader(*f.Buf))

	newBuf, err :=  ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = walkFn(File{
		FileName: TrimSuffix(f.FileName),
		Buf: &newBuf,
	})
	return err
}

func (s *Snappy) Close() error { return nil }

// NewSnappy returns a new, default instance ready to be customized and used.
func NewSnappy() *Snappy {
	return new(Snappy)
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(Snappy))
	_ = Closeable(new(Snappy))
)
