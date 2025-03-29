package types

import (
	"fmt"
	"sync/atomic"
)

type ParsingValueDataType string

const (
	INT    ParsingValueDataType = "int"
	STRING ParsingValueDataType = "string"
	DOUBLE ParsingValueDataType = "double"
	NULL ParsingValueDataType = "NULL"
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
	for len(o.mapDatas) < useIdx + 1 {
		o.mapDatas = append(o.mapDatas, make(map[string]_ParsingMapPair))
	}
}
func (o *ParsingMap) OverReadFrom(src *ParsingMap) error {
	if atomic.LoadInt32(&o.state) != _ParsingMapStateNone  || atomic.LoadInt32(&src.state) != _ParsingMapStateNone{
		return NewAppError(ErrorAppLock, "")
	}

	atomic.StoreInt32(&o.state, _ParsingMapStatePush)
	atomic.StoreInt32(&src.state, _ParsingMapPairGet)
	o.chkMemAndallocColumnMap(len(src.mapDatas))
	for idx := range src.mapDatas {
		for key, pair := range src.mapDatas[idx] {
			src.mapDatas[idx][key] = pair
		}
	}

	atomic.StoreInt32(&src.state, _ParsingMapStateNone)
	atomic.StoreInt32(&o.state, _ParsingMapStatePush)

	return nil
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
	if atomic.LoadInt32(&o.state) != _ParsingMapStateNone {
		return NewAppError(ErrorAppLock, "Map is using fetch state")
	}

	atomic.StoreInt32(&o.state, _ParsingMapStatePush)
	
	o.chkMemAndallocColumnMap(idx)

	if len(o.mapDatas) < idx {
		panic(fmt.Sprintf("ParsingMap idx[%d] is over mapData[%d]", idx, len(o.mapDatas)))
	}

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

func (o *ParsingMap) FetchOne(idx int) (*ParsingMapFetch,error) {
	if mapLen := len(o.mapDatas); mapLen >= idx {
		atomic.StoreInt32(&o.state, _ParsingMapStateNone)
		return nil, NewAppError(ErrorAppSys, "Get Idx(%d) overflow map length(%d)", idx, mapLen)
	}
	
	if state := atomic.LoadInt32(&o.state); state != _ParsingMapStateNone  {
		return nil, NewAppError(ErrorAppSys, "ParsingMap Fetch failed state:%d", state)
	}
	atomic.StoreInt32(&o.state, _ParsingMapStateFetch)
	
	return &ParsingMapFetch{
		mapPtr: o,
		currentIdx: idx,
		endIdx: idx + 1,
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

func (f *ParsingMapFetch) Next() (isEnd bool) {
	if f.currentIdx >= f.endIdx {
		isEnd = true
		return 
	}
	isEnd = false
	f.currentIdx += 1
	return 
}

func (f *ParsingMapFetch) IsEnd() (isEnd bool) {
	if f.currentIdx >= f.endIdx {
		isEnd = true
		return 
	}
	isEnd = false
	return 
}

func (f *ParsingMapFetch) Reset(idx int) (isEnd bool) {
	if idx >= f.endIdx {
		isEnd = true
		return 
	}
	isEnd = false
	f.currentIdx = idx
	return 
}



func (f *ParsingMapFetch) GetData() (key []string, val []any, valType []ParsingValueDataType, err error) {
	if f.currentIdx >= f.endIdx {
		return nil, nil, nil, nil
	}

	if state := atomic.LoadInt32(&f.mapPtr.state);state != _ParsingMapStateFetch {
		return nil, nil, nil, NewAppError(ErrorAppSys, "state not fetch, state=%d", state)
	}
	allocSize := f.endIdx + 1
	key = make([]string, allocSize)
	val = make([]any, allocSize)
	valType = make([]ParsingValueDataType, allocSize)

	fetchIdx := 0

	if f.currentIdx >= len(f.mapPtr.mapDatas) {
		panic(fmt.Sprintf("overflow currentIdx[%d] mapData size[%d]", f.currentIdx, len(f.mapPtr.mapDatas)))
	}
	
	for k, v := range f.mapPtr.mapDatas[f.currentIdx] {
		key[fetchIdx] = k
		val[fetchIdx] = v.Val
		valType[fetchIdx] = v.ValType
		fetchIdx += 1
	}
	err = nil
	return 
}
