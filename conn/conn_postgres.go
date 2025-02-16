package conn

import (
	"context"
	"database/sql"

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

func (pc *pgConn) Run(ctx context.Context, cmd string, param *types.ParsingValue) (*types.ParsingResultSet, error) {
	realQuery, realParam, changeErr := utils.ChangeSqlToNumBindSupportSql(cmd, param)

	if changeErr != nil {
		return nil, types.NewAppError(changeErr, "failed change real querys")
	}

	rows, rowsErr := pc.db.QueryContext(ctx, realQuery, realParam...)
	if rowsErr != nil {
		return nil, types.NewAppError(rowsErr, "failed pg connection run query")
	}
	defer rows.Close()

	cName, cErr := rows.Columns() 
	if cErr != nil {
		return nil, types.NewAppError(cErr, "get failed cols names")
	}

	ct, ctErr := rows.ColumnTypes()
	if ctErr != nil {
		return nil, types.NewAppError(ctErr, "get failed column types")
	}
	
	sysType, convertErr := convertPgTypeToSysType(ct)

	if convertErr != nil {
		return nil, types.NewAppError(convertErr, "failed convert pg type to sys type")
	}

	colBuffer := utils.NewColOutBuffer(sysType)
	output := types.NewParsingResultSet(sysType);

	rowIdx := 0
	for rows.Next() {
		if err := rows.Scan(colBuffer.GetPtrs()...); err != nil {
			return nil, types.NewAppError(err, "failed scan row data")
		}

		datas := colBuffer.GetDatas()

		for idx := range datas {
			if err := output.Set(cName[idx], rowIdx, datas[idx], sysType[idx]); err != nil {
				return nil, types.NewAppError(err, "failed result set push data [name:%s] [idx:%d]", cName[idx], rowIdx)
			}
		}
		rowIdx += 1	
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
