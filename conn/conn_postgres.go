package conn

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"mycfgrest/conn/utils"
	"mycfgrest/types"

	_ "github.com/lib/pq"
)

type pgConn struct {
	db *sql.DB
}

func convertPgTypeToSysType(cts []*sql.ColumnType) ([]types.ParsingValueDataType, *types.AppError) {
	ret := make([]types.ParsingValueDataType, len(cts))

	for i, ct := range cts {
		switch ct.Name() {
		case "INT2", "INT4", "INT8":
			ret[i]=types.INT
		case "CHAR", "VARCHAR":
			ret[i]=types.STRING
		case "FLOAT4", "FLOAT8":
			ret[i] =types.DOUBLE
		default:
			return nil, types.NewAppError(types.ErrorAppNoData, "not support this pg type [%s]", ct.Name())
		}
	}

	return ret, nil
}

func (pc *pgConn) runEach(ctx context.Context, cmd string, prefix string, fetch *types.ParsingMapFetch, output *types.ParsingMap) error {
	realQuery, realParam, changeErr := utils.ChangeSqlToNumBindSupportSql(cmd, fetch)

	if changeErr != nil {
		return types.NewAppError(changeErr, "failed change real querys")
	}

	rows, rowsErr := pc.db.QueryContext(ctx, realQuery, realParam...)
	if rowsErr != nil {
		return types.NewAppError(rowsErr, "failed pg connection run query")
	}
	defer rows.Close()

	cName, cErr := rows.Columns() 
	if cErr != nil {
		return types.NewAppError(cErr, "get failed cols names")
	}

	ct, ctErr := rows.ColumnTypes()
	if ctErr != nil {
		return types.NewAppError(ctErr, "get failed column types")
	}
	
	sysType, convertErr := convertPgTypeToSysType(ct)

	if convertErr != nil {
		return types.NewAppError(convertErr, "failed convert pg type to sys type")
	}

	colBuffer := utils.NewColOutBuffer(sysType)

	rowIdx := 0
	for rows.Next() {
		if err := rows.Scan(colBuffer.GetPtrs()...); err != nil {
			return types.NewAppError(err, "failed scan row data")
		}

		datas := colBuffer.GetDatas()

		for idx := range datas {
			if err := output.Set(rowIdx, strings.Join([]string{prefix, strconv.Itoa(idx), cName[idx]},"."), datas[idx], sysType[idx]); err != nil {
				return types.NewAppError(err, "failed result set push data [name:%s] [idx:%d]", cName[idx], rowIdx)
			}
		}
		rowIdx += 1	
	}

	return nil
}

func (pc *pgConn) Run(ctx context.Context, cmd string, param *types.ParsingMap) (*types.ParsingMap, error) {
	output := types.NewParsingMap()
	var fetch *types.ParsingMapFetch = nil
	var err error = nil

	if fetch, err = param.Fetch(); err != nil {
		return nil, types.NewAppError(types.ErrorAppSys, "pgConn Run Failed")
	}

	idx := 0
	for isEnd := fetch.IsEnd(); !isEnd; isEnd = fetch.Next() {
		if err = pc.runEach(ctx, cmd, strconv.Itoa(idx), fetch, output); err != nil {
			return nil, types.NewAppError(err, "")
		}
		idx += 1
	}
	
	return output, nil
}

func NewPgConn() (Conn, error) {
	if db, err := sql.Open("postgres", ""); err != nil {
		return nil, types.NewAppError(err, "failed pg connection")
	} else {
		db.SetMaxIdleConns(1)
		db.SetMaxOpenConns(1)
		return &pgConn{db: db}, nil
	}
}
