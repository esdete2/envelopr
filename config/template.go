package config

const DefaultConfigTemplate = `# Directory paths for templates, partials and output
paths:
  # Main directory containing your MJML templates
  documents: documents
  # Directory containing partial templates that can be included
  partials: partials
  # Output directory for compiled HTML files
  output: output

# MJML compilation settings
mjml:
  # Validation level: "strict", "soft", or "skip"
  # - strict: Validates and fails on any error
  # - soft: Shows warnings but continues
  # - skip: Skips validation entirely
  validationLevel: soft
  # Keep comments in output HTML
  keepComments: false
  # Beautify the output HTML
  beautify: true
  # Minify the output HTML
  minify: true
  # Custom fonts to include
  fonts:
    # Roboto: https://fonts.googleapis.com/css?family=Roboto

# Template processing settings
template:
  # By default, mjml url-encodes the value of href attributes.
  # This preserves template expressions in href attributes, like href="{{ expression .url }}".
  # Useful for two-stage template processing.
  preserveHrefExpressions: true

  # Global static variables available to all templates
  variables:
    # companyName: ACME Corp

  # Per-document variables
  documents:
    # Static variables for a template named newsletter.mjml
    # newsletter:
      # shopUrl: https://example.shop
`
