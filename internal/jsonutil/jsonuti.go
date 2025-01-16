package jsonutil

import (
	"encoding/json"
	"errors"
)

// ToJSON - универсальная функция для преобразования объекта в JSON.
func ToJSON(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, errors.New("cannot marshal nil value")
	}
	return json.Marshal(v)
}

// FromJSON - универсальная функция для десериализации JSON в объект.
func FromJSON(data []byte, v interface{}) error {
	if len(data) <= 1 {
		return errors.New("cannot unmarshal empty data")
	}

	if v == nil {
		return errors.New("cannot unmarshal to nil object")
	}
	return json.Unmarshal(data, v)
}
