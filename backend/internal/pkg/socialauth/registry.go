package socialauth

import "fmt"

type Registry struct {
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{
		providers: map[string]Provider{
			"google":   NewGoogleProvider(),
			"facebook": NewFacebookProvider(),
			"x":        NewXProvider(),
		},
	}
}

func (r *Registry) GetProvider(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return p, nil
}
