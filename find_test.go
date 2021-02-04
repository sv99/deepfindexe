package deepfindexe

import (
	"fmt"
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
	}
}

func TestFindBat(t *testing.T) {
	for _, tc := range testSuffix {
		testFindItem(t, "test.bat", tc)
	}
}

func testFindItem(t *testing.T, item string, suffix string) {
	fn := TestDir + item + suffix
	opts := DefOptions(fn)
	res, err := Find(opts)
	if err != nil {
		t.Errorf("Test Find error on %s %s", fn, err.Error())
	}
	fmt.Println(fn, res, opts)
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

func TestFindPasswordedArchive(t *testing.T) {
	for _, tc := range []string{
		"rar_with_pass.rar",
		"zip_with_pass.zip",
	} {
		fn := TestDir + tc
		opts := DefOptions(fn)
		res, err := Find(opts)
		if err != nil {
			t.Errorf("Test Find error on %s %s", fn, err.Error())
		}
		fmt.Println(fn, res, opts)
		//if filepath.Base(res) != item {
		//	t.Errorf("Find not detect executable in %s: %s", fn, res)
		//}
	}
}
