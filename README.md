# 📧 Envelopr - Live reload MJML email development

[![Go Report Card](https://goreportcard.com/badge/github.com/esdete2/envelopr)](https://goreportcard.com/report/github.com/esdete2/envelopr)
[![License](https://img.shields.io/github/license/esdete2/envelopr)](https://github.com/esdete2/envelopr/blob/main/LICENSE)

Envelopr is a live-reloading development tool for building MJML email templates with instant preview. Write your MJML templates and see changes in real-time, with support for Go template syntax for dynamic content.

## Features

- 🔄 Live preview with hot reload
- 🧩 Two-stage template processing (Go templates + MJML)
- 📂 Support for partials and nested directory structures
- 📱 Responsive preview with mobile/tablet/desktop views
- 🎨 Static template variables via YAML configuration

## Installation

### Via Install Script (Recommended)

```bash
curl -sSfL https://raw.githubusercontent.com/esdete2/envelopr/master/install.sh | sh
```

### Via Go Install

With Go 1.23 or higher:
```bash
go install github.com/esdete2/envelopr@latest
```

## Quick Start

1. First enter into your project
```bash
cd /path/to/your_project
```

2. Initialize a new project in current directory:
```bash
envelopr init
```

3. Create your first template in `documents/welcome.mjml`:
```html
<mjml>
  <mj-body>
    <mj-section>
      <mj-column>
        <mj-text>Hello World!</mj-text>
      </mj-column>
    </mj-section>
  </mj-body>
</mjml>
```

4. Start the development server:
```bash
envelopr watch
```

5. Open http://localhost:3600 in your browser

## Configuration

Envelopr uses a YAML configuration file (`envelopr.yaml`):

```yaml
paths:
  documents: documents   # MJML templates directory
  partials: partials     # Reusable components
  output: output         # Compiled HTML output

mjml:
  validationLevel: soft  # strict/soft/skip
  keepComments: false    # Preserve HTML comments
  beautify: true         # Pretty print HTML
  minify: false          # Minify output
  fonts:                 # Custom web fonts
    Roboto: https://fonts.googleapis.com/css?family=Roboto

template:
  variables:             # Global variables
    companyName: ACME Corp
    supportEmail: support@acme.com
    
  documents:             # Template-specific variables
    welcome:
      name: John Doe
      activationLink: https://example.com/activate
```

## Templates

Envelopr uses Go's template language with additional features for email template development.

### Go Templates

Basic Go template syntax is supported with variables from your config file:
```html
# With config:
template:
  variables:
    name: John
    showHeader: true
    items:
      - name: Product
        price: $99

# In template:
<!-- Variables -->
<mj-text>Hello {{ .name }}</mj-text>

<!-- Conditionals -->
{{ if .showHeader }}
  {{ template "header" . }}
{{ end}}

<!-- Loops -->
{{ range .items}}
  <mj-text>Item: {{ .name }} - {{ .price }}</mj-text>
{{ end}}
```

Template variables are defined in your `envelopr.yaml` configuration.

### Expression Preservation

Use `expression` (or its shorter alias `exp`) to preserve Go template expressions in the output HTML:

```html
<mj-button href="{{ expression ".profileUrl" }}">
  Visit Profile
</mj-button>

<!-- Same using shorter syntax -->
<mj-button href="{{ exp ".profileUrl" }}">
  Visit Profile
</mj-button>

<!-- Output HTML (both variants) -->
<a href="{{ .profileUrl }}">Visit Profile</a>
```

Here's another example preserving loop and conditional expressions:
```html

<mj-section>
    <mj-column>
        {{ exp "range .products" }}
            <mj-text>{{ exp ".name" }}</mj-text>
            {{ exp "if .onSale" }}
                <mj-text color="red">Sale!</mj-text>
            {{ exp "end" }}
        {{ exp "end" }}
    </mj-column>
</mj-section>

<!-- Output HTML -->
<div>
    {{ range .products }}
        <p>{{ .name }}</p>
        {{ if .onSale }}
            <p style="color: red">Sale!</p>
        {{ end }}
    {{ end }}
</div>
```

This is useful when the final HTML needs to be processed by another template engine.

### Partials and Layouts

Create reusable components in the `partials` directory:

`partials/layout.mjml`:
```html

<mjml>
    <mj-head>
        <mj-title>{{ .title }}</mj-title>
    </mj-head>
    <mj-body>
        {{ template "header" . }}
        {{ template "content" . }}  <!-- Main content injection -->
        {{ template "footer" . }}
    </mj-body>
</mjml>
```

`partials/button.mjml`:
```html

<mj-button
    href="{{.url}}"
    background-color="#2563eb"
    border-radius="6px"
>
    {{ .label }}
</mj-button>
```

Use them in your templates:
```html
{{ template "layout" . }}

{{ define "content" }}
<mj-section>
    <mj-column>
        <mj-text>Welcome {{ .name }}!</mj-text>

        {{ template "button" dict
            "url" (exp ".campaignUrl")
            "label" "Shop Now"
        }}
    </mj-column>
</mj-section>
{{ end }}
```

## Commands

```sh
# Initialize new project
envelopr init

# Compile templates
envelopr build

# Start development server
envelopr watch
```

### Command Options

```bash
# Build with custom config file
envelopr build -c custom-config.yaml

# Watch with custom host and port
envelopr watch --host 127.0.0.1 --port 8080

# Initialize without interactive prompts
envelopr init -y

# Set log verbosity (1=error to 5=trace)
envelopr -v 4 watch
```
For a full list of commands and options, run `envelopr --help`.

## Contributing

Pull requests are welcome! Feel free to:

- Report bugs
- Suggest new features
- Submit pull requests

## License

Apache-2.0 license - see [LICENSE](LICENSE) for details.