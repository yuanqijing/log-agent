package kube

// Map2Json
//  m := map[string]string{
//		"app": "redis-t",
//	} -> `{"app":"redis-t"}`
func Map2Json(m map[string]string) string {
	var s string
	for k, v := range m {
		s += `"` + k + `":"` + v + `",`
	}
	return `{` + s[:len(s)-1] + `}`
}

// Map2Str
//  m := map[string]string{
//		"app": "redis-t",
//	} -> "app=redis-t"
func Map2Str(m map[string]string) string {
	var s string
	for k, v := range m {
		s += k + "=" + v + ","
	}
	return s[:len(s)-1]
}
