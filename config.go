// Package config provides simple persistent storage for application configuration
package config

import (
	"errors"
	"fmt"
	"github.com/LastPossum/kamino"
	"sync"
)

type Config[T any] struct {
	mu         sync.Mutex
	storage    Storage
	serializer Serializer[T]

	Data T
}

func NewConfig[T any](defaultData T, storage Storage, serializer Serializer[T]) *Config[T] {
	return &Config[T]{
		storage:    storage,
		serializer: serializer,
		Data:       defaultData,
	}
}

// NewConfigMinimal creates config which uses YAML format
// and stores data in config.yaml file inside current working directory
func NewConfigMinimal[T any](defaultData T) *Config[T] {
	return NewConfig(
		defaultData,
		FileStorage{Filepath: "config.yaml"},
		YamlSerializer[T]{},
	)
}

// Load will load data from storage into Data.
// If storage reports that data is missing (e.g. file not found), Data will not be changed
func (c *Config[T]) Load() error {
	rawData, missing, err := c.storage.Load()

	if err != nil {
		return err
	}

	if missing {
		return nil
	}

	data, err := c.serializer.DeserializeData(rawData)

	if err != nil {
		return err
	}
	c.Data = *data
	return nil
}

func (c *Config[T]) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.serializeAndStore()
}

func (c *Config[T]) serializeAndStore() error {
	data, err := c.serializer.SerializeData(&c.Data)
	if err != nil {
		return err
	}
	return c.storage.Save(data)
}

// Transaction simplifies atomic config changes. It runs transaction, then Save.
// If error or panic occurs inside transaction or Save, config changes are rolled back.
func (c *Config[T]) Transaction(transaction func(data *T) error) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	backup, err := kamino.Clone(c.Data)

	if err != nil {
		return
	}

	panicked := true
	defer func() {
		if panicked || err != nil {
			c.Data = backup
		}
	}()

	err = transaction(&c.Data)

	if err != nil {
		return
	}

	err = c.serializeAndStore()

	panicked = false

	return
}

func panicToError(err *error) {
	r := recover()
	if r == nil {
		return
	}
	switch e := r.(type) {
	case error:
		*err = e
	case string:
		*err = errors.New(e)
	default:
		*err = fmt.Errorf("unknown panic: %v", e)
	}
}
