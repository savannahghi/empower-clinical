package domain

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/savannahghi/scalarutils"
)

// FHIRTask is a FHIR task resource data class
type FHIRTask struct {
	ID                    *string                  `json:"id,omitempty"`
	Meta                  *FHIRMeta                `json:"meta,omitempty"`
	ImplicitRules         *string                  `json:"implicitRules,omitempty"`
	Language              *string                  `json:"language,omitempty"`
	Text                  *FHIRNarrative           `json:"text,omitempty"`
	Extension             []*FHIRExtension         `json:"extension,omitempty"`
	ModifierExtension     []*FHIRExtension         `json:"modifierExtension,omitempty"`
	Identifier            []*FHIRIdentifier        `json:"identifier,omitempty"`
	InstantiatesCanonical *scalarutils.Canonical   `json:"instantiatesCanonical,omitempty"`
	InstantiatesUri       *scalarutils.URI         `json:"instantiatesUri,omitempty"`
	BasedOn               []*FHIRReference         `json:"basedOn,omitempty"`
	GroupIdentifier       *FHIRIdentifier          `json:"groupIdentifier,omitempty"`
	PartOf                []*FHIRReference         `json:"partOf,omitempty"`
	Status                *scalarutils.Code        `json:"status"`
	StatusReason          *FHIRCodeableReference   `json:"statusReason,omitempty"`
	BusinessStatus        *FHIRCodeableConcept     `json:"businessStatus,omitempty"`
	Intent                string                   `json:"intent"`
	Priority              *scalarutils.Code        `json:"priority,omitempty"`
	Code                  *FHIRCodeableConcept     `json:"code,omitempty"`
	Description           string                   `json:"description,omitempty"`
	Focus                 *FHIRReference           `json:"focus,omitempty"`
	For                   *FHIRReference           `json:"for,omitempty"`
	Encounter             *FHIRReference           `json:"encounter,omitempty"`
	ExecutionPeriod       *FHIRPeriod              `json:"executionPeriod,omitempty"`
	AuthoredOn            *string                  `json:"authoredOn,omitempty"`
	LastModified          *string                  `json:"lastModified,omitempty"`
	Requester             *FHIRReference           `json:"requester,omitempty"`
	PerformerType         []*FHIRCodeableConcept   `json:"performerType,omitempty"`
	Owner                 *FHIRReference           `json:"owner,omitempty"`
	Location              *FHIRReference           `json:"location,omitempty"`
	Reason                []*FHIRCodeableReference `json:"reason,omitempty"`
	ReasonReference       *FHIRReference           `json:"reasonReference,omitempty"`
	Insurance             []*FHIRReference         `json:"insurance,omitempty"`
	Note                  []*FHIRAnnotation        `json:"note,omitempty"`
	RelevantHistory       []*FHIRReference         `json:"relevantHistory,omitempty"`
	Restriction           *TaskRestriction         `json:"restriction,omitempty"`
	Input                 []*TaskInput             `json:"input,omitempty"`
	Output                []*TaskOutput            `json:"output,omitempty"`
}

// GetServiceRequestIDFromTask is used to extract the referral (service request) ID from a task
func (t *FHIRTask) GetServiceRequestIDFromTask() (string, error) {
	if t == nil {
		return "", errors.New("task is nil")
	}

	var referralID string

	for _, serviceRequest := range t.BasedOn {
		if serviceRequest.Type != nil && string(*serviceRequest.Type) == ReferralServiceRequestType.String() {
			referralID = fmt.Sprintf("ServiceRequest/%s", *serviceRequest.ID)
			break
		}
	}

	return referralID, nil
}

// TaskRestriction models the constraints on fulfillment tasks
type TaskRestriction struct {
	ID                *string          `json:"id,omitempty"`
	Extension         []*FHIRExtension `json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension `json:"modifierExtension,omitempty"`
	Repetitions       *int             `json:"repetitions,omitempty"`
	Period            *FHIRPeriod      `json:"period,omitempty"`
	Recipient         []*FHIRReference `json:"recipient,omitempty"`
}

// TaskInput models the information needed to perform a task
type TaskInput struct {
	ID                     *string                `json:"id,omitempty"`
	Extension              []*FHIRExtension       `json:"extension,omitempty"`
	ModifierExtension      []*FHIRExtension       `json:"modifierExtension,omitempty"`
	Type                   *FHIRCodeableConcept   `json:"type,omitempty"`
	ValueBase64Binary      string                 `json:"valueBase64Binary,omitempty"`
	ValueBoolean           bool                   `json:"valueBoolean,omitempty"`
	ValueCanonical         string                 `json:"valueCanonical,omitempty"`
	ValueCode              string                 `json:"valueCode,omitempty"`
	ValueDate              string                 `json:"valueDate,omitempty"`
	ValueDateTime          string                 `json:"valueDateTime,omitempty"`
	ValueDecimal           json.Number            `json:"valueDecimal,omitempty"`
	ValueId                string                 `json:"valueId,omitempty"`
	ValueInstant           string                 `json:"valueInstant,omitempty"`
	ValueInteger           int                    `json:"valueInteger,omitempty"`
	ValueMarkdown          string                 `json:"valueMarkdown,omitempty"`
	ValueOid               string                 `json:"valueOid,omitempty"`
	ValuePositiveInt       int                    `json:"valuePositiveInt,omitempty"`
	ValueString            string                 `json:"valueString,omitempty"`
	ValueTime              string                 `json:"valueTime,omitempty"`
	ValueUnsignedInt       int                    `json:"valueUnsignedInt,omitempty"`
	ValueReference         *FHIRReference         `json:"valueReference,omitempty"`
	ValueUri               string                 `json:"valueUri,omitempty"`
	ValueUrl               string                 `json:"valueUrl,omitempty"`
	ValueCodeableReference *FHIRCodeableReference `json:"valueCodeableReference"`
}

// TaskOutput models the information produced as part of task
type TaskOutput struct {
	ID                *string              `json:"id,omitempty"`
	Extension         []*FHIRExtension     `json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension     `json:"modifierExtension,omitempty"`
	Type              *FHIRCodeableConcept `json:"type"`
	ValueBase64Binary string               `json:"valueBase64Binary"`
	ValueBoolean      bool                 `json:"valueBoolean"`
	ValueCanonical    string               `json:"valueCanonical"`
	ValueCode         string               `json:"valueCode"`
	ValueDate         string               `json:"valueDate"`
	ValueDateTime     string               `json:"valueDateTime"`
	ValueDecimal      json.Number          `json:"valueDecimal"`
	ValueId           string               `json:"valueId"`
	ValueInstant      string               `json:"valueInstant"`
	ValueInteger      int                  `json:"valueInteger"`
	ValueMarkdown     string               `json:"valueMarkdown"`
	ValueOid          string               `json:"valueOid"`
	ValuePositiveInt  int                  `json:"valuePositiveInt"`
	ValueString       string               `json:"valueString"`
	ValueTime         string               `json:"valueTime"`
	ValueUnsignedInt  int                  `json:"valueUnsignedInt"`
	ValueUri          string               `json:"valueUri"`
	ValueUrl          string               `json:"valueUrl"`
	ValueUuid         string               `json:"valueUuid"`
}

// FHIRTaskInput models the payload used to create a task
type FHIRTaskInput struct {
	ResourceType       string                      `json:"resourceType,omitempty"`
	ID                 *string                     `json:"id,omitempty"`
	Meta               *FHIRMetaInput              `json:"meta,omitempty"`
	Text               *FHIRNarrative              `json:"text,omitempty"`
	Status             *scalarutils.Code           `json:"status,omitempty"`
	Reason             []*FHIRCodeableReference    `json:"reason,omitempty"`
	BusinessStatus     *FHIRCodeableConceptInput   `json:"businessStatus,omitempty"`
	PartOf             []FHIRReference             `json:"partOf,omitempty"`
	StatusReason       *FHIRCodeableReferenceInput `json:"statusReason,omitempty"`
	Intent             *scalarutils.Code           `json:"intent,omitempty"`
	Priority           *scalarutils.Code           `json:"priority,omitempty"`
	Code               *FHIRCodeableConceptInput   `json:"code,omitempty"`
	Description        *string                     `json:"description,omitempty"`
	BasedOn            []*FHIRReferenceInput       `json:"basedOn,omitempty"`
	For                *FHIRReferenceInput         `json:"for,omitempty"`
	Focus              *FHIRReferenceInput         `json:"focus,omitempty"`
	Encounter          *FHIRReferenceInput         `json:"encounter,omitempty"`
	AuthoredOn         *scalarutils.DateTime       `json:"authoredOn,omitempty"`
	Requester          *FHIRReferenceInput         `json:"requester,omitempty"`
	Owner              *FHIRReferenceInput         `json:"owner,omitempty"`
	ReasonReference    *FHIRReferenceInput         `json:"reasonReference,omitempty"`
	Note               []*FHIRAnnotationInput      `json:"note,omitempty"`
	RequestedPeriod    *FHIRPeriodInput            `json:"requestedPeriod,omitempty"`
	ExecutionPeriod    *FHIRPeriodInput            `json:"executionPeriod,omitempty"`
	LastModified       scalarutils.DateTime        `json:"lastModified,omitempty"`
	RequestedPerformer []*FHIRCodeableReference    `json:"requestedPerformer,omitempty"`
	Input              []*TaskInput                `json:"input,omitempty"`
}

// PagedFHIRTask is used to return paginated list of tasks
type PagedFHIRTask struct {
	Tasks           []FHIRTask `mapstructure:"Task" json:"tasks"`
	HasNextPage     bool       `mapstructure:"hasNextPage" json:"hasNextPage"`
	NextCursor      string     `mapstructure:"nextCursor" json:"nextCursor"`
	HasPreviousPage bool       `mapstructure:"hasPreviousPage" json:"hasPreviousPage"`
	PreviousCursor  string     `mapstructure:"previousCursor" json:"previousCursor"`
	TotalCount      *int       `mapstructure:"totalCount" json:"totalCount"`
}

// FHIRTaskRelayPayload is used to return single instances of Task
type FHIRTaskRelayPayload struct {
	Resource *FHIRTask `json:"resource,omitempty"`
}
