// Package config provides simple persistent storage for application configuration
package config

type Config[T any] struct {
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
	data, err := c.serializer.SerializeData(&c.Data)
	if err != nil {
		return err
	}
	return c.storage.Save(data)
}
