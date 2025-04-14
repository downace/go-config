package config

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"strings"
)

type Serializer[T any] interface {
	SerializeData(data *T) ([]byte, error)
	DeserializeData(rawData []byte) (*T, error)
}

// YamlSerializer uses YAML format to store config
type YamlSerializer[T any] struct{}

func (s YamlSerializer[T]) SerializeData(data *T) (result []byte, err error) {
	defer panicToError(&err)
	return yaml.Marshal(data)
}

func (s YamlSerializer[T]) DeserializeData(rawData []byte) (*T, error) {
	var data T

	err := yaml.Unmarshal(rawData, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

// JsonSerializer uses JSON format to store config
type JsonSerializer[T any] struct {
	Indent int
}

func (s JsonSerializer[T]) SerializeData(data *T) ([]byte, error) {
	if s.Indent > 0 {
		return json.MarshalIndent(data, "", strings.Repeat(" ", s.Indent))
	} else {
		return json.Marshal(data)
	}
}

func (s JsonSerializer[T]) DeserializeData(rawData []byte) (*T, error) {
	var data T

	err := json.Unmarshal(rawData, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}
