package handle

import (
	"os"
	"path"

	"golang.org/x/exp/maps"

	"mycfgrest/types"
)

type HandleMetaType int

const (
	HandleTypeToml = iota
)

type HandleLoaderUtil struct {
	fileList   []string
	currentIdx int
	loaderType HandleMetaType
}

func NewLoaderUtils(dir string, loaderType HandleMetaType) (*HandleLoaderUtil, *types.AppError) {
	sub_file, err := os.ReadDir(dir)
	if err != nil {
		return nil, types.NewAppError(err, "loader utils read failed [dir:%s]", dir)
	}

	lu := new(HandleLoaderUtil)
	for _, f := range sub_file {
		info, infoErr := f.Info()
		if infoErr != nil {
			return nil, types.NewAppError(infoErr, "loader utils read file info failed [dir:%s]", f.Name())
		}

		if info.IsDir() {
			continue
		}

		lu.fileList = append(lu.fileList, path.Join(dir, info.Name()))
	}
	lu.currentIdx = 0
	lu.loaderType = loaderType
	return lu, nil
}

func (lu *HandleLoaderUtil) Size() int {
	return len(lu.fileList)
}

func (lu *HandleLoaderUtil) Seek(pos int) {
	lu.currentIdx = pos
}

func (lu *HandleLoaderUtil) Cur() int {
	return lu.currentIdx
}

func (*HandleLoaderUtil) checkSymbolDuplicate(m *HandleMeta) (bool, string, string, string) {
	for qk, qval := range m.Data.Request.QueryString {
		for bk, bval := range m.Data.Request.Body {
			if qval.Symbol == bval.Symbol {
				return true, qval.Symbol, qk, bk
			}
		}
	}

	for lk, lval := range m.Data.Load {
		tup := make(map[string]struct{})
		for gl, gval := range lval.GetData {
			_, exists := tup[gval.Symbol]

			if exists {
				return true, lk, gl, gval.Symbol
			}

			tup[gval.Symbol] = struct{}{}
		}
		maps.Clear(tup)
	}
	return false, "", "", ""
}

func (lu *HandleLoaderUtil) Next() (*HandleMeta, *types.AppError) {
	if lu.currentIdx >= len(lu.fileList) {
		return nil, nil
	}

	p := lu.fileList[lu.currentIdx]

	data, readErr := os.ReadFile(p)
	if readErr != nil {
		return nil, types.NewAppError(readErr, "read meta data failed [file:%s] [type:%d]", p, lu.loaderType)
	}
	var meta *HandleMeta
	var metaErr error

	if lu.loaderType == HandleTypeToml {
		meta, metaErr = readTomlHandleCfg(data)
	}

	if metaErr != nil {
		return nil, types.NewAppError(metaErr, "parsing meta data failed [file:%s] [type:%d]", p, lu.loaderType)
	}

	if is_dup, s1, s2, s3 := lu.checkSymbolDuplicate(meta); is_dup {
		return nil, types.NewAppError(types.ErrorAppDuplicate,
			"duplicate [s1:%s, s2:%s, s3:%s]", s1, s2, s3)

	}

	return meta, nil
}
