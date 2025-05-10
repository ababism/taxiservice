package app

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestInitAppError(t *testing.T) {

	t.Run("NilConfig", func(t *testing.T) {
		t.Parallel()

		err := InitAppError(nil)
		require.Error(t, err)
		assert.Equal(t, "application config is nil", err.Error())
	})

	t.Run("ValidConfig", func(t *testing.T) {

		cfg := &Config{
			Name:        "TestApp",
			Environment: ProductionEnv,
			Version:     "1.0.0",
		}
		err := InitAppError(cfg)
		require.NoError(t, err)
	})
}

func TestNewAppError(t *testing.T) {
	t.Parallel()

	err := NewError(http.StatusNotFound, "Not Found", "Resource not found", errors.New("not found"))
	require.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Equal(t, "Not Found", err.Message)
	assert.Equal(t, "Resource not found", err.DevMessage)
	assert.NotNil(t, err.Err)
	assert.Equal(t, "not found", err.Err.Error())
}

func TestGetLastMessage(t *testing.T) {

	const (
		message    = "Internal Server Error"
		devMessage = "Unexpected error occurred"
	)

	testNestedErr := errors.New("internal error")
	testErr := NewError(http.StatusInternalServerError, message, devMessage, testNestedErr)

	t.Run("AppErrorProduction", func(t *testing.T) {

		appConfig = &Config{
			Environment: ProductionEnv,
		}
		require.Equal(t, appConfig.IsProduction(), true)

		resultMessage := GetLastMessage(testErr)
		assert.Equal(t, message, resultMessage)
	})

	t.Run("AppErrorLocal", func(t *testing.T) {

		appConfig = &Config{
			Environment: DevelopmentEnv,
		}
		require.Equal(t, appConfig.IsDevelopment(), true)

		resultMessage := GetLastMessage(testErr)
		assert.Equal(t, devMessage+": "+testNestedErr.Error(), resultMessage)
	})

	t.Run("NonAppErrorLocal", func(t *testing.T) {
		const errorText = "random error"

		appConfig = &Config{
			Environment: DevelopmentEnv,
		}
		require.Equal(t, appConfig.IsDevelopment(), true)

		resultMessage := GetLastMessage(errors.New(errorText))
		assert.Equal(t, errorText, resultMessage)
	})

	t.Run("NonAppErrorProduction", func(t *testing.T) {
		const errorText = "random error"

		appConfig = &Config{
			Environment: ProductionEnv,
		}
		require.Equal(t, appConfig.IsProduction(), true)

		resultMessage := GetLastMessage(errors.New(errorText))
		assert.Equal(t, ErrUnknown, resultMessage)
	})

	t.Run("DoubleWrappedProduction", func(t *testing.T) {
		insideAppErr := NewError(http.StatusBadRequest, "inside "+message, "inside "+devMessage,
			errors.New("double inside text"))

		doubleErr := NewError(http.StatusInternalServerError, message, devMessage, insideAppErr)
		const errorText = "random error"

		appConfig = &Config{
			Environment: ProductionEnv,
		}
		require.Equal(t, appConfig.IsProduction(), true)

		resultMessage := GetLastMessage(doubleErr)
		assert.Equal(t, message, resultMessage)
	})

	t.Run("DoubleWrappedLocal", func(t *testing.T) {
		insideAppErr := NewError(http.StatusBadRequest, "nested "+message, "nested "+devMessage,
			errors.New("double-nested error text"))

		doubleErr := NewError(http.StatusInternalServerError, message, devMessage, insideAppErr)

		appConfig = &Config{
			Environment: DevelopmentEnv,
		}
		require.True(t, appConfig.IsDevelopment())

		resultMessage := GetLastMessage(doubleErr)
		assert.Equal(t, devMessage+": "+insideAppErr.Error(), resultMessage)
	})

	t.Run("NilErrProduction", func(t *testing.T) {
		appConfig = &Config{
			Environment: ProductionEnv,
		}
		require.Equal(t, appConfig.IsProduction(), true)

		resultMessage := GetLastMessage(nil)
		assert.Equal(t, "", resultMessage)
	})

	t.Run("NilErrLocal", func(t *testing.T) {
		appConfig = &Config{
			Environment: DevelopmentEnv,
		}
		require.Equal(t, appConfig.IsDevelopment(), true)

		resultMessage := GetLastMessage(nil)
		assert.Equal(t, "", resultMessage)
	})
}

func TestGetCode(t *testing.T) {
	t.Parallel()

	err := NewError(http.StatusNotFound, "Not Found", "Resource not found", errors.New("not found"))

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		code := GetCode(err)
		assert.Equal(t, http.StatusNotFound, code)
	})

	t.Run("NonAppError", func(t *testing.T) {
		t.Parallel()

		code := GetCode(errors.New("unknown error"))
		assert.Equal(t, http.StatusInternalServerError, code)
	})
}

func TestAppErrorUnwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	appErr := NewError(http.StatusInternalServerError, "Internal Server Error", "Unexpected error occurred", innerErr)

	unwrappedErr := appErr.Unwrap()
	assert.Equal(t, innerErr, unwrappedErr)
}

func TestAppErrorError(t *testing.T) {
	const (
		message    = "Internal Server Error"
		devMessage = "Unexpected error occurred"

		errorText = "inner error"
	)

	t.Run("GoodScenario", func(t *testing.T) {
		innerErr := errors.New(errorText)
		appErr := NewError(http.StatusInternalServerError, message, devMessage, innerErr)

		expectedErrorMessage := fmt.Sprintf("[%d %s]: %s: %s", http.StatusInternalServerError, message, devMessage, errorText)
		assert.Equal(t, expectedErrorMessage, appErr.Error())
	})
	t.Run("NilErrScenario", func(t *testing.T) {

		appErr := NewError(http.StatusInternalServerError, message, devMessage,
			nil)

		expectedErrorMessage := fmt.Sprintf("[%d %s]: %s", http.StatusInternalServerError, message, devMessage)
		assert.Equal(t, expectedErrorMessage, appErr.Error())
	})

}
