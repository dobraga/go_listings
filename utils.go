package main

import (
	"bytes"
	"encoding/gob"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func GetKeys(input map[string]interface{}) []string {
	var keys []string
	for k := range input {
		keys = append(keys, k)
	}
	return keys
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

func GetFirst[T int64 | int | float64 | float32 | string](listaValores []T, url string, variable string) T {
	var value T
	switch qtdElements := len(listaValores); qtdElements {
	case 0:
		log.Debug(fmt.Sprintf(`Property "%s" without %s values`, url, variable))
	case 1:
		value = listaValores[0]
	default:
		log.Warn(fmt.Sprintf(`Property "%s" with %v %s`, url, listaValores, variable))
		value = listaValores[0]
	}

	return value
}
