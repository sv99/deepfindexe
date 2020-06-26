package deepfindexe

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nwaples/rardecode"
)

// Rar provides facilities for reading RAR archives.
// See https://www.rarlab.com/technote.htm.
type Rar struct {
	// If true, errors encountered during reading
	// or writing a single file will be logged and
	// the operation will continue on remaining files.
	ContinueOnError bool

	// The password to open archives (optional).
	Password string

	rr *rardecode.Reader     // underlying stream reader
	rc *rardecode.ReadCloser // supports multi-volume archives (files only)
}

// Open opens t for reading an archive from
// in. The size parameter is not used.
func (r *Rar) Open(buf *[]byte) error {
	if r.rr != nil {
		return fmt.Errorf("rar archive is already open for reading")
	}
	var err error
	r.rr, err = rardecode.NewReader(bytes.NewReader(*buf), r.Password)
	return err
}

// Read reads the next file from t, which must have
// already been opened for reading. If there are no
// more files, the error is io.EOF. The File must
// be closed when finished reading from it.
func (r *Rar) Read() (File, error) {
	if r.rr == nil {
		return File{}, fmt.Errorf("rar archive is not open")
	}

	hdr, err := r.rr.Next()
	if err != nil {
		return File{}, err // don't wrap error; preserve io.EOF
	}

	buf, err := ioutil.ReadAll(r.rr)
	if err != nil {
		return File{}, err // don't wrap error; preserve io.EOF
	}

	file := File{
		FileName:   hdr.Name,
		Buf: 		&buf,
	}

	return file, nil
}

// Close closes the rar archive(s) opened by Create and Open.
func (r *Rar) Close() error {
	var err error
	if r.rc != nil {
		rc := r.rc
		r.rc = nil
		err = rc.Close()
	}
	if r.rr != nil {
		r.rr = nil
	}
	return err
}

// WalkByMime calls walkFn for each visited item in archive.
func (r *Rar) WalkByMime(f File, walkFn WalkFunc) error {
	err := r.Open(f.Buf)
	if err != nil {
		return fmt.Errorf("opening archive file: %v", err)
	}

	for {
		f, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if r.ContinueOnError {
				log.Printf("[ERROR] Opening next file: %v", err)
				continue
			}
			return fmt.Errorf("opening next file: %v", err)
		}
		err = walkFn(f)
		if err != nil {
			if err == ErrStopWalk {
				break
			}
			if r.ContinueOnError {
				log.Printf("[ERROR] Walking %s: %v", f.FileName, err)
				continue
			}
			return err
		}
	}

	return nil
}

// NewRar returns a new, default instance ready to be customized and used.
func NewRar() *Rar {
	return &Rar{
	}
}

type rarFileInfo struct {
	fh *rardecode.FileHeader
}

func (rfi rarFileInfo) Name() string       { return rfi.fh.Name }
func (rfi rarFileInfo) Size() int64        { return rfi.fh.UnPackedSize }
func (rfi rarFileInfo) Mode() os.FileMode  { return rfi.fh.Mode() }
func (rfi rarFileInfo) ModTime() time.Time { return rfi.fh.ModificationTime }
func (rfi rarFileInfo) IsDir() bool        { return rfi.fh.IsDir }
func (rfi rarFileInfo) Sys() interface{}   { return nil }

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(Rar))
	_ = Closeable(new(Rar))
)
