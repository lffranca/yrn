package ybase

type (
	EventType    int
	DeployStatus string
)

const (
	EventTypeScheduler EventType = iota + 1
	EventTypeApi
	EventTypeHandler

	DeployStatusPending     DeployStatus = "PENDING"
	DeployStatusInProgress  DeployStatus = "IN_PROGRESS"
	DeployStatusInOperation DeployStatus = "IN_OPERATION"
	DeployStatusFailed      DeployStatus = "FAILED"
	DeployStatusRollback    DeployStatus = "ROLLBACK"
	DeployStatusCanceled    DeployStatus = "CANCELED"
)

type ConnectorParent struct {
	Connector         *Plugin     `json:"connector"`
	TransferCondition map[int]int `json:"transfer_condition"`
	Template          string      `json:"template,omitempty"`
	PositionData      string      `json:"position_data"`
}

type Event struct {
	Type EventType `json:"type,omitempty"`
	Body any       `json:"body,omitempty"`
}

type Flow struct {
	DeployStatus DeployStatus            `json:"deploy_status"`
	Version      string                  `json:"version"`
	TriggerEvent *Event                  `json:"trigger_event,omitempty"`
	Connectors   map[int]ConnectorParent `json:"connectors"`
}
