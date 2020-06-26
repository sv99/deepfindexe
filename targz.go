package deepfindexe

import (
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/pgzip"
	"io"
)

// TarGz facilitates gzip compression
// (RFC 1952) of tarball archives.
type TarGz struct {
	*Tar

	// The compression level to use, as described
	// in the compress/gzip package.
	CompressionLevel int

	// Disables parallel gzip.
	SingleThreaded bool
}

// Walk calls walkFn for each visited item in archive.
func (tgz *TarGz) WalkByMime(f File, walkFn WalkFunc) error {
	tgz.wrapReader()
	return tgz.Tar.WalkByMime(f, walkFn)
}

// Open opens t for reading a compressed archive from
// in. The size parameter is not used.
func (tgz *TarGz) Open(buf *[]byte, size int64) error {
	tgz.wrapReader()
	return tgz.Tar.Open(buf)
}

func (tgz *TarGz) Close() error {
	return nil
}

func (tgz *TarGz) wrapReader() {
	var gzr io.ReadCloser
	tgz.Tar.readerWrapFn = func(r io.Reader) (io.Reader, error) {
		var err error
		if tgz.SingleThreaded {
			gzr, err = gzip.NewReader(r)
		} else {
			gzr, err = pgzip.NewReader(r)
		}
		return gzr, err
	}
	tgz.Tar.cleanupWrapFn = func() {
		gzr.Close()
	}
}

// NewTarGz returns a new, default instance ready to be customized and used.
func NewTarGz() *TarGz {
	return &TarGz{
		CompressionLevel: gzip.DefaultCompression,
		Tar:              NewTar(),
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(TarGz))
	_ = Closeable(new(TarGz))
)

