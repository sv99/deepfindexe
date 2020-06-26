package deepfindexe

import (
	"github.com/klauspost/compress/zstd"
	"io"
)

// TarZstd facilitates Zstandard compression
// (RFC 8478) of tarball archives.
type TarZstd struct {
	*Tar
}

// Walk calls walkFn for each visited item in archive.
func (tzst *TarZstd) WalkByMime(f File, walkFn WalkFunc) error {
	tzst.wrapReader()
	return tzst.Tar.WalkByMime(f, walkFn)
}

// Open opens t for reading a compressed archive from
// in. The size parameter is not used.
func (tzst *TarZstd) Open(buf *[]byte, size int64) error {
	tzst.wrapReader()
	return tzst.Tar.Open(buf)
}

func (tzst *TarZstd) wrapReader() {
	var zstdr *zstd.Decoder
	tzst.Tar.readerWrapFn = func(r io.Reader) (io.Reader, error) {
		var err error
		zstdr, err = zstd.NewReader(r)
		return zstdr, err
	}
	tzst.Tar.cleanupWrapFn = func() {
		zstdr.Close()
	}
}

func (tzst *TarZstd) Close() error { return nil }

// NewTarZstd returns a new, default instance ready to be customized and used.
func NewTarZstd() *TarZstd {
	return &TarZstd{
		Tar: NewTar(),
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(TarZstd))
	_ = Closeable(new(TarZstd))
)
