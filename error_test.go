package yerror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateNewErrorWithMerrorPackage(t *testing.T) {
	err := New("Hello World",
		Code(200),
		zap.String("STRING", "TEST"),
		zap.Int("INT", 24),
		zap.Error(errors.New("ERROR")),
		zap.Bool("BOOL", true),
		zap.Float64("TEST", 24.1),
	)

	assert.NotNil(t, err)
}

func TestGetCodeFunctionWithMerrorAndNotMerror(t *testing.T) {
	err1 := New("Hello World",
		Code(200),
		zap.String("STRING", "TEST"),
		zap.Int("INT", 24),
		zap.Error(errors.New("ERROR")),
		zap.Bool("BOOL", true),
		zap.Float64("TEST", 24.1),
	)

	merr := AsMerror(err1)

	assert.Equal(t, merr.GetCode(), 200)

}

func TestSetCodeFuntionWhenErrorIsMerror(t *testing.T) {
	err := New("Hello World",
		Code(200),
		zap.String("STRING", "TEST"),
		zap.Int("INT", 24),
		zap.Error(errors.New("ERROR")),
		zap.Bool("BOOL", true),
		zap.Float64("TEST", 24.1),
	)

	err.SetCode(500)

	assert.Equal(t, err.GetCode(), 500)
}

func TestSetCodeFuntionWhenErrorIsMerrorAndCodeNotImplemented(t *testing.T) {
	err := New("Hello World",
		zap.String("STRING", "TEST"),
		zap.Int("INT", 24),
		zap.Error(errors.New("ERROR")),
		zap.Bool("BOOL", true),
		zap.Float64("TEST", 24.1),
	)

	err.SetCode(400)

	assert.Equal(t, err.GetCode(), 400)
}

func TestWhenErrorIsMerrorThenShouldGetFields(t *testing.T) {
	err := New("Hello World",
		Code(200),
		zap.String("STRING", "TEST"),
		zap.Int("INT", 24),
		zap.Error(errors.New("ERROR")),
		zap.Bool("BOOL", true),
		zap.Float64("TEST", 24.1),
	)

	assert.Equal(t, err.GetFields(), []zap.Field{
		Code(200),
		zap.String("STRING", "TEST"),
		zap.Int("INT", 24),
		zap.Error(errors.New("ERROR")),
		zap.Bool("BOOL", true),
		zap.Float64("TEST", 24.1),
	})
}

func TestWhenErrorIsMerrorThenIsMerrorFunctionShouldReturnTrue(t *testing.T) {
	err := New("Hello World",
		Code(200),
		zap.String("STRING", "TEST"),
		zap.Int("INT", 24),
		zap.Error(errors.New("ERROR")),
		zap.Bool("BOOL", true),
		zap.Float64("TEST", 24.1),
	)

	assert.True(t, IsMerror(err))

}

func TestWhenErrorIsNotMerrorThenIsMerrorFunctionShouldReturnFalse(t *testing.T) {
	err := errors.New("Hello World")

	assert.False(t, IsMerror(err))

}

func TestWhenIsFunctionUsedThenShouldReturnExpected(t *testing.T) {
	err := New("Hello World",
		Code(200),
		zap.String("STRING", "TEST"),
		zap.Int("INT", 24),
		zap.Error(errors.New("ERROR")),
		zap.Bool("BOOL", true),
		zap.Float64("TEST", 24.1),
	)

	assert.True(t, Is(err, err))
	assert.False(t, Is(err, errors.New("ERROR")))
}

func TestWhenAsFunctionUsedThenShouldReturnExpected(t *testing.T) {
	err := New("Hello World",
		Code(200),
		zap.String("STRING", "TEST"),
		zap.Int("INT", 24),
		zap.Error(errors.New("ERROR")),
		zap.Bool("BOOL", true),
		zap.Float64("TEST", 24.1),
	)

	expectedStatus := 200
	expectedFields := []zap.Field{
		Code(200),
		zap.String("STRING", "TEST"),
		zap.Int("INT", 24),
		zap.Error(errors.New("ERROR")),
		zap.Bool("BOOL", true),
		zap.Float64("TEST", 24.1),
	}

	merr := AsMerror(err)

	assert.Equal(t, merr.GetCode(), expectedStatus)
	assert.Equal(t, merr.GetFields(), expectedFields)
}

func TestGivenErrorWhenWrapThenShouldBeSameMerror(t *testing.T) {
	err := New("this is an error", zap.String("key", "value"))
	wrappedErr := Wrap(err)
	assert.Equal(t, wrappedErr, err)
}

func TestGivenNilToWrapWhenWrapThenShouldBeNil(t *testing.T) {
	wrappedErr := Wrap(nil)
	assert.Nil(t, wrappedErr)
}

func TestGivenNilToWrapAsMerrorWhenWrapThenShouldBeNil(t *testing.T) {
	err := Wrap(nil)
	merr := AsMerror(err)
	assert.Nil(t, merr)
}
