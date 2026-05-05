// Package sysfs provides functions to read/write special files in sysfs, such
// as "/sys/bus/iio/devices/iio:device0".
// This can be used for IIO, GPIO etc.
package sysfs

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//-------------------------------------------------------------------------------------------------
// Seams for configuration & testing

// SysFsRoot sets the root for all relative node paths.
var SysFsRoot string

var GetString = func(node string) (string, error) {
	return readTextFile(fullPath(node))
}

var SetString = func(node, value string) error {
	return os.WriteFile(fullPath(node), []byte(value), 0777)
}

var HexPrefix = "0x"

//-------------------------------------------------------------------------------------------------

func fullPath(node string) string {
	if strings.HasPrefix(node, "/") || SysFsRoot == "" {
		return node
	}

	return filepath.Join(SysFsRoot, node)
}

func readTextFile(path string) (string, error) {
	bs, e2 := os.ReadFile(path)
	if e2 != nil {
		return "", e2
	}

	return strings.TrimSpace(string(bs)), nil
}

//-------------------------------------------------------------------------------------------------

// GetValue gets an arbitrary value of any type given a parsing function for
// that type.
func GetValue[T any](node string, parse func(string) (T, error)) (T, error) {
	s, err := GetString(node)
	if err != nil {
		var zero T
		return zero, err
	}
	return parse(s)
}

// GetInt64 gets an integer as a signed int64.
func GetInt64(node string) (int64, error) {
	return GetValue(node, parseInt64)
}

// GetHex64 gets a hexadecimal integer, discarding any "0x" prefix if present.
func GetHex64(node string) (int64, error) {
	return GetValue(node, parseHex64)
}

// GetInt gets an integer as an int.
func GetInt(node string) (int, error) {
	return GetValue(node, strconv.Atoi)
}

// GetBool gets a boolean value. True values are 1, t, or true.
// False values are 0, f, or false. The case is ignored.
func GetBool(node string) (bool, error) {
	return GetValue(node, strconv.ParseBool)
}

// GetStrings gets a slice of space-separated strings. If there is a
// surrounding "[" and "]", these are discarded.
func GetStrings(node string) ([]string, error) {
	s, err := GetString(node)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		s = s[1 : len(s)-1]
	}
	return strings.Split(s, " "), nil
}

//-------------------------------------------------------------------------------------------------

func parseInt64(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }

func parseHex64(s string) (int64, error) {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	return strconv.ParseInt(s, 16, 64)
}

//-------------------------------------------------------------------------------------------------

// SetInt64 writes an integer in base 10.
func SetInt64(node string, value int64) error {
	return SetString(node, strconv.FormatInt(value, 10))
}

// SetHex64 writes an integer in base 16 prefixed by [HexPrefix].
func SetHex64(node string, value int64) error {
	return SetString(node, HexPrefix+strconv.FormatInt(value, 16))
}

// SetInt writes an integer in base 10.
func SetInt(node string, value int) error {
	return SetString(node, strconv.Itoa(value))
}

// SetBool writes a boolean value as 1 or 0.
func SetBool(node string, value bool) error {
	if value {
		return SetString(node, "1")
	}
	return SetString(node, "0")
}
