package types

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var (
	ErrorAppSys = errors.New("sys error")
	ErrorAppNoData = errors.New("no data")
	ErrorAppDuplicate = errors.New("duplicate")
	ErrorAppHttpBadRequest = errors.New("http bad request")
	ErrorAppLock = errors.New("lock error")
)

type AppError struct {
	Origin error
	msgs []string
}

func (a *AppError)Error() string {
	buf := strings.Join(a.msgs, "\n")
	return buf
}

func makeErrMsg(origin error, fn string, file string, line int, msgFmt string, args ...any) string {
	msg := "[nil]"
	if origin != nil {
		msg = origin.Error()
	}

	return fmt.Sprintf("(func:%.128s,file:%.256s:%d, origin:%.60s): %s",
		fn, file, line, msg, fmt.Sprintf(msgFmt, args...))
}

func (a *AppError)PushError(origin error, subMsgFmt string, args ...any) {
	pc, file,line, ok := runtime.Caller(1)
	
	msg := ""

	if !ok {
		msg = makeErrMsg(origin, "", "", 0, subMsgFmt, args...)
	} else {
		fn := runtime.FuncForPC(pc)
		msg = makeErrMsg(origin, fn.Name(), file, line, subMsgFmt, args...)
	}

	a.msgs = append(a.msgs, msg)
}

func NewAppError(origin error, subMsgFmt string, args ...any) *AppError {
	pc, file,line, ok := runtime.Caller(1)
	
	msg_list := make([]string, 0)
	msg := ""

	if !ok {
		msg = makeErrMsg(origin, "", "", 0, subMsgFmt, args...)
	} else {
		fn := runtime.FuncForPC(pc)
		msg = makeErrMsg(origin, fn.Name(), file, line, subMsgFmt, args...)
	}

	return &AppError {
		Origin: origin,
		msgs: append(msg_list, msg),
	}
}