package intel

import (
	"fmt"
	"net/http"
)

func (c *Client) Header() http.Header {
	protoBase := c.Proto() + c.Base
	cookie := fmt.Sprintf(
		"csrftoken=%s; sessionid=%s",
		c.CSRF, c.SessID,
	)

	hdr := http.Header{}
	hdr.Set("User-Agent", c.UA)
	hdr.Set("Accept", "application/json")
	hdr.Set("Content-Type", "application/json")
	hdr.Set("X-CSRFToken", c.CSRF)
	hdr.Set("Referer", protoBase)
	hdr.Set("Cookie", cookie)

	return hdr
}
