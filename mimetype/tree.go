package mimetype

import "deepfindexe/mimetype/matchers"

// root is a matcher which passes for any slice of bytes.
// When a matcher passes the check, the children matchers
// are tried in order to find a more accurate MIME type.
var root = newMIME("application/octet-stream", "", func([]byte) bool { return true },
	sevenZ, zip,
	exe, elf, ar, tar, xar, bz2, fits, gzip, class, swf, crx, wasm,
	dbf, dcm, rar, nes, macho,
	mrc, zstd, cab,
	rpm, xz, lzip, lz4,
	warc, snappy,
)

// The list of nodes appended to the root node.
var (
	xz   = newMIME("application/x-xz", ".xz", matchers.Xz)
	gzip = newMIME("application/gzip", ".gz", matchers.Gzip).
		alias("application/x-gzip", "application/x-gunzip", "application/gzipped", "application/gzip-compressed", "application/x-gzip-compressed", "gzip/document")
	sevenZ = newMIME("application/x-7z-compressed", ".7z", matchers.SevenZ)
	zip    = newMIME("application/zip", ".zip", matchers.Zip).
		alias("application/x-zip", "application/x-zip-compressed")
	tar = newMIME("application/x-tar", ".tar", matchers.Tar)
	xar = newMIME("application/x-xar", ".xar", matchers.Xar)
	bz2 = newMIME("application/x-bzip2", ".bz2", matchers.Bz2)
	fits = newMIME("application/fits", ".fits", matchers.Fits)
	class = newMIME("application/x-java-applet; charset=binary", ".class", matchers.Class)
	swf   = newMIME("application/x-shockwave-flash", ".swf", matchers.Swf)
	crx   = newMIME("application/x-chrome-extension", ".crx", matchers.Crx)
	wasm    = newMIME("application/wasm", ".wasm", matchers.Wasm)
	dbf     = newMIME("application/x-dbf", ".dbf", matchers.Dbf)
	exe     = newMIME("application/vnd.microsoft.portable-executable", ".exe", matchers.Exe)
	elf     = newMIME("application/x-elf", "", matchers.Elf, elfObj, elfExe, elfLib, elfDump)
	elfObj  = newMIME("application/x-object", "", matchers.ElfObj)
	elfExe  = newMIME("application/x-executable", "", matchers.ElfExe)
	elfLib  = newMIME("application/x-sharedlib", ".so", matchers.ElfLib)
	elfDump = newMIME("application/x-coredump", "", matchers.ElfDump)
	ar      = newMIME("application/x-archive", ".a", matchers.Ar, deb).
		alias("application/x-unix-archive")
	deb = newMIME("application/vnd.debian.binary-package", ".deb", matchers.Deb)
	rpm = newMIME("application/x-rpm", ".rpm", matchers.Rpm)
	dcm = newMIME("application/dicom", ".dcm", matchers.Dcm)
	rar = newMIME("application/x-rar-compressed", ".rar", matchers.Rar).
		alias("application/x-rar")
	warc    = newMIME("application/warc", ".warc", matchers.Warc)
	nes     = newMIME("application/vnd.nintendo.snes.rom", ".nes", matchers.Nes)
	macho   = newMIME("application/x-mach-binary", ".macho", matchers.MachO)
	mrc     = newMIME("application/marc", ".mrc", matchers.Marc)
	zstd    = newMIME("application/zstd", ".zst", matchers.Zstd)
	cab     = newMIME("application/vnd.ms-cab-compressed", ".cab", matchers.Cab)
	lzip    = newMIME("application/lzip", ".lz", matchers.Lzip)
	snappy  = newMIME("application/x-snappy-framed", ".sz", matchers.Snappy)
	lz4     = newMIME("application/x-lz4", ".lz4", matchers.Lz4)
)
