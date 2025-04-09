package config

import "os"

type Storage interface {
	Load() (data []byte, missing bool, err error)
	Save(data []byte) error
}

type FileStorage struct {
	Filepath string
}

func (p FileStorage) Load() ([]byte, bool, error) {
	rawData, err := os.ReadFile(p.Filepath)

	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return nil, true, err
	}

	return rawData, false, nil
}

func (p FileStorage) Save(data []byte) error {
	return os.WriteFile(p.Filepath, data, 0664)
}
