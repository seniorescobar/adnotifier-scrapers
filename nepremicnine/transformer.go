package nepremicnine

import "net/url"

type Transformer struct{}

func (t *Transformer) Transform(uo *url.URL) (*url.URL, error) {
	ut := *uo

	ut.Scheme = "https"
	ut.Host = "www.nepremicnine.net"

	qVals := ut.Query()
	qVals.Set("s", "16")
	ut.RawQuery = qVals.Encode()

	return &ut, nil
}
