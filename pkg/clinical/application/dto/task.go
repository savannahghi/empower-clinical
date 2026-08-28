package dto

import (
	"github.com/savannahghi/scalarutils"
)

// TaskInput modes the data class used to create a task
type TaskInput struct {
	EncounterID string                `json:"encounterID"`
	Task        string                `json:"task"`
	Workflow    ScreeningTypeEnum     `json:"workflow"`
	Description string                `json:"description"`
	DueDate     *scalarutils.DateTime `json:"dueDate"`
}

// TaskOutput is used to display a task result
type TaskOutput struct {
	ID           string                `json:"id,omitempty"`
	EncounterID  string                `json:"encounterID,omitempty"`
	Task         string                `json:"task,omitempty"`
	Description  string                `json:"description,omitempty"`
	Status       TaskStatus            `json:"status,omitempty"`
	StatusReason string                `json:"statusReason,omitempty"`
	Workflow     string                `json:"workflow,omitempty"`
	AuthoredOn   *scalarutils.DateTime `json:"authoredOn,omitempty" swaggertype:"primitive,string"`
	DueDate      *scalarutils.DateTime `json:"dueDate,omitempty" swaggertype:"primitive,string"`
	Priority     string                `json:"priority,omitempty"`
	Attachment   []*AttachmentResponse `json:"attachment,omitempty"`
	Notes        *Annotation           `json:"notes,omitempty"`
	Subject      *Reference            `json:"subject,omitempty"`
	LastUpdated  *scalarutils.DateTime `json:"lastUpdated,omitempty" swaggertype:"primitive,string"`
	UsageContext ScreeningTypeEnum     `json:"usageContext,omitempty"`
}

// AttachmentResponse is a custom response model to mention what kind of attachment it represents
type AttachmentResponse struct {
	Type       string      `json:"type,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

// TaskEdge is an task edge
type TaskEdge struct {
	Node   TaskOutput `json:"node,omitempty"`
	Cursor string     `json:"cursor,omitempty"`
}

// TaskConnection  is an TaskConnection Connection type
type TaskOutputConnection struct {
	TotalCount *int       `json:"totalCount"`
	Edges      []TaskEdge `json:"edges,omitempty"`
	PageInfo   PageInfo   `json:"pageInfo,omitempty"`
}

// CreateTaskConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateTaskConnection(tasks []*TaskOutput, pageInfo PageInfo, total *int) TaskOutputConnection {
	connection := TaskOutputConnection{
		TotalCount: total,
		Edges:      []TaskEdge{},
		PageInfo:   pageInfo,
	}

	for _, task := range tasks {
		edge := TaskEdge{
			Node:   *task,
			Cursor: task.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}
