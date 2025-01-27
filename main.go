package main

import (
	"net/http"
	"os"

	"mycfgrest/global"
	"mycfgrest/loader/handle"

	"golang.org/x/exp/slog"
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
		ConnConf: args.ConnCfgPath,
	})

	if globalErr != nil {
		slog.Error(globalErr.Error())
		os.Exit(2)
	}
	
	http.HandleFunc("/", httpRootHandleFunc)
	http.ListenAndServe(args.BindAddrPort, nil)
}
