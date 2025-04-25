package flowmanager

type (
	Flow struct {
		Id               string       `json:"id"`
		Name             string       `json:"name"`
		Description      string       `json:"description"`
		Tenant           string       `json:"tenant"`
		FirstPluginToRun string       `json:"first_plugin_to_run"`
		Plugins          []FlowPlugin `json:"plugins"`
		Version          int          `json:"version"`
	}

	FlowPlugin struct {
		Id                          string   `json:"id"`
		Slug                        string   `json:"slug"`
		Name                        string   `json:"name"`
		Description                 string   `json:"description"`
		Version                     int      `json:"version"`
		SchemaInput                 string   `json:"schema_input"`
		ContinueEvenWithError       bool     `json:"continue_even_with_error"`
		ShareResponseWithAllPlugins bool     `json:"share_response_with_all_plugins"`
		NextToBeExecuted            []string `json:"next_to_be_executed"`
	}
)
