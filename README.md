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