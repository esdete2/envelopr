// config/init.go
package config

import (
	"os"
	"path/filepath"

	"github.com/friendsofgo/errors"
)

// CreateConfig writes the default config template to the specified path
func CreateDefaultConfig(path string, force bool) error {
	// Check if file exists
	if _, err := os.Stat(path); err == nil && !force {
		return errors.New("config file already exists, use --force to overwrite")
	}

	// Create directory if needed
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrap(err, "creating config directory")
		}
	}

	// Write config file
	if err := os.WriteFile(path, []byte(DefaultConfigTemplate), 0600); err != nil {
		return errors.Wrap(err, "writing config file")
	}

	return nil
}
