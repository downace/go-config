# Go Config

Simple persistent storage for application configuration

### Example

```go
type AppConfig struct {
    Host string `yaml:"host"`
    Port uint16 `yaml:"port"`
}

// NewConfigMinimal uses YAML format
// and stores data in config.yaml file inside current working directory
conf := config.NewConfigMinimal(AppConfig{
    Host: "192.168.10.1",
    Port: 8888,
})

// Load config from file
conf.Load()
// Use config
fmt.Printf("Address is %s:%d\n", conf.Data.Host, conf.Data.Port)
// Change config
conf.Data.Port = 8080
// Save config to file
conf.Save()
```

You can customize serializer and storage:

```go
type AppConfig struct {
    AuthToken string `json:"authToken"`
}

conf := config.NewConfig(
    AppConfig{},
    config.FileStorage{Filepath: "my-config.json"},
    config.JsonSerializer[AppConfig]{},
)

// ...
```

Using `Transaction` method

```go
type AppConfig struct {
    Host string `yaml:"host"`
    Port uint16 `yaml:"port"`
}

conf := NewConfigMinimal(AppConfig{})

err := conf.Transaction(func(data *AppConfig) error {
    data.Host = "192.168.10.1"
    data.Port = 8888
    // Use panic or return non-nil error to rollback changes
    return nil
})

if err != nil {
    // ...
}

// Changes are saved automatically if there were no error nor panic
```
