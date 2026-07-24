package graph

import (
	"strings"

	"aether/model"
)

// FilterOptions specifies query parameters for filtering graph nodes.
type FilterOptions struct {
	Type       model.EntityType
	Language   string
	PathPrefix string
	NameSearch string
}

// QueryNodes filters nodes in the graph matching specified criteria.
func (e *Engine) QueryNodes(opts FilterOptions) []model.Node {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var matches []model.Node
	for _, node := range e.nodes {
		if opts.Type != "" && node.Type != opts.Type {
			continue
		}
		if opts.Language != "" && !strings.EqualFold(node.Language, opts.Language) {
			continue
		}
		if opts.PathPrefix != "" && !strings.HasPrefix(node.Path, opts.PathPrefix) {
			continue
		}
		if opts.NameSearch != "" && !strings.Contains(strings.ToLower(node.Name), strings.ToLower(opts.NameSearch)) {
			continue
		}
		matches = append(matches, node)
	}

	return matches
}
