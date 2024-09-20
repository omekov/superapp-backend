package external

import "net/http"

type Client struct {
	kgdURL      string
	exchangeURL string
	httpClient  http.Client
}

func NewExternatClient(kgdURL, exchangeURL string) Client {
	return Client{
		kgdURL:      kgdURL,
		exchangeURL: exchangeURL,
		httpClient:  *http.DefaultClient,
	}
}
