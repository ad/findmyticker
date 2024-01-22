package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
)

// Config ...
type Config struct {
	Token        string   `json:"token"`
	URL          string   `json:"url"`
	Period       int      `json:"period"`
	Ignore       []string `json:"ignore"`
	AllowItems   bool     `json:"allowItems"`
	AllowDevices bool     `json:"allowDevices"`
}

func InitConfig() (*Config, error) {
	config := &Config{}

	u, err := user.Lookup(os.ExpandEnv("$USER"))
	if err != nil {
		return nil, fmt.Errorf("user lookup: %+v", err)
	}

	configFileName := fmt.Sprintf("%s/findmytickerconfig.json", u.HomeDir)

	var initFromFile = false

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

	if !initFromFile {
		flag.StringVar(&config.Token, "TOKEN", lookupEnvOrString("TOKEN", config.Token), "homeassistant token")
		flag.StringVar(&config.URL, "URL", lookupEnvOrString("URL", config.URL), "homeassistant url")
		flag.Parse()
	}

	if config.Token == "" {
		_ = openConfigEditor(configFileName)

		return nil, fmt.Errorf("TOKEN env var not set")
	}

	if config.URL == "" {
		_ = openConfigEditor(configFileName)

		return nil, fmt.Errorf("URL env var not set")
	}

	if config.Period == 0 {
		config.Period = 60
	}

	return config, nil
}

func lookupEnvOrString(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
}

func openConfigEditor(path string) error {
	if _, err := os.Stat(path); err == nil {
		// path exists
	} else if errors.Is(err, os.ErrNotExist) {
		// path does *not* exist
		initialConfig := `{
	"url": "homeassistant url",
	"token": "homeassistant token",
	"period": 60,
	"ignore": [],
	"allowItems": true,
	"allowDevices": true
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

	cmd := exec.Command(`open`, "-t", path)
	stderr, err := cmd.StderrPipe()
	log.SetOutput(os.Stderr)

	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := io.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return nil
}
