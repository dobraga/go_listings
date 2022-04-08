package main

import (
	"bytes"
	"encoding/gob"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetBytes(key interface{}) ([]byte, error) {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}
