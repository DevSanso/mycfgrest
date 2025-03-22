package global

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	c "mycfgrest/conn"
	"mycfgrest/loader/conn"
	"mycfgrest/loader/handle"
	"mycfgrest/types"
)

var global struct {
	initOnce sync.Once
	handleMetasIdx map[string]int
	handleMetas   []handle.HandleMeta
	sqlPool map[string]c.SQLConn
}

type GlobalOptions struct {
	HandleDir  string
	ConnConf   string
	HandleType handle.HandleMetaType
}

func GetHandleMetaRerf(url string) *handle.HandleMeta {
	idx, ok := global.handleMetasIdx[url]
	if !ok {
		return nil
	}

	return &global.handleMetas[idx]
}

func GetSqlPool(name string, ctx context.Context) (c.SQLConn) {
	p, ok := global.sqlPool[name]
	if !ok {
		return nil
	}
	return p
}

func initSqlPool(meta *conn.ConnMeta) *types.AppError {
	for k, v := range meta.Sql.Postgres {
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", v.User, v.Password, v.Addr, v.Port, v.Dbname)
		db, dbErr := sql.Open("postgres", connStr)

		if dbErr != nil {
			return types.NewAppError(dbErr, "init postgresql pool failed [ip:%s]", v.Addr)
		}
		if _, ok := global.sqlPool[k]; ok {
			return types.NewAppError(types.ErrorAppDuplicate,"duplicate sqlpool name=%s", k)
		}
		global.sqlPool[k] = c.NewPgConn(db)
		
	}

	for k, v := range meta.Sql.Sqlite {
		db, dbErr := sql.Open("sqlite", v.Dbname)
		if dbErr != nil {
			types.NewAppError(dbErr, "init sqlite pool failed [ip:%s]", v.Addr)
		}
		if _, ok := global.sqlPool[k]; ok {
			return types.NewAppError(types.ErrorAppDuplicate,"duplicate sqlpool name=%s", k)
		}
		global.sqlPool[k] = c.NewPgConn(db)
	}

	return nil
}

func Init(opt *GlobalOptions) error {
	var err error = nil

	global.initOnce.Do(func() {
		global.handleMetas = make([]handle.HandleMeta, 0, 10)

		connMeta, connErr := conn.ReadTomlConnCfg(opt.ConnConf)
	
		if connErr != nil {
			err = types.NewAppError(connErr, "read failed connection info file")
			return
		}
	
		if pErr := initSqlPool(connMeta); pErr != nil {
			err = types.NewAppError(pErr, "init pools failed")
			return
		}
	
		utils, utilsErr := handle.NewLoaderUtils(opt.HandleDir, opt.HandleType)
		if utilsErr != nil {
			err = utilsErr
			return
		}
	
		if utils.Size() <= 0 {
			err = types.NewAppError(types.ErrorAppNoData, "")
			return
		}
		
		for {
			m, mErr := utils.Next()
			if mErr != nil {
				err = types.NewAppError(mErr, "utils Next method error")
				return
			}
	
			if m == nil {
				break
			}
	
			_, exists := global.handleMetasIdx[m.Data.Url]
			if exists {
				err = types.NewAppError(types.ErrorAppDuplicate, "exists handle url %s", m.Data.Url)
				return
			}
			
			global.handleMetas = append(global.handleMetas, *m)
			global.handleMetasIdx[m.Data.Url] = len(global.handleMetas)
		}
	})

	return err
}
