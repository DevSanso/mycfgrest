package types

import (
	"sync/atomic"
)

const (
	_ParsingMapStateNone  = iota << 1
	_ParsingMapStateFetch
	_ParsingMapStatePush
	_ParsingMapPairGet
)

type _ParsingMapPair struct{Val any; ValType ParsingValueDataType}

type ParsingMap struct {
	mapDatas []map[string]_ParsingMapPair
	state int32
}

func NewParsingMap() *ParsingMap {
	o := &ParsingMap{}
	o.mapDatas = make([]map[string]_ParsingMapPair, 0, 10)
	
	return o
}

func (o *ParsingMap) chkMemAndallocColumnMap(useIdx int) {
	for len(o.mapDatas) < useIdx {
		o.mapDatas = append(o.mapDatas, make(map[string]_ParsingMapPair))
	}
}

func (o *ParsingMap) Get(idx int, key string) (data any, dataType ParsingValueDataType, err error) {
	if atomic.LoadInt32(&o.state) != _ParsingMapStateNone {
		return nil, NULL, NewAppError(ErrorAppLock, "Map is using state")
	}

	atomic.StoreInt32(&o.state, _ParsingMapPairGet)

	if mapLen := len(o.mapDatas); mapLen >= idx {
		atomic.StoreInt32(&o.state, _ParsingMapStateNone)
		return nil, NULL, NewAppError(ErrorAppSys, "Get Idx(%d) overflow map length(%d)", idx, mapLen)
	}

	val,ok := o.mapDatas[idx][key]
	if !ok {
		atomic.StoreInt32(&o.state, _ParsingMapStateNone)
		return nil, NULL, NewAppError(ErrorAppNoData, "Get Failed, callOk(%t)", ok)
	}

	atomic.StoreInt32(&o.state, _ParsingMapStateNone)
	return val.Val, val.ValType, nil
}

func (o *ParsingMap) Set(idx int, key string, val any, valType ParsingValueDataType) error {
	if atomic.LoadInt32(&o.state) == _ParsingMapStateNone {
		return NewAppError(ErrorAppLock, "Map is using fetch state")
	}

	atomic.StoreInt32(&o.state, _ParsingMapStatePush)
	
	o.chkMemAndallocColumnMap(idx)

	if chk,ok := o.mapDatas[idx][key]; ok {
		atomic.StoreInt32(&o.state, _ParsingMapStateNone)
		return NewAppError(ErrorAppDuplicate,"already exists data, key(%s) idx(%d) type(%s)", key, idx, chk.ValType)
	}


	o.mapDatas[idx][key] = _ParsingMapPair{Val: val, ValType: valType}
	atomic.StoreInt32(&o.state, _ParsingMapStateNone)
	return nil
}

func (o *ParsingMap) Fetch() (*ParsingMapFetch,error) {
	if state := atomic.LoadInt32(&o.state); state != _ParsingMapStateNone  {
		return nil, NewAppError(ErrorAppSys, "ParsingMap Fetch failed state:%d", state)
	}
	atomic.StoreInt32(&o.state, _ParsingMapStateFetch)
	
	return &ParsingMapFetch{
		mapPtr: o,
		currentIdx: 0,
		endIdx: len(o.mapDatas),
	}, nil
}

type ParsingMapFetch struct {
	mapPtr *ParsingMap

	currentIdx int
	endIdx int
}

func (f *ParsingMapFetch) Close() error {
	atomic.StoreInt32(&f.mapPtr.state, _ParsingMapStateNone)
	f.mapPtr = nil

	return nil
}

func (f *ParsingMapFetch) Next() (key []string, val []any, valType []ParsingValueDataType, err error) {
	if f.currentIdx >= f.endIdx {
		return nil, nil, nil, nil
	}

	if state := atomic.LoadInt32(&f.mapPtr.state);state != _ParsingMapStateFetch {
		return nil, nil, nil, NewAppError(ErrorAppSys, "state not fetch, state=%d", state)
	}

	key = make([]string, f.endIdx)
	val = make([]any, f.endIdx)
	valType = make([]ParsingValueDataType, f.endIdx)

	fetchIdx := 0
	for k, v := range f.mapPtr.mapDatas[f.currentIdx] {
		key[fetchIdx] = k
		val[fetchIdx] = v.Val
		valType[fetchIdx] = v.ValType
		fetchIdx += 1
	}
	err = nil
	f.currentIdx += 1
	return 
}

