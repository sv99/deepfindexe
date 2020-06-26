package deepfindexe

import (
	"bytes"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/pgzip"
	"io/ioutil"
)

// Gz facilitates gzip compression.
type Gz struct {
	CompressionLevel int
}

// WalkByMime calls walkFn for each visited item in archive.
func (gz *Gz) WalkByMime(f File, walkFn WalkFunc) error {
	r, err := pgzip.NewReader(bytes.NewReader(*f.Buf))
	if err != nil {
		return err
	}
	defer r.Close()

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

func (gz *Gz) Close() error { return nil }

// NewGz returns a new, default instance ready to be customized and used.
func NewGz() *Gz {
	return &Gz{
		CompressionLevel: gzip.DefaultCompression,
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(Gz))
	_ = Closeable(new(Gz))
)
