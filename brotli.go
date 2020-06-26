package deepfindexe

import (
	"bytes"
	"github.com/andybalholm/brotli"
	"io/ioutil"
)

// Bz2 facilitates bzip2 compression.
type Brotli struct {
	Quality int
}

// WalkByMime calls walkFn for each visited item in archive.
func (bz *Brotli) WalkByMime(f File, walkFn WalkFunc) error {
	r := brotli.NewReader(bytes.NewReader(*f.Buf))

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

func (bz *Brotli) Close() error { return nil }

// NewBz2 returns a new, default instance ready to be customized and used.
func NewBrotli() *Brotli {
	return &Brotli{
		Quality: brotli.DefaultCompression,
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(Bz2))
	_ = Closeable(new(Bz2))
)
