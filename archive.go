package deepfindexe

import (
	"deepfindexe/mimetype"
	"fmt"
	"path/filepath"
	"strings"
)

// File provides methods for accessing information about
// or contents of a file within an archive.
type File struct {
	FileName string
	// File data
	Buf *[]byte
}

/// Detect file mime with store position in the file
func (f File) Detect() *mimetype.MIME {
	return mimetype.Detect(*f.Buf)
}

type Archive struct{
	WalkerByMime
	Closeable
}

type Closeable interface {
	Close() error
}

// WalkFunc is called at each item visited by Walk.
// If an error is returned, the walk may continue
// if the Walker is configured to continue on error.
// The sole exception is the error value ErrStopWalk,
// which stops the walk without an actual error.
type WalkFunc func(f File) error

// Walker can walk an archive file and return information
// about each item in the archive.
type WalkerByMime interface {
	WalkByMime(f File, walkFn WalkFunc) error
}

// ErrStopWalk signals Walk to break without error.
var ErrStopWalk = fmt.Errorf("walk stopped")

// ErrFormatNotRecognized is an error that will be
// returned if the file is not a valid archive format.
var ErrFormatNotRecognized = fmt.Errorf("format not recognized")

func TrimSuffix(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// ArchiveByMime returns an archiver and unarchiver, or compressor
// and decompressor, based on the mime of the file.
func ArchiveByMime(mime *mimetype.MIME) (interface{}, error) {
	if mime.Is("application/x-rar") {
		return NewRar(), nil
	} else if mime.Is("application/zip") {
		return NewZip(), nil
	} else if mime.Is("application/x-tar") {
		return NewTar(), nil
	} else if mime.Is("application/gzip") {
		return NewGz(), nil
	} else if mime.Is("application/x-bzip2") {
		return NewBz2(), nil
	} else if mime.Is("application/x-lz4") {
		return NewLz4(), nil
	} else if mime.Is("application/x-xz") {
		return NewXz(), nil
	} else if mime.Is("application/zstd") {
		return NewZstd(), nil
	} else if mime.Is("application/x-snappy-framed") {
		return NewSnappy(), nil
	}
	return nil, fmt.Errorf("format unrecognized by mime: %s", mime.String())
}

// ByExtension returns an archiver and unarchiver, or compressor
// and decompressor, based on the extension of the filename.
func ArchiveByExtension(filename string) (interface{}, error) {
	if strings.HasSuffix(filename, ".tar.gz") {
		return NewTarGz(), nil
	} else if
	strings.HasSuffix(filename, ".tgz") {
		return NewTarGz(), nil
	} else if strings.HasSuffix(filename, ".tar.bz2") {
		return NewTarBz2(), nil
	} else if strings.HasSuffix(filename, ".tbz2") {
		return NewTarBz2(), nil
	} else if strings.HasSuffix(filename, ".tar.br") {
		return NewTarBrotli(), nil
	} else if strings.HasSuffix(filename, ".tbr") {
		return NewTarBrotli(), nil
	} else if strings.HasSuffix(filename, ".br") {
		return NewBrotli(), nil
	} else if strings.HasSuffix(filename, ".tar.lz4") {
		return NewTarLz4(), nil
	} else if strings.HasSuffix(filename, ".tlz4") {
		return NewTarLz4(), nil
	} else if strings.HasSuffix(filename, ".tar.sz") {
		return NewTarSz(), nil
	} else if strings.HasSuffix(filename, ".tsz") {
		return NewTarSz(), nil
	} else if strings.HasSuffix(filename, ".tar.xz") {
		return NewTarXz(), nil
	} else if strings.HasSuffix(filename, ".txz") {
		return NewTarXz(), nil
	} else if strings.HasSuffix(filename, ".tar.zst") {
		return NewTarZstd(), nil
	}
	return nil, fmt.Errorf("format unrecognized by filename: %s", filename)
}

func ExtEqualsAny(filename string, extensions []string) bool {
	ext := filepath.Ext(filename)
	// check if ext exists
	if len(ext) == 0 {
		return false
	}
	ext1 := ext[1:]
	for _, e := range extensions {
		if e == ext1 {
			return true
		}
	}
	return false
}
