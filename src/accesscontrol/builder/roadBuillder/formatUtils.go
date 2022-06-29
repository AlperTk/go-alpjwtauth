package roadBuillder

import "strings"

func FormatEndpoint(endpoint string) string {
	res := strings.HasSuffix(endpoint, "/")
	if !res {
		endpoint += "/"
	}

	res = strings.HasPrefix(endpoint, "/")
	if !res {
		endpoint = "/" + endpoint
	}
	return endpoint
}
