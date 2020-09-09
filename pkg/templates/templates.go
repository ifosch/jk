package templates

import (
	"fmt"
	"strings"
	"text/template"
)

// trim is a helper function for templates to remove surrounding
// newlines.
func trim(in string) string {
	return strings.TrimPrefix(
		strings.TrimSuffix(in, "\n"),
		"\n",
	)
}

// Describe ...
func Describe() (t *template.Template, err error) {
	tmpl := "`{{ .Raw.Name }}`: "
	tmpl += "{{ .Raw.Description }}\n"
	tmpl += "Parameters:\n"
	tmpl += "{{ with .Raw.Property }}"
	tmpl += "{{ range . }}"
	tmpl += "{{ with .ParameterDefinitions }}"
	tmpl += "{{ range . }}"
	tmpl += "- *{{ .Name }}* ({{ .Type }}): "
	tmpl += "{{ trim .Description }} "
	tmpl += "(Default: `{{ .DefaultParameterValue.Value }}`)\n"
	tmpl += "{{ end }}"
	tmpl += "{{ end }}"
	tmpl += "{{ end }}"
	tmpl += "{{ end }}"
	t, err = template.New("Job").Funcs(template.FuncMap{
		"trim": trim,
	}).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("Template Creation error %v", err)
	}
	return t, nil
}
