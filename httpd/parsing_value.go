package httpd

import (
	"sync/atomic"

	"mycfgrest/app_error"
)

type ParsingValueDataType string

const (
	INT ParsingValueDataType = "int"
	STRING = "string"
	DOUBLE = "double"
)

type ParsingValue struct {
	size int

	keys []string
	vals []any

    dataTypes []ParsingValueDataType

	is_fetching atomic.Bool
	is_pushing atomic.Bool
}

type parsingValueFetch struct {
	current int
	
	pm *ParsingValue
}


func NewParsingValue() *ParsingValue {
	return &ParsingValue{
		size : 0,
		keys : make([]string, 10),
		vals : make([]any, 10),
		dataTypes: make([]ParsingValueDataType, 10),
	}
}

func (pm *ParsingValue)Push(key string, val any, dataType ParsingValueDataType) *app_error.AppError {
	pm.is_pushing.Store(true)
	
	if pm.is_fetching.Load() {
		pm.is_pushing.Store(false)
		return app_error.NewError(app_error.ErrorLock, "fetch locking")
	}

	pm.keys = append(pm.keys, key)
	pm.vals = append(pm.vals, val)
	pm.dataTypes = append(pm.dataTypes, dataType)

	pm.is_pushing.Store(false)
	return nil
}

func (pm *ParsingValue) Fetch() (parsingValueFetch, *app_error.AppError ) {
	if pm.is_pushing.Load() {
		return parsingValueFetch{}, app_error.NewError(app_error.ErrorLock, "push locking")
	}

	if pm.is_fetching.Load() {
		return parsingValueFetch{}, app_error.NewError(app_error.ErrorSys, "already fetching")
	}

	pm.is_fetching.Store(true)

	return parsingValueFetch {
		current: 0,
		pm : pm,
	}, nil
}
func (pmf *parsingValueFetch)Close() {
	pmf.pm.is_fetching.Store(false)
}
func (pmf *parsingValueFetch) Next() (is_end bool, key string, val any, dataType ParsingValueDataType) {
	if pmf.current >= pmf.pm.size {
		is_end = true
		return 
	}
	is_end = false
	key = pmf.pm.keys[pmf.current]
	val = pmf.pm.vals[pmf.current]
	dataType = pmf.pm.dataTypes[pmf.current]

	pmf.current += 1
	return
}