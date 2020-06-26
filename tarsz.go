package deepfindexe

import (
	"github.com/golang/snappy"
	"io"
)

// TarSz facilitates Snappy compression
// (https://github.com/google/snappy)
// of tarball archives.
type TarSz struct {
	*Tar
}

// Walk calls walkFn for each visited item in archive.
func (tsz *TarSz) WalkByMime(f File, walkFn WalkFunc) error {
	tsz.wrapReader()
	return tsz.Tar.WalkByMime(f, walkFn)
}

// Open opens t for reading a compressed archive from
// in. The size parameter is not used.
func (tsz *TarSz) Open(buf *[]byte) error {
	tsz.wrapReader()
	return tsz.Tar.Open(buf)
}

func (tsz *TarSz) wrapReader() {
	tsz.Tar.readerWrapFn = func(r io.Reader) (io.Reader, error) {
		return snappy.NewReader(r), nil
	}
}

func (tsz *TarSz) Close() error { return nil }

// NewTarSz returns a new, default instance ready to be customized and used.
func NewTarSz() *TarSz {
	return &TarSz{
		Tar: NewTar(),
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(TarSz))
	_ = Closeable(new(TarSz))
)
