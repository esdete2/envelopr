package handler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/esdete2/mjml-dev/handler"
)

func TestFileLoader_LoadDocuments(t *testing.T) {
	t.Run("empty directory path", func(t *testing.T) {
		r := require.New(t)

		loader := handler.NewFileLoader("", "")
		docs, err := loader.LoadDocuments()
		r.NoError(err)
		r.Nil(docs)
	})

	t.Run("non-existent directory", func(t *testing.T) {
		r := require.New(t)

		loader := handler.NewFileLoader("/does/not/exist", "")
		_, err := loader.LoadDocuments()
		r.Error(err)
		r.Contains(err.Error(), "directory does not exist")
	})

	t.Run("empty but existing directory", func(t *testing.T) {
		r := require.New(t)

		// Create empty temp directory
		tmpDir, err := os.MkdirTemp("", "mjml-dev-test-empty")
		r.NoError(err)
		defer os.RemoveAll(tmpDir)

		loader := handler.NewFileLoader(tmpDir, "")
		docs, err := loader.LoadDocuments()
		r.NoError(err)
		r.Empty(docs)
	})

	t.Run("directory with files", func(t *testing.T) {
		r := require.New(t)

		tmpDir, err := os.MkdirTemp("", "mjml-dev-test")
		r.NoError(err)
		defer os.RemoveAll(tmpDir)

		files := map[string]string{
			"welcome.mjml":              "<mjml>1</mjml>",
			"marketing/newsletter.mjml": "<mjml>2</mjml>",
			"not-a-template.txt":        "ignored",
		}

		for path, content := range files {
			fullPath := filepath.Join(tmpDir, path)
			err := os.MkdirAll(filepath.Dir(fullPath), 0755)
			r.NoError(err)
			err = os.WriteFile(fullPath, []byte(content), 0644)
			r.NoError(err)
		}

		loader := handler.NewFileLoader(tmpDir, "")
		docs, err := loader.LoadDocuments()
		r.NoError(err)
		r.Len(docs, 2) // Only .mjml files

		// Create map for easier lookup during testing
		docMap := make(map[string]string)
		for _, doc := range docs {
			docMap[doc.Name] = doc.Content
		}

		r.Equal("<mjml>1</mjml>", docMap["welcome"])
		r.Equal("<mjml>2</mjml>", docMap["marketing/newsletter"])
	})
}
func TestFileLoader_LoadPartials(t *testing.T) {
	t.Run("empty directory path", func(t *testing.T) {
		r := require.New(t)

		loader := handler.NewFileLoader("", "")
		partials, err := loader.LoadPartials()
		r.NoError(err)
		r.Nil(partials)
	})

	t.Run("non-existent directory", func(t *testing.T) {
		r := require.New(t)

		loader := handler.NewFileLoader("", "/does/not/exist")
		_, err := loader.LoadPartials()
		r.Error(err)
		r.Contains(err.Error(), "directory does not exist")
	})

	t.Run("empty but existing directory", func(t *testing.T) {
		r := require.New(t)

		tmpDir, err := os.MkdirTemp("", "mjml-dev-test-empty")
		r.NoError(err)
		defer os.RemoveAll(tmpDir)

		loader := handler.NewFileLoader("", tmpDir)
		partials, err := loader.LoadPartials()
		r.NoError(err)
		r.Empty(partials)
	})

	t.Run("directory with files", func(t *testing.T) {
		r := require.New(t)

		tmpDir, err := os.MkdirTemp("", "mjml-dev-test")
		r.NoError(err)
		defer os.RemoveAll(tmpDir)

		files := map[string]string{
			"header.mjml":        "<mj-section>header</mj-section>",
			"shared/footer.mjml": "<mj-section>footer</mj-section>",
			"not-a-partial.txt":  "ignored",
		}

		for path, content := range files {
			fullPath := filepath.Join(tmpDir, path)
			err := os.MkdirAll(filepath.Dir(fullPath), 0755)
			r.NoError(err)
			err = os.WriteFile(fullPath, []byte(content), 0644)
			r.NoError(err)
		}

		loader := handler.NewFileLoader("", tmpDir)
		partials, err := loader.LoadPartials()
		r.NoError(err)
		r.Len(partials, 2) // Only .mjml files

		partialsMap := make(map[string]string)
		for _, partial := range partials {
			partialsMap[partial.Name] = partial.Content
		}

		r.Equal("<mj-section>header</mj-section>", partialsMap["header"])
		r.Equal("<mj-section>footer</mj-section>", partialsMap["shared/footer"])
	})
}
