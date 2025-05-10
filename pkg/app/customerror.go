package app

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	appUnknown = "unknown_app"
	verUnknown = "unknown_version"
	ErrUnknown = "unknown_error"
)

var appConfig *Config

func init() {
	appConfig = &Config{
		Name:        appUnknown,
		Environment: ProductionEnv,
		Version:     verUnknown,
	}
}

func InitAppError(cfg *Config) error {
	if cfg == nil {
		appConfig = &Config{
			Name:        appUnknown,
			Environment: ProductionEnv,
			Version:     verUnknown,
		}
		return errors.New("application config is nil")
	}
	cfgCopy := *cfg
	appConfig = &cfgCopy
	return nil
}

type Error struct {
	Code       int    // HTTP status code
	Message    string // Safe message for HTTP response
	DevMessage string // Developer/debugging message
	Err        error  // Underlying error from the repository or other layers
}

func NewError(code int, message string, devMessage string, err error) Error {
	return Error{Code: code, Message: message, DevMessage: devMessage, Err: err}
}

func (a Error) Unwrap() error {
	return a.Err
}

func (a Error) Error() string {
	if a.Err == nil {
		return fmt.Sprintf("[%d %s]: %s", a.Code, a.Message, a.DevMessage)
	}
	return fmt.Sprintf("[%d %s]: %s: %s", a.Code, a.Message, a.DevMessage, a.Err.Error())
}

func GetLastMessage(err error) string {
	if err == nil {
		return ""
	}
	var myErr Error
	if errors.As(err, &myErr) {
		if appConfig.IsProduction() {
			return myErr.Message
		} else if appConfig.IsDevelopment() {
			if myErr.Err != nil {
				return myErr.DevMessage + ": " + myErr.Unwrap().Error()
			}
			return myErr.DevMessage
		}
		return myErr.Message
	} else {
		if appConfig.IsDevelopment() {
			return err.Error()
		}
		return ErrUnknown
	}
}

func GetCode(err error) int {
	var myErr Error
	if errors.As(err, &myErr) {
		return myErr.Code
	}
	return http.StatusInternalServerError
}
