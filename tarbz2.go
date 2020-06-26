package deepfindexe

import (
	"github.com/dsnet/compress/bzip2"
	"io"
)

// TarBz2 facilitates bzip2 compression
// (https://github.com/dsnet/compress/blob/master/doc/bzip2-format.pdf)
// of tarball archives.
type TarBz2 struct {
	*Tar

	CompressionLevel int
}

// Walk calls walkFn for each visited item in archive.
func (tbz2 *TarBz2) WalkByMime(f File, walkFn WalkFunc) error {
	tbz2.wrapReader()
	return tbz2.Tar.WalkByMime(f, walkFn)
}

// Open opens t for reading a compressed archive from
// in. The size parameter is not used.
func (tbz2 *TarBz2) Open(buf *[]byte) error {
	tbz2.wrapReader()
	return tbz2.Tar.Open(buf)
}

func (tbz2 *TarBz2) wrapReader() {
	var bz2r *bzip2.Reader
	tbz2.Tar.readerWrapFn = func(r io.Reader) (io.Reader, error) {
		var err error
		bz2r, err = bzip2.NewReader(r, nil)
		return bz2r, err
	}
	tbz2.Tar.cleanupWrapFn = func() {
		bz2r.Close()
	}
}

func (tbz2 *TarBz2) Close() error { return nil }

// NewTarBz2 returns a new, default instance ready to be customized and used.
func NewTarBz2() *TarBz2 {
	return &TarBz2{
		CompressionLevel: bzip2.DefaultCompression,
		Tar:              NewTar(),
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = Closeable(new(TarBz2))
	_ = WalkerByMime(new(TarBz2))
)
