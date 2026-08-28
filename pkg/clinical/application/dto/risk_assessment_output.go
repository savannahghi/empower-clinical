package dto

import (
	"strings"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
)

// RiskAssessment ...
type RiskAssessment struct {
	ID           *string                    `json:"id,omitempty" mapstructure:"id"`
	Subject      Reference                  `json:"subject,omitempty"`
	Encounter    *Reference                 `json:"encounter,omitempty"`
	Prediction   []RiskAssessmentPrediction `json:"prediction,omitempty"`
	Note         []Annotation               `json:"note,omitempty"`
	Date         *scalarutils.DateTime      `json:"occurrenceDateTime,omitempty" swaggertype:"primitive,string"`
	UsageContext string                     `json:"usageContext,omitempty"`
	Basis        []Reference                `json:"basis"`
}

func (ra *RiskAssessment) AppendQuestionnaireResponseReferenceID() {
	for i := range ra.Basis {
		if ra.Basis[i].ID == "" {
			basisID := strings.TrimPrefix(ra.Basis[i].Reference, "QuestionnaireResponse/")
			ra.Basis[i].ID = basisID
		}
	}
}

// RiskAssessmentPrediction describes the predicted outcome
type RiskAssessmentPrediction struct {
	ID                 *string          `json:"id,omitempty"`
	Outcome            *CodeableConcept `json:"qualitativeRisk,omitempty"`
	ProbabilityDecimal *float64         `json:"probabilityDecimal,omitempty"`
}

// NodeID implements the NodeID() method that returns ID of a Node
func (r RiskAssessment) NodeID() string {
	return *r.ID
}

// ConvertRiskAssessmentEdges converts the generic edges to a risk assessment edges (strict go types ⛓️ )
func ConvertRiskAssessmentEdges(edges []serverutils.Edge[RiskAssessment]) []RiskAssessmentEdge {
	result := make([]RiskAssessmentEdge, len(edges))
	for i, edge := range edges {
		result[i] = RiskAssessmentEdge{
			Cursor: edge.Cursor,
			Node:   edge.Node,
		}
	}

	return result
}

// RiskAssessmentEdge is an risk assessment edge
type RiskAssessmentEdge struct {
	Node   RiskAssessment `json:"node,omitempty"`
	Cursor string         `json:"cursor,omitempty"`
}

// RiskAssessmentConnection  is an RiskAssessmentConnection Connection Type
type RiskAssessmentConnection struct {
	SearchID   string               `json:"searchID"`
	TotalCount int                  `json:"totalCount"`
	Edges      []RiskAssessmentEdge `json:"edges,omitempty"`
	PageInfo   serverutils.PageInfo `json:"pageInfo,omitempty"`
}

// CreateRiskAssessmentConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateRiskAssessmentConnection(assessments []*RiskAssessment, pageInfo serverutils.PageInfo, total int) RiskAssessmentConnection {
	connection := RiskAssessmentConnection{
		TotalCount: total,
		Edges:      []RiskAssessmentEdge{},
		PageInfo:   pageInfo,
	}

	for _, assessment := range assessments {
		edge := RiskAssessmentEdge{
			Node:   *assessment,
			Cursor: *assessment.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}

// ObservationOutput serves a sole purpose of serving screening report with the data model
// that distinctively states categorizes Examinations and Tests
type ObservationOutput struct {
	Examinations []Observation    `json:"examinations,omitempty"`
	Tests        []Observation    `json:"tests,omitempty"`
	ReferredTest []ReferralDetail `json:"referredTests,omitempty"`
}

// ScreeningReport is used to get the screening report
type ScreeningReport struct {
	Consent         []*Consent        `json:"consent,omitempty"`
	RiskAssessment  []*RiskAssessment `json:"riskAssessment,omitempty"`
	Observation     ObservationOutput `json:"observation,omitempty"`
	ReferralDetails []ReferralDetail  `json:"referralDetails,omitempty"`
}
