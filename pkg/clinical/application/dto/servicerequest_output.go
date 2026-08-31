package dto

import "github.com/savannahghi/scalarutils"

// ServiceRequest is a record of a request for a procedure or diagnostic or other service to be planned
type ServiceRequest struct {
	ID                string                `json:"id,omitempty"`
	Status            ServiceRequestStatus  `json:"status,omitempty"`
	Intent            string                `json:"intent,omitempty"`
	Priority          string                `json:"priority,omitempty"`
	Note              []Annotation          `json:"note,omitempty"`
	Subject           Reference             `json:"subject,omitempty"`
	Encounter         *Reference            `json:"encounter,omitempty"`
	ReceivingFacility string                `json:"receivingFacility,omitempty"`
	Category          string                `json:"category,omitempty"`
	OrderDetails      OrderDetails          `json:"orderDetails,omitempty"`
	Date              *scalarutils.DateTime `json:"date,omitempty" swaggertype:"primitive,string"`
	Results           []Observation         `json:"results,omitempty"`
	UsageContext      ScreeningTypeEnum     `json:"usageContext,omitempty"`
}

type OrderDetails struct {
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}

// ServiceRequestEdge is an service request edge
type ServiceRequestEdge struct {
	Node   ServiceRequest `json:"node,omitempty"`
	Cursor string         `json:"cursor,omitempty"`
}

// ServiceRequestConnection  is an service request Connection type
type ServiceRequestOutputConnection struct {
	TotalCount *int                 `json:"totalCount"`
	Edges      []ServiceRequestEdge `json:"edges,omitempty"`
	PageInfo   PageInfo             `json:"pageInfo,omitempty"`
}

// CreateServiceRequestConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateServiceRequestConnection(serviceRequests []*ServiceRequest, pageInfo PageInfo, total *int) ServiceRequestOutputConnection {
	connection := ServiceRequestOutputConnection{
		TotalCount: total,
		Edges:      []ServiceRequestEdge{},
		PageInfo:   pageInfo,
	}

	for _, ServiceRequest := range serviceRequests {
		edge := ServiceRequestEdge{
			Node:   *ServiceRequest,
			Cursor: ServiceRequest.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}

// IsEntity marks a service request as an entity
func (s ServiceRequest) IsEntity() {}
