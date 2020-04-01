package intel

import "net/http"

type Client struct {
	Client *http.Client

	Base    string
	Secure  bool
	UA      string
	Version string
	CSRF    string
	SessID  string
}

func NewClient() (*Client, error) {
	c := &Client{}

	c.Client = &http.Client{}

	return c, nil
}

func (c *Client) Proto() string {
	if c.Secure {
		return "https://"
	}
	return "http://"
}
