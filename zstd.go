package deepfindexe

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/klauspost/compress/zstd"
)

// Zstd facilitates Zstandard compression.
type Zstd struct {
}

// Decompress reads in, decompresses it, and writes it to out.
func (zs *Zstd) Decompress(in io.Reader, out io.Writer) error {
	r, err := zstd.NewReader(in)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(out, r)
	return err
}

// WalkByMime calls walkFn for each visited item in archive.
func (gz *Zstd) WalkByMime(f File, walkFn WalkFunc) error {
	r, err := zstd.NewReader(bytes.NewReader(*f.Buf))
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

func (zs *Zstd) Close() error { return nil }

// NewZstd returns a new, default instance ready to be customized and used.
func NewZstd() *Zstd {
	return new(Zstd)
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(Zstd))
	_ = Closeable(new(Zstd))
)
