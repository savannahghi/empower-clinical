package domain

import (
	"time"

	"github.com/savannahghi/scalarutils"
)

// FHIRObservation definition: measurements and simple assertions made about a patient, device or other subject.
type FHIRObservation struct {

	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// A unique identifier assigned to this observation.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// A plan, proposal or order that is fulfilled in whole or in part by this event.  For example, a MedicationRequest may require a patient to have laboratory test performed before  it is dispensed.
	BasedOn []*FHIRReference `json:"basedOn,omitempty"`

	// A larger event of which this particular Observation is a component or step.  For example,  an observation as part of a procedure.
	PartOf []*FHIRReference `json:"partOf,omitempty"`

	// The status of the result value.
	Status *ObservationStatusEnum `json:"status,omitempty"`

	// A code that classifies the general type of observation being made.
	Category []*FHIRCodeableConcept `json:"category,omitempty"`

	// Describes what was observed. Sometimes this is called the observation "name".
	Code *FHIRCodeableConcept `json:"code,omitempty"`

	// The patient, or group of patients, location, or device this observation is about and into whose record the observation is placed. If the actual focus of the observation is different from the subject (or a sample of, part, or region of the subject), the `focus` element or the `code` itself specifies the actual focus of the observation.
	Subject *FHIRReference `json:"subject,omitempty"`

	// The actual focus of an observation when it is not the patient of record representing something or someone associated with the patient such as a spouse, parent, fetus, or donor. For example, fetus observations in a mother's record.  The focus of an observation could also be an existing condition,  an intervention, the subject's diet,  another observation of the subject,  or a body structure such as tumor or implanted device.   An example use case would be using the Observation resource to capture whether the mother is trained to change her child's tracheostomy tube. In this example, the child is the patient of record and the mother is the focus.
	Focus []*FHIRReference `json:"focus,omitempty"`

	// The healthcare event  (e.g. a patient and healthcare provider interaction) during which this observation is made.
	Encounter *FHIRReference `json:"encounter,omitempty"`

	// The time or time-period the observed value is asserted as being true. For biological subjects - e.g. human patients - this is usually called the "physiologically relevant time". This is usually either the time of the procedure or of specimen collection, but very often the source of the date/time is not known, only the date/time itself.
	EffectiveDateTime *scalarutils.DateTime `json:"effectiveDateTime,omitempty"`

	// The time or time-period the observed value is asserted as being true. For biological subjects - e.g. human patients - this is usually called the "physiologically relevant time". This is usually either the time of the procedure or of specimen collection, but very often the source of the date/time is not known, only the date/time itself.
	EffectivePeriod *FHIRPeriod `json:"effectivePeriod,omitempty"`

	// The time or time-period the observed value is asserted as being true. For biological subjects - e.g. human patients - this is usually called the "physiologically relevant time". This is usually either the time of the procedure or of specimen collection, but very often the source of the date/time is not known, only the date/time itself.
	EffectiveTiming *FHIRTiming `json:"effectiveTiming,omitempty"`

	// The time or time-period the observed value is asserted as being true. For biological subjects - e.g. human patients - this is usually called the "physiologically relevant time". This is usually either the time of the procedure or of specimen collection, but very often the source of the date/time is not known, only the date/time itself.
	EffectiveInstant *scalarutils.Instant `json:"effectiveInstant,omitempty"`

	// The date and time this version of the observation was made available to providers, typically after the results have been reviewed and verified.
	Issued *scalarutils.Instant `json:"issued,omitempty"`

	// Who was responsible for asserting the observed value as "true".
	Performer []*FHIRReference `json:"performer,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueQuantity *FHIRQuantity `json:"valueQuantity,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueCodeableConcept *FHIRCodeableConcept `json:"valueCodeableConcept,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueString *string `json:"valueString,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueBoolean *bool `json:"valueBoolean,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueInteger *string `json:"valueInteger,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueRange *FHIRRange `json:"valueRange,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueRatio *FHIRRatio `json:"valueRatio,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueSampledData *FHIRSampledData `json:"valueSampledData,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueTime *time.Time `json:"valueTime,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueDateTime *scalarutils.Date `json:"valueDateTime,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValuePeriod *FHIRPeriod `json:"valuePeriod,omitempty"`

	// Provides a reason why the expected value in the element Observation.value[x] is missing.
	DataAbsentReason *FHIRCodeableConcept `json:"dataAbsentReason,omitempty"`

	// A categorical assessment of an observation value.  For example, high, low, normal.
	Interpretation []*FHIRCodeableConcept `json:"interpretation,omitempty"`

	// Comments about the observation or the results.
	Note []*FHIRAnnotation `json:"note,omitempty"`

	// Indicates the site on the subject's body where the observation was made (i.e. the target site).
	BodySite *FHIRCodeableConcept `json:"bodySite,omitempty"`

	// Indicates the mechanism used to perform the observation.
	Method *FHIRCodeableConcept `json:"method,omitempty"`

	// The specimen that was used when this observation was made.
	Specimen *FHIRReference `json:"specimen,omitempty"`

	// The device used to generate the observation data.
	Device *FHIRReference `json:"device,omitempty"`

	// Guidance on how to interpret the value by comparison to a normal or recommended range.  Multiple reference ranges are interpreted as an "OR".   In other words, to represent two distinct target populations, two `referenceRange` elements would be used.
	ReferenceRange []*FHIRObservationReferencerange `json:"referenceRange,omitempty"`

	// This observation is a group observation (e.g. a battery, a panel of tests, a set of vital sign measurements) that includes the target as a member of the group.
	HasMember []*FHIRReference `json:"hasMember,omitempty"`

	// The target resource that represents a measurement from which this observation value is derived. For example, a calculated anion gap or a fetal measurement based on an ultrasound image.
	DerivedFrom []*FHIRReference `json:"derivedFrom,omitempty"`

	// Some observations have multiple component observations.  These component observations are expressed as separate code value pairs that share the same attributes.  Examples include systolic and diastolic component observations for blood pressure measurement and multiple component observations for genetics observations.
	Component []*FHIRObservationComponent `json:"component,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMeta `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

func (o *FHIRObservation) GetBasedOn() *FHIRReference {
	if o.BasedOn != nil {
		for _, basedOn := range o.BasedOn {
			return &FHIRReference{
				Reference: basedOn.Reference,
				ID:        basedOn.ID,
				Display:   basedOn.Display,
			}
		}
	}

	return nil
}

func (o *FHIRObservation) Coding() *FHIRCoding {
	if o.Code != nil {
		for _, code := range o.Code.Coding {
			return &FHIRCoding{
				System:  code.System,
				Code:    code.Code,
				Display: code.Display,
			}
		}
	}

	return nil
}

// FHIRObservationComponent definition: measurements and simple assertions made about a patient, device or other subject.
type FHIRObservationComponent struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Describes what was observed. Sometimes this is called the observation "code".
	Code FHIRCodeableConcept `json:"code,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueQuantity *FHIRQuantity `json:"valueQuantity,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueCodeableConcept *scalarutils.Code `json:"valueCodeableConcept,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueString *string `json:"valueString,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueBoolean *bool `json:"valueBoolean,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueInteger *string `json:"valueInteger,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueRange *FHIRRange `json:"valueRange,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueRatio *FHIRRatio `json:"valueRatio,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueSampledData *FHIRSampledData `json:"valueSampledData,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueTime *time.Time `json:"valueTime,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueDateTime *scalarutils.Date `json:"valueDateTime,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValuePeriod *FHIRPeriod `json:"valuePeriod,omitempty"`

	// Provides a reason why the expected value in the element Observation.component.value[x] is missing.
	DataAbsentReason *FHIRCodeableConcept `json:"dataAbsentReason,omitempty"`

	// A categorical assessment of an observation value.  For example, high, low, normal.
	Interpretation []*FHIRCodeableConcept `json:"interpretation,omitempty"`

	// Guidance on how to interpret the value by comparison to a normal or recommended range.
	ReferenceRange []*FHIRObservationReferencerange `json:"referenceRange,omitempty"`
}

// FHIRObservationComponentInput is the input type for ObservationComponent
type FHIRObservationComponentInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Describes what was observed. Sometimes this is called the observation "code".
	Code FHIRCodeableConceptInput `json:"code,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueQuantity *FHIRQuantityInput `json:"valueQuantity,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueCodeableConcept *scalarutils.Code `json:"valueCodeableConcept,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueString *string `json:"valueString,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueBoolean *bool `json:"valueBoolean,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueInteger *string `json:"valueInteger,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueRange *FHIRRangeInput `json:"valueRange,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueRatio *FHIRRatioInput `json:"valueRatio,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueSampledData *FHIRSampledDataInput `json:"valueSampledData,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueTime *time.Time `json:"valueTime,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueDateTime *scalarutils.Date `json:"valueDateTime,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValuePeriod *FHIRPeriodInput `json:"valuePeriod,omitempty"`

	// Provides a reason why the expected value in the element Observation.component.value[x] is missing.
	DataAbsentReason *FHIRCodeableConceptInput `json:"dataAbsentReason,omitempty"`

	// A categorical assessment of an observation value.  For example, high, low, normal.
	Interpretation []*FHIRCodeableConceptInput `json:"interpretation,omitempty"`

	// Guidance on how to interpret the value by comparison to a normal or recommended range.
	ReferenceRange []*FHIRObservationReferencerangeInput `json:"referenceRange,omitempty"`
}

// FHIRObservationInput is the input type for Observation
type FHIRObservationInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	Text *FHIRNarrative `json:"text,omitempty"`

	// A unique identifier assigned to this observation.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	// A plan, proposal or order that is fulfilled in whole or in part by this event.  For example, a MedicationRequest may require a patient to have laboratory test performed before  it is dispensed.
	BasedOn []*FHIRReferenceInput `json:"basedOn,omitempty"`

	// A larger event of which this particular Observation is a component or step.  For example,  an observation as part of a procedure.
	PartOf []*FHIRReferenceInput `json:"partOf,omitempty"`

	// The status of the result value.
	Status *ObservationStatusEnum `json:"status,omitempty"`

	// A code that classifies the general type of observation being made.
	Category []*FHIRCodeableConceptInput `json:"category,omitempty"`

	// Describes what was observed. Sometimes this is called the observation "name".
	Code *FHIRCodeableConceptInput `json:"code,omitempty"`

	// The patient, or group of patients, location, or device this observation is about and into whose record the observation is placed. If the actual focus of the observation is different from the subject (or a sample of, part, or region of the subject), the `focus` element or the `code` itself specifies the actual focus of the observation.
	Subject *FHIRReferenceInput `json:"subject,omitempty"`

	// The actual focus of an observation when it is not the patient of record representing something or someone associated with the patient such as a spouse, parent, fetus, or donor. For example, fetus observations in a mother's record.  The focus of an observation could also be an existing condition,  an intervention, the subject's diet,  another observation of the subject,  or a body structure such as tumor or implanted device.   An example use case would be using the Observation resource to capture whether the mother is trained to change her child's tracheostomy tube. In this example, the child is the patient of record and the mother is the focus.
	Focus []*FHIRReferenceInput `json:"focus,omitempty"`

	// The healthcare event  (e.g. a patient and healthcare provider interaction) during which this observation is made.
	Encounter *FHIRReferenceInput `json:"encounter,omitempty"`

	// The time or time-period the observed value is asserted as being true. For biological subjects - e.g. human patients - this is usually called the "physiologically relevant time". This is usually either the time of the procedure or of specimen collection, but very often the source of the date/time is not known, only the date/time itself.
	EffectiveDateTime *scalarutils.DateTime `json:"effectiveDateTime,omitempty"`

	// The time or time-period the observed value is asserted as being true. For biological subjects - e.g. human patients - this is usually called the "physiologically relevant time". This is usually either the time of the procedure or of specimen collection, but very often the source of the date/time is not known, only the date/time itself.
	EffectivePeriod *FHIRPeriodInput `json:"effectivePeriod,omitempty"`

	// The time or time-period the observed value is asserted as being true. For biological subjects - e.g. human patients - this is usually called the "physiologically relevant time". This is usually either the time of the procedure or of specimen collection, but very often the source of the date/time is not known, only the date/time itself.
	EffectiveTiming *FHIRTimingInput `json:"effectiveTiming,omitempty"`

	// The time or time-period the observed value is asserted as being true. For biological subjects - e.g. human patients - this is usually called the "physiologically relevant time". This is usually either the time of the procedure or of specimen collection, but very often the source of the date/time is not known, only the date/time itself.
	EffectiveInstant *scalarutils.Instant `json:"effectiveInstant,omitempty"`

	// The date and time this version of the observation was made available to providers, typically after the results have been reviewed and verified.
	Issued *scalarutils.Instant `json:"issued,omitempty"`

	// Who was responsible for asserting the observed value as "true".
	Performer []*FHIRReferenceInput `json:"performer,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueQuantity *FHIRQuantityInput `json:"valueQuantity,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueCodeableConcept *FHIRCodeableConceptInput `json:"valueCodeableConcept,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueString *string `json:"valueString,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueBoolean *bool `json:"valueBoolean,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueInteger *string `json:"valueInteger,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueRange *FHIRRangeInput `json:"valueRange,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueRatio *FHIRRatioInput `json:"valueRatio,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueSampledData *FHIRSampledDataInput `json:"valueSampledData,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueTime *time.Time `json:"valueTime,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValueDateTime *scalarutils.Date `json:"valueDateTime,omitempty"`

	// The information determined as a result of making the observation, if the information has a simple value.
	ValuePeriod *FHIRPeriodInput `json:"valuePeriod,omitempty"`

	// Provides a reason why the expected value in the element Observation.value[x] is missing.
	DataAbsentReason *FHIRCodeableConceptInput `json:"dataAbsentReason,omitempty"`

	// A categorical assessment of an observation value.  For example, high, low, normal.
	Interpretation []*FHIRCodeableConceptInput `json:"interpretation,omitempty"`

	// Comments about the observation or the results.
	Note []*FHIRAnnotationInput `json:"note,omitempty"`

	// Indicates the site on the subject's body where the observation was made (i.e. the target site).
	BodySite *FHIRCodeableConceptInput `json:"bodySite,omitempty"`

	// Indicates the mechanism used to perform the observation.
	Method *FHIRCodeableConceptInput `json:"method,omitempty"`

	// The specimen that was used when this observation was made.
	Specimen *FHIRReferenceInput `json:"specimen,omitempty"`

	// The device used to generate the observation data.
	Device *FHIRReferenceInput `json:"device,omitempty"`

	// Guidance on how to interpret the value by comparison to a normal or recommended range.  Multiple reference ranges are interpreted as an "OR".   In other words, to represent two distinct target populations, two `referenceRange` elements would be used.
	ReferenceRange []*FHIRObservationReferencerangeInput `json:"referenceRange,omitempty"`

	// This observation is a group observation (e.g. a battery, a panel of tests, a set of vital sign measurements) that includes the target as a member of the group.
	HasMember []*FHIRReferenceInput `json:"hasMember,omitempty"`

	// The target resource that represents a measurement from which this observation value is derived. For example, a calculated anion gap or a fetal measurement based on an ultrasound image.
	DerivedFrom []*FHIRReferenceInput `json:"derivedFrom,omitempty"`

	// Some observations have multiple component observations.  These component observations are expressed as separate code value pairs that share the same attributes.  Examples include systolic and diastolic component observations for blood pressure measurement and multiple component observations for genetics observations.
	Component []*FHIRObservationComponentInput `json:"component,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

// FHIRObservationReferencerange definition: measurements and simple assertions made about a patient, device or other subject.
type FHIRObservationReferencerange struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The value of the low bound of the reference range.  The low bound of the reference range endpoint is inclusive of the value (e.g.  reference range is >=5 - <=9). If the low bound is omitted,  it is assumed to be meaningless (e.g. reference range is <=2.3).
	Low *FHIRQuantity `json:"low,omitempty"`

	// The value of the high bound of the reference range.  The high bound of the reference range endpoint is inclusive of the value (e.g.  reference range is >=5 - <=9). If the high bound is omitted,  it is assumed to be meaningless (e.g. reference range is >= 2.3).
	High *FHIRQuantity `json:"high,omitempty"`

	// Codes to indicate the what part of the targeted reference population it applies to. For example, the normal or therapeutic range.
	Type *FHIRCodeableConcept `json:"type,omitempty"`

	// Codes to indicate the target population this reference range applies to.  For example, a reference range may be based on the normal population or a particular sex or race.  Multiple `appliesTo`  are interpreted as an "AND" of the target populations.  For example, to represent a target population of African American females, both a code of female and a code for African American would be used.
	AppliesTo []*FHIRCodeableConcept `json:"appliesTo,omitempty"`

	// The age at which this reference range is applicable. This is a neonatal age (e.g. number of weeks at term) if the meaning says so.
	Age *FHIRRange `json:"age,omitempty"`

	// Text based reference range in an observation which may be used when a quantitative range is not appropriate for an observation.  An example would be a reference value of "Negative" or a list or table of "normals".
	Text *string `json:"text,omitempty"`
}

// FHIRObservationReferencerangeInput is the input type for ObservationReferencerange
type FHIRObservationReferencerangeInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The value of the low bound of the reference range.  The low bound of the reference range endpoint is inclusive of the value (e.g.  reference range is >=5 - <=9). If the low bound is omitted,  it is assumed to be meaningless (e.g. reference range is <=2.3).
	Low *FHIRQuantityInput `json:"low,omitempty"`

	// The value of the high bound of the reference range.  The high bound of the reference range endpoint is inclusive of the value (e.g.  reference range is >=5 - <=9). If the high bound is omitted,  it is assumed to be meaningless (e.g. reference range is >= 2.3).
	High *FHIRQuantityInput `json:"high,omitempty"`

	// Codes to indicate the what part of the targeted reference population it applies to. For example, the normal or therapeutic range.
	Type *FHIRCodeableConceptInput `json:"type,omitempty"`

	// Codes to indicate the target population this reference range applies to.  For example, a reference range may be based on the normal population or a particular sex or race.  Multiple `appliesTo`  are interpreted as an "AND" of the target populations.  For example, to represent a target population of African American females, both a code of female and a code for African American would be used.
	AppliesTo []*FHIRCodeableConceptInput `json:"appliesTo,omitempty"`

	// The age at which this reference range is applicable. This is a neonatal age (e.g. number of weeks at term) if the meaning says so.
	Age *FHIRRangeInput `json:"age,omitempty"`

	// Text based reference range in an observation which may be used when a quantitative range is not appropriate for an observation.  An example would be a reference value of "Negative" or a list or table of "normals".
	Text *string `json:"text,omitempty"`
}

// PagedFHIRAllergy is an allergy's paginaton dataclass
type PagedFHIRObservations struct {
	BundleID        string
	Observations    []FHIRObservation
	HasNextPage     bool
	NextPageURL     string
	HasPreviousPage bool
	PreviousPageURL string
	TotalCount      int
}

// FHIRObservationRelayPayload is used to return single instance of Observation
type FHIRObservationRelayPayload struct {
	Resource *FHIRObservation `json:"resource,omitempty"`
}
