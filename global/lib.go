package global

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"mycfgrest/app_error"
	"mycfgrest/loader/conn"
	"mycfgrest/loader/handle"
)

var global struct {
	handles []*handle.HandleMeta

	sqlPool map[string]*sql.DB
}

type GlobalOptions struct {
	HandleDir string
	ConnConf string
	HandleType handle.HandleMetaType
}

func initSqlPool(meta *conn.ConnMeta) *app_error.AppError {
	for k,v := range meta.Sql.Postgres {
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", v.User, v.Password, v.Addr, v.Port, v.Dbname)
		db, dbErr := sql.Open("postgres", connStr)
		
		if dbErr != nil {
			app_error.NewError(dbErr, "init postgresql pool failed [ip:%s]", v.Addr)
		}

		global.sqlPool[fmt.Sprintf("sql.postgres.%s", k)] = db
	}

	for k,v := range meta.Sql.Sqlite {
		db, dbErr := sql.Open("sqlite", v.Dbname)
		if dbErr != nil {
			app_error.NewError(dbErr, "init sqlite pool failed [ip:%s]", v.Addr)
		}

		global.sqlPool[fmt.Sprintf("sql.sqlite.%s", k)] = db
	}

	return nil
}

func Init(opt *GlobalOptions) *app_error.AppError {
	connMeta, connErr := conn.ReadTomlConnCfg(opt.ConnConf)

	if connErr != nil {
		return app_error.NewError(connErr, "read failed connection info file")
	}

	if pErr := initSqlPool(connMeta); pErr != nil {
		return app_error.NewError(pErr, "init pools failed")
	}

	utils, utilsErr := handle.NewLoaderUtils(opt.HandleDir, opt.HandleType)
	if utilsErr != nil {
		return utilsErr
	}

	if utils.Size() <= 0 {
		return app_error.NewError(app_error.ErrorNoData,"")
	}

	for {
		m, mErr := utils.Next()
		if mErr != nil {
			return app_error.NewError(mErr, "utils Next method error")
		}

		if m == nil {
			break
		}

		global.handles = append(global.handles, m)
	}

	return nil
}