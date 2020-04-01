package intel

import "fmt"

type Portal struct {
	Team  string
	Lat   float64
	Lng   float64
	Image string
	Name  string
}

func (c *Client) GetPortal(guid string) (*Portal, error) {
	p := &Portal{}

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
			return nil, err
		}

		keys := make([]string, 0)
		for k := range res {
			keys = append(keys, k)
		}

		var ok bool
		result, ok = res["result"].([]interface{})
		if ok {
			break
		}

		tries++
		if tries > 10 {
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
		return nil, fmt.Errorf("assert image")
	}
	p.Image = image

	name, ok := result[8].(string)
	if !ok {
		return nil, fmt.Errorf("assert name")
	}
	p.Name = name

	return p, nil
}
