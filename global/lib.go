package global

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"mycfgrest/loader/conn"
	"mycfgrest/loader/handle"
	"mycfgrest/types"
)

var global struct {
	handles []*handle.HandleMeta

	sqlPool map[string]*sql.DB
}

type GlobalOptions struct {
	HandleDir  string
	ConnConf   string
	HandleType handle.HandleMetaType
}

func initSqlPool(meta *conn.ConnMeta) *types.AppError {
	for k, v := range meta.Sql.Postgres {
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", v.User, v.Password, v.Addr, v.Port, v.Dbname)
		db, dbErr := sql.Open("postgres", connStr)

		if dbErr != nil {
			types.NewAppError(dbErr, "init postgresql pool failed [ip:%s]", v.Addr)
		}

		global.sqlPool[fmt.Sprintf("sql.postgres.%s", k)] = db
	}

	for k, v := range meta.Sql.Sqlite {
		db, dbErr := sql.Open("sqlite", v.Dbname)
		if dbErr != nil {
			types.NewAppError(dbErr, "init sqlite pool failed [ip:%s]", v.Addr)
		}

		global.sqlPool[fmt.Sprintf("sql.sqlite.%s", k)] = db
	}

	return nil
}

func Init(opt *GlobalOptions) *types.AppError {
	connMeta, connErr := conn.ReadTomlConnCfg(opt.ConnConf)

	if connErr != nil {
		return types.NewAppError(connErr, "read failed connection info file")
	}

	if pErr := initSqlPool(connMeta); pErr != nil {
		return types.NewAppError(pErr, "init pools failed")
	}

	utils, utilsErr := handle.NewLoaderUtils(opt.HandleDir, opt.HandleType)
	if utilsErr != nil {
		return utilsErr
	}

	if utils.Size() <= 0 {
		return types.NewAppError(types.ErrorAppNoData, "")
	}

	for {
		m, mErr := utils.Next()
		if mErr != nil {
			return types.NewAppError(mErr, "utils Next method error")
		}

		if m == nil {
			break
		}

		global.handles = append(global.handles, m)
	}

	return nil
}
