package avtonet

import "net/url"

type Transformer struct{}

func (t *Transformer) Transform(uo *url.URL) (*url.URL, error) {
	ut := *uo

	ut.Scheme = "https"
	ut.Host = "www.avto.net"

	return &ut, nil
}
