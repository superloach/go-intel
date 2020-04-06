package intel

import "fmt"

type Portal struct {
	ID    string  `json:"id" rethinkdb:"id"`
	Team  string  `json:"ingr" rethinkdb:"ingr"`
	Lat   float64 `json:"lat" rethinkdb:"lat"`
	Lng   float64 `json:"lng" rethinkdb:"lng"`
	Image string  `json:"image" rethinkdb:"image"`
	Name  string  `json:"name" rethinkdb:"name"`
}

func (c *Client) GetPortal(guid string) (*Portal, error) {
	p := &Portal{}

	p.ID = guid

	var result []interface{}
	tries := 0
	for {
		res, err := c.jsonPost(
			"/r/getPortalDetails",
			obj{
				"guid": guid,
				"v":    c.Version,
			},
		)
		if err != nil {
			continue
		}

		var ok bool
		result, ok = res["result"].([]interface{})
		if ok {
			break
		}

		tries++
		if c.MaxTries > 0 && tries > c.MaxTries {
			return nil, fmt.Errorf("max tries portal details")
		}
	}

	if len(result) < 9 {
		return nil, fmt.Errorf("short result")
	}

	team, ok := result[1].(string)
	if !ok {
		return nil, fmt.Errorf("assert team")
	}
	p.Team = team

	lat, ok := result[2].(float64)
	if !ok {
		return nil, fmt.Errorf("assert lat")
	}
	p.Lat = lat / 1000000

	lng, ok := result[3].(float64)
	if !ok {
		return nil, fmt.Errorf("assert lng")
	}
	p.Lng = lng / 1000000

	image, ok := result[7].(string)
	if !ok {
		if result[7] == nil {
			image = ""
		} else {
			return nil, fmt.Errorf("assert image")
		}
	}
	p.Image = image

	name, ok := result[8].(string)
	if !ok {
		return nil, fmt.Errorf("assert name")
	}
	p.Name = name

	return p, nil
}
