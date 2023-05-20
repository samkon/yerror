package yerror

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"go.uber.org/zap"
)

type merror struct {
	msg    message
	fields []zap.Field
}

const _code = "code"

var maxStack = 5

type message struct {
	Message string `json:"message"`
	Stack   string `json:"stack"`
}

// SetStackSize Set stack size of errors. Default is 5
func SetStackSize(size int) {
	maxStack = size
}

func extractStack(frames *runtime.Frames) string {
	lines := ""
	for frame, hasNext := frames.Next(); hasNext; frame, hasNext = frames.Next() {
		if !strings.Contains(frame.Function, "modanisa") {
			break
		}
		functionParts := strings.Split(frame.Function, ".")
		funcName := functionParts[len(functionParts)-1]
		fileParts := strings.Split(frame.File, "/")
		fileName := fileParts[len(fileParts)-1]
		packageName := fileParts[len(fileParts)-2]
		lines = fmt.Sprintf("%s/%s:%s():%d -> ", packageName, fileName, funcName, frame.Line) + lines
	}
	if len(lines) >= 4 {
		lines = lines[:len(lines)-4]
	}
	return lines
}

func clearKey(key string) string {
	fieldKey := strings.ReplaceAll(key, "_", "")
	fieldKey = strings.ReplaceAll(fieldKey, "-", "")
	return strings.ToLower(fieldKey)
}

func isKeysEqual(source, target string) bool {
	return clearKey(source) == clearKey(target)
}

func (m *merror) isFieldExists(key string) bool {
	for _, field := range m.fields {
		if isKeysEqual(field.Key, key) {
			return true
		}
	}
	return false
}

func (m *merror) getIndexOfField(key string) int {
	for idx := range m.fields {
		if clearKey(m.fields[idx].Key) == clearKey(key) {
			return idx
		}
	}
	return -1
}

func New(msg string, fields ...zap.Field) *merror {
	pcs := make([]uintptr, maxStack+1)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	lines := extractStack(frames)
	return &merror{
		msg: message{
			Message: msg,
			Stack:   lines,
		},
		fields: fields,
	}
}

// GetMessage Returns message struct belongs to this error.This is used to print error message in logger.
func (m *merror) GetMessage() interface{} {
	return m.msg
}

// Wrap function wraps the error passed as an argument to this function.
func Wrap(err error, fields ...zap.Field) error {
	if err == nil {
		return nil
	}
	if IsMerror(err) {
		merr := AsMerror(err)
		for _, field := range fields {
			if !merr.isFieldExists(field.Key) {
				merr.fields = append(merr.fields, field)
			} else {
				idx := merr.getIndexOfField(field.Key)
				merr.fields[idx].Key = field.Key
				merr.fields[idx].String = field.String
				merr.fields[idx].Integer = field.Integer
				merr.fields[idx].Interface = field.Interface
				merr.fields[idx].Type = field.Type
			}
		}
		return merr
	}
	pcs := make([]uintptr, maxStack+1)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	lines := extractStack(frames)
	return &merror{
		msg: message{
			Message: err.Error(),
			Stack:   lines,
		},
		fields: fields,
	}
}

func (m *merror) Error() string {
	return m.msg.Message
}

func (m *merror) GetFields() []zap.Field {
	return m.fields
}

func (m *merror) SetCode(code int) bool {
	for i, field := range m.fields {
		if field.Key == _code {
			m.fields[i] = Code(code)
			return true
		}
	}

	m.AddFields(Code(code))
	return true
}

func (m *merror) AddFields(field ...zap.Field) {
	m.fields = append(m.fields, field...)
}

func (m *merror) GetCode() int {
	for _, field := range m.fields {
		if field.Key == _code {
			return int(field.Integer)
		}
	}
	return -1
}

// IsMerror Custom Functions
func IsMerror(err error) bool {
	_, ok := err.(*merror)
	return ok
}

func AsMerror(err error) *merror {
	merr, ok := err.(*merror)
	if !ok {
		wrapped := Wrap(err)
		merr2, _ := wrapped.(*merror)
		return merr2
	}
	return merr
}

// Is Default Wrapped Functions
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
