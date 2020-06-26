package deepfindexe

import (
	fastxz "github.com/xi2/xz"
	"io"
)

// TarXz facilitates xz compression
// (https://tukaani.org/xz/format.html)
// of tarball archives.
type TarXz struct {
	*Tar
}

// Walk calls walkFn for each visited item in archive.
func (txz *TarXz) WalkByMime(f File, walkFn WalkFunc) error {
	txz.wrapReader()
	return txz.Tar.WalkByMime(f, walkFn)
}

// Open opens t for reading a compressed archive from
// in. The size parameter is not used.
func (txz *TarXz) Open(buf *[]byte) error {
	txz.wrapReader()
	return txz.Tar.Open(buf)
}

func (txz *TarXz) wrapReader() {
	var xzr *fastxz.Reader
	txz.Tar.readerWrapFn = func(r io.Reader) (io.Reader, error) {
		var err error
		xzr, err = fastxz.NewReader(r, 0)
		return xzr, err
	}
}

func (txz *TarXz) Close() error { return nil }

// NewTarXz returns a new, default instance ready to be customized and used.
func NewTarXz() *TarXz {
	return &TarXz{
		Tar: NewTar(),
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(TarXz))
	_ = Closeable(new(TarXz))
)
