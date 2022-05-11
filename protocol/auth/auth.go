package auth

type HTTPAuth interface {
	HTTPHeaders(url string, appID string, secret string) map[string][]string
}
