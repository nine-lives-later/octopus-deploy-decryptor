package html

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/projectExport"
	"html/template"
	"log"
	"strings"
)

var (
	//go:embed html.tmpl
	htmlContent string
)

func RenderHTML(entities projectExport.EntityMap, key []byte, sensitiveOnly bool) (string, error) {
	// render the template
	t, err := template.New("html").Parse(htmlContent)
	if err != nil {
		return "", fmt.Errorf("failed to create HTML template: %w", err)
	}

	vm := buildViewModel(entities, key, sensitiveOnly)

	var buf bytes.Buffer
	err = t.Execute(&buf, vm)
	if err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return buf.String(), nil
}

func buildViewModel(entities projectExport.EntityMap, key []byte, sensitiveOnly bool) map[string]any {
	root := map[string]any{
		"Projects": buildViewModel_Projects(entities, key, sensitiveOnly),
	}

	return root
}

func buildViewModel_Projects(entities projectExport.EntityMap, key []byte, sensitiveOnly bool) []map[string]any {
	projects := entities.Projects()
	ret := make([]map[string]any, 0, len(projects))

	for _, p := range projects {
		ret = append(ret, map[string]any{
			"ID":                  p.EntityID(),
			"Name":                p.EntityName(),
			"SpaceID":             p.EntitySpaceID(),
			"ProjectVariableSet":  buildViewModel_ProjectVariableSet(entities, p, key, sensitiveOnly),
			"LibraryVariableSets": buildViewModel_ProjectLibraryVariableSets(entities, p, key, sensitiveOnly),
		})
	}

	return ret
}

func buildViewModel_ProjectVariableSet(entities projectExport.EntityMap, p *projectExport.Project, key []byte, sensitiveOnly bool) map[string]any {
	s := entities.VariableSetByOwner(p.EntityID())

	variables := make([]any, 0, len(s.Variables))

	for _, v := range s.Variables {
		if sensitiveOnly && v.Type != "Sensitive" {
			continue
		}

		variables = append(variables, map[string]any{
			"Name":  v.Name,
			"Value": decryptedValue(v, key),
			"Scope": buildViewModel_VariableScope(entities, &v.Scope),
		})
	}

	return map[string]any{
		"ID":        s.EntityID(),
		"Name":      s.EntityName(),
		"SpaceID":   s.EntitySpaceID(),
		"Variables": variables,
	}
}

func buildViewModel_ProjectLibraryVariableSets(entities projectExport.EntityMap, p *projectExport.Project, key []byte, sensitiveOnly bool) []map[string]any {
	ret := make([]map[string]any, 0, len(p.IncludedLibraryVariableSetIds))

	for _, id := range p.IncludedLibraryVariableSetIds {
		ls := entities[id].(*projectExport.LibraryVariableSet)
		if ls == nil {
			log.Printf("Missing library variable set %s in project %s", id, p.Id)
			continue
		}

		s := entities[ls.VariableSetId].(*projectExport.VariableSet)
		if s == nil {
			log.Printf("Missing variable set %s for library variable set %s", ls.VariableSetId, id)
			continue
		}

		variables := make([]any, 0, len(s.Variables))

		for _, v := range s.Variables {
			if sensitiveOnly && v.Type != "Sensitive" {
				continue
			}

			variables = append(variables, map[string]any{
				"Name":  v.Name,
				"Value": decryptedValue(v, key),
				"Scope": buildViewModel_VariableScope(entities, &v.Scope),
			})
		}

		ret = append(ret, map[string]any{
			"ID":        s.EntityID(),
			"Name":      s.EntityName(),
			"SpaceID":   s.EntitySpaceID(),
			"Variables": variables,
		})
	}

	return ret
}

func buildViewModel_VariableScope(entities projectExport.EntityMap, scope *projectExport.VariableScope) string {
	if scope == nil {
		return ""
	}

	var ret strings.Builder

	if len(scope.EnvironmentIDs) > 0 {
		//ret.WriteString("Environment: ")

		for i, id := range scope.EnvironmentIDs {
			if i > 0 {
				ret.WriteString(", ")
			}

			name := id
			if e := entities[id].(*projectExport.Environment); e != nil {
				name = e.EntityName()
			}

			ret.WriteString(name)
		}
	}

	return ret.String()
}
