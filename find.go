package deepfindexe

import (
	"deepfindexe/mimetype"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var archiveMimes = []string{
	"application/x-rar-compressed",
	"application/x-rar",
	"application/zip",
	"application/x-zip",
	"application/x-zip-compressed",
	"application/x-tar",
	"application/gzip",
	"application/x-gzip",
	"application/x-bzip2",
	"application/x-lz4",
	"application/x-xz",
	"application/zstd",
	"application/x-snappy-framed",
	"application/vnd.ms-cab-compressed",
	"application/x-7z-compressed",
}

var execMimes = []string{
	"application/vnd.microsoft.portable-executable",
	"application/x-executable",
	"application/x-elf",
}

var execExtensions = []string{
	"ade", "adp", "asd",
	"bas", "bat",
	"cab", "chm", "cmd", "com", "cpl", "crt",
	"dll", "exe", "hlp", "hta", "inf", "ins", "isp",
	"jse", "jar", "lib", "lnk",
	"mdb", "mde", "mdz", "msc", "msi", "msp", "mst",
	"ole", "ocx",
	"pcd", "pif",
	"reg",
	"scr", "sct", "shs", "shb", "sys",
	"url",
	"vbe", "vbs", "vxd",
	"wsc", "wsf", "wsh",
}

// Find first executable in the archive.
// Return error with file name executable in the message.
func Find(opts *Options) (string, error) {
	// check file extension
	if opts.Positional.Filename != "" {
		// detected executable attachments
		det := ExtEqualsAny(opts.Positional.Filename, opts.ExtensionArray)
		if det {
			return opts.Positional.Filename, nil
		}
	}

	file, err := os.Open(opts.Positional.Filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	res, err := Walk(File{
		FileName: opts.Positional.Filename,
		Buf: &buf,
	}, 0, opts)
	return res, err
}

// pretty indent for internal archives
func indent(level int) string {
	return strings.Repeat("  ", level)
}

// check file and detect executable
func Walk(f File, level int, opts *Options) (string, error) {
	// first step detect executable by ext
	// base check - file must have extension executable in the Explorer!!
	det := ExtEqualsAny(f.FileName, opts.ExtensionArray)
	if det {
		if opts.Verbose {
			fmt.Println(indent(level) + "Executable:", f.FileName)
		}
		return f.FileName, nil
	}
	// second step detect executable by mime
	mime := mimetype.Detect(*f.Buf)
	if mimetype.EqualsAny(mime.String(), execMimes...) {
		if opts.Verbose {
			fmt.Println(indent(level) + "Executable:", f.FileName, mime.String())
		}
		return f.FileName, nil
	}

	a, err := ArchiveByExtension(f.FileName)
	if err != nil {
		if mimetype.EqualsAny(mime.String(), archiveMimes...) {
			a, err = ArchiveByMime(mime)
			if err != nil {
				if opts.Verbose {
					fmt.Println(indent(level) + "File:", f.FileName, mime.String())
				}
				return "", err
			}
		} else {
			if opts.Verbose {
				fmt.Println(indent(level) + "File:", f.FileName, mime.String())
			}
			return "", nil
		}
	}
	defer a.(Closeable).Close()

	if opts.Verbose {
		fmt.Println(indent(level) + "Archive:", f.FileName, mime.String())
	}
	w, ok := a.(WalkerByMime)
	if !ok {
		return "", fmt.Errorf(
			"format specified by archive filename is not a walker format: (%T)", a)
	}
	var res string
	err = w.WalkByMime(f, func(f File) error {
		mime = mimetype.Detect(*f.Buf)
		if mimetype.EqualsAny(mime.String(), archiveMimes...) {
			// find internal archive
			res, err = Walk(f, level + 1, opts)
			if err != nil {
				return err
			}
		} else {
			// archive item detect executable
			if ExtEqualsAny(f.FileName, opts.ExtensionArray) || mimetype.EqualsAny(
				mime.String(), execMimes...) {
				res = f.FileName
				if opts.Verbose {
					fmt.Println(indent(level + 1) + "Executable", f.FileName, mime.String())
				}
				// stop walking
				return nil
			} else {
				if opts.Verbose {
					fmt.Println(indent(level + 1) + "File:", f.FileName, mime.String())
				}
			}
		}
		return nil
	})
	return res, err
}
