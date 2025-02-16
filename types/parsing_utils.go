package types

import (
	"sync/atomic"
)

type ParsingValueDataType string

const (
	INT    ParsingValueDataType = "int"
	STRING ParsingValueDataType = "string"
	DOUBLE ParsingValueDataType = "double"
)

type ParsingValue struct {
	size int

	keys []string
	vals []any

	dataTypes []ParsingValueDataType

	is_fetching atomic.Bool
	is_pushing  atomic.Bool
}

type parsingValueFetch struct {
	current int

	pm *ParsingValue
}

func NewParsingValue() *ParsingValue {
	return &ParsingValue{
		size:      0,
		keys:      make([]string, 10),
		vals:      make([]any, 10),
		dataTypes: make([]ParsingValueDataType, 10),
	}
}

func (pm *ParsingValue) Push(key string, val any, dataType ParsingValueDataType) *AppError {
	pm.is_pushing.Store(true)

	if pm.is_fetching.Load() {
		pm.is_pushing.Store(false)
		return NewAppError(ErrorAppLock, "fetch locking")
	}

	pm.keys = append(pm.keys, key)
	pm.vals = append(pm.vals, val)
	pm.dataTypes = append(pm.dataTypes, dataType)

	pm.is_pushing.Store(false)
	return nil
}

func (pm *ParsingValue) Fetch() (parsingValueFetch, *AppError) {
	if pm.is_pushing.Load() {
		return parsingValueFetch{}, NewAppError(ErrorAppLock, "push locking")
	}

	if pm.is_fetching.Load() {
		return parsingValueFetch{}, NewAppError(ErrorAppSys, "already fetching")
	}

	pm.is_fetching.Store(true)

	return parsingValueFetch{
		current: 0,
		pm:      pm,
	}, nil
}
func (pmf *parsingValueFetch) Close() {
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

func (pmf *parsingValueFetch) FullSize() int {
	return pmf.pm.size
}

func (pmf *parsingValueFetch) Keys(out []string) {
	copy(out, pmf.pm.keys)
}

func (pmf *parsingValueFetch) Values(out []any) {
	copy(out, pmf.pm.vals)
}

func (pmf *parsingValueFetch) Types(out []ParsingValueDataType) {
	copy(out, pmf.pm.dataTypes)
}


type ParsingResultSet struct {
	sets map[struct{idx int; colName string}]any
	colTypes []ParsingValueDataType
}

func NewParsingResultSet(ts []ParsingValueDataType) *ParsingResultSet {
	copyTs := make([]ParsingValueDataType, len(ts))
	copy(copyTs, ts)
	
	return &ParsingResultSet{
		sets : make(map[struct{idx int; colName string}]any),
		colTypes : copyTs,
	}
}

func (prs *ParsingResultSet) Set(name string, idx int, data any, t ParsingValueDataType) *AppError {
	if t != prs.colTypes[idx] {
		return NewAppError(ErrorAppSys, "%s is not matching type", name)
	}

	if data := prs.sets[struct{idx int; colName string}{idx : idx, colName: name}]; data != nil {
		return NewAppError(ErrorAppDuplicate, "%s is duplicate", name)
	}

	prs.sets[struct{idx int; colName string}{idx:idx, colName: name}] = data

	return nil
}

func (prs *ParsingResultSet) Get(name string, idx int) (data any) {
	data = prs.sets[struct{idx int; colName string}{idx: idx, colName: name}]
	return 
}