package domain

import (
	"fmt"

	"github.com/savannahghi/scalarutils"
)

// FHIRCondition definition: a clinical condition, problem, diagnosis, or other event, situation, issue, or clinical concept that has risen to a level of concern.
type FHIRCondition struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// Business identifiers assigned to this condition by the performer or other systems which remain constant as the resource is updated and propagates from server to server.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// The clinical status of the condition.
	ClinicalStatus *FHIRCodeableConcept `json:"clinicalStatus,omitempty"`

	// The verification status to support the clinical status of the condition.
	VerificationStatus *FHIRCodeableConcept `json:"verificationStatus,omitempty"`

	// A category assigned to the condition.
	Category []*FHIRCodeableConcept `json:"category,omitempty"`

	// A subjective assessment of the severity of the condition as evaluated by the clinician.
	Severity *FHIRCodeableConcept `json:"severity,omitempty"`

	// Identification of the condition, problem or diagnosis.
	Code *FHIRCodeableConcept `json:"code,omitempty"`

	// The anatomical location where this condition manifests itself.
	BodySite []*FHIRCodeableConcept `json:"bodySite,omitempty"`

	// Indicates the patient or group who the condition record is associated with.
	Subject *FHIRReference `json:"subject,omitempty"`

	// The Encounter during which this Condition was created or to which the creation of this record is tightly associated.
	Encounter *FHIRReference `json:"encounter,omitempty"`

	// Estimated or actual date or date-time  the condition began, in the opinion of the clinician.
	OnsetDateTime *scalarutils.Date `json:"onsetDateTime,omitempty"`

	// Estimated or actual date or date-time  the condition began, in the opinion of the clinician.
	OnsetAge *FHIRAge `json:"onsetAge,omitempty"`

	// Estimated or actual date or date-time  the condition began, in the opinion of the clinician.
	OnsetPeriod *FHIRPeriod `json:"onsetPeriod,omitempty"`

	// Estimated or actual date or date-time  the condition began, in the opinion of the clinician.
	OnsetRange *FHIRRange `json:"onsetRange,omitempty"`

	// Estimated or actual date or date-time  the condition began, in the opinion of the clinician.
	OnsetString *string `json:"onsetString,omitempty"`

	// The date or estimated date that the condition resolved or went into remission. This is called "abatement" because of the many overloaded connotations associated with "remission" or "resolution" - Conditions are never really resolved, but they can abate.
	AbatementDateTime *scalarutils.Date `json:"abatementDateTime,omitempty"`

	// The date or estimated date that the condition resolved or went into remission. This is called "abatement" because of the many overloaded connotations associated with "remission" or "resolution" - Conditions are never really resolved, but they can abate.
	AbatementAge *FHIRAge `json:"abatementAge,omitempty"`

	// The date or estimated date that the condition resolved or went into remission. This is called "abatement" because of the many overloaded connotations associated with "remission" or "resolution" - Conditions are never really resolved, but they can abate.
	AbatementPeriod *FHIRPeriod `json:"abatementPeriod,omitempty"`

	// The date or estimated date that the condition resolved or went into remission. This is called "abatement" because of the many overloaded connotations associated with "remission" or "resolution" - Conditions are never really resolved, but they can abate.
	AbatementRange *FHIRRange `json:"abatementRange,omitempty"`

	// The date or estimated date that the condition resolved or went into remission. This is called "abatement" because of the many overloaded connotations associated with "remission" or "resolution" - Conditions are never really resolved, but they can abate.
	AbatementString *string `json:"abatementString,omitempty"`

	// The recordedDate represents when this particular Condition record was created in the system, which is often a system-generated date.
	RecordedDate *scalarutils.Date `json:"recordedDate,omitempty"`

	// Individual who recorded the record and takes responsibility for its content.
	Recorder *FHIRReference `json:"recorder,omitempty"`

	// Individual who is making the condition statement.
	Asserter *FHIRReference `json:"asserter,omitempty"`

	// Clinical stage or grade of a condition. May include formal severity assessments.
	Stage []*FHIRConditionStage `json:"stage,omitempty"`

	// Supporting evidence / manifestations that are the basis of the Condition's verification status, such as evidence that confirmed or refuted the condition.
	Evidence []*FHIRConditionEvidence `json:"evidence,omitempty"`

	// Additional information about the Condition. This is a general notes/comments entry  for description of the Condition, its diagnosis and prognosis.
	Note []*FHIRAnnotation `json:"note,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMeta `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension   []*FHIRExtension           `json:"extension,omitempty"`
	Participant []FHIRConditionParticipant `json:"participant,omitempty"`
}

func (c *FHIRCondition) GetStage() string {
	if len(c.Stage) > 0 {
		for _, stage := range c.Stage {
			return stage.Summary.Text
		}
	}

	return ""
}

type FHIRConditionParticipant struct {
	ID                *string              `json:"id,omitempty"`
	Extension         []Extension          `json:"extension,omitempty"`
	ModifierExtension []Extension          `json:"modifierExtension,omitempty"`
	Function          *FHIRCodeableConcept `json:"function,omitempty"`
	Actor             *FHIRReference       `json:"actor"`
}

// FHIRConditionEvidence definition: a clinical condition, problem, diagnosis, or other event, situation, issue, or clinical concept that has risen to a level of concern.
type FHIRConditionEvidence struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A manifestation or symptom that led to the recording of this condition.
	Code *FHIRCodeableConcept `json:"code,omitempty"`

	// Links to other relevant information, including pathology reports.
	Detail []*FHIRReference `json:"detail,omitempty"`
}

// FHIRConditionEvidenceInput is the input type for ConditionEvidence
type FHIRConditionEvidenceInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A manifestation or symptom that led to the recording of this condition.
	Code *FHIRCodeableConceptInput `json:"code,omitempty"`

	// Links to other relevant information, including pathology reports.
	Detail []*FHIRReferenceInput `json:"detail,omitempty"`
}

// FHIRConditionInput is the input type for Condition
type FHIRConditionInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// Business identifiers assigned to this condition by the performer or other systems which remain constant as the resource is updated and propagates from server to server.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrativeInput `json:"text,omitempty"`

	// The clinical status of the condition.
	ClinicalStatus *FHIRCodeableConceptInput `json:"clinicalStatus,omitempty"`

	// The verification status to support the clinical status of the condition.
	VerificationStatus *FHIRCodeableConceptInput `json:"verificationStatus,omitempty"`

	// A category assigned to the condition.
	Category []*FHIRCodeableConceptInput `json:"category,omitempty"`

	// A subjective assessment of the severity of the condition as evaluated by the clinician.
	Severity *FHIRCodeableConceptInput `json:"severity,omitempty"`

	// Identification of the condition, problem or diagnosis.
	Code *FHIRCodeableConceptInput `json:"code,omitempty"`

	// The anatomical location where this condition manifests itself.
	BodySite []*FHIRCodeableConceptInput `json:"bodySite,omitempty"`

	// Indicates the patient or group who the condition record is associated with.
	Subject *FHIRReferenceInput `json:"subject,omitempty"`

	// The Encounter during which this Condition was created or to which the creation of this record is tightly associated.
	Encounter *FHIRReferenceInput `json:"encounter,omitempty"`

	// Estimated or actual date or date-time  the condition began, in the opinion of the clinician.
	OnsetDateTime *scalarutils.Date `json:"onsetDateTime,omitempty"`

	// Estimated or actual date or date-time  the condition began, in the opinion of the clinician.
	OnsetAge *FHIRAgeInput `json:"onsetAge,omitempty"`

	// Estimated or actual date or date-time  the condition began, in the opinion of the clinician.
	OnsetPeriod *FHIRPeriodInput `json:"onsetPeriod,omitempty"`

	// Estimated or actual date or date-time  the condition began, in the opinion of the clinician.
	OnsetRange *FHIRRangeInput `json:"onsetRange,omitempty"`

	// Estimated or actual date or date-time  the condition began, in the opinion of the clinician.
	OnsetString *string `json:"onsetString,omitempty"`

	// The date or estimated date that the condition resolved or went into remission. This is called "abatement" because of the many overloaded connotations associated with "remission" or "resolution" - Conditions are never really resolved, but they can abate.
	AbatementDateTime *scalarutils.Date `json:"abatementDateTime,omitempty"`

	// The date or estimated date that the condition resolved or went into remission. This is called "abatement" because of the many overloaded connotations associated with "remission" or "resolution" - Conditions are never really resolved, but they can abate.
	AbatementAge *FHIRAgeInput `json:"abatementAge,omitempty"`

	// The date or estimated date that the condition resolved or went into remission. This is called "abatement" because of the many overloaded connotations associated with "remission" or "resolution" - Conditions are never really resolved, but they can abate.
	AbatementPeriod *FHIRPeriodInput `json:"abatementPeriod,omitempty"`

	// The date or estimated date that the condition resolved or went into remission. This is called "abatement" because of the many overloaded connotations associated with "remission" or "resolution" - Conditions are never really resolved, but they can abate.
	AbatementRange *FHIRRangeInput `json:"abatementRange,omitempty"`

	// The date or estimated date that the condition resolved or went into remission. This is called "abatement" because of the many overloaded connotations associated with "remission" or "resolution" - Conditions are never really resolved, but they can abate.
	AbatementString *string `json:"abatementString,omitempty"`

	// The recordedDate represents when this particular Condition record was created in the system, which is often a system-generated date.
	RecordedDate *scalarutils.Date `json:"recordedDate,omitempty"`

	// Individual who recorded the record and takes responsibility for its content.
	Recorder *FHIRReferenceInput `json:"recorder,omitempty"`

	// Individual who is making the condition statement.
	Asserter *FHIRReferenceInput `json:"asserter,omitempty"`

	// Clinical stage or grade of a condition. May include formal severity assessments.
	Stage []*FHIRConditionStageInput `json:"stage,omitempty"`

	// Supporting evidence / manifestations that are the basis of the Condition's verification status, such as evidence that confirmed or refuted the condition.
	Evidence []*FHIRConditionEvidenceInput `json:"evidence,omitempty"`

	// Additional information about the Condition. This is a general notes/comments entry  for description of the Condition, its diagnosis and prognosis.
	Note []*FHIRAnnotationInput `json:"note,omitempty"`

	// Meta stores more information about the resource
	Meta FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`

	Participant []*FHIRConditionParticipantInput `json:"participant,omitempty"`
}

type FHIRConditionParticipantInput struct {
	ID                *string                   `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []Extension               `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []Extension               `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Function          *FHIRCodeableConceptInput `bson:"function,omitempty" json:"function,omitempty"`
	Actor             *FHIRReferenceInput       `bson:"actor" json:"actor"`
}

// ComposeConditionCodeField explicitly composes the Code field of a FHIR Condition resource.
// It strictly validates the provided TerminologySource and returns an error if invalid.
func (ci *FHIRConditionInput) ComposeConditionCodeField(code, displayName, system string, allowedSources map[string]string) error {
	systemURL, ok := allowedSources[system]
	if !ok {
		return fmt.Errorf("unsupported terminology source '%s'. Supported sources are ICD-11 and ICHI", system)
	}

	if code == "" || displayName == "" {
		return fmt.Errorf("code and displayName must both be provided")
	}

	ci.Code = &FHIRCodeableConceptInput{
		Coding: []*FHIRCodingInput{
			{
				Code:    scalarutils.Code(code),
				Display: displayName,
				System:  (*scalarutils.URI)(&systemURL),
			},
		},
		Text: displayName,
	}

	return nil
}

// FHIRConditionStage definition: a clinical condition, problem, diagnosis, or other event, situation, issue, or clinical concept that has risen to a level of concern.
type FHIRConditionStage struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A simple summary of the stage such as "Stage 3". The determination of the stage is disease-specific.
	Summary *FHIRCodeableConcept `json:"summary,omitempty"`

	// Reference to a formal record of the evidence on which the staging assessment is based.
	Assessment []*FHIRReference `json:"assessment,omitempty"`

	// The kind of staging, such as pathological or clinical staging.
	Type *FHIRCodeableConcept `json:"type,omitempty"`
}

// FHIRConditionStageInput is the input type for ConditionStage
type FHIRConditionStageInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	Extension []*Extension `json:"extension,omitempty"`

	// A simple summary of the stage such as "Stage 3". The determination of the stage is disease-specific.
	Summary *FHIRCodeableConceptInput `json:"summary,omitempty"`

	// Reference to a formal record of the evidence on which the staging assessment is based.
	Assessment []*FHIRReferenceInput `json:"assessment,omitempty"`

	// The kind of staging, such as pathological or clinical staging.
	Type *FHIRCodeableConceptInput `json:"type,omitempty"`
}

// PagedFHIRCondition ...
type PagedFHIRCondition struct {
	Conditions      []FHIRCondition
	HasNextPage     bool
	NextCursor      string
	HasPreviousPage bool
	PreviousCursor  string
	TotalCount      int
}

// FHIRConditionRelayPayload is used to return single instances of Condition
type FHIRConditionRelayPayload struct {
	Resource *FHIRCondition `json:"resource,omitempty"`
}
