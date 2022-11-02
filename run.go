package main

import "time"

// Based on https://www.terraform.io/docs/cloud/api/notification-configurations.html#sample-payload
type Run struct {
	ID             string         `json:"run_id"`
	PayloadVersion int            `json:"payload_version"`
	RunURL         string         `json:"run_url"`
	WorkspaceID    string         `json:"workspace_id"`
	WorkspaceName  string         `json:"workspace_name"`
	Notifications  []Notification `json:"notifications"`
	Message        string         `json:"run_message"`
	CreatedAt      time.Time      `json:"run_created_at"`
	CreatedBy      string         `json:"run_created_by"`
}

type Notification struct {
	Message      string    `json:"message"`
	Trigger      string    `json:"trigger"`
	RunStatus    string    `json:"run_status"`
	RunUpdatedAt time.Time `json:"run_updated_at"`
}
