package pgarray

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// Int32Slice can be scanned out of the DB
type Int32Slice []int32

// Scan Implements sql.Scanner for the String slice type
// Scanners take the database value (in this case as a byte slice)
// and sets the value of the type.  Here we cast to a string and
// do a regexp based parse
func (s *Int32Slice) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	asBytes, ok := src.([]byte)
	if !ok {
		return error(fmt.Errorf("Scan source was not []bytes: %v", src))
	}

	asString := string(asBytes)
	parsed := s.parseArray(asString)
	(*s) = Int32Slice(parsed)

	return nil
}

// Value implements driver.Valuer interface which allows StringSlices to be automatically parsed
// when passed into sql functions
func (s Int32Slice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	var buffer bytes.Buffer

	buffer.WriteString("{")
	last := len(s) - 1
	for i, val := range s {
		buffer.WriteString(fmt.Sprintf("%d", val))
		if i != last {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")

	return string(buffer.Bytes()), nil
}

// Parse the output string from the array type.
// Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
func (s *Int32Slice) parseArray(array string) []int32 {
	var results []int32
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		i64, _ := strconv.ParseInt(s, 10, 32)
		results = append(results, int32(i64))
	}
	return results
}

// Int64Slice can be scanned out of the DB
type Int64Slice []int64

// Scan Implements sql.Scanner for the String slice type
// Scanners take the database value (in this case as a byte slice)
// and sets the value of the type.  Here we cast to a string and
// do a regexp based parse
func (s *Int64Slice) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	asBytes, ok := src.([]byte)
	if !ok {
		return error(fmt.Errorf("Scan source was not []bytes: %v", src))
	}

	asString := string(asBytes)
	parsed := s.parseArray(asString)
	(*s) = Int64Slice(parsed)

	return nil
}

// Value implements driver.Valuer interface which allows StringSlices to be automatically parsed
// when passed into sql functions
func (s Int64Slice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	var buffer bytes.Buffer

	buffer.WriteString("{")
	last := len(s) - 1
	for i, val := range s {
		buffer.WriteString(fmt.Sprintf("%d", val))
		if i != last {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")

	return string(buffer.Bytes()), nil
}

// Parse the output string from the array type.
// Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
func (s *Int64Slice) parseArray(array string) []int64 {
	var results []int64
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		i64, _ := strconv.ParseInt(s, 10, 64)
		results = append(results, i64)
	}
	return results
}
