package utils

import "github.com/google/go-querystring/query"

func BuildUrlWithQueryString(url string, params interface{}) string {
	queryString := getQueryString(params)

	if queryString != "" {
		url += "?" + queryString
	}

	return url
}

func getQueryString(params interface{}) string {
	v, err := query.Values(params)

	if err != nil {
		return ""
	}

	return v.Encode()
}
