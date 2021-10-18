// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: configs/config.proto

package configs

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
)

// Validate checks the field values on Config with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Config) Validate() error {
	if m == nil {
		return nil
	}

	if _, ok := _Config_Mode_InLookup[m.GetMode()]; !ok {
		return ConfigValidationError{
			field:  "Mode",
			reason: "value must be in list [debug normal]",
		}
	}

	if utf8.RuneCountInString(m.GetGrpcAddr()) < 3 {
		return ConfigValidationError{
			field:  "GrpcAddr",
			reason: "value length must be at least 3 runes",
		}
	}

	if utf8.RuneCountInString(m.GetHttpAddr()) < 3 {
		return ConfigValidationError{
			field:  "HttpAddr",
			reason: "value length must be at least 3 runes",
		}
	}

	if utf8.RuneCountInString(m.GetPprofAddr()) < 3 {
		return ConfigValidationError{
			field:  "PprofAddr",
			reason: "value length must be at least 3 runes",
		}
	}

	// no validation rules for ApmUrl

	// no validation rules for MysqlConf

	// no validation rules for TenSecretId

	// no validation rules for TenSecretKey

	// no validation rules for TenSmsConf

	// no validation rules for WechatMiniAppid

	// no validation rules for WechatMiniSecret

	// no validation rules for JwtKey

	return nil
}

// ConfigValidationError is the validation error returned by Config.Validate if
// the designated constraints aren't met.
type ConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ConfigValidationError) ErrorName() string { return "ConfigValidationError" }

// Error satisfies the builtin error interface
func (e ConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ConfigValidationError{}

var _Config_Mode_InLookup = map[string]struct{}{
	"debug":  {},
	"normal": {},
}