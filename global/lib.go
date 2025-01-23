package global

import (
	"database/sql"
	"mycfgrest/loader/handle"
	"mycfgrest/app_error"
)

var global struct {
	handles []handle.HandleMeta

	sqlPool map[string]*sql.DB
}

type GlobalOptions struct {
	HandleDir string
	HandleType handle.HandleMetaType
}

func Init(opt *GlobalOptions) *app_error.AppError {
	utils, utilsErr := handle.NewLoaderUtils(opt.HandleDir, opt.HandleType)
	if utilsErr != nil {
		return utilsErr
	}

	if utils.Size() <= 0 {
		return app_error.NewError(app_error.ErrorNoData,"")
	}

	return nil
}