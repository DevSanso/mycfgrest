package conn

import (
	"os"
	"github.com/BurntSushi/toml"

	"mycfgrest/app_error"
)

func ReadTomlConnCfg(filePath string) (*ConnMeta, *app_error.AppError) {
	connByte, connErr := os.ReadFile(filePath)
	
	if connErr != nil {
		return nil, app_error.NewError(connErr, "read failed, conn toml [file:%s]", filePath)
	}

	meta := new(ConnMeta)
	_, err := toml.Decode(string(connByte), meta)

	if err != nil {
		return nil, app_error.NewError(err, "toml decode is failed")
	}

	return meta, nil
}


