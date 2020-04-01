package intel

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type obj = map[string]interface{}

func (c *Client) jsonPost(endpoint string, data obj) (obj, error) {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(data)
	if err != nil {
		return nil, err
	}

	url := c.Proto() + c.Base + endpoint

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header = c.Header()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	res := obj{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
