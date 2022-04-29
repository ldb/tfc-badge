package main

// Based on https://www.terraform.io/docs/cloud/api/notification-configurations.html#sample-payload
type Run struct {
	PayloadVersion int            `payload_version`
	RunURL         string         `json:"run_url"`
	WorkspaceID    string         `json:"workspace_id"`
	WorkspaceName  string         `json:"workspace_name"`
	Notifications  []Notification `json:"notifications"`
}

type Notification struct {
	Message      string `json:"message"`
	Trigger      string `json:"trigger"`
	RunStatus    string `json:"run_status"`
	RunUpdatedAt string `json:"run_updated_at"`
}
