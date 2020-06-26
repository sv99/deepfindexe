package deepfindexe

import (
	"path/filepath"
	"testing"
)

var testSuffix = []string{
	"",
	".br",
	".bz2",
	".gz",
	".lz4",
	".rar",
	".sz",
	".tar",
	".tar.br",
	".tar.bz2",
	".tar.gz",
	".tar.lz4",
	".tar.rar",
	".tar.sz",
	".tar.xz",
	".tar.zip",
	".tar.zst",
	".xz",
	".zip",
	".zst",
}

func TestFindExe(t *testing.T) {
	for _, tc := range testSuffix {
		testFindItem(t, "test.exe", tc)
		//fn := TestDir + "test.exe" + tc
		//res, err := Find(fn, false)
		//if err != nil {
		//	t.Errorf("Test Find error on %s %s", fn, err.Error())
		//}
		//if filepath.Base(res) != "test.exe" {
		//	t.Errorf("Find not detect executable in %s: %s", fn, res)
		//}
	}
}

func TestFindBat(t *testing.T) {
	for _, tc := range testSuffix {
		testFindItem(t, "test.bat", tc)
	}
}

func testFindItem(t *testing.T, item string, suffix string) {
	fn := TestDir + item + suffix
	res, err := Find(fn, false)
	if err != nil {
		t.Errorf("Test Find error on %s %s", fn, err.Error())
	}
	if filepath.Base(res) != item {
		t.Errorf("Find not detect executable in %s: %s", fn, res)
	}
}

func TestFindRecursiveArchive(t *testing.T) {
	for _, tc := range []string{
		".zip.rar",
		".rar.zip",
	} {
		testFindItem(t, "test.exe", tc)
	}
}

