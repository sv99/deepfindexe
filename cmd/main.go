package main

import (
	"deepfindexe"
	"fmt"
	"github.com/sv99/go-flags"
	"os"
	"path/filepath"
	"strings"
)

func main()  {
	var opts deepfindexe.Options
	var parser = flags.NewParser(&opts, flags.Default)

	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(2)
		} else {
			if opts.Verbose {
				fmt.Println(err.Error())
			}
			os.Exit(2)
		}
	}

	// if not provided optional filename - make it from path
	if opts.Positional.Filename == "" {
		opts.Positional.Filename = filepath.Base(opts.Positional.Filepath)
	}
	// prepare ExtensionArray
	opts.ExtensionArray = strings.Split(opts.Extensions, "|")

	res, err := deepfindexe.Find(&opts)
	if err != nil {
		if opts.Verbose {
			fmt.Println(err.Error())
		}
		os.Exit(2)
	}
	if res != "" {
		fmt.Println(res)
		os.Exit(1)
	}
}

