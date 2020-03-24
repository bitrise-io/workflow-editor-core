package core

// YmlToJSONKeyTypeConversion ...
func YmlToJSONKeyTypeConversion(i interface{}) interface{} {
	switch x := i.(type) {
	case map[string]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k] = YmlToJSONKeyTypeConversion(v)
		}
		return m2
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = YmlToJSONKeyTypeConversion(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = YmlToJSONKeyTypeConversion(v)
		}
	}
	return i
}
