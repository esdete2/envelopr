package web

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/networkteam/slogutils"

	"github.com/esdete2/mjml-dev/web/views"
)

type IndexData struct {
	Templates []TemplateInfo
}

type TemplateInfo struct {
	Name string
	Path string
}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tree, err := s.listTemplates()
		if err != nil {
			slog.Error("failed to list templates", slogutils.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Debug("rendering index view")

		err = views.IndexView(tree).Render(r.Context(), w)
		if err != nil {
			slog.Error("failed to render index view", slogutils.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templatePath := chi.URLParam(r, "*")

		_, err := os.ReadFile(filepath.Join(s.options.Output, templatePath))
		if err != nil {
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}

		tmpl := views.TemplateContent{
			Path: templatePath,
			Name: strings.TrimSuffix(filepath.Base(templatePath), ".html"),
		}

		err = views.TemplateView(tmpl).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleRawTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templatePath := chi.URLParam(r, "*")

		content, err := os.ReadFile(filepath.Join(s.options.Output, templatePath))
		if err != nil {
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}

		// Set proper content type for HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listTemplates() (views.TreeNode, error) {
	root := &views.TreeNode{
		Name:     "/",
		IsDir:    true,
		Children: make([]*views.TreeNode, 0),
	}

	err := filepath.Walk(s.options.Output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == s.options.Output {
			return nil
		}

		// Only process directories and HTML files
		if !info.IsDir() && filepath.Ext(path) != ".html" {
			return nil
		}

		relPath, err := filepath.Rel(s.options.Output, path)
		if err != nil {
			return err
		}

		// Split path into components
		parts := strings.Split(filepath.ToSlash(relPath), "/")

		// Traverse or create nodes
		current := root
		for i, part := range parts {
			isLast := i == len(parts)-1
			var child *views.TreeNode

			// Find or create the child node
			for _, c := range current.Children {
				if c.Name == part {
					child = c
					break
				}
			}

			if child == nil {
				child = &views.TreeNode{
					Name:     strings.TrimSuffix(part, ".html"),
					IsDir:    info.IsDir() || !isLast,
					Children: make([]*views.TreeNode, 0),
					Path:     relPath,
				}
				current.Children = append(current.Children, child)
			}
			current = child
		}

		return nil
	})

	if err != nil {
		return views.TreeNode{}, err
	}

	// Convert root back to value type (if needed)
	return *root, nil
}
