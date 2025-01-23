package main

import (
	"flag"
	"fmt"
)

type Args struct {
	LogFile string
	LogLevel string
	HandleDir string

	ConnCfgPath string
}

func NewArgs() (*Args, error) {
	args := new(Args)

	flag.StringVar(&args.LogFile, "log-file", "", "logfile path")
	flag.StringVar(&args.LogLevel, "log-level", "", "log level")
	flag.StringVar(&args.HandleDir, "handle_dir", "", "handle dir path")
	flag.StringVar(&args.ConnCfgPath, "conn-cfg-path", "", "connection cfg file path")

	flag.Parse()

	const errFmt = "process param not setting %s"
	if args.ConnCfgPath == "" {
		return nil, fmt.Errorf(errFmt, "conn-cfg-path")
	}

	if args.LogFile == "" {
		return nil, fmt.Errorf(errFmt, "log-file")
	}

	if args.LogLevel == "" {
		return nil, fmt.Errorf(errFmt, "log-level")
	}

	if args.HandleDir == "" {
		return nil, fmt.Errorf(errFmt, "handle_dir")
	}

	return args, nil
}