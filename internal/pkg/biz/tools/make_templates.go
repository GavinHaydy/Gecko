package tools

import "Gecko/internal/pkg/dal/rao"

func MakeTemplates(functions ...func() rao.Template) []rao.Template {
	var result []rao.Template
	for _, f := range functions {
		result = append(result, f())
	}
	return result
}
