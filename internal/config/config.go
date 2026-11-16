package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Repos []string `yaml:"repos"`
}

func filePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "statboard", "config.yaml"), nil
}

// Load reads the config file. A missing file means "no repos yet", not an error.
func Load() (*Config, error) {
	p, err := filePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	p, err := filePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

func (c *Config) Add(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return fmt.Errorf("cannot access %s: %w", abs, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", abs)
	}
	if _, err := os.Stat(filepath.Join(abs, ".git")); err != nil {
		return fmt.Errorf("%s is not a Git repository (no .git found)", abs)
	}
	for _, r := range c.Repos {
		if r == abs {
			return fmt.Errorf("%s is already tracked", abs)
		}
	}
	c.Repos = append(c.Repos, abs)
	return nil
}

func (c *Config) Remove(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	for i, r := range c.Repos {
		if r == abs {
			c.Repos = append(c.Repos[:i], c.Repos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%s is not tracked", abs)
}
