package deepfindexe

import (
	"github.com/pierrec/lz4/v4"
	"io"
)

// TarLz4 facilitates lz4 compression
// (https://github.com/lz4/lz4/tree/master/doc)
// of tarball archives.
type TarLz4 struct {
	*Tar

	// The compression level to use when writing.
	// Minimum 0 (fast compression), maximum 12
	// (most space savings).
	CompressionLevel int
}

// Walk calls walkFn for each visited item in archive.
func (tlz4 *TarLz4) WalkByMime(f File, walkFn WalkFunc) error {
	tlz4.wrapReader()
	return tlz4.Tar.WalkByMime(f, walkFn)
}

// Open opens t for reading a compressed archive from
// in. The size parameter is not used.
func (tlz4 *TarLz4) Open(buf *[]byte) error {
	tlz4.wrapReader()
	return tlz4.Tar.Open(buf)
}

func (tlz4 *TarLz4) wrapReader() {
	tlz4.Tar.readerWrapFn = func(r io.Reader) (io.Reader, error) {
		return lz4.NewReader(r), nil
	}
}

func (tlz4 *TarLz4) Close() error { return nil }

// NewTarLz4 returns a new, default instance ready to be customized and used.
func NewTarLz4() *TarLz4 {
	return &TarLz4{
		CompressionLevel: 9, // https://github.com/lz4/lz4/blob/1b819bfd633ae285df2dfe1b0589e1ec064f2873/lib/lz4hc.h#L48
		Tar:              NewTar(),
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(TarLz4))
	_ = Closeable(new(TarLz4))
)
