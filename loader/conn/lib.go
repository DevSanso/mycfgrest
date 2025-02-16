package conn

import (
	"os"

	"github.com/BurntSushi/toml"

	"mycfgrest/types"
)

func ReadTomlConnCfg(filePath string) (*ConnMeta, *types.AppError) {
	connByte, connErr := os.ReadFile(filePath)

	if connErr != nil {
		return nil, types.NewAppError(connErr, "read failed, conn toml [file:%s]", filePath)
	}

	meta := new(ConnMeta)
	_, err := toml.Decode(string(connByte), meta)

	if err != nil {
		return nil, types.NewAppError(err, "toml decode is failed")
	}

	return meta, nil
}
