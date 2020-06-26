package mimetype

import "testing"

const TestDir = "../testdata/"

func TestArchiveByExt(t *testing.T) {
	for _, tc := range []struct {
		path, mime string
	}{
		{
			path:  "test.bat",
			mime:  "application/octet-stream",
		},
		{
			path:  "test.exe",
			mime:  "application/vnd.microsoft.portable-executable",
		},
		{
			path:  "test.bat.br",
			mime:  "application/octet-stream",
		},
		{
			path:  "test.bat.bz2",
			mime:  "application/x-bzip2",
		},
		{
			path:  "test.bat.gz",
			mime:  "application/gzip",
		},
		{
			path:  "test.bat.lz4",
			mime:  "application/x-lz4",
		},
		{
			path:  "test.bat.rar",
			mime:  "application/x-rar-compressed",
		},
		{
			path:  "test.bat.sz",
			mime:  "application/x-snappy-framed",
		},
		{
			path:  "test.bat.tar",
			mime:  "application/x-tar",
		},
	}{
		mime, err := DetectFile(TestDir + tc.path)
		if err != nil {
			t.Errorf("DetectFile error: %s", err.Error())
			return
		}
		detectedMime := mime.String()
		if mime.String() != tc.mime {
			t.Errorf("Test %s: Expected %s but got %s", tc.path, tc.mime, detectedMime)
		}
	}
}


