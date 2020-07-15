package deepfindexe

import (
	fp "path/filepath"
	"strings"
)


var executables = "ade|adp|asd|bas|bat|cab|chm|cmd|com|cpl|crt|" +
		"dll|exe|hlp|hta|inf|ins|isp|jse|jar|lib|lnk|" +
		"mdb|mde|mdz|msc|msi|msp|mst|ole|ocx|pcd|pif|reg|" +
		"scr|sct|shs|shb|sys|url|vbe|vbs|vxd|wsc|wsf|wsh"

type Options struct {
	Verbose		bool	`short:"v" long:"verbose" description:"Show verbose debug information"`
	Extensions	string 	`short:"e" long:"ext" description:"Executable extensions" default:"ade|adp|asd|bas|bat|cab|chm|cmd|com|cpl|crt|dll|exe|hlp|hta|inf|ins|isp|jse|jar|lib|lnk|mdb|mde|mdz|msc|msi|msp|mst|ole|ocx|pcd|pif|reg|scr|sct|shs|shb|sys|url|vbe|vbs|vxd|wsc|wsf|wsh"`
	// stripped by | Extensions param
	ExtensionArray []string
	Positional struct {
		Filepath string `positional-arg-name:"filepath" required:"yes"`
		Filename string `positional-arg-name:"filename"`
	} `positional-args:"yes"`
}

func DefOptions(filepath string) *Options {
	opts := Options{
		Verbose:    false,
		Extensions: executables,
	}
	opts.Positional.Filepath = filepath
	opts.Positional.Filename = fp.Base(filepath)
	// prepare ExtensionArray
	opts.ExtensionArray = strings.Split(opts.Extensions, "|")
	return &opts
}
