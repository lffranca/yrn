package pluginhttp

type HTTPSchema struct {
	Request HTTPRequest  `json:"request"`
	Retry   *RetryConfig `json:"retry,omitempty"`
}

type HTTPRequest struct {
	Method      string            `json:"method"` // GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers,omitempty"`
	QueryParams map[string]string `json:"queryParams,omitempty"`
	Body        interface{}       `json:"body,omitempty"` // string, object, array, or null
	Timeout     int               `json:"timeout,omitempty"`
}

type RetryConfig struct {
	MaxAttempts int `json:"maxAttempts"`
	Delay       int `json:"delay"` // milliseconds
}
