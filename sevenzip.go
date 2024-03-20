package deepfindexe

import (
	"bytes"
	"compress/flate"
	"fmt"
	"github.com/bodgit/sevenzip"
	"io"
	"io/ioutil"
	"log"
)

// SevenZip provides facilities for operating 7Z archives.
type SevenZip struct {
	// The compression level to use, as described
	// in the compress/flate package.
	CompressionLevel int

	// Whether to overwrite existing files; if false,
	// an error is returned if the file exists.
	OverwriteExisting bool

	// Whether to make all the directories necessary
	// to create a 7z archive in the desired path.
	MkdirAll bool

	// If enabled, selective compression will only
	// compress files which are not already in a
	// compressed format; this is decided based
	// simply on file extension.
	SelectiveCompression bool

	// A single top-level folder can be implicitly
	// created by the Archive or Unarchive methods
	// if the files to be added to the archive
	// or the files to be extracted from the archive
	// do not all have a common root. This roughly
	// mimics the behavior of archival tools integrated
	// into OS file browsers which create a subfolder
	// to avoid unexpectedly littering the destination
	// folder with potentially many files, causing a
	// problematic cleanup/organization situation.
	// This feature is available for both creation
	// and extraction of archives, but may be slightly
	// inefficient with lots and lots of files,
	// especially on extraction.
	ImplicitTopLevelFolder bool

	// If true, errors encountered during reading
	// or writing a single file will be logged and
	// the operation will continue on remaining files.
	ContinueOnError bool

	zr   *sevenzip.Reader
	ridx int
}

// Open opens z for reading an archive from in,
// which is expected to have the given size and
// which must be an io.ReaderAt.
func (z *SevenZip) Open(buf *[]byte) error {
	if z.zr != nil {
		return fmt.Errorf("zip archive is already open for reading")
	}
	reader := bytes.NewReader(*buf)
	var err error
	z.zr, err = sevenzip.NewReader(reader, int64(len(*buf)))
	if err != nil {
		return fmt.Errorf("creating reader: %v", err)
	}
	z.ridx = 0
	return nil
}

// Read reads the next file from z, which must have
// already been opened for reading. If there are no
// more files, the error is io.EOF. The File must
// be closed when finished reading from it.
func (z *SevenZip) Read() (File, error) {
	if z.zr == nil {
		return File{}, fmt.Errorf("zip archive is not open")
	}
	if z.ridx >= len(z.zr.File) {
		return File{}, io.EOF
	}

	// access the file and increment counter so that
	// if there is an error processing this file, the
	// caller can still iterate to the next file
	zf := z.zr.File[z.ridx]
	z.ridx++

	rc, err := zf.Open()
	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return File{}, err // don't wrap error; preserve io.EOF
	}

	file := File{
		FileName: zf.Name,
		Buf:      &buf,
	}
	return file, nil
}

// Close closes the zip archive(s) opened by Create and Open.
func (z *SevenZip) Close() error {
	if z.zr != nil {
		z.zr = nil
	}
	return nil
}

// Walk calls walkFn for each visited item in archive.
func (z *SevenZip) WalkByMime(f File, walkFn WalkFunc) error {
	err := z.Open(f.Buf)
	if err != nil {
		return fmt.Errorf("opening 7z reader: %v", err)
	}

	for _, zf := range z.zr.File {
		zfrc, err := zf.Open()
		if err != nil {
			if z.ContinueOnError {
				log.Printf("[ERROR] Opening %s: %v", zf.Name, err)
				continue
			}
			return fmt.Errorf("opening %s: %v", zf.Name, err)
		}

		buf, err := ioutil.ReadAll(zfrc)
		if err != nil {
			// ignore error for passworded zip
			if err == flate.CorruptInputError(4) {
				log.Printf("[ERROR] read passworded zip file: %v", err)
				continue
			}
			return err // don't wrap error; preserve io.EOF
		}
		_ = zfrc.Close()

		err = walkFn(File{
			FileName: zf.Name,
			Buf:      &buf,
		})
		if err != nil {
			if err == ErrStopWalk {
				break
			}
			if z.ContinueOnError {
				log.Printf("[ERROR] Walking %s: %v", zf.Name, err)
				continue
			}
			return err
		}
	}

	return nil
}

// NewZip returns a new, default instance ready to be customized and used.
func NewSevenZip() *SevenZip {
	return &SevenZip{
		CompressionLevel:     flate.DefaultCompression,
		MkdirAll:             true,
		SelectiveCompression: true,
	}
}

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(SevenZip))
	_ = Closeable(new(SevenZip))
)
