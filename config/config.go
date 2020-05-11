package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
)

type server struct {
	Port string `yaml:"port"`
	ReadTimeout uint64 `yaml:"readTimeout"`
	WriteTimeout uint64 `yaml:"writeTimeout"`
}

type EmailCredentials struct {
	Address string `yaml:"address"`
	Password string `yaml:"password"`
}

type Config struct {
	Original string
	Server server `yaml:"server"`
	EmailCredentials EmailCredentials `yaml:"emailCredentials"`
	RecaptchaSecretKey string `yaml:"recaptchaSecretKey"`
}

var (
	DefaultConfig = Config{
		Original: "",
		Server: server{
			Port: "8080",
			ReadTimeout: 5000,
			WriteTimeout: 5000,
		},
		EmailCredentials: EmailCredentials{
			Address: "",
			Password: "",
		},
		RecaptchaSecretKey: "",
	}
)

// this is used for production, in which Heroku forces config vars through environment
func LoadEnv() (*Config, error) {
	cfg := &Config{}
	port := os.Getenv("PORT")
	cfg.Server.Port = port; if cfg.Server.Port == "" {
		return nil, errors.New("PORT missing in environment")
	}

	readTimeout := os.Getenv("SERVER_READ_TIMEOUT")
	if readTimeout == "" {
		return nil, errors.New("SERVER_READ_TIMEOUT missing in environment")
	}
	var err error
	cfg.Server.ReadTimeout, err = strconv.ParseUint(readTimeout, 10, 16)
	if err != nil {
		return nil, err
	}

	writeTimeout := os.Getenv("SERVER_WRITE_TIMEOUT")
	if writeTimeout == "" {
		return nil, errors.New("SERVER_WRITE_TIMEOUT missing in environment")
	}
	cfg.Server.WriteTimeout, err = strconv.ParseUint(writeTimeout, 10, 16)
	if err != nil {
		return nil, err
	}

	emailAddress := os.Getenv("EMAIL_ADDRESS")
	if emailAddress == "" {
		return nil, errors.New("EMAIL_ADDRESS missing in environment")
	}
	cfg.EmailCredentials.Address = emailAddress

	emailPassword := os.Getenv("EMAIL_PASSWORD")
	if emailPassword == "" {
		return nil, errors.New("EMAIL_PASSWORD missing in environment")
	}
	cfg.EmailCredentials.Password = emailPassword

	recaptchaSecretKey := os.Getenv("RECAPTCHA_SECRET_KEY")
	if recaptchaSecretKey == "" {
		return nil, errors.New("RECAPTCHA_SECRET_KEY missing in environment")
	}
	cfg.RecaptchaSecretKey = recaptchaSecretKey

	return cfg, nil
}

// Load parses the YAML input s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}
	// if the entire config body is empty the UnmarshalYAML method is
	// never called. We thus have to set the DefaultConfig at the entry
	// point as well.
	*cfg = DefaultConfig

	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	cfg.Original = s
	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
