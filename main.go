package main

import (
	"os"

	"golang.org/x/exp/slog"
	"mycfgrest/global"
	"mycfgrest/loader/handle"
)

func main() {
	args,argsErr := NewArgs()

	if argsErr != nil {
		println(argsErr.Error())
		os.Exit(2)
	}

	f, fErr := os.OpenFile(args.LogFile, os.O_APPEND | os.O_WRONLY | os.O_CREATE, os.FileMode(0640))
	if fErr != nil {
		println(fErr.Error())
		os.Exit(2)
	}
	defer f.Close()
	slog.SetDefault(NewLogger(f, args.LogLevel))
	
	globalErr := global.Init(&global.GlobalOptions{
		HandleDir: args.HandleDir,
		HandleType: handle.HandleTypeToml,
	})

	if globalErr != nil {
		slog.Error(globalErr.Error())
		os.Exit(2)
	}
	
}
