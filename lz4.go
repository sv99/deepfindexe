package deepfindexe

import (
	"bytes"
	"github.com/pierrec/lz4/v3"
	"io"
	"io/ioutil"
)

// Lz4 facilitates LZ4 compression.
type Lz4 struct {
	CompressionLevel int
}

// Decompress reads in, decompresses it, and writes it to out.
func (lz *Lz4) Decompress(in io.Reader, out io.Writer) error {
	r := lz4.NewReader(in)
	_, err := io.Copy(out, r)
	return err
}

// WalkByMime calls walkFn for each visited item in archive.
func (gz *Lz4) WalkByMime(f File, walkFn WalkFunc) error {
	r := lz4.NewReader(bytes.NewReader(*f.Buf))

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

func (lz *Lz4) Close() error { return nil }

// NewLz4 returns a new, default instance ready to be customized and used.
func NewLz4() *Lz4 {
	return &Lz4{
		CompressionLevel: 9, // https://github.com/lz4/lz4/blob/1b819bfd633ae285df2dfe1b0589e1ec064f2873/lib/lz4hc.h#L48
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(Lz4))
	_ = Closeable(new(Lz4))
)
