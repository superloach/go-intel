package intel

import "fmt"

func (c *Client) PortalIDs(tileKeys []string) ([]string, error) {
	ids := make([]string, 0)

	res, err := c.jsonPost(
		"/r/getEntities",
		obj{
			"tileKeys": tileKeys,
			"v":        c.Version,
		},
	)
	if err != nil {
		return nil, err
	}

	result, ok := res["result"].(obj)
	if !ok {
		return nil, fmt.Errorf("assert portal ids result")
	}

	_map, ok := result["map"].(obj)
	if !ok {
		return nil, fmt.Errorf("assert _map")
	}

	for _, ichunk := range _map {
		chunk, ok := ichunk.(obj)
		if !ok {
			return nil, fmt.Errorf("assert chunk")
		}

		ents, ok := chunk["gameEntities"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("assert ents")
		}

		for _, ient := range ents {
			ent, ok := ient.([]interface{})
			if !ok {
				return nil, fmt.Errorf("assert ent")
			}

			if len(ent) < 3 {
				return nil, fmt.Errorf("short ent")
			}

			parts, ok := ent[2].([]interface{})
			if !ok {
				return nil, fmt.Errorf("assert parts")
			}

			if len(parts) < 1 {
				return nil, fmt.Errorf("short parts")
			}

			typ, ok := parts[0].(string)
			if !ok {
				return nil, fmt.Errorf("assert typ")
			}

			switch typ {
			case "p": // portal
				id, ok := ent[0].(string)
				if !ok {
					return nil, fmt.Errorf("assert portal id")
				}

				ids = append(ids, id)
			case "e": // line
				if len(parts) < 6 {
					return nil, fmt.Errorf("short line parts")
				}

				id1, ok := parts[2].(string)
				if !ok {
					return nil, fmt.Errorf("assert line id1")
				}

				ids = append(ids, id1)

				id2, ok := parts[5].(string)
				if !ok {
					return nil, fmt.Errorf("assert line id1")
				}

				ids = append(ids, id2)
			case "r": // poly
				if len(parts) < 3 {
					return nil, fmt.Errorf("short poly parts")
				}

				portals, ok := parts[2].([]interface{})
				if !ok {
					return nil, fmt.Errorf("assert poly portals")
				}

				for _, iportal := range portals {
					portal, ok := iportal.([]interface{})
					if !ok {
						return nil, fmt.Errorf("assert poly portal")
					}

					if len(portal) < 1 {
						return nil, fmt.Errorf("short poly portal")
					}

					id, ok := portal[0].(string)
					if !ok {
						return nil, fmt.Errorf("assert poly id")
					}

					ids = append(ids, id)
				}
			default:
				return nil, fmt.Errorf("unknown typ " + typ)
			}
		}
	}

	return ids, nil
}
