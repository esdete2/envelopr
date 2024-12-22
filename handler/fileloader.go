package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/friendsofgo/errors"

	"github.com/esdete2/mjml-dev/template"
)

type FileLoader struct {
	documentsPath string
	partialsPath  string
}

func NewFileLoader(documentsPath, partialsPath string) *FileLoader {
	return &FileLoader{
		documentsPath: documentsPath,
		partialsPath:  partialsPath,
	}
}

func (l *FileLoader) LoadDocuments() ([]template.Template, error) {
	if l.documentsPath == "" {
		return nil, nil
	}

	return l.loadTemplates(l.documentsPath)
}

func (l *FileLoader) LoadDocument(name string) ([]template.Template, error) {
	if l.documentsPath == "" {
		return nil, nil
	}

	// Handle both with and without .mjml extension
	if !strings.HasSuffix(name, ".mjml") {
		name = name + ".mjml"
	}

	// Build the full path
	fullPath := filepath.Join(l.documentsPath, name)

	// Check if file exists
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template does not exist: %s", name)
		}
		return nil, errors.Wrap(err, "checking template file")
	}

	// Read the file
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, errors.Wrap(err, "reading template file")
	}

	// Create template
	name = strings.TrimSuffix(filepath.ToSlash(name), ".mjml")
	return []template.Template{{
		Name:    name,
		Content: string(content),
	}}, nil
}

func (l *FileLoader) LoadPartials() ([]template.Template, error) {
	if l.partialsPath == "" {
		return nil, nil
	}

	return l.loadTemplates(l.partialsPath)
}

func (l *FileLoader) loadTemplates(dir string) ([]template.Template, error) {
	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory does not exist: %s", dir)
		}
		return nil, errors.Wrap(err, "checking directory")
	}

	templates := make([]template.Template, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "walking directory")
		}
		if info.IsDir() || filepath.Ext(path) != ".mjml" {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return errors.Wrap(err, "getting relative path")
		}

		name := strings.TrimSuffix(filepath.ToSlash(relPath), ".mjml")

		content, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, "reading file")
		}

		templates = append(templates, template.Template{
			Name:    name,
			Content: string(content),
		})
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "walking directory for templates")
	}

	return templates, nil
}
