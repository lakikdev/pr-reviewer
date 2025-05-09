package utils

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RequestSource string

const (
	// prs-sc-lp = 1a32aa551a32aa4af41a44d1a32aa43cd0e1a32a
	// prs-sc-tr = 18cc134318cc17c14118cc178f18cc178773c318
	// prs-sc-rp = c5729c572ab12ac57096c5720f8e82f196c5720f
	CMSRequestSource     RequestSource = "CMS"
	CMSRequestSourceSHA1 string        = "55af444dd0e343c141f773c39ab12a709682f196"

	// prs-sc-lp = 81381afc81381ab4a68195f81381ab0302a81381
	// prs-sc-tr = e123d8a5e123dd3a12e123dd19e123dd1dcb9ae1
	// prs-sc-rp = e90ace90ae91e9e91ee1e90a168be7feb9e90a16
	AppRequestSource     RequestSource = "App"
	AppRequestSourceSHA1 string        = "fc4a695f02a8a53a129dcb9ace91e91ee1e7feb9"

	// prs-sc-lp = 5b2917e65b2917e7b85bb2b5b2917eb55235b291
	// prs-sc-tr = 00c2670a00c269491a00c269bb00c269bebc4300
	// prs-sc-rp = 289642896f25be2831722896733efbce60289673
	PostmanRequestSource     RequestSource = "Postman"
	PostmanRequestSourceSHA1 string        = "e67b8b2b52370a491abebc434f25be3172fbce60"

	//Headers not much or not exist in request
	OtherRequestSource RequestSource = "Other"

	NotSetRequestSource RequestSource = "NotSet"
)

func ParseRequestSource(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context ) error {
		req, _ := parseSourceHeaders(c.Request())
		c.SetRequest(req)

		return next(c)
	}
}

type requestSourceContextKeyType struct{}

var requestSourceContextKey requestSourceContextKeyType

func parseSourceHeaders(r *http.Request) (*http.Request, error) {
	lpHeader := r.Header.Get("prs-sc-lp")
	trHeader := r.Header.Get("prs-sc-tr")
	rpHeader := r.Header.Get("prs-sc-rp")

	var combinedHash string
	if len(lpHeader) > 0 && len(trHeader) > 0 && len(rpHeader) > 0 {
		//every hash was blended bt using specific pattern
		//DO NOT CHANGE PATTERN - unless pattern in generator was changed
		//and we changed headers in all sources
		hashPart1 := unblendText(lpHeader, []int{6, 2, 7, 3, 2, 3, 9, 3, 5})
		hashPart2 := unblendText(trHeader, []int{5, 3, 6, 4, 7, 1, 7, 5, 2})
		hashPart3 := unblendText(rpHeader, []int{4, 1, 4, 5, 2, 4, 8, 6, 6})
		combinedHash = hashPart1 + hashPart2 + hashPart3
	}
	var requestSource RequestSource

	switch combinedHash {
	case CMSRequestSourceSHA1:
		requestSource = CMSRequestSource
	case AppRequestSourceSHA1:
		requestSource = AppRequestSource
	case PostmanRequestSourceSHA1:
		requestSource = PostmanRequestSource
	default:
		requestSource = OtherRequestSource
	}

	req := r.WithContext(context.WithValue(r.Context(), requestSourceContextKey, requestSource))

	return req, nil
}

func unblendText(blendedText string, pattern []int) string {
	result := ""
	isRandom := true
	for _, p := range pattern {
		if isRandom {
			blendedText = blendedText[p:]
			isRandom = false
			continue
		}

		result += blendedText[:p]
		blendedText = blendedText[p:]
		isRandom = true
	}
	return result
}

func GetRequestSource(r *http.Request) RequestSource {
	if requestSource, ok := r.Context().Value(requestSourceContextKey).(RequestSource); ok {
		return requestSource
	}
	return NotSetRequestSource
}
