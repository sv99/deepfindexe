package deepfindexe

import (
	"bytes"
	"github.com/dsnet/compress/bzip2"
	"io/ioutil"
)

// Bz2 facilitates bzip2 compression.
type Bz2 struct {
	CompressionLevel int
}

// WalkByMime calls walkFn for each visited item in archive.
func (bz *Bz2) WalkByMime(f File, walkFn WalkFunc) error {
	r, err := bzip2.NewReader(bytes.NewReader(*f.Buf), nil)
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

func (bz *Bz2) Close() error { return nil }

// NewBz2 returns a new, default instance ready to be customized and used.
func NewBz2() *Bz2 {
	return &Bz2{
		CompressionLevel: bzip2.DefaultCompression,
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(Bz2))
	_ = Closeable(new(Bz2))
)
