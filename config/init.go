package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/friendsofgo/errors"
)

func CreateDefaultConfig(path, documentsDir, partialsDir, outputDir string, force bool) error {
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

	configTemplate := fmt.Sprintf(DefaultConfigTemplate, documentsDir, partialsDir, outputDir)

	// Write config file
	if err := os.WriteFile(path, []byte(configTemplate), 0600); err != nil {
		return errors.Wrap(err, "writing config file")
	}

	return nil
}
