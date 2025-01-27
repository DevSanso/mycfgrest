package handle

import (
	"github.com/BurntSushi/toml"
)

func readTomlHandleCfg(data []byte) (*HandleMeta, error) {
	meta := new(HandleMeta)
	_, err := toml.Decode(string(data), meta)

	if err != nil {
		return nil, err
	}

	return meta, nil
}
