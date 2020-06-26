package main

import (
	"deepfindexe"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

func main()  {
	var opts deepfindexe.Options
	var parser = flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS] file"

	args, err := parser.Parse()
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

	if len(args) == 0 {
		parser.WriteHelp(os.Stderr)
		os.Exit(2)
	}
	res, err := deepfindexe.Find(args[0], opts.Verbose)
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

