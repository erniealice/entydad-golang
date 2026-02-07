package entydad

import "context"

// DataSource provides technology-agnostic data access for views.
// Consumer apps satisfy this interface by wrapping their database adapter.
// espyna's DatabaseAdapter already matches this signature directly.
type DataSource interface {
	ListSimple(ctx context.Context, collection string) ([]map[string]any, error)
}
