package deepfindexe

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

// Tar provides facilities for operating TAR archives.
// See http://www.gnu.org/software/tar/manual/html_node/Standard.html.
type Tar struct {
	// Whether to overwrite existing files; if false,
	// an error is returned if the file exists.
	OverwriteExisting bool

	// Whether to make all the directories necessary
	// to create a tar archive in the desired path.
	MkdirAll bool

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

	tw *tar.Writer
	tr *tar.Reader

	readerWrapFn  func(io.Reader) (io.Reader, error)
	writerWrapFn  func(io.Writer) (io.Writer, error)
	cleanupWrapFn func()
}

// Open opens t for reading an archive from
// in. The size parameter is not used.
func (t *Tar) Open(buf *[]byte) error {
	if t.tr != nil {
		return fmt.Errorf("tar archive is already open for reading")
	}
	in := bytes.NewReader(*buf)
	// wrapping readers allows us to open compressed tarballs
	if t.readerWrapFn != nil {
		var err error
		wrappedIn, err := t.readerWrapFn(in)
		if err != nil {
			return fmt.Errorf("wrapping file reader: %v", err)
		}
		t.tr = tar.NewReader(wrappedIn)
	} else {
		t.tr = tar.NewReader(in)
	}
	return nil
}

// Read reads the next file from t, which must have
// already been opened for reading. If there are no
// more files, the error is io.EOF. The File must
// be closed when finished reading from it.
func (t *Tar) Read() (File, error) {
	if t.tr == nil {
		return File{}, fmt.Errorf("tar archive is not open")
	}

	hdr, err := t.tr.Next()
	if err != nil {
		return File{}, err // don't wrap error; preserve io.EOF
	}

	buf, err := ioutil.ReadAll(t.tr)
	if err != nil {
		return File{}, err // don't wrap error; preserve io.EOF
	}

	file := File{
		FileName:   hdr.Name,
		Buf:		&buf,
	}
	return file, nil
}

// Close closes the tar archive(s) opened by Create and Open.
func (t *Tar) Close() error {
	var err error
	if t.tr != nil {
		t.tr = nil
	}
	if t.tw != nil {
		tw := t.tw
		t.tw = nil
		err = tw.Close()
	}
	// make sure cleanup of "Reader/Writer wrapper"
	// (say that ten times fast) happens AFTER the
	// underlying stream is closed
	if t.cleanupWrapFn != nil {
		t.cleanupWrapFn()
	}
	return err
}

// Walk calls walkFn for each visited item in archive.
func (t *Tar) WalkByMime(f File, walkFn WalkFunc) error {
	err := t.Open(f.Buf)
	if err != nil {
		return fmt.Errorf("opening archive: %v", err)
	}

	for {
		f, err := t.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if t.ContinueOnError {
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
			if t.ContinueOnError {
				log.Printf("[ERROR] Walking %s: %v", f.FileName, err)
				continue
			}
			return err
		}
	}

	return nil
}

// NewTar returns a new, default instance ready to be customized and used.
func NewTar() *Tar {
	return &Tar{
		MkdirAll: true,
	}
}

const tarBlockSize = 512

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ = WalkerByMime(new(Tar))
	_ = Closeable(new(Tar))
)
