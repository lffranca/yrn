package flowmanager

type (
	Flow struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Tenant      string   `json:"tenant"`
		Plugins     []Plugin `json:"plugins"`
		Version     int      `json:"version"`
	}

	Plugin struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Schema      string   `json:"schema"`
		DiagramData string   `json:"diagram_data"`
		FlowId      string   `json:"flow_id"`
		Tenant      string   `json:"tenant"`
		NextSteps   []string `json:"next_steps"`
	}
)
