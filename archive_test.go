package deepfindexe

import (
	"deepfindexe/mimetype"
	"reflect"
	"testing"
)

const TestDir = "testdata/"

func TestArchiveByExt(t *testing.T) {
	for _, tc := range []struct {
		path string
		archive interface{}
	}{
		{
			path:  "test.tar.gz",
			archive:  &TarGz{},
		},
		{
			path:  "test.tgz",
			archive:  &TarGz{},
		},
		{
			path:  "test.tar.bz2",
			archive:  &TarBz2{},
		},
		{
			path:  "test.tbz2",
			archive:  &TarBz2{},
		},
		{
			path:  "test.tar.br",
			archive:  &TarBrotli{},
		},
		{
			path:  "test.tbr",
			archive:  &TarBrotli{},
		},
		{
			path:  "test.tar.lz4",
			archive:  &TarLz4{},
		},
		{
			path:  "test.tlz4",
			archive:  &TarLz4{},
		},
		{
			path:  "test.tar.sz",
			archive:  &TarSz{},
		},
		{
			path:  "test.tsz",
			archive:  &TarSz{},
		},
		{
			path:  "test.tar.xz",
			archive:  &TarXz{},
		},
		{
			path:  "test.txz",
			archive:  &TarXz{},
		},
		{
			path:  "test.tar.zst",
			archive:  &TarZstd{},
		},
	}{
		a, err := ArchiveByExtension(tc.path)
		if err != nil {
			t.Errorf("ArchiveByExtension error: %s", err.Error())
			return
		}

		if reflect.TypeOf(a) != reflect.TypeOf(tc.archive) {
			t.Errorf("Test %s: Expected %T but got %T", tc.path, tc.archive, a)
		}
	}
}

func TestArchiveByExtUnrecognized(t *testing.T) {
	path := "test.exe"
	_, err := ArchiveByExtension(path)
	if err == nil {
		t.Errorf("ArchiveByExtension must error on %s", path)
	}
}

func TestArchiveByMime(t *testing.T) {
	for _, tc := range []struct {
		path string
		archive interface{}
	}{
		{
			path:    "test.bat.rar",
			archive: &Rar{},
		},
		{
			path:    "test.bat.tar",
			archive: &Tar{},
		},
		{
			path:    "test.bat.gz",
			archive: &Gz{},
		},
		{
			path:    "test.bat.bz2",
			archive: &Bz2{},
		},
		{
			path:    "test.bat.lz4",
			archive: &Lz4{},
		},
		{
			path:    "test.bat.xz",
			archive: &Xz{},
		},
		{
			path:    "test.bat.zst",
			archive: &Zstd{},
		},
		{
			path:    "test.bat.sz",
			archive: &Snappy{},
		},
	}{
		mime, err := mimetype.DetectFile(TestDir + tc.path)
		if err != nil {
			t.Errorf("DetectFile error %s: %s", tc.path, err.Error())
			return
		}

		a, err := ArchiveByMime(mime)
		if err != nil {
			t.Errorf("ArchiveByMime error %s mime %s: %s", tc.path, mime.String(), err.Error())
			return
		}

		if reflect.TypeOf(a) != reflect.TypeOf(tc.archive) {
			t.Errorf("Test %s: Expected %T but got %T", tc.path, tc.archive, a)
		}

	}
}

func TestArchiveByMimeUnrecognized(t *testing.T) {
	path := "test.exe"
	mime, err := mimetype.DetectFile(TestDir + path)
	if err != nil {
		t.Errorf("DetectFile error %s: %s", path, err.Error())
		return
	}

	_, err = ArchiveByMime(mime)
	if err == nil {
		t.Errorf("ArchiveByExtension must error on %s", path)
	}
}

func TestExtEqualsAny(t *testing.T) {

	batName := "test.bat"
	res := ExtEqualsAny(batName, execExtensions)
	if !res {
		t.Errorf("ExtEqualsAny not detect %s", batName)
	}

	docName := "test.doc"
	res = ExtEqualsAny(docName, execExtensions)
	if res {
		t.Errorf("ExtEqualsAny detect %s", docName)
	}

}

var archiveFormats = []interface{}{
	NewZip(),
	NewTar(),
	NewRar(),
	NewBrotli(),
	NewTarBrotli(),
	NewBz2(),
	NewTarBz2(),
	NewGz(),
	NewTarGz(),
	NewLz4(),
	NewTarLz4(),
	NewSnappy(),
	NewTarSz(),
	NewXz(),
	NewTarXz(),
	NewZstd(),
	NewTarZstd(),
}

