package utils

import (
	"mycfgrest/types"
)

type ColOutBuffer struct {
	realData []any
	dataPtr []any
}

func NewColOutBuffer(cols []types.ParsingValueDataType) ColOutBuffer {
	realSlice := make([]any, len(cols))
	dataPtrSlice := make([]any, len(cols))

	for idx := range cols {
		switch cols[idx] {
		case types.DOUBLE:
			realSlice[idx] = 0.0
		case types.INT:
			realSlice[idx] = 0
		case types.STRING:
			realSlice[idx] = ""
		}

		dataPtrSlice[idx] = &realSlice[idx]
	}

	return ColOutBuffer{
		realData: realSlice,
		dataPtr: dataPtrSlice,
	}
}

func (buf *ColOutBuffer) GetDatas() []any {
	return buf.realData
}

func (buf *ColOutBuffer) GetPtrs() []any {
	return buf.dataPtr
}