package parse

type ContentParser interface {
	Parse(configContent interface{}) (map[string]interface{}, error)
}
