package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

func (c *Config) OpenTime() (time.Time, error) {
	t, err := time.Parse("15:04:05", c.OpenAt)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid OpenAt format: %w", err)
	}
	return t, nil
}

func (c *Config) CloseTime() (time.Time, error) {
	open, err := c.OpenTime()
	if err != nil {
		return time.Time{}, err
	}
	return open.Add(time.Duration(c.Duration) * time.Hour), nil
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Floors <= 0 {
		return fmt.Errorf("floors must be > 0")
	}
	if c.Monsters <= 0 {
		return fmt.Errorf("monsters must be > 0")
	}
	if c.Duration <= 0 {
		return fmt.Errorf("duration must be > 0")
	}
	if _, err := c.OpenTime(); err != nil {
		return err
	}
	return nil
}
