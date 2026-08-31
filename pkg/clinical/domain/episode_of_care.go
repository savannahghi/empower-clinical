package domain

import (
	"github.com/savannahghi/firebasetools"
)

// FHIREpisodeOfCare definition: an association between a patient and an organization / healthcare provider(s) during which time encounters may occur. the managing organization assumes a level of responsibility for the patient during this time.
type FHIREpisodeOfCare struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// The EpisodeOfCare may be known by different identifiers for different contexts of use, such as when an external agency is tracking the Episode for funding purposes.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// planned | waitlist | active | onhold | finished | cancelled.
	Status *EpisodeOfCareStatusEnum `json:"status,omitempty"`

	// The history of statuses that the EpisodeOfCare has been through (without requiring processing the history of the resource).
	StatusHistory []*FHIREpisodeofcareStatushistory `json:"statusHistory,omitempty"`

	// A classification of the type of episode of care; e.g. specialist referral, disease management, type of funded care.
	Type []*FHIRCodeableConcept `json:"type,omitempty"`

	// The list of diagnosis relevant to this episode of care.
	Diagnosis []*FHIREpisodeofcareDiagnosis `json:"diagnosis,omitempty"`

	// The patient who is the focus of this episode of care.
	Patient *FHIRReference `json:"patient,omitempty"`

	// The organization that has assumed the specific responsibilities for the specified duration.
	ManagingOrganization *FHIRReference `json:"managingOrganization,omitempty"`

	// The interval during which the managing organization assumes the defined responsibility.
	Period *FHIRPeriod `json:"period,omitempty"`

	// Referral Request(s) that are fulfilled by this EpisodeOfCare, incoming referrals.
	ReferralRequest []*FHIRReference `json:"referralRequest,omitempty"`

	// The practitioner that is the care manager/care coordinator for this patient.
	CareManager *FHIRReference `json:"careManager,omitempty"`

	// The list of practitioners that may be facilitating this episode of care for specific purposes.
	Team []*FHIRReference `json:"team,omitempty"`

	// The set of accounts that may be used for billing for this EpisodeOfCare.
	Account []*FHIRReference `json:"account,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMeta `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

// FHIREpisodeOfCareInput is the input type for EpisodeOfCare
type FHIREpisodeOfCareInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// The EpisodeOfCare may be known by different identifiers for different contexts of use, such as when an external agency is tracking the Episode for funding purposes.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	// planned | waitlist | active | onhold | finished | cancelled.
	Status *EpisodeOfCareStatusEnum `json:"status,omitempty"`

	// The history of statuses that the EpisodeOfCare has been through (without requiring processing the history of the resource).
	StatusHistory []*FHIREpisodeofcareStatushistoryInput `json:"statusHistory,omitempty"`

	// A classification of the type of episode of care; e.g. specialist referral, disease management, type of funded care.
	Type []*FHIRCodeableConceptInput `json:"type,omitempty"`

	// The list of diagnosis relevant to this episode of care.
	Diagnosis []*FHIREpisodeofcareDiagnosisInput `json:"diagnosis,omitempty"`

	// The patient who is the focus of this episode of care.
	Patient *FHIRReferenceInput `json:"patient,omitempty"`

	// The organization that has assumed the specific responsibilities for the specified duration.
	ManagingOrganization *FHIRReferenceInput `json:"managingOrganization,omitempty"`

	// The interval during which the managing organization assumes the defined responsibility.
	Period *FHIRPeriodInput `json:"period,omitempty"`

	// Referral Request(s) that are fulfilled by this EpisodeOfCare, incoming referrals.
	ReferralRequest []*FHIRReferenceInput `json:"referralRequest,omitempty"`

	// The practitioner that is the care manager/care coordinator for this patient.
	CareManager *FHIRReferenceInput `json:"careManager,omitempty"`

	// The list of practitioners that may be facilitating this episode of care for specific purposes.
	Team []*FHIRReferenceInput `json:"team,omitempty"`

	// The set of accounts that may be used for billing for this EpisodeOfCare.
	Account []*FHIRReferenceInput `json:"account,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

// FHIREpisodeOfCareRelayConnection is a Relay connection for EpisodeOfCare
type FHIREpisodeOfCareRelayConnection struct {
	Edges []*FHIREpisodeOfCareRelayEdge `json:"edges,omitempty"`

	PageInfo *firebasetools.PageInfo `json:"pageInfo,omitempty"`
}

// FHIREpisodeOfCareRelayEdge is a Relay edge for EpisodeOfCare
type FHIREpisodeOfCareRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIREpisodeOfCare `json:"node,omitempty"`
}

// FHIREpisodeOfCareRelayPayload is used to return single instances of EpisodeOfCare
type FHIREpisodeOfCareRelayPayload struct {
	Resource *FHIREpisodeOfCare `json:"resource,omitempty"`
}

// FHIREpisodeofcareDiagnosis definition: an association between a patient and an organization / healthcare provider(s) during which time encounters may occur. the managing organization assumes a level of responsibility for the patient during this time.
type FHIREpisodeofcareDiagnosis struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A list of conditions/problems/diagnoses that this episode of care is intended to be providing care for.
	Condition *FHIRReference `json:"condition,omitempty"`

	// Role that this diagnosis has within the episode of care (e.g. admission, billing, discharge …).
	Role *FHIRCodeableConcept `json:"role,omitempty"`

	// Ranking of the diagnosis (for each role type).
	Rank *string `json:"rank,omitempty"`
}

// FHIREpisodeofcareDiagnosisInput is the input type for EpisodeofcareDiagnosis
type FHIREpisodeofcareDiagnosisInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A list of conditions/problems/diagnoses that this episode of care is intended to be providing care for.
	Condition *FHIRReferenceInput `json:"condition,omitempty"`

	// Role that this diagnosis has within the episode of care (e.g. admission, billing, discharge …).
	Role *FHIRCodeableConceptInput `json:"role,omitempty"`

	// Ranking of the diagnosis (for each role type).
	Rank *string `json:"rank,omitempty"`
}

// FHIREpisodeofcareStatushistory definition: an association between a patient and an organization / healthcare provider(s) during which time encounters may occur. the managing organization assumes a level of responsibility for the patient during this time.
type FHIREpisodeofcareStatushistory struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// planned | waitlist | active | onhold | finished | cancelled.
	Status *EpisodeOfCareStatusHistoryStatusEnum `json:"status,omitempty"`

	// The period during this EpisodeOfCare that the specific status applied.
	Period *FHIRPeriod `json:"period,omitempty"`
}

// FHIREpisodeofcareStatushistoryInput is the input type for EpisodeofcareStatushistory
type FHIREpisodeofcareStatushistoryInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// planned | waitlist | active | onhold | finished | cancelled.
	Status *EpisodeOfCareStatusHistoryStatusEnum `json:"status,omitempty"`

	// The period during this EpisodeOfCare that the specific status applied.
	Period *FHIRPeriodInput `json:"period,omitempty"`
}

// EpisodeOfCarePayload is used to return the results after creation of
// episodes of care
type EpisodeOfCarePayload struct {
	EpisodeOfCare *FHIREpisodeOfCare `json:"episodeOfCare"`
}
