package domain

import (
	"github.com/savannahghi/scalarutils"
)

// FHIRAllergyIntolerance definition: risk of harmful or undesirable, physiological response which is unique to an individual and associated with exposure to a substance.
type FHIRAllergyIntolerance struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// Business identifiers assigned to this AllergyIntolerance by the performer or other systems which remain constant as the resource is updated and propagates from server to server.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// The clinical status of the allergy or intolerance.
	ClinicalStatus FHIRCodeableConcept `json:"clinicalStatus,omitempty"`

	// Assertion about certainty associated with the propensity, or potential risk, of a reaction to the identified substance (including pharmaceutical product).
	VerificationStatus FHIRCodeableConcept `json:"verificationStatus,omitempty"`

	// Identification of the underlying physiological mechanism for the reaction risk.
	Type *FHIRCodeableConcept `json:"type,omitempty"`

	// Category of the identified substance.
	Category []*AllergyIntoleranceCategoryEnum `json:"category,omitempty"`

	// Estimate of the potential clinical harm, or seriousness, of the reaction to the identified substance.
	Criticality AllergyIntoleranceCriticalityEnum `json:"criticality,omitempty"`

	// Code for an allergy or intolerance statement (either a positive or a negated/excluded statement).  This may be a code for a substance or pharmaceutical product that is considered to be responsible for the adverse reaction risk (e.g., "Latex"), an allergy or intolerance condition (e.g., "Latex allergy"), or a negated/excluded code for a specific substance or class (e.g., "No latex allergy") or a general or categorical negated statement (e.g.,  "No known allergy", "No known drug allergies").  Note: the substance for a specific reaction may be different from the substance identified as the cause of the risk, but it must be consistent with it. For instance, it may be a more specific substance (e.g. a brand medication) or a composite product that includes the identified substance. It must be clinically safe to only process the 'code' and ignore the 'reaction.substance'.  If a receiving system is unable to confirm that AllergyIntolerance.reaction.substance falls within the semantic scope of AllergyIntolerance.code, then the receiving system should ignore AllergyIntolerance.reaction.substance.
	Code *FHIRCodeableConcept `json:"code,omitempty"`

	// The patient who has the allergy or intolerance.
	Patient *FHIRReference `json:"patient,omitempty"`

	// The encounter when the allergy or intolerance was asserted.
	Encounter *FHIRReference `json:"encounter,omitempty"`

	// Estimated or actual date,  date-time, or age when allergy or intolerance was identified.
	OnsetDateTime *scalarutils.Date `json:"onsetDateTime,omitempty"`

	// Estimated or actual date,  date-time, or age when allergy or intolerance was identified.
	OnsetAge *FHIRAge `json:"onsetAge,omitempty"`

	// Estimated or actual date,  date-time, or age when allergy or intolerance was identified.
	OnsetPeriod *FHIRPeriod `json:"onsetPeriod,omitempty"`

	// Estimated or actual date,  date-time, or age when allergy or intolerance was identified.
	OnsetRange *FHIRRange `json:"onsetRange,omitempty"`

	// Estimated or actual date,  date-time, or age when allergy or intolerance was identified.
	OnsetString *string `json:"onsetString,omitempty"`

	// The recordedDate represents when this particular AllergyIntolerance record was created in the system, which is often a system-generated date.
	RecordedDate *scalarutils.DateTime `json:"recordedDate,omitempty"`

	// Individual who recorded the record and takes responsibility for its content.
	Recorder *FHIRReference `json:"recorder,omitempty"`

	// The source of the information about the allergy that is recorded.
	Asserter *FHIRReference `json:"asserter,omitempty"`

	// Represents the date and/or time of the last known occurrence of a reaction event.
	LastOccurrence *scalarutils.DateTime `json:"lastOccurrence,omitempty"`

	// Additional narrative about the propensity for the Adverse Reaction, not captured in other fields.
	Note []*FHIRAnnotation `json:"note,omitempty"`

	// Details about each adverse reaction event linked to exposure to the identified substance.
	Reaction []*FHIRAllergyintoleranceReaction `json:"reaction,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMeta `json:"meta,omitempty"`

	Participant []*AllergyIntoleranceParticipant `json:"participant,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

type AllergyIntoleranceParticipant struct {
	ID                *string              `json:"id,omitempty"`
	Extension         []FHIRExtension      `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension      `json:"modifierExtension,omitempty"`
	Function          *FHIRCodeableConcept `json:"function,omitempty"`
	Actor             *FHIRReference       `json:"actor"`
}

// FHIRAllergyIntoleranceInput is the input type for AllergyIntolerance
type FHIRAllergyIntoleranceInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	Text *FHIRNarrativeInput `json:"text,omitempty"`

	// Business identifiers assigned to this AllergyIntolerance by the performer or other systems which remain constant as the resource is updated and propagates from server to server.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	// The clinical status of the allergy or intolerance.
	ClinicalStatus FHIRCodeableConceptInput `json:"clinicalStatus,omitempty"`

	// Assertion about certainty associated with the propensity, or potential risk, of a reaction to the identified substance (including pharmaceutical product).
	VerificationStatus FHIRCodeableConceptInput `json:"verificationStatus,omitempty"`

	// Identification of the underlying physiological mechanism for the reaction risk.
	Type *FHIRCodeableConceptInput `json:"type,omitempty"`

	// Category of the identified substance.
	Category []*AllergyIntoleranceCategoryEnum `json:"category,omitempty"`

	// Estimate of the potential clinical harm, or seriousness, of the reaction to the identified substance.
	Criticality AllergyIntoleranceCriticalityEnum `json:"criticality,omitempty"`

	// Code for an allergy or intolerance statement (either a positive or a negated/excluded statement).  This may be a code for a substance or pharmaceutical product that is considered to be responsible for the adverse reaction risk (e.g., "Latex"), an allergy or intolerance condition (e.g., "Latex allergy"), or a negated/excluded code for a specific substance or class (e.g., "No latex allergy") or a general or categorical negated statement (e.g.,  "No known allergy", "No known drug allergies").  Note: the substance for a specific reaction may be different from the substance identified as the cause of the risk, but it must be consistent with it. For instance, it may be a more specific substance (e.g. a brand medication) or a composite product that includes the identified substance. It must be clinically safe to only process the 'code' and ignore the 'reaction.substance'.  If a receiving system is unable to confirm that AllergyIntolerance.reaction.substance falls within the semantic scope of AllergyIntolerance.code, then the receiving system should ignore AllergyIntolerance.reaction.substance.
	Code FHIRCodeableConceptInput `json:"code,omitempty"`

	// The patient who has the allergy or intolerance.
	Patient *FHIRReferenceInput `json:"patient,omitempty"`

	// The encounter when the allergy or intolerance was asserted.
	Encounter *FHIRReferenceInput `json:"encounter,omitempty"`

	// Estimated or actual date,  date-time, or age when allergy or intolerance was identified.
	OnsetDateTime *scalarutils.DateTime `json:"onsetDateTime,omitempty"`

	// Estimated or actual date,  date-time, or age when allergy or intolerance was identified.
	OnsetAge *FHIRAgeInput `json:"onsetAge,omitempty"`

	// Estimated or actual date,  date-time, or age when allergy or intolerance was identified.
	OnsetPeriod *FHIRPeriodInput `json:"onsetPeriod,omitempty"`

	// Estimated or actual date,  date-time, or age when allergy or intolerance was identified.
	OnsetRange *FHIRRangeInput `json:"onsetRange,omitempty"`

	// Estimated or actual date,  date-time, or age when allergy or intolerance was identified.
	OnsetString *string `json:"onsetString,omitempty"`

	// The recordedDate represents when this particular AllergyIntolerance record was created in the system, which is often a system-generated date.
	RecordedDate *scalarutils.DateTime `json:"recordedDate,omitempty"`

	// Individual who recorded the record and takes responsibility for its content.
	Recorder *FHIRReferenceInput `json:"recorder,omitempty"`

	// The source of the information about the allergy that is recorded.
	Asserter *FHIRReferenceInput `json:"asserter,omitempty"`

	// Represents the date and/or time of the last known occurrence of a reaction event.
	LastOccurrence *scalarutils.DateTime `json:"lastOccurrence,omitempty"`

	// Additional narrative about the propensity for the Adverse Reaction, not captured in other fields.
	Note []*FHIRAnnotationInput `json:"note,omitempty"`

	// Details about each adverse reaction event linked to exposure to the identified substance.
	Reaction []*FHIRAllergyintoleranceReactionInput `json:"reaction,omitempty"`

	// Meta stores more information about the resource
	Meta FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`

	Participant []*AllergyIntoleranceParticipantInput `bson:"participant,omitempty" json:"participant,omitempty"`
}

type AllergyIntoleranceParticipantInput struct {
	ID                *string             `json:"id,omitempty"`
	Extension         []*FHIRExtension    `json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension    `json:"modifierExtension,omitempty"`
	Actor             *FHIRReferenceInput `json:"actor,omitempty"`
}

// FHIRAllergyintoleranceReaction definition: risk of harmful or undesirable, physiological response which is unique to an individual and associated with exposure to a substance.
type FHIRAllergyintoleranceReaction struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Identification of the specific substance (or pharmaceutical product) considered to be responsible for the Adverse Reaction event. Note: the substance for a specific reaction may be different from the substance identified as the cause of the risk, but it must be consistent with it. For instance, it may be a more specific substance (e.g. a brand medication) or a composite product that includes the identified substance. It must be clinically safe to only process the 'code' and ignore the 'reaction.substance'.  If a receiving system is unable to confirm that AllergyIntolerance.reaction.substance falls within the semantic scope of AllergyIntolerance.code, then the receiving system should ignore AllergyIntolerance.reaction.substance.
	Substance *FHIRCodeableConcept `json:"substance,omitempty"`

	// Clinical symptoms and/or signs that are observed or associated with the adverse reaction event.
	Manifestation []*FHIRCodeableReference `json:"manifestation,omitempty"`

	// Text description about the reaction as a whole, including details of the manifestation if required.
	Description *string `json:"description,omitempty"`

	// Record of the date and/or time of the onset of the Reaction.
	Onset *scalarutils.DateTime `json:"onset,omitempty"`

	// Clinical assessment of the severity of the reaction event as a whole, potentially considering multiple different manifestations.
	Severity *AllergyIntoleranceReactionSeverityEnum `json:"severity,omitempty"`

	// Identification of the route by which the subject was exposed to the substance.
	ExposureRoute *FHIRCodeableConcept `json:"exposureRoute,omitempty"`

	// Additional text about the adverse reaction event not captured in other fields.
	Note []*FHIRAnnotation `json:"note,omitempty"`
}

// FHIRAllergyintoleranceReactionInput is the input type for AllergyintoleranceReaction
type FHIRAllergyintoleranceReactionInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Identification of the specific substance (or pharmaceutical product) considered to be responsible for the Adverse Reaction event. Note: the substance for a specific reaction may be different from the substance identified as the cause of the risk, but it must be consistent with it. For instance, it may be a more specific substance (e.g. a brand medication) or a composite product that includes the identified substance. It must be clinically safe to only process the 'code' and ignore the 'reaction.substance'.  If a receiving system is unable to confirm that AllergyIntolerance.reaction.substance falls within the semantic scope of AllergyIntolerance.code, then the receiving system should ignore AllergyIntolerance.reaction.substance.
	Substance *FHIRCodeableConceptInput `json:"substance,omitempty"`

	// Clinical symptoms and/or signs that are observed or associated with the adverse reaction event.
	Manifestation []*FHIRCodeableReferenceInput `json:"manifestation,omitempty"`

	// Text description about the reaction as a whole, including details of the manifestation if required.
	Description *string `json:"description,omitempty"`

	// Record of the date and/or time of the onset of the Reaction.
	Onset *scalarutils.DateTime `json:"onset,omitempty"`

	// Clinical assessment of the severity of the reaction event as a whole, potentially considering multiple different manifestations.
	Severity *AllergyIntoleranceReactionSeverityEnum `json:"severity,omitempty"`

	// Identification of the route by which the subject was exposed to the substance.
	ExposureRoute *FHIRCodeableConceptInput `json:"exposureRoute,omitempty"`

	// Additional text about the adverse reaction event not captured in other fields.
	Note []*FHIRAnnotationInput `json:"note,omitempty"`
}

// FHIRAllergyIntoleranceRelayPayload is used to return single instances of AllergyIntolerance
type FHIRAllergyIntoleranceRelayPayload struct {
	Resource *FHIRAllergyIntolerance `json:"resource,omitempty"`
}

// PagedFHIRAllergy is an allergy's paginaton dataclass
type PagedFHIRAllergy struct {
	Allergies       []FHIRAllergyIntolerance
	HasNextPage     bool
	NextCursor      string
	HasPreviousPage bool
	PreviousCursor  string
	TotalCount      int
}

func (f *FHIRAllergyIntolerance) GetConceptData() *FHIRCoding {
	if f.Code != nil {
		for _, item := range f.Code.Coding {
			if item.Code != nil && item.System != nil && item.Display != "" {
				return &FHIRCoding{
					Code:    item.Code,
					Display: item.Display,
					System:  item.System,
				}
			}
		}
	}

	return nil
}
