package jsonpath

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func ExtractString(jsonStr []byte, path string) (string, error) {
	iface, err := Extract(jsonStr, path)
	if err != nil {
		return "", err
	}

	str, ok := iface.(string)

	if !ok {
		return "", fmt.Errorf("Element [%v] is not a string.", iface)
	}

	return str, nil
}

func ExtractFloat64(jsonStr []byte, path string) (float64, error) {
	iface, err := Extract(jsonStr, path)
	if err != nil {
		return 0, err
	}

	flt, ok := iface.(float64)

	if !ok {
		return 0, fmt.Errorf("Element [%v] is not a number.", iface)
	}

	return flt, nil
}

// Extract unmarshals the given json object or array and extracts a subobject, array,
// or value based on the given dot-separated path. Elements in the path can be a
// field name or array index. Array indices can be negative to be counted backward
// from the end of the array.
// As an example: {"a": [{"val":0}, {"val":7}]}  "a.[-1].val" would return 7.
func Extract(jsonStr []byte, path string) (interface{}, error) {
	var rtn interface{}

	err := json.Unmarshal(jsonStr, &rtn)

	if err != nil {
		return nil, err
	}

	pieces := strings.Split(path, ".")

	for _, piece := range pieces {
		if len(piece) < 1 {
			return nil, fmt.Errorf("Empty field in path")
		}

		if piece[0] == '[' {
			rtn, err = extractArray(rtn, piece)

			if err != nil {
				return nil, err
			}
		} else {
			rtn, err = extractObject(rtn, piece)

			if err != nil {
				return nil, err
			}
		}
	}

	return rtn, nil
}

func extractObject(jsonIface interface{}, piece string) (interface{}, error) {
	m, ok := jsonIface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("No object found when looking for field %v", piece)
	}

	rtn, ok := m[piece]
	if !ok {
		return nil, fmt.Errorf("Field not found when looking for field %v", piece)
	}

	return rtn, nil
}

func extractArray(jsonIface interface{}, piece string) (interface{}, error) {
	if piece[len(piece)-1] != ']' {
		return nil, fmt.Errorf("Array index missing ']' in %v", piece)
	}

	a, ok := jsonIface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("No array found when looking for %v", piece)
	}

	piece = piece[1 : len(piece)-1]

	index, err := strconv.Atoi(piece)
	if err != nil {
		return nil, fmt.Errorf("Parsing array index %v gave %v", piece, err.Error())
	}

	if index >= len(a) || (index < -len(a)) {
		return nil, fmt.Errorf("Array index %v out of range %v", index, len(a))
	}

	if index >= 0 {
		return a[index], nil
	} else {
		return a[len(a)+index], nil
	}
}
