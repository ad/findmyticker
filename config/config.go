package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
)

// Config ...
type Config struct {
	Homeassistant struct {
		Token string `json:"token"`
		URL   string `json:"url"`
	} `json:"homeassistant"`
	Update struct {
		Period          int      `json:"period"`
		Ignore          []string `json:"ignore"`
		AllowItems      bool     `json:"allowItems"`
		AllowDevices    bool     `json:"allowDevices"`
		MinimalAccuracy int      `json:"minimalAccuracy"`
	} `json:"update"`
	FindMyApp struct {
		BringToFrontOnIdle bool `json:"bringToFrontOnIdle"`
		BringToFronDelay   int  `json:"bringToFrontDelay"`
		OpenOnStartup      bool `json:"openOnStartup"`
	} `json:"findMyApp"`
}

func InitConfig() (*Config, error) {
	config := &Config{}

	configFileName, errGetConfigPath := GetConfigPath()

	var initFromFile = false

	if errGetConfigPath == nil {
		if _, err := os.Stat(configFileName); err == nil {
			jsonFile, err := os.Open(configFileName)
			if err == nil {
				byteValue, _ := io.ReadAll(jsonFile)
				if err = json.Unmarshal(byteValue, &config); err != nil {
					return nil, fmt.Errorf("error on unmarshal config from file %s", err.Error())
				} else {
					initFromFile = true
				}
			}
		}
	}

	if !initFromFile {
		flag.StringVar(&config.Homeassistant.Token, "TOKEN", lookupEnvOrString("TOKEN", config.Homeassistant.Token), "homeassistant token")
		flag.StringVar(&config.Homeassistant.URL, "URL", lookupEnvOrString("URL", config.Homeassistant.URL), "homeassistant url")
		flag.Parse()
	}

	if config.Homeassistant.Token == "" {
		_ = OpenConfigEditor()

		return nil, fmt.Errorf("TOKEN env var not set")
	}

	if config.Homeassistant.URL == "" {
		_ = OpenConfigEditor()

		return nil, fmt.Errorf("URL env var not set")
	}

	if config.Update.Period <= 0 {
		config.Update.Period = 60
	}

	if config.FindMyApp.BringToFronDelay <= 0 {
		config.FindMyApp.BringToFronDelay = 60
	}

	if config.Update.MinimalAccuracy <= 0 {
		config.Update.MinimalAccuracy = 200
	}

	return config, nil
}

func lookupEnvOrString(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
}

func OpenConfigEditor() error {
	path, errGetConfigPath := GetConfigPath()
	if errGetConfigPath != nil {
		return errGetConfigPath
	}

	if _, err := os.Stat(path); err == nil {
		// path exists
	} else if errors.Is(err, os.ErrNotExist) {
		// path does *not* exist
		initialConfig := `{
	"homeassistant": {
		"url": "homeassistant url",
		"token": "homeassistant token"
	},
	"update": {
		"period": 60,
		"ignore": [],
		"allowItems": true,
		"allowDevices": true,
		"minimalAccuracy": 200
	},
	"findMyApp": {
		"openOnStartup": true,
		"bringToFrontOnIdle": true,
		"bringToFrontDelay": 60
	}
}`

		f, err := os.Create(path)
		if err != nil {
			return err
		}

		// Create a new writer.
		w := bufio.NewWriter(f)

		// Write a string to the file.
		_, _ = w.WriteString(initialConfig)

		// Flush.
		w.Flush()
	}

	return exec.Command(`open`, "-e", path).Run()
}

func GetConfigPath() (string, error) {
	u, err := user.Lookup(os.ExpandEnv("$USER"))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/findmytickerconfig.json", u.HomeDir), nil
}
