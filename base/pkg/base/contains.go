package base

func Contains(arr []string, v string) bool {
	for _, k := range arr {
		if k == v {
			return true
		}
	}
	return false
}

func MapContainsAny(p map[string]interface{}, keys ...string) bool {
	for _, k := range keys {
		if _, ok := p[k]; ok {
			return true
		}
	}
	return false
}

func MapContainsAll(p map[string]interface{}, keys ...string) bool {
	for _, k := range keys {
		if _, ok := p[k]; !ok {
			return false
		}
	}
	return true
}
