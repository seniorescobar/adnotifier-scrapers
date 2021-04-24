package avtonet

import "net/url"

// Transformer represents an AvtoNet Transformer.
type Transformer struct{}

// NewTransformer creates a new instance of AvtoNet Transformer.
func NewTransformer() *Transformer {
	return &Transformer{}
}

// Transform transforms the url into a normalized form.
func (t *Transformer) Transform(uo *url.URL) (*url.URL, error) {
	ut := *uo

	ut.Scheme = "https"
	ut.Host = "www.avto.net"

	return &ut, nil
}
