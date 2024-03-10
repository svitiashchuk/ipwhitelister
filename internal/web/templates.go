package web

import (
	"embed"
	"strings"
	"sync"
	"text/template"
)

//go:embed assets/templates/*
var templateFS embed.FS

// TemplateCache stores the compiled templates to avoid reparsing text templates on each render call.
var TemplateCache = make(map[string]*template.Template)
var cacheMutex = &sync.Mutex{}

// Template parses the templates from the embedded file system and caches the result.
func Template(tmpl ...string) (*template.Template, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	for i, t := range tmpl {
		tmpl[i] = "assets/templates/" + t
	}

	// Create a unique key for the set of templates
	tmplKey := strings.Join(tmpl, ",")

	t, ok := TemplateCache[tmplKey]
	if !ok {
		var err error

		t, err = template.ParseFS(templateFS, tmpl...)
		if err != nil {
			return nil, err
		}
		TemplateCache[tmplKey] = t
	}

	return t, nil
}
