package catalog

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"strings"
)

const indexTemplateBody = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <style>
    :root {
      color-scheme: light dark;
      font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      max-width: 960px;
      margin: 0 auto;
      padding: 2rem 1rem 4rem;
      line-height: 1.5;
      background-color: #f7f7f7;
    }

    body {
      margin: 0;
    }

    header {
      margin-bottom: 2rem;
    }

    h1 {
      font-size: clamp(2rem, 5vw, 3rem);
      margin-bottom: 0.25rem;
    }

    .utilities {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
      gap: 1rem;
    }

    .utility {
      border-radius: 0.75rem;
      background: rgba(255, 255, 255, 0.85);
      backdrop-filter: blur(8px);
      box-shadow: 0 0.5rem 1.5rem rgba(15, 15, 15, 0.1);
      padding: 1.5rem;
      transition: transform 160ms ease, box-shadow 160ms ease;
    }

    .utility:focus-within,
    .utility:hover {
      transform: translateY(-4px);
      box-shadow: 0 1rem 2.5rem rgba(15, 15, 15, 0.18);
    }

    .utility h2 {
      margin: 0 0 0.5rem;
      font-size: 1.25rem;
    }

    .utility p {
      margin: 0 0 1rem;
      color: rgba(20, 20, 20, 0.75);
    }

    .utility a {
      display: inline-flex;
      gap: 0.5rem;
      align-items: center;
      color: #1a4eb0;
      font-weight: 600;
      text-decoration: none;
    }

    .utility a:hover {
      text-decoration: underline;
    }
  </style>
</head>
<body>
  <header>
    <h1>{{ .Title }}</h1>
    <p>{{ .Intro }}</p>
  </header>
  <main>
    <section class="utilities">
      {{- range .Utilities }}
      <article class="utility">
        <h2>{{ .Name }}</h2>
        <p>{{ .Description }}</p>
        <a href="{{ .URL }}">Open utility →</a>
      </article>
      {{- end }}
    </section>
  </main>
</body>
</html>`

var indexTemplate = template.Must(template.New("index").Parse(indexTemplateBody))

// IndexData feeds the index page template.
type IndexData struct {
	Title      string
	Intro      string
	Utilities  []IndexUtility
	BaseDomain string
}

// IndexUtility represents a utility entry on the index page.
type IndexUtility struct {
	Name        string
	Description string
	URL         string
}

// RenderIndex builds the catalogue index markup with the provided utilities.
func RenderIndex(baseURL string, utils []Utility) ([]byte, error) {
	data := IndexData{
		Title: "Rafael Selbach — Utilities",
		Intro: "Shared tools deployed under utils.rselbach.com.",
	}

	data.Utilities = make([]IndexUtility, 0, len(utils))

	for _, util := range utils {
		data.Utilities = append(data.Utilities, IndexUtility{
			Name:        util.Name,
			Description: util.Description,
			URL:         linkFor(baseURL, util.Slug),
		})
	}

	var buf bytes.Buffer
	if err := indexTemplate.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	return buf.Bytes(), nil
}

func linkFor(baseURL, slug string) string {
	clean := strings.Trim(path.Clean("/"+slug), "/")

	if clean == "." || clean == "" {
		clean = ""
	}

	if baseURL == "" {
		if clean == "" {
			return "./"
		}

		return "./" + clean + "/"
	}

	base := strings.TrimSuffix(baseURL, "/")
	if clean == "" {
		return base + "/"
	}

	return base + "/" + clean + "/"
}
