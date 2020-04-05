package intel

func Dedup(orig []string) []string {
	keys := make(map[string]struct{})
	list := make([]string, 0)
	for _, item := range orig {
		if _, ok := keys[item]; !ok {
			keys[item] = struct{}{}
			list = append(list, item)
		}
	}
	return list
}
