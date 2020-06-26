package deepfindexe

import (
	"bytes"
	fastxz "github.com/xi2/xz"
	"io"
	"io/ioutil"
)

// Xz facilitates XZ compression.
type Xz struct{}

// Decompress reads in, decompresses it, and writes it to out.
func (x *Xz) Decompress(in io.Reader, out io.Writer) error {
	r, err := fastxz.NewReader(in, 0)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, r)
	return err
}

// WalkByMime calls walkFn for each visited item in archive.
func (gz *Xz) WalkByMime(f File, walkFn WalkFunc) error {
	r, err := fastxz.NewReader(bytes.NewReader(*f.Buf), 0)
	if err != nil {
		return err
	}

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

func (x *Xz) Close() error { return nil }

// NewXz returns a new, default instance ready to be customized and used.
func NewXz() *Xz {
	return new(Xz)
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(Xz))
	_ = Closeable(new(Xz))
)
