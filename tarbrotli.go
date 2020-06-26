package deepfindexe

import (
	"github.com/andybalholm/brotli"
	"io"
)

// TarBrotli facilitates brotli compression of tarball archives.
type TarBrotli struct {
	*Tar
	Quality int
}

// Walk calls walkFn for each visited item in archive.
func (tbr *TarBrotli) WalkByMime(f File, walkFn WalkFunc) error {
	tbr.wrapReader()
	return tbr.Tar.WalkByMime(f, walkFn)
}

// Open opens t for reading a compressed archive from
// in. The size parameter is not used.
func (tbr *TarBrotli) Open(buf *[]byte) error {
	tbr.wrapReader()
	return tbr.Tar.Open(buf)
}

func (tbr *TarBrotli) wrapReader() {
	tbr.Tar.readerWrapFn = func(r io.Reader) (io.Reader, error) {
		return brotli.NewReader(r), nil
	}
}

// NewTarBrotli returns a new, default instance ready to be customized and used.
func NewTarBrotli() *TarBrotli {
	return &TarBrotli{
		Tar:     NewTar(),
		Quality: brotli.DefaultCompression,
	}
}

func (tbz2 *TarBrotli) Close() error { return nil }

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(TarBrotli))
	_ = Closeable(new(TarBrotli))
)
