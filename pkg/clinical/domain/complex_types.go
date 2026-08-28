package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/savannahghi/scalarutils"
)

// FHIRAddress definition: an address expressed using postal conventions (as opposed to gps or other location definition formats).  this data type may be used to convey addresses for use in delivering mail as well as for visiting locations which might not be valid for mail delivery.  there are a variety of postal address formats defined around the world.
type FHIRAddress struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The purpose of this address.
	Use *AddressUseEnum `json:"use,omitempty"`

	// Distinguishes between physical addresses (those you can visit) and mailing addresses (e.g. PO Boxes and care-of addresses). Most addresses are both.
	Type *AddressTypeEnum `json:"type,omitempty"`

	// Specifies the entire address as it should be displayed e.g. on a postal label. This may be provided instead of or as well as the specific parts.
	Text string `json:"text,omitempty"`

	// This component contains the house number, apartment number, street name, street direction,  P.O. Box number, delivery hints, and similar address information.
	Line []*string `json:"line,omitempty"`

	// The name of the city, town, suburb, village or other community or delivery center.
	City *string `json:"city,omitempty"`

	// The name of the administrative area (county).
	District *string `json:"district,omitempty"`

	// Sub-unit of a country with limited sovereignty in a federally organized country. A code may be used if codes are in common use (e.g. US 2 letter state codes).
	State *string `json:"state,omitempty"`

	// A postal code designating a region defined by the postal service.
	PostalCode *scalarutils.Code `json:"postalCode,omitempty"`

	// Country - a nation as commonly understood or generally accepted.
	Country *string `json:"country,omitempty"`

	// Time period when address was/is in use.
	Period *FHIRPeriod `json:"period,omitempty"`
}

// FHIRAddressInput is the input type for Address
type FHIRAddressInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The purpose of this address.
	Use *AddressUseEnum `json:"use,omitempty"`

	// Distinguishes between physical addresses (those you can visit) and mailing addresses (e.g. PO Boxes and care-of addresses). Most addresses are both.
	Type *AddressTypeEnum `json:"type,omitempty"`

	// Specifies the entire address as it should be displayed e.g. on a postal label. This may be provided instead of or as well as the specific parts.
	Text string `json:"text,omitempty"`

	// This component contains the house number, apartment number, street name, street direction,  P.O. Box number, delivery hints, and similar address information.
	Line []*string `json:"line,omitempty"`

	// The name of the city, town, suburb, village or other community or delivery center.
	City *string `json:"city,omitempty"`

	// The name of the administrative area (county).
	District *string `json:"district,omitempty"`

	// Sub-unit of a country with limited sovereignty in a federally organized country. A code may be used if codes are in common use (e.g. US 2 letter state codes).
	State *string `json:"state,omitempty"`

	// A postal code designating a region defined by the postal service.
	PostalCode *scalarutils.Code `json:"postalCode,omitempty"`

	// Country - a nation as commonly understood or generally accepted.
	Country *string `json:"country,omitempty"`

	// Time period when address was/is in use.
	Period *FHIRPeriodInput `json:"period,omitempty"`
}

// FHIRAge definition: a duration of time during which an organism (or a process) has existed.
type FHIRAge struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The value of the measured amount. The value includes an implicit precision in the presentation of the value.
	Value *scalarutils.Decimal `json:"value,omitempty"`

	// How the value should be understood and represented - whether the actual value is greater or less than the stated value due to measurement issues; e.g. if the comparator is "<" , then the real value is < stated value.
	Comparator *AgeComparatorEnum `json:"comparator,omitempty"`

	// A human-readable form of the unit.
	Unit *string `json:"unit,omitempty"`

	// The identification of the system that provides the coded form of the unit.
	System *scalarutils.URI `json:"system,omitempty"`

	// A computer processable form of the unit in some unit representation system.
	Code *scalarutils.Code `json:"code,omitempty"`
}

// FHIRAgeInput is the input type for Age
type FHIRAgeInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The value of the measured amount. The value includes an implicit precision in the presentation of the value.
	Value *scalarutils.Decimal `json:"value,omitempty"`

	// How the value should be understood and represented - whether the actual value is greater or less than the stated value due to measurement issues; e.g. if the comparator is "<" , then the real value is < stated value.
	Comparator *AgeComparatorEnum `json:"comparator,omitempty"`

	// A human-readable form of the unit.
	Unit *string `json:"unit,omitempty"`

	// The identification of the system that provides the coded form of the unit.
	System *scalarutils.URI `json:"system,omitempty"`

	// A computer processable form of the unit in some unit representation system.
	Code *scalarutils.Code `json:"code,omitempty"`
}

// FHIRAnnotation definition: a  text note which also  contains information about who made the statement and when.
type FHIRAnnotation struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The individual responsible for making the annotation.
	AuthorReference *FHIRReference `json:"authorReference,omitempty"`

	// The individual responsible for making the annotation.
	AuthorString *string `json:"authorString,omitempty"`

	// Indicates when this particular annotation was made.
	Time *scalarutils.DateTime `json:"time,omitempty"`

	// The text of the annotation in markdown format.
	Text *scalarutils.Markdown `json:"text,omitempty"`
}

// FHIRAnnotationInput is the input type for Annotation
type FHIRAnnotationInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The individual responsible for making the annotation.
	AuthorReference *FHIRReferenceInput `json:"authorReference,omitempty"`

	// The individual responsible for making the annotation.
	AuthorString *string `json:"authorString,omitempty"`

	// Indicates when this particular annotation was made.
	Time *scalarutils.DateTime `json:"time,omitempty"`

	// The text of the annotation in markdown format.
	Text *scalarutils.Markdown `json:"text,omitempty"`
}

// FHIRAttachment definition: for referring to data content defined in other formats.
type FHIRAttachment struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Identifies the type of the data in the attachment and allows a method to be chosen to interpret or render the data. Includes mime type parameters such as charset where appropriate.
	ContentType *scalarutils.Code `json:"contentType,omitempty"`

	// The human language of the content. The value can be any valid value according to BCP 47.
	Language *scalarutils.Code `json:"language,omitempty"`

	// The actual data of the attachment - a sequence of bytes, base64 encoded.
	Data *scalarutils.Base64Binary `json:"data,omitempty"`

	// A location where the data can be accessed.
	URL *scalarutils.URL `json:"url,omitempty"`

	// The number of bytes of data that make up this attachment (before base64 encoding, if that is done).
	Size *int `json:"size,omitempty"`

	// The calculated hash of the data using SHA-1. Represented using base64.
	Hash *scalarutils.Base64Binary `json:"hash,omitempty"`

	// A label or set of text to display in place of the data.
	Title *string `json:"title,omitempty"`

	// The date that the attachment was first created.
	Creation *scalarutils.DateTime `json:"creation,omitempty"`
}

// FHIRAttachmentInput is the input type for Attachment
type FHIRAttachmentInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Identifies the type of the data in the attachment and allows a method to be chosen to interpret or render the data. Includes mime type parameters such as charset where appropriate.
	ContentType *scalarutils.Code `json:"contentType,omitempty"`

	// The human language of the content. The value can be any valid value according to BCP 47.
	Language *scalarutils.Code `json:"language,omitempty"`

	// The actual data of the attachment - a sequence of bytes, base64 encoded.
	Data *scalarutils.Base64Binary `json:"data,omitempty"`

	// A location where the data can be accessed.
	URL *scalarutils.URL `json:"url,omitempty"`

	// The number of bytes of data that make up this attachment (before base64 encoding, if that is done).
	Size *int `json:"size,omitempty"`

	// The calculated hash of the data using SHA-1. Represented using base64.
	Hash *scalarutils.Base64Binary `json:"hash,omitempty"`

	// A label or set of text to display in place of the data.
	Title *string `json:"title,omitempty"`

	// The date that the attachment was first created.
	Creation *scalarutils.DateTime `json:"creation,omitempty"`
}

// FHIRCodeableConcept definition: a concept that may be defined by a formal reference to a terminology or ontology or may be provided by text.
type FHIRCodeableConcept struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A reference to a code defined by a terminology system.
	Coding []*FHIRCoding `json:"coding,omitempty"`

	// A human language representation of the concept as seen/selected/uttered by the user who entered the data and/or which represents the intended meaning of the user.
	Text string `json:"text,omitempty"`
}

// FHIRCodeableReference is documented here http://hl7.org/fhir/StructureDefinition/CodeableReference
type FHIRCodeableReference struct {
	ID        *string              `json:"id,omitempty"`
	Concept   *FHIRCodeableConcept `json:"concept,omitempty"`
	Reference *FHIRReference       `json:"reference,omitempty"`
}

// FHIRCodeableReference is documented here http://hl7.org/fhir/StructureDefinition/CodeableReference
type FHIRCodeableReferenceInput struct {
	ID        *string                   `json:"id,omitempty"`
	Concept   *FHIRCodeableConceptInput `json:"concept,omitempty"`
	Reference *FHIRReferenceInput       `json:"reference,omitempty"`
}

// FHIRCodeableConceptInput is the input type for CodeableConcept
type FHIRCodeableConceptInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A reference to a code defined by a terminology system.
	Coding []*FHIRCodingInput `json:"coding,omitempty"`

	// A human language representation of the concept as seen/selected/uttered by the user who entered the data and/or which represents the intended meaning of the user.
	Text string `json:"text,omitempty"`
}

// FHIRCoding definition: a reference to a code defined by a terminology system.
type FHIRCoding struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The identification of the code system that defines the meaning of the symbol in the code.
	System *scalarutils.URI `json:"system,omitempty" swaggertype:"primitive,string"`

	// The version of the code system which was used when choosing this code. Note that a well-maintained code system does not need the version reported, because the meaning of codes is consistent across versions. However this cannot consistently be assured, and when the meaning is not guaranteed to be consistent, the version SHOULD be exchanged.
	Version *string `json:"version,omitempty"`

	// A symbol in syntax defined by the system. The symbol may be a predefined code or an expression in a syntax defined by the coding system (e.g. post-coordination).
	Code *scalarutils.Code `json:"code,omitempty" swaggertype:"primitive,string"`

	// A representation of the meaning of the code in the system, following the rules of the system.
	Display string `json:"display,omitempty"`

	// Indicates that this coding was chosen by a user directly - e.g. off a pick list of available items (codes or displays).
	UserSelected *bool `json:"userSelected,omitempty"`
}

// ToString returns the Display field of the Coding struct as a string
func (c *FHIRCoding) ToString() string {
	return c.Display
}

// FHIRCodingInput is the input type for Coding
type FHIRCodingInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID string `json:"id,omitempty" mapstructure:"id"`

	// The identification of the code system that defines the meaning of the symbol in the code.
	System *scalarutils.URI `json:"system,omitempty" mapstructure:"system"`

	// The version of the code system which was used when choosing this code. Note that a well-maintained code system does not need the version reported, because the meaning of codes is consistent across versions. However this cannot consistently be assured, and when the meaning is not guaranteed to be consistent, the version SHOULD be exchanged.
	Version *string `json:"version,omitempty" mapstructure:"version"`

	// A symbol in syntax defined by the system. The symbol may be a predefined code or an expression in a syntax defined by the coding system (e.g. post-coordination).
	Code scalarutils.Code `json:"code,omitempty" mapstructure:"code"`

	// A representation of the meaning of the code in the system, following the rules of the system.
	Display string `json:"display,omitempty" mapstructure:"display"`

	// Indicates that this coding was chosen by a user directly - e.g. off a pick list of available items (codes or displays).
	UserSelected *bool `json:"userSelected,omitempty" mapstructure:"userSelected"`
}

// FHIRContactPoint definition: details for all kinds of technology mediated contact points for a person or organization, including telephone, email, etc.
type FHIRContactPoint struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Telecommunications form for contact point - what communications system is required to make use of the contact.
	System *ContactPointSystemEnum `json:"system,omitempty"`

	// The actual contact point details, in a form that is meaningful to the designated communication system (i.e. phone number or email address).
	Value *string `json:"value,omitempty"`

	// Identifies the purpose for the contact point.
	Use *ContactPointUseEnum `json:"use,omitempty"`

	// Specifies a preferred order in which to use a set of contacts. ContactPoints with lower rank values are more preferred than those with higher rank values.
	Rank *int64 `json:"rank,omitempty"`

	// Time period when the contact point was/is in use.
	Period *FHIRPeriod `json:"period,omitempty"`
}

// FHIRContactPointInput is the input type for ContactPoint
type FHIRContactPointInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Telecommunications form for contact point - what communications system is required to make use of the contact.
	System *ContactPointSystemEnum `json:"system,omitempty"`

	// The actual contact point details, in a form that is meaningful to the designated communication system (i.e. phone number or email address).
	Value *string `json:"value,omitempty"`

	// Identifies the purpose for the contact point.
	Use *ContactPointUseEnum `json:"use,omitempty"`

	// Specifies a preferred order in which to use a set of contacts. ContactPoints with lower rank values are more preferred than those with higher rank values.
	Rank *int64 `json:"rank,omitempty"`

	// Time period when the contact point was/is in use.
	Period *FHIRPeriodInput `json:"period,omitempty"`

	Extension []Extension `json:"extension,omitempty"`
}

// FHIRDosage definition: indicates how the medication is/was taken or should be taken by the patient.
type FHIRDosage struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Indicates the order in which the dosage instructions should be applied or interpreted.
	Sequence *json.Number `json:"sequence,omitempty"`

	// Free text dosage instructions e.g. SIG.
	Text *string `json:"text,omitempty"`

	// Supplemental instructions to the patient on how to take the medication  (e.g. "with meals" or"take half to one hour before food") or warnings for the patient about the medication (e.g. "may cause drowsiness" or "avoid exposure of skin to direct sunlight or sunlamps").
	AdditionalInstruction []*FHIRCodeableConcept `json:"additionalInstruction,omitempty"`

	// Instructions in terms that are understood by the patient or consumer.
	PatientInstruction *string `json:"patientInstruction,omitempty"`

	// When medication should be administered.
	Timing *FHIRTiming `json:"timing,omitempty"`

	// Indicates whether the Medication is only taken when needed within a specific dosing schedule (Boolean option), or it indicates the precondition for taking the Medication (CodeableConcept).
	AsNeeded *bool `json:"asNeeded,omitempty"`

	// Indicates whether the Medication is only taken when needed within a specific dosing schedule (Boolean option), or it indicates the precondition for taking the Medication (CodeableConcept).
	AsNeededCodeableConcept *scalarutils.Code `json:"asNeededCodeableConcept,omitempty"`

	// Body site to administer to.
	Site *FHIRCodeableConcept `json:"site,omitempty"`

	// How drug should enter body.
	Route *FHIRCodeableConcept `json:"route,omitempty"`

	// Technique for administering medication.
	Method *FHIRCodeableConcept `json:"method,omitempty"`

	// The amount of medication administered.
	DoseAndRate []*FHIRDosageDoseandrate `json:"doseAndRate,omitempty"`

	// Upper limit on medication per unit of time.
	MaxDosePerPeriod *FHIRRatio `json:"maxDosePerPeriod,omitempty"`

	// Upper limit on medication per administration.
	MaxDosePerAdministration *FHIRQuantity `json:"maxDosePerAdministration,omitempty"`

	// Upper limit on medication per lifetime of the patient.
	MaxDosePerLifetime *FHIRQuantity `json:"maxDosePerLifetime,omitempty"`
}

// FHIRDosageDoseandrate definition: indicates how the medication is/was taken or should be taken by the patient.
type FHIRDosageDoseandrate struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The kind of dose or rate specified, for example, ordered or calculated.
	Type *FHIRCodeableConcept `json:"type,omitempty"`

	// Amount of medication per dose.
	DoseRange *FHIRRange `json:"doseRange,omitempty"`

	// Amount of medication per dose.
	DoseQuantity *FHIRQuantity `json:"doseQuantity,omitempty"`

	// Amount of medication per unit of time.
	RateRatio *FHIRRatio `json:"rateRatio,omitempty"`

	// Amount of medication per unit of time.
	RateRange *FHIRRange `json:"rateRange,omitempty"`

	// Amount of medication per unit of time.
	RateQuantity *FHIRQuantity `json:"rateQuantity,omitempty"`
}

// FHIRDosageDoseandrateInput is the input type for DosageDoseandrate
type FHIRDosageDoseandrateInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The kind of dose or rate specified, for example, ordered or calculated.
	Type *FHIRCodeableConceptInput `json:"type,omitempty"`

	// Amount of medication per dose.
	DoseRange *FHIRRangeInput `json:"doseRange,omitempty"`

	// Amount of medication per dose.
	DoseQuantity *FHIRQuantityInput `json:"doseQuantity,omitempty"`

	// Amount of medication per unit of time.
	RateRatio *FHIRRatioInput `json:"rateRatio,omitempty"`

	// Amount of medication per unit of time.
	RateRange *FHIRRangeInput `json:"rateRange,omitempty"`

	// Amount of medication per unit of time.
	RateQuantity *FHIRQuantityInput `json:"rateQuantity,omitempty"`
}

// FHIRDosageInput is the input type for Dosage
type FHIRDosageInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Indicates the order in which the dosage instructions should be applied or interpreted.
	Sequence *string `json:"sequence,omitempty"`

	// Free text dosage instructions e.g. SIG.
	Text *string `json:"text,omitempty"`

	// Supplemental instructions to the patient on how to take the medication  (e.g. "with meals" or"take half to one hour before food") or warnings for the patient about the medication (e.g. "may cause drowsiness" or "avoid exposure of skin to direct sunlight or sunlamps").
	AdditionalInstruction []*FHIRCodeableConceptInput `json:"additionalInstruction,omitempty"`

	// Instructions in terms that are understood by the patient or consumer.
	PatientInstruction *string `json:"patientInstruction,omitempty"`

	// When medication should be administered.
	Timing *FHIRTimingInput `json:"timing,omitempty"`

	// Indicates whether the Medication is only taken when needed within a specific dosing schedule (Boolean option), or it indicates the precondition for taking the Medication (CodeableConcept).
	AsNeeded *bool `json:"asNeeded,omitempty"`

	// Indicates whether the Medication is only taken when needed within a specific dosing schedule (Boolean option), or it indicates the precondition for taking the Medication (CodeableConcept).
	AsNeededCodeableConcept *scalarutils.Code `json:"asNeededCodeableConcept,omitempty"`

	// Body site to administer to.
	Site *FHIRCodeableConceptInput `json:"site,omitempty"`

	// How drug should enter body.
	Route *FHIRCodeableConceptInput `json:"route,omitempty"`

	// Technique for administering medication.
	Method *FHIRCodeableConceptInput `json:"method,omitempty"`

	// The amount of medication administered.
	DoseAndRate []*FHIRDosageDoseandrateInput `json:"doseAndRate,omitempty"`

	// Upper limit on medication per unit of time.
	MaxDosePerPeriod *FHIRRatioInput `json:"maxDosePerPeriod,omitempty"`

	// Upper limit on medication per administration.
	MaxDosePerAdministration *FHIRQuantityInput `json:"maxDosePerAdministration,omitempty"`

	// Upper limit on medication per lifetime of the patient.
	MaxDosePerLifetime *FHIRQuantityInput `json:"maxDosePerLifetime,omitempty"`
}

// FHIRDuration definition: a length of time.
type FHIRDuration struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The value of the measured amount. The value includes an implicit precision in the presentation of the value.
	Value *scalarutils.Decimal `json:"value,omitempty"`

	// How the value should be understood and represented - whether the actual value is greater or less than the stated value due to measurement issues; e.g. if the comparator is "<" , then the real value is < stated value.
	Comparator *DurationComparatorEnum `json:"comparator,omitempty"`

	// A human-readable form of the unit.
	Unit *string `json:"unit,omitempty"`

	// The identification of the system that provides the coded form of the unit.
	System *scalarutils.URI `json:"system,omitempty"`

	// A computer processable form of the unit in some unit representation system.
	Code *scalarutils.Code `json:"code,omitempty"`
}

// FHIRDurationInput is the input type for Duration
type FHIRDurationInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The value of the measured amount. The value includes an implicit precision in the presentation of the value.
	Value *scalarutils.Decimal `json:"value,omitempty"`

	// How the value should be understood and represented - whether the actual value is greater or less than the stated value due to measurement issues; e.g. if the comparator is "<" , then the real value is < stated value.
	Comparator *DurationComparatorEnum `json:"comparator,omitempty"`

	// A human-readable form of the unit.
	Unit *string `json:"unit,omitempty"`

	// The identification of the system that provides the coded form of the unit.
	System *scalarutils.URI `json:"system,omitempty"`

	// A computer processable form of the unit in some unit representation system.
	Code *scalarutils.Code `json:"code,omitempty"`
}

// FHIRHumanName definition: a human's name with the ability to identify parts and usage.
type FHIRHumanName struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Identifies the purpose for this name.
	Use HumanNameUseEnum `json:"use,omitempty"`

	// Specifies the entire name as it should be displayed e.g. on an application UI. This may be provided instead of or as well as the specific parts.
	Text string `json:"text,omitempty"`

	// The part of a name that links to the genealogy. In some cultures (e.g. Eritrea) the family name of a son is the first name of his father.
	Family *string `json:"family,omitempty"`

	// Given name.
	Given []*string `json:"given,omitempty"`

	// Part of the name that is acquired as a title due to academic, legal, employment or nobility status, etc. and that appears at the start of the name.
	Prefix []*string `json:"prefix,omitempty"`

	// Part of the name that is acquired as a title due to academic, legal, employment or nobility status, etc. and that appears at the end of the name.
	Suffix []*string `json:"suffix,omitempty"`

	// Indicates the period of time when this name was valid for the named person.
	Period *FHIRPeriod `json:"period,omitempty"`
}

// FHIRHumanNameInput is the input type for HumanName
type FHIRHumanNameInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Identifies the purpose for this name.
	Use HumanNameUseEnum `json:"use,omitempty"`

	// Specifies the entire name as it should be displayed e.g. on an application UI. This may be provided instead of or as well as the specific parts.
	Text string `json:"text,omitempty"`

	// The part of a name that links to the genealogy. In some cultures (e.g. Eritrea) the family name of a son is the first name of his father.
	Family string `json:"family,omitempty"`

	// Given name.
	Given []string `json:"given,omitempty"`

	// Part of the name that is acquired as a title due to academic, legal, employment or nobility status, etc. and that appears at the start of the name.
	Prefix []*string `json:"prefix,omitempty"`

	// Part of the name that is acquired as a title due to academic, legal, employment or nobility status, etc. and that appears at the end of the name.
	Suffix []*string `json:"suffix,omitempty"`

	// Indicates the period of time when this name was valid for the named person.
	Period *FHIRPeriodInput `json:"period,omitempty"`
}

// FHIRIdentifier definition: an identifier - identifies some entity uniquely and unambiguously. typically this is used for business identifiers.
type FHIRIdentifier struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The purpose of this identifier.
	Use IdentifierUseEnum `json:"use,omitempty"`

	// Allows for the inclusion of data that are supported by FHIR
	Extension []FHIRExtension `json:"extension,omitempty"`

	// A coded type for the identifier that can be used to determine which identifier to use for a specific purpose.
	Type FHIRCodeableConcept `json:"type,omitempty"`

	// Establishes the namespace for the value - that is, a URL that describes a set values that are unique.
	System *scalarutils.URI `json:"system,omitempty" swaggertype:"primitive,string"`

	// The portion of the identifier typically relevant to the user and which is unique within the context of the system.
	Value string `json:"value,omitempty"`

	// Time period during which identifier is/was valid for use.
	Period *FHIRPeriod `json:"period,omitempty"`

	// Organization that issued/manages the identifier.
	Assigner *FHIRReference `json:"assigner,omitempty"`
}

// FHIRIdentifierInput is the input type for Identifier
type FHIRIdentifierInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The purpose of this identifier.
	Use IdentifierUseEnum `json:"use,omitempty"`

	// A coded type for the identifier that can be used to determine which identifier to use for a specific purpose.
	Type *FHIRCodeableConceptInput `json:"type,omitempty"`

	// Establishes the namespace for the value - that is, a URL that describes a set values that are unique.
	System *scalarutils.URI `json:"system,omitempty"`

	// The portion of the identifier typically relevant to the user and which is unique within the context of the system.
	Value string `json:"value,omitempty"`

	// Time period during which identifier is/was valid for use.
	Period *FHIRPeriodInput `json:"period,omitempty"`

	// Organization that issued/manages the identifier.
	Assigner *FHIRReferenceInput `json:"assigner,omitempty"`
}

func (id *FHIRIdentifierInput) AddAssigner(organizationID string) {
	reference := "Organization/" + organizationID

	id.Assigner = &FHIRReferenceInput{
		Reference: &reference,
		Display:   organizationID,
	}
}

// FHIRNarrative definition: a human-readable summary of the resource conveying the essential clinical and business information for the resource.
type FHIRNarrative struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The status of the narrative - whether it's entirely generated (from just the defined data or the extensions too), or whether a human authored it and it may contain additional data.
	Status *NarrativeStatusEnum `json:"status,omitempty"`

	// The actual narrative content, a stripped down version of XHTML.
	Div scalarutils.XHTML `json:"div,omitempty" swaggertype:"primitive,string"`
}

// FHIRNarrativeInput is the input type for Narrative
type FHIRNarrativeInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The status of the narrative - whether it's entirely generated (from just the defined data or the extensions too), or whether a human authored it and it may contain additional data.
	Status *NarrativeStatusEnum `json:"status,omitempty"`

	// The actual narrative content, a stripped down version of XHTML.
	Div scalarutils.XHTML `json:"div,omitempty"`
}

// FHIRPeriod definition: a time period defined by a start and end date and optionally time.
type FHIRPeriod struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The start of the period. The boundary is inclusive.
	Start scalarutils.DateTime `json:"start,omitempty"`

	// The end of the period. If the end of the period is missing, it means no end was known or planned at the time the instance was created. The start may be in the past, and the end date in the future, which means that period is expected/planned to end at that time.
	End scalarutils.DateTime `json:"end,omitempty"`
}

// FHIRPeriodInput is the input type for Period
type FHIRPeriodInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The start of the period. The boundary is inclusive.
	Start *scalarutils.DateTime `json:"start,omitempty"`

	// The end of the period. If the end of the period is missing, it means no end was known or planned at the time the instance was created. The start may be in the past, and the end date in the future, which means that period is expected/planned to end at that time.
	End *scalarutils.DateTime `json:"end,omitempty"`
}

// FHIRQuantity definition: a measured amount (or an amount that can potentially be measured). note that measured amounts include amounts that are not precisely quantified, including amounts involving arbitrary units and floating currencies.
type FHIRQuantity struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The value of the measured amount. The value includes an implicit precision in the presentation of the value.
	Value float64 `json:"value"`

	// How the value should be understood and represented - whether the actual value is greater or less than the stated value due to measurement issues; e.g. if the comparator is "<" , then the real value is < stated value.
	Comparator *QuantityComparatorEnum `json:"comparator,omitempty"`

	// A human-readable form of the unit.
	Unit string `json:"unit"`

	// The identification of the system that provides the coded form of the unit.
	System scalarutils.URI `json:"system"`

	// A computer processable form of the unit in some unit representation system.
	Code *scalarutils.Code `json:"code"`
}

// FHIRQuantityInput is the input type for Quantity
type FHIRQuantityInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The value of the measured amount. The value includes an implicit precision in the presentation of the value.
	Value float64 `json:"value"`

	// How the value should be understood and represented - whether the actual value is greater or less than the stated value due to measurement issues; e.g. if the comparator is "<" , then the real value is < stated value.
	Comparator *QuantityComparatorEnum `json:"comparator,omitempty"`

	// A human-readable form of the unit.
	Unit string `json:"unit"`

	// The identification of the system that provides the coded form of the unit.
	System scalarutils.URI `json:"system"`

	// A computer processable form of the unit in some unit representation system.
	Code *scalarutils.Code `json:"code"`
}

// FHIRRange definition: a set of ordered quantities defined by a low and high limit.
type FHIRRange struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The low limit. The boundary is inclusive.
	Low FHIRQuantity `json:"low,omitempty"`

	// The high limit. The boundary is inclusive.
	High FHIRQuantity `json:"high,omitempty"`
}

// FHIRRangeInput is the input type for Range
type FHIRRangeInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The low limit. The boundary is inclusive.
	Low FHIRQuantityInput `json:"low,omitempty"`

	// The high limit. The boundary is inclusive.
	High FHIRQuantityInput `json:"high,omitempty"`
}

// FHIRRatio definition: a relationship of two quantity values - expressed as a numerator and a denominator.
type FHIRRatio struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The value of the numerator.
	Numerator FHIRQuantity `json:"numerator,omitempty"`

	// The value of the denominator.
	Denominator FHIRQuantity `json:"denominator,omitempty"`
}

// FHIRRatioInput is the input type for Ratio
type FHIRRatioInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The value of the numerator.
	Numerator FHIRQuantityInput `json:"numerator,omitempty"`

	// The value of the denominator.
	Denominator FHIRQuantityInput `json:"denominator,omitempty"`
}

// FHIRReference definition: a reference from one resource to another.
type FHIRReference struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A reference to a location at which the other resource is found. The reference may be a relative reference, in which case it is relative to the service base URL, or an absolute URL that resolves to the location where the resource is found. The reference may be version specific or not. If the reference is not to a FHIR RESTful server, then it should be assumed to be version specific. Internal fragment references (start with '#') refer to contained resources.
	Reference *string `json:"reference,omitempty"`

	//     The expected type of the target of the reference. If both Reference.type and Reference.reference are populated and Reference.reference is a FHIR URL, both SHALL be consistent.
	//
	// The type is the Canonical URL of Resource Definition that is the type this reference refers to. References are URLs that are relative to http://hl7.org/fhir/StructureDefinition/ e.g. "Patient" is a reference to http://hl7.org/fhir/StructureDefinition/Patient. Absolute URLs are only allowed for logical models (and can only be used in references in logical models, not resources).
	Type *scalarutils.URI `json:"type,omitempty" swaggertype:"primitive,string"`

	// An identifier for the target resource. This is used when there is no way to reference the other resource directly, either because the entity it represents is not available through a FHIR server, or because there is no way for the author of the resource to convert a known identifier to an actual location. There is no requirement that a Reference.identifier point to something that is actually exposed as a FHIR instance, but it SHALL point to a business concept that would be expected to be exposed as a FHIR instance, and that instance would need to be of a FHIR resource type allowed by the reference.
	Identifier *FHIRIdentifier `json:"identifier,omitempty"`

	// Plain text narrative that identifies the resource in addition to the resource reference.
	Display string `json:"display,omitempty"`
}

// ResourceID safely extracts the logical ID from the Reference.reference field (e.g., from "ResourceType/12345").
//
// Note:
// Since we dont want the codebase to be full of such code,
//
// resourceID, err := activity.Reference.GetResourceID()
//
//	if err != nil {
//		return nil, err
//	}
//
// We shall leave the responsibility of handling nil references to the function caller.
// e.g (something like below).
//
//	if r == nil || r.Reference == nil || *r.Reference == "" {
//		return "", fmt.Errorf("cannot get resource ID from a nil or empty reference")
//	}.
func (r *FHIRReference) ResourceID() string {
	if r == nil || r.Reference == nil || *r.Reference == "" {
		return ""
	}

	refStr := *r.Reference

	lastSlashIndex := strings.LastIndex(refStr, "/")
	if lastSlashIndex == -1 || lastSlashIndex == len(refStr)-1 {
		return ""
	}

	return refStr[lastSlashIndex+1:]
}

// FHIRReferenceInput is the input type for Reference
type FHIRReferenceInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A reference to a location at which the other resource is found. The reference may be a relative reference, in which case it is relative to the service base URL, or an absolute URL that resolves to the location where the resource is found. The reference may be version specific or not. If the reference is not to a FHIR RESTful server, then it should be assumed to be version specific. Internal fragment references (start with '#') refer to contained resources.
	Reference *string `json:"reference,omitempty"`

	//     The expected type of the target of the reference. If both Reference.type and Reference.reference are populated and Reference.reference is a FHIR URL, both SHALL be consistent.
	//
	// The type is the Canonical URL of Resource Definition that is the type this reference refers to. References are URLs that are relative to http://hl7.org/fhir/StructureDefinition/ e.g. "Patient" is a reference to http://hl7.org/fhir/StructureDefinition/Patient. Absolute URLs are only allowed for logical models (and can only be used in references in logical models, not resources).
	Type *scalarutils.URI `json:"type,omitempty"`

	// An identifier for the target resource. This is used when there is no way to reference the other resource directly, either because the entity it represents is not available through a FHIR server, or because there is no way for the author of the resource to convert a known identifier to an actual location. There is no requirement that a Reference.identifier point to something that is actually exposed as a FHIR instance, but it SHALL point to a business concept that would be expected to be exposed as a FHIR instance, and that instance would need to be of a FHIR resource type allowed by the reference.
	Identifier *FHIRIdentifierInput `json:"identifier,omitempty"`

	// Plain text narrative that identifies the resource in addition to the resource reference.
	Display string `json:"display,omitempty"`
}

// FHIRSampledData definition: a series of measurements taken by a device, with upper and lower limits. there may be more than one dimension in the data.
type FHIRSampledData struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The base quantity that a measured value of zero represents. In addition, this provides the units of the entire measurement series.
	Origin *FHIRQuantity `json:"origin,omitempty"`

	// The length of time between sampling times, measured in milliseconds.
	Period *scalarutils.Decimal `json:"period,omitempty"`

	// A correction factor that is applied to the sampled data points before they are added to the origin.
	Factor *scalarutils.Decimal `json:"factor,omitempty"`

	// The lower limit of detection of the measured points. This is needed if any of the data points have the value "L" (lower than detection limit).
	LowerLimit *scalarutils.Decimal `json:"lowerLimit,omitempty"`

	// The upper limit of detection of the measured points. This is needed if any of the data points have the value "U" (higher than detection limit).
	UpperLimit *scalarutils.Decimal `json:"upperLimit,omitempty"`

	// The number of sample points at each time point. If this value is greater than one, then the dimensions will be interlaced - all the sample points for a point in time will be recorded at once.
	Dimensions *string `json:"dimensions,omitempty"`

	// A series of data points which are decimal values separated by a single space (character u20). The special values "E" (error), "L" (below detection limit) and "U" (above detection limit) can also be used in place of a decimal value.
	Data *string `json:"data,omitempty"`
}

// FHIRSampledDataInput is the input type for SampledData
type FHIRSampledDataInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The base quantity that a measured value of zero represents. In addition, this provides the units of the entire measurement series.
	Origin *FHIRQuantityInput `json:"origin,omitempty"`

	// The length of time between sampling times, measured in milliseconds.
	Period *scalarutils.Decimal `json:"period,omitempty"`

	// A correction factor that is applied to the sampled data points before they are added to the origin.
	Factor *scalarutils.Decimal `json:"factor,omitempty"`

	// The lower limit of detection of the measured points. This is needed if any of the data points have the value "L" (lower than detection limit).
	LowerLimit *scalarutils.Decimal `json:"lowerLimit,omitempty"`

	// The upper limit of detection of the measured points. This is needed if any of the data points have the value "U" (higher than detection limit).
	UpperLimit *scalarutils.Decimal `json:"upperLimit,omitempty"`

	// The number of sample points at each time point. If this value is greater than one, then the dimensions will be interlaced - all the sample points for a point in time will be recorded at once.
	Dimensions *string `json:"dimensions,omitempty"`

	// A series of data points which are decimal values separated by a single space (character u20). The special values "E" (error), "L" (below detection limit) and "U" (above detection limit) can also be used in place of a decimal value.
	Data *string `json:"data,omitempty"`
}

// FHIRTiming definition: specifies an event that may occur multiple times. timing schedules are used to record when things are planned, expected or requested to occur. the most common usage is in dosage instructions for medications. they are also used when planning care of various kinds, and may be used for reporting the schedule to which past regular activities were carried out.
type FHIRTiming struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Identifies specific times when the event occurs.
	Event []*scalarutils.DateTime `json:"event,omitempty"`

	// A set of rules that describe when the event is scheduled.
	Repeat *FHIRTimingRepeat `json:"repeat,omitempty"`

	// A code for the timing schedule (or just text in code.text). Some codes such as BID are ubiquitous, but many institutions define their own additional codes. If a code is provided, the code is understood to be a complete statement of whatever is specified in the structured timing data, and either the code or the data may be used to interpret the Timing, with the exception that .repeat.bounds still applies over the code (and is not contained in the code).
	Code scalarutils.Code `json:"code,omitempty"`
}

// FHIRTimingInput is the input type for Timing
type FHIRTimingInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Identifies specific times when the event occurs.
	Event *scalarutils.DateTime `json:"event,omitempty"`

	// A set of rules that describe when the event is scheduled.
	Repeat *FHIRTimingRepeatInput `json:"repeat,omitempty"`

	// A code for the timing schedule (or just text in code.text). Some codes such as BID are ubiquitous, but many institutions define their own additional codes. If a code is provided, the code is understood to be a complete statement of whatever is specified in the structured timing data, and either the code or the data may be used to interpret the Timing, with the exception that .repeat.bounds still applies over the code (and is not contained in the code).
	Code scalarutils.Code `json:"code,omitempty"`
}

// FHIRTimingRepeat definition: specifies an event that may occur multiple times. timing schedules are used to record when things are planned, expected or requested to occur. the most common usage is in dosage instructions for medications. they are also used when planning care of various kinds, and may be used for reporting the schedule to which past regular activities were carried out.
type FHIRTimingRepeat struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Either a duration for the length of the timing schedule, a range of possible length, or outer bounds for start and/or end limits of the timing schedule.
	BoundsDuration *FHIRDuration `json:"boundsDuration,omitempty"`

	// Either a duration for the length of the timing schedule, a range of possible length, or outer bounds for start and/or end limits of the timing schedule.
	BoundsRange *FHIRRange `json:"boundsRange,omitempty"`

	// Either a duration for the length of the timing schedule, a range of possible length, or outer bounds for start and/or end limits of the timing schedule.
	BoundsPeriod *FHIRPeriod `json:"boundsPeriod,omitempty"`

	// A total count of the desired number of repetitions across the duration of the entire timing specification. If countMax is present, this element indicates the lower bound of the allowed range of count values.
	Count *int `json:"count,omitempty"`

	// If present, indicates that the count is a range - so to perform the action between [count] and [countMax] times.
	CountMax *string `json:"countMax,omitempty"`

	// How long this thing happens for when it happens. If durationMax is present, this element indicates the lower bound of the allowed range of the duration.
	Duration *json.Number `json:"duration,omitempty"`

	// If present, indicates that the duration is a range - so to perform the action between [duration] and [durationMax] time length.
	DurationMax *scalarutils.Decimal `json:"durationMax,omitempty"`

	// The units of time for the duration, in UCUM units.
	DurationUnit *UnitsOfTimeEnum `json:"durationUnit,omitempty"`

	// The number of times to repeat the action within the specified period. If frequencyMax is present, this element indicates the lower bound of the allowed range of the frequency.
	Frequency *int `json:"frequency,omitempty"`

	// If present, indicates that the frequency is a range - so to repeat between [frequency] and [frequencyMax] times within the period or period range.
	FrequencyMax *string `json:"frequencyMax,omitempty"`

	// Indicates the duration of time over which repetitions are to occur; e.g. to express "3 times per day", 3 would be the frequency and "1 day" would be the period. If periodMax is present, this element indicates the lower bound of the allowed range of the period length.
	Period *json.Number `json:"period,omitempty"`

	// If present, indicates that the period is a range from [period] to [periodMax], allowing expressing concepts such as "do this once every 3-5 days.
	PeriodMax *scalarutils.Decimal `json:"periodMax,omitempty"`

	// The units of time for the period in UCUM units.
	PeriodUnit *UnitsOfTimeEnum `json:"periodUnit,omitempty"`

	// If one or more days of week is provided, then the action happens only on the specified day(s).
	DayOfWeek []*scalarutils.Code `json:"dayOfWeek,omitempty"`

	// Specified time of day for action to take place.
	TimeOfDay *time.Time `json:"timeOfDay,omitempty"`

	// An approximate time period during the day, potentially linked to an event of daily living that indicates when the action should occur.
	When []TimingRepeatWhenEnum `json:"when,omitempty"`

	// The number of minutes from the event. If the event code does not indicate whether the minutes is before or after the event, then the offset is assumed to be after the event.
	Offset *int `json:"offset,omitempty"`
}

// FHIRTimingRepeatInput is the input type for TimingRepeat
type FHIRTimingRepeatInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Either a duration for the length of the timing schedule, a range of possible length, or outer bounds for start and/or end limits of the timing schedule.
	BoundsDuration *FHIRDurationInput `json:"boundsDuration,omitempty"`

	// Either a duration for the length of the timing schedule, a range of possible length, or outer bounds for start and/or end limits of the timing schedule.
	BoundsRange *FHIRRangeInput `json:"boundsRange,omitempty"`

	// Either a duration for the length of the timing schedule, a range of possible length, or outer bounds for start and/or end limits of the timing schedule.
	BoundsPeriod *FHIRPeriodInput `json:"boundsPeriod,omitempty"`

	// A total count of the desired number of repetitions across the duration of the entire timing specification. If countMax is present, this element indicates the lower bound of the allowed range of count values.
	Count *string `json:"count,omitempty"`

	// If present, indicates that the count is a range - so to perform the action between [count] and [countMax] times.
	CountMax *string `json:"countMax,omitempty"`

	// How long this thing happens for when it happens. If durationMax is present, this element indicates the lower bound of the allowed range of the duration.
	Duration *json.Number `json:"duration,omitempty"`

	// If present, indicates that the duration is a range - so to perform the action between [duration] and [durationMax] time length.
	DurationMax *scalarutils.Decimal `json:"durationMax,omitempty"`

	// The units of time for the duration, in UCUM units.
	DurationUnit *UnitsOfTimeEnum `json:"durationUnit,omitempty"`

	// The number of times to repeat the action within the specified period. If frequencyMax is present, this element indicates the lower bound of the allowed range of the frequency.
	Frequency *int `json:"frequency,omitempty"`

	// If present, indicates that the frequency is a range - so to repeat between [frequency] and [frequencyMax] times within the period or period range.
	FrequencyMax *string `json:"frequencyMax,omitempty"`

	// Indicates the duration of time over which repetitions are to occur; e.g. to express "3 times per day", 3 would be the frequency and "1 day" would be the period. If periodMax is present, this element indicates the lower bound of the allowed range of the period length.
	Period *json.Number `json:"period,omitempty"`

	// If present, indicates that the period is a range from [period] to [periodMax], allowing expressing concepts such as "do this once every 3-5 days.
	PeriodMax *scalarutils.Decimal `json:"periodMax,omitempty"`

	// The units of time for the period in UCUM units.
	PeriodUnit *UnitsOfTimeEnum `json:"periodUnit,omitempty"`

	// If one or more days of week is provided, then the action happens only on the specified day(s).
	DayOfWeek *scalarutils.Code `json:"dayOfWeek,omitempty"`

	// Specified time of day for action to take place.
	TimeOfDay *time.Time `json:"timeOfDay,omitempty"`

	// An approximate time period during the day, potentially linked to an event of daily living that indicates when the action should occur.
	When *TimingRepeatWhenEnum `json:"when,omitempty"`

	// The number of minutes from the event. If the event code does not indicate whether the minutes is before or after the event, then the offset is assumed to be after the event.
	Offset *int `json:"offset,omitempty"`
}

// FHIRExpression is documented here http://hl7.org/fhir/StructureDefinition/Expression
type FHIRExpression struct {
	ID          *string     `json:"id,omitempty"`
	Extension   []Extension `json:"extension,omitempty"`
	Description *string     `json:"description,omitempty"`
	Name        *string     `json:"name,omitempty"`
	Language    string      `json:"language,omitempty"`
	Expression  *string     `json:"expression,omitempty"`
	Reference   *string     `json:"reference,omitempty"`
}

// AddressTypeEnum is a FHIR enum
type AddressTypeEnum string

const (
	// AddressTypeEnumPostal ...
	AddressTypeEnumPostal AddressTypeEnum = "postal"
	// AddressTypeEnumPhysical ...
	AddressTypeEnumPhysical AddressTypeEnum = "physical"
	// AddressTypeEnumBoth ...
	AddressTypeEnumBoth AddressTypeEnum = "both"
)

// AllAddressTypeEnum ...
var AllAddressTypeEnum = []AddressTypeEnum{
	AddressTypeEnumPostal,
	AddressTypeEnumPhysical,
	AddressTypeEnumBoth,
}

// IsValid ...
func (e AddressTypeEnum) IsValid() bool {
	switch e {
	case AddressTypeEnumPostal, AddressTypeEnumPhysical, AddressTypeEnumBoth:
		return true
	}

	return false
}

// String ...
func (e AddressTypeEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *AddressTypeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AddressTypeEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AddressTypeEnum", str)
	}

	return nil
}

// MarshalGQL writes the address type enum to the supplied writer as a quoted string
func (e AddressTypeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// AddressUseEnum is a FHIR enum
type AddressUseEnum string

const (
	// AddressUseEnumHome ...
	AddressUseEnumHome AddressUseEnum = "home"
	// AddressUseEnumWork ...
	AddressUseEnumWork AddressUseEnum = "work"
	// AddressUseEnumTemp ...
	AddressUseEnumTemp AddressUseEnum = "temp"
	// AddressUseEnumOld ...
	AddressUseEnumOld AddressUseEnum = "old"
	// AddressUseEnumBilling ...
	AddressUseEnumBilling AddressUseEnum = "billing"
)

// AllAddressUseEnum ...
var AllAddressUseEnum = []AddressUseEnum{
	AddressUseEnumHome,
	AddressUseEnumWork,
	AddressUseEnumTemp,
	AddressUseEnumOld,
	AddressUseEnumBilling,
}

// IsValid ...
func (e AddressUseEnum) IsValid() bool {
	switch e {
	case AddressUseEnumHome, AddressUseEnumWork, AddressUseEnumTemp, AddressUseEnumOld, AddressUseEnumBilling:
		return true
	}

	return false
}

// String ...
func (e AddressUseEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *AddressUseEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AddressUseEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AddressUseEnum", str)
	}

	return nil
}

// MarshalGQL writes the address use enum to the supplied writer as a quoted string
func (e AddressUseEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// AgeComparatorEnum is a FHIR enum
type AgeComparatorEnum string

const (
	// AgeComparatorEnumLessThan ...
	AgeComparatorEnumLessThan AgeComparatorEnum = "less_than"
	// AgeComparatorEnumLessThanOrEqualTo ...
	AgeComparatorEnumLessThanOrEqualTo AgeComparatorEnum = "less_than_or_equal_to"
	// AgeComparatorEnumGreaterThanOrEqualTo ...
	AgeComparatorEnumGreaterThanOrEqualTo AgeComparatorEnum = "greater_than_or_equal_to"
	// AgeComparatorEnumGreaterThan ...
	AgeComparatorEnumGreaterThan AgeComparatorEnum = "greater_than"
)

// AllAgeComparatorEnum ...
var AllAgeComparatorEnum = []AgeComparatorEnum{
	AgeComparatorEnumLessThan,
	AgeComparatorEnumLessThanOrEqualTo,
	AgeComparatorEnumGreaterThanOrEqualTo,
	AgeComparatorEnumGreaterThan,
}

// IsValid ...
func (e AgeComparatorEnum) IsValid() bool {
	switch e {
	case AgeComparatorEnumLessThan, AgeComparatorEnumLessThanOrEqualTo, AgeComparatorEnumGreaterThanOrEqualTo, AgeComparatorEnumGreaterThan:
		return true
	}

	return false
}

// String renders an age comparator enum as a string
func (e AgeComparatorEnum) String() string {
	switch e {
	case AgeComparatorEnumLessThan:
		return "<"
	case AgeComparatorEnumLessThanOrEqualTo:
		return "<="
	case AgeComparatorEnumGreaterThanOrEqualTo:
		return ">="
	case AgeComparatorEnumGreaterThan:
		return ">"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *AgeComparatorEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AgeComparatorEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AgeComparatorEnum", str)
	}

	return nil
}

// MarshalGQL writes the age comparator to the supplied writer as a quoted string
func (e AgeComparatorEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// ContactPointSystemEnum is a FHIR enum
type ContactPointSystemEnum string

const (
	// ContactPointSystemEnumPhone ...
	ContactPointSystemEnumPhone ContactPointSystemEnum = "phone"
	// ContactPointSystemEnumFax ...
	ContactPointSystemEnumFax ContactPointSystemEnum = "fax"
	// ContactPointSystemEnumEmail ...
	ContactPointSystemEnumEmail ContactPointSystemEnum = "email"
	// ContactPointSystemEnumPager ...
	ContactPointSystemEnumPager ContactPointSystemEnum = "pager"
	// ContactPointSystemEnumURL ...
	ContactPointSystemEnumURL ContactPointSystemEnum = "url"
	// ContactPointSystemEnumSms ...
	ContactPointSystemEnumSms ContactPointSystemEnum = "sms"
	// ContactPointSystemEnumOther ...
	ContactPointSystemEnumOther ContactPointSystemEnum = "other"
)

// AllContactPointSystemEnum ...
var AllContactPointSystemEnum = []ContactPointSystemEnum{
	ContactPointSystemEnumPhone,
	ContactPointSystemEnumFax,
	ContactPointSystemEnumEmail,
	ContactPointSystemEnumPager,
	ContactPointSystemEnumURL,
	ContactPointSystemEnumSms,
	ContactPointSystemEnumOther,
}

// IsValid ...
func (e ContactPointSystemEnum) IsValid() bool {
	switch e {
	case ContactPointSystemEnumPhone, ContactPointSystemEnumFax, ContactPointSystemEnumEmail, ContactPointSystemEnumPager, ContactPointSystemEnumURL, ContactPointSystemEnumSms, ContactPointSystemEnumOther:
		return true
	}

	return false
}

// String ...
func (e ContactPointSystemEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *ContactPointSystemEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ContactPointSystemEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ContactPointSystemEnum", str)
	}

	return nil
}

// MarshalGQL writes the given enum to the supplied writer as a quoted string
func (e ContactPointSystemEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// ContactPointUseEnum is a FHIR enum
type ContactPointUseEnum string

const (
	// ContactPointUseEnumHome ...
	ContactPointUseEnumHome ContactPointUseEnum = "home"
	// ContactPointUseEnumWork ...
	ContactPointUseEnumWork ContactPointUseEnum = "work"
	// ContactPointUseEnumTemp ...
	ContactPointUseEnumTemp ContactPointUseEnum = "temp"
	// ContactPointUseEnumOld ...
	ContactPointUseEnumOld ContactPointUseEnum = "old"
	// ContactPointUseEnumMobile ...
	ContactPointUseEnumMobile ContactPointUseEnum = "mobile"
)

// AllContactPointUseEnum ...
var AllContactPointUseEnum = []ContactPointUseEnum{
	ContactPointUseEnumHome,
	ContactPointUseEnumWork,
	ContactPointUseEnumTemp,
	ContactPointUseEnumOld,
	ContactPointUseEnumMobile,
}

// IsValid ...
func (e ContactPointUseEnum) IsValid() bool {
	switch e {
	case ContactPointUseEnumHome, ContactPointUseEnumWork, ContactPointUseEnumTemp, ContactPointUseEnumOld, ContactPointUseEnumMobile:
		return true
	}

	return false
}

// String ...
func (e ContactPointUseEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *ContactPointUseEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ContactPointUseEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ContactPointUseEnum", str)
	}

	return nil
}

// MarshalGQL writes the contact point use to the supplied writer as a quoted string
func (e ContactPointUseEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// DurationComparatorEnum is a FHIR enum
type DurationComparatorEnum string

const (
	// DurationComparatorEnumLessThan ...
	DurationComparatorEnumLessThan DurationComparatorEnum = "less_than"
	// DurationComparatorEnumLessThanOrEqualTo ...
	DurationComparatorEnumLessThanOrEqualTo DurationComparatorEnum = "less_than_or_equal_to"
	// DurationComparatorEnumGreaterThanOrEqualTo ...
	DurationComparatorEnumGreaterThanOrEqualTo DurationComparatorEnum = "greater_than_or_equal_to"
	// DurationComparatorEnumGreaterThan ...
	DurationComparatorEnumGreaterThan DurationComparatorEnum = "greater_than"
)

// AllDurationComparatorEnum ...
var AllDurationComparatorEnum = []DurationComparatorEnum{
	DurationComparatorEnumLessThan,
	DurationComparatorEnumLessThanOrEqualTo,
	DurationComparatorEnumGreaterThanOrEqualTo,
	DurationComparatorEnumGreaterThan,
}

// IsValid ...
func (e DurationComparatorEnum) IsValid() bool {
	switch e {
	case DurationComparatorEnumLessThan, DurationComparatorEnumLessThanOrEqualTo, DurationComparatorEnumGreaterThanOrEqualTo, DurationComparatorEnumGreaterThan:
		return true
	}

	return false
}

// String ...
func (e DurationComparatorEnum) String() string {
	switch e {
	case DurationComparatorEnumLessThan:
		return "<"
	case DurationComparatorEnumLessThanOrEqualTo:
		return "<="
	case DurationComparatorEnumGreaterThan:
		return ">"
	case DurationComparatorEnumGreaterThanOrEqualTo:
		return ">="
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *DurationComparatorEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DurationComparatorEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DurationComparatorEnum", str)
	}

	return nil
}

// MarshalGQL writes the duration comparator to the supplied writer as a quoted string
func (e DurationComparatorEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// HumanNameUseEnum is a FHIR enum
type HumanNameUseEnum string

const (
	// HumanNameUseEnumUsual ...
	HumanNameUseEnumUsual HumanNameUseEnum = "usual"
	// HumanNameUseEnumOfficial ...
	HumanNameUseEnumOfficial HumanNameUseEnum = "official"
	// HumanNameUseEnumTemp ...
	HumanNameUseEnumTemp HumanNameUseEnum = "temp"
	// HumanNameUseEnumNickname ...
	HumanNameUseEnumNickname HumanNameUseEnum = "nickname"
	// HumanNameUseEnumAnonymous ...
	HumanNameUseEnumAnonymous HumanNameUseEnum = "anonymous"
	// HumanNameUseEnumOld ...
	HumanNameUseEnumOld HumanNameUseEnum = "old"
	// HumanNameUseEnumMaiden ...
	HumanNameUseEnumMaiden HumanNameUseEnum = "maiden"
)

// AllHumanNameUseEnum ...
var AllHumanNameUseEnum = []HumanNameUseEnum{
	HumanNameUseEnumUsual,
	HumanNameUseEnumOfficial,
	HumanNameUseEnumTemp,
	HumanNameUseEnumNickname,
	HumanNameUseEnumAnonymous,
	HumanNameUseEnumOld,
	HumanNameUseEnumMaiden,
}

// IsValid ...
func (e HumanNameUseEnum) IsValid() bool {
	switch e {
	case HumanNameUseEnumUsual, HumanNameUseEnumOfficial, HumanNameUseEnumTemp, HumanNameUseEnumNickname, HumanNameUseEnumAnonymous, HumanNameUseEnumOld, HumanNameUseEnumMaiden:
		return true
	}

	return false
}

// String ...
func (e HumanNameUseEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *HumanNameUseEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HumanNameUseEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HumanNameUseEnum", str)
	}

	return nil
}

// MarshalGQL writes the given enum to the supplied writer as a quoted string
func (e HumanNameUseEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// IdentifierUseEnum is a FHIR enum
type IdentifierUseEnum string

const (
	// IdentifierUseEnumUsual ...
	IdentifierUseEnumUsual IdentifierUseEnum = "usual"
	// IdentifierUseEnumOfficial ...
	IdentifierUseEnumOfficial IdentifierUseEnum = "official"
	// IdentifierUseEnumTemp ...
	IdentifierUseEnumTemp IdentifierUseEnum = "temp"
	// IdentifierUseEnumSecondary ...
	IdentifierUseEnumSecondary IdentifierUseEnum = "secondary"
	// IdentifierUseEnumOld ...
	IdentifierUseEnumOld IdentifierUseEnum = "old"
)

// AllIdentifierUseEnum ...
var AllIdentifierUseEnum = []IdentifierUseEnum{
	IdentifierUseEnumUsual,
	IdentifierUseEnumOfficial,
	IdentifierUseEnumTemp,
	IdentifierUseEnumSecondary,
	IdentifierUseEnumOld,
}

// IsValid ...
func (e IdentifierUseEnum) IsValid() bool {
	switch e {
	case IdentifierUseEnumUsual, IdentifierUseEnumOfficial, IdentifierUseEnumTemp, IdentifierUseEnumSecondary, IdentifierUseEnumOld:
		return true
	}

	return false
}

// String ...
func (e IdentifierUseEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *IdentifierUseEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IdentifierUseEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IdentifierUseEnum", str)
	}

	return nil
}

// MarshalGQL writes the identifier use to the supplied writer as a quoted string
func (e IdentifierUseEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// NarrativeStatusEnum is a FHIR enum
type NarrativeStatusEnum string

const (
	// NarrativeStatusEnumGenerated ...
	NarrativeStatusEnumGenerated NarrativeStatusEnum = "generated"
	// NarrativeStatusEnumExtensions ...
	NarrativeStatusEnumExtensions NarrativeStatusEnum = "extensions"
	// NarrativeStatusEnumAdditional ...
	NarrativeStatusEnumAdditional NarrativeStatusEnum = "additional"
	// NarrativeStatusEnumEmpty ...
	NarrativeStatusEnumEmpty NarrativeStatusEnum = "empty"
)

// AllNarrativeStatusEnum ...
var AllNarrativeStatusEnum = []NarrativeStatusEnum{
	NarrativeStatusEnumGenerated,
	NarrativeStatusEnumExtensions,
	NarrativeStatusEnumAdditional,
	NarrativeStatusEnumEmpty,
}

// IsValid ...
func (e NarrativeStatusEnum) IsValid() bool {
	switch e {
	case NarrativeStatusEnumGenerated, NarrativeStatusEnumExtensions, NarrativeStatusEnumAdditional, NarrativeStatusEnumEmpty:
		return true
	}

	return false
}

// String ...
func (e NarrativeStatusEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *NarrativeStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NarrativeStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NarrativeStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the given enum to the supplied writer as a quoted string
func (e NarrativeStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// QuantityComparatorEnum is a FHIR enum
type QuantityComparatorEnum string

const (
	// QuantityComparatorEnumLessThan ...
	QuantityComparatorEnumLessThan QuantityComparatorEnum = "less_than"
	// QuantityComparatorEnumLessThanOrEqualTo ...
	QuantityComparatorEnumLessThanOrEqualTo QuantityComparatorEnum = "less_than_or_equal_to"
	// QuantityComparatorEnumGreaterThanOrEqualTo ...
	QuantityComparatorEnumGreaterThanOrEqualTo QuantityComparatorEnum = "greater_than_or_equal_to"
	// QuantityComparatorEnumGreaterThan ...
	QuantityComparatorEnumGreaterThan QuantityComparatorEnum = "greater_than"
)

// AllQuantityComparatorEnum ...
var AllQuantityComparatorEnum = []QuantityComparatorEnum{
	QuantityComparatorEnumLessThan,
	QuantityComparatorEnumLessThanOrEqualTo,
	QuantityComparatorEnumGreaterThanOrEqualTo,
	QuantityComparatorEnumGreaterThan,
}

// IsValid ...
func (e QuantityComparatorEnum) IsValid() bool {
	switch e {
	case QuantityComparatorEnumLessThan, QuantityComparatorEnumLessThanOrEqualTo, QuantityComparatorEnumGreaterThanOrEqualTo, QuantityComparatorEnumGreaterThan:
		return true
	}

	return false
}

// String ...
func (e QuantityComparatorEnum) String() string {
	switch e {
	case QuantityComparatorEnumLessThan:
		return "<"
	case QuantityComparatorEnumLessThanOrEqualTo:
		return "<="
	case QuantityComparatorEnumGreaterThanOrEqualTo:
		return ">="
	case QuantityComparatorEnumGreaterThan:
		return ">"
	default:
		return string(e)
	}
}

// UnmarshalGQL ...
func (e *QuantityComparatorEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = QuantityComparatorEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid QuantityComparatorEnum", str)
	}

	return nil
}

// MarshalGQL writes the quality comparator to the supplied writer as a quoted string
func (e QuantityComparatorEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// UnitsOfTimeEnum is a FHIR enum
type UnitsOfTimeEnum string

const (
	// UnitsOfTimeEnumS ...
	UnitsOfTimeEnumS UnitsOfTimeEnum = "s"
	// UnitsOfTimeEnumMin ...
	UnitsOfTimeEnumMin UnitsOfTimeEnum = "min"
	// UnitsOfTimeEnumH ...
	UnitsOfTimeEnumH UnitsOfTimeEnum = "h"
	// UnitsOfTimeEnumD ...
	UnitsOfTimeEnumD UnitsOfTimeEnum = "d"
	// UnitsOfTimeEnumWk ...
	UnitsOfTimeEnumWk UnitsOfTimeEnum = "wk"
	// UnitsOfTimeEnumMo ...
	UnitsOfTimeEnumMo UnitsOfTimeEnum = "mo"
	// UnitsOfTimeEnumA ...
	UnitsOfTimeEnumA UnitsOfTimeEnum = "a"
)

// AllUnitsOfTimeEnum ...
var AllUnitsOfTimeEnum = []UnitsOfTimeEnum{
	UnitsOfTimeEnumS,
	UnitsOfTimeEnumMin,
	UnitsOfTimeEnumH,
	UnitsOfTimeEnumD,
	UnitsOfTimeEnumWk,
	UnitsOfTimeEnumMo,
	UnitsOfTimeEnumA,
}

// IsValid ...
func (e UnitsOfTimeEnum) IsValid() bool {
	switch e {
	case UnitsOfTimeEnumS, UnitsOfTimeEnumMin, UnitsOfTimeEnumH, UnitsOfTimeEnumD, UnitsOfTimeEnumWk, UnitsOfTimeEnumMo, UnitsOfTimeEnumA:
		return true
	}

	return false
}

// String ...
func (e UnitsOfTimeEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *UnitsOfTimeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UnitsOfTimeEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Timing_RepeatPeriodUnitEnum", str)
	}

	return nil
}

// MarshalGQL writes the timing repeat period to the supplied writer as a quoted string
func (e UnitsOfTimeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// TimingRepeatWhenEnum is a FHIR enum
type TimingRepeatWhenEnum string

const (
	// TimingRepeatWhenEnumMorn ...
	TimingRepeatWhenEnumMorn TimingRepeatWhenEnum = "MORN"
	// TimingRepeatWhenEnumMornEarly ...
	TimingRepeatWhenEnumMornEarly TimingRepeatWhenEnum = "MORN_early"
	// TimingRepeatWhenEnumMornLate ...
	TimingRepeatWhenEnumMornLate TimingRepeatWhenEnum = "MORN_late"
	// TimingRepeatWhenEnumNoon ...
	TimingRepeatWhenEnumNoon TimingRepeatWhenEnum = "NOON"
	// TimingRepeatWhenEnumAft ...
	TimingRepeatWhenEnumAft TimingRepeatWhenEnum = "AFT"
	// TimingRepeatWhenEnumAftEarly ...
	TimingRepeatWhenEnumAftEarly TimingRepeatWhenEnum = "AFT_early"
	// TimingRepeatWhenEnumAftLate ...
	TimingRepeatWhenEnumAftLate TimingRepeatWhenEnum = "AFT_late"
	// TimingRepeatWhenEnumEve ...
	TimingRepeatWhenEnumEve TimingRepeatWhenEnum = "EVE"
	// TimingRepeatWhenEnumEveEarly ...
	TimingRepeatWhenEnumEveEarly TimingRepeatWhenEnum = "EVE_early"
	// TimingRepeatWhenEnumEveLate ...
	TimingRepeatWhenEnumEveLate TimingRepeatWhenEnum = "EVE_late"
	// TimingRepeatWhenEnumNight ...
	TimingRepeatWhenEnumNight TimingRepeatWhenEnum = "NIGHT"
	// TimingRepeatWhenEnumPhs ...
	TimingRepeatWhenEnumPhs TimingRepeatWhenEnum = "PHS"
	// TimingRepeatWhenEnumHs ...
	TimingRepeatWhenEnumHs TimingRepeatWhenEnum = "HS"
	// TimingRepeatWhenEnumWake ...
	TimingRepeatWhenEnumWake TimingRepeatWhenEnum = "WAKE"
	// TimingRepeatWhenEnumC ...
	TimingRepeatWhenEnumC TimingRepeatWhenEnum = "C"
	// TimingRepeatWhenEnumCm ...
	TimingRepeatWhenEnumCm TimingRepeatWhenEnum = "CM"
	// TimingRepeatWhenEnumCd ...
	TimingRepeatWhenEnumCd TimingRepeatWhenEnum = "CD"
	// TimingRepeatWhenEnumCv ...
	TimingRepeatWhenEnumCv TimingRepeatWhenEnum = "CV"
	// TimingRepeatWhenEnumAc ...
	TimingRepeatWhenEnumAc TimingRepeatWhenEnum = "AC"
	// TimingRepeatWhenEnumAcm ...
	TimingRepeatWhenEnumAcm TimingRepeatWhenEnum = "ACM"
	// TimingRepeatWhenEnumAcd ...
	TimingRepeatWhenEnumAcd TimingRepeatWhenEnum = "ACD"
	// TimingRepeatWhenEnumAcv ...
	TimingRepeatWhenEnumAcv TimingRepeatWhenEnum = "ACV"
	// TimingRepeatWhenEnumPc ...
	TimingRepeatWhenEnumPc TimingRepeatWhenEnum = "PC"
	// TimingRepeatWhenEnumPcm ...
	TimingRepeatWhenEnumPcm TimingRepeatWhenEnum = "PCM"
	// TimingRepeatWhenEnumPcd ...
	TimingRepeatWhenEnumPcd TimingRepeatWhenEnum = "PCD"
	// TimingRepeatWhenEnumPcv ...
	TimingRepeatWhenEnumPcv       TimingRepeatWhenEnum = "PCV"
	TimingRepeatWhenEnumImmediate TimingRepeatWhenEnum = "IMD"
)

// AllTimingRepeatWhenEnum ...
var AllTimingRepeatWhenEnum = []TimingRepeatWhenEnum{
	TimingRepeatWhenEnumMorn,
	TimingRepeatWhenEnumMornEarly,
	TimingRepeatWhenEnumMornLate,
	TimingRepeatWhenEnumNoon,
	TimingRepeatWhenEnumAft,
	TimingRepeatWhenEnumAftEarly,
	TimingRepeatWhenEnumAftLate,
	TimingRepeatWhenEnumEve,
	TimingRepeatWhenEnumEveEarly,
	TimingRepeatWhenEnumEveLate,
	TimingRepeatWhenEnumNight,
	TimingRepeatWhenEnumPhs,
	TimingRepeatWhenEnumHs,
	TimingRepeatWhenEnumWake,
	TimingRepeatWhenEnumC,
	TimingRepeatWhenEnumCm,
	TimingRepeatWhenEnumCd,
	TimingRepeatWhenEnumCv,
	TimingRepeatWhenEnumAc,
	TimingRepeatWhenEnumAcm,
	TimingRepeatWhenEnumAcd,
	TimingRepeatWhenEnumAcv,
	TimingRepeatWhenEnumPc,
	TimingRepeatWhenEnumPcm,
	TimingRepeatWhenEnumPcd,
	TimingRepeatWhenEnumPcv,
	TimingRepeatWhenEnumImmediate,
}

// IsValid ...
func (e TimingRepeatWhenEnum) IsValid() bool {
	switch e {
	case TimingRepeatWhenEnumMorn, TimingRepeatWhenEnumMornEarly, TimingRepeatWhenEnumMornLate,
		TimingRepeatWhenEnumNoon, TimingRepeatWhenEnumAft, TimingRepeatWhenEnumAftEarly, TimingRepeatWhenEnumAftLate,
		TimingRepeatWhenEnumEve, TimingRepeatWhenEnumEveEarly, TimingRepeatWhenEnumEveLate, TimingRepeatWhenEnumNight,
		TimingRepeatWhenEnumPhs, TimingRepeatWhenEnumHs, TimingRepeatWhenEnumWake, TimingRepeatWhenEnumC, TimingRepeatWhenEnumCm,
		TimingRepeatWhenEnumCd, TimingRepeatWhenEnumCv, TimingRepeatWhenEnumAc, TimingRepeatWhenEnumAcm, TimingRepeatWhenEnumAcd,
		TimingRepeatWhenEnumAcv, TimingRepeatWhenEnumPc, TimingRepeatWhenEnumPcm, TimingRepeatWhenEnumPcd, TimingRepeatWhenEnumPcv, TimingRepeatWhenEnumImmediate:
		return true
	}

	return false
}

// String ...
func (e TimingRepeatWhenEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *TimingRepeatWhenEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TimingRepeatWhenEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Timing_RepeatWhenEnum", str)
	}

	return nil
}

// MarshalGQL writes when timings repeat to the supplied writer as a quoted string
func (e TimingRepeatWhenEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// AllergyIntoleranceCategoryEnum is a FHIR enum
type AllergyIntoleranceCategoryEnum string

const (
	// AllergyIntoleranceCategoryEnumFood ...
	AllergyIntoleranceCategoryEnumFood AllergyIntoleranceCategoryEnum = "food"
	// AllergyIntoleranceCategoryEnumMedication ...
	AllergyIntoleranceCategoryEnumMedication AllergyIntoleranceCategoryEnum = "medication"
	// AllergyIntoleranceCategoryEnumEnvironment ...
	AllergyIntoleranceCategoryEnumEnvironment AllergyIntoleranceCategoryEnum = "environment"
	// AllergyIntoleranceCategoryEnumBiologic ...
	AllergyIntoleranceCategoryEnumBiologic AllergyIntoleranceCategoryEnum = "biologic"
)

// AllAllergyIntoleranceCategoryEnum ...
var AllAllergyIntoleranceCategoryEnum = []AllergyIntoleranceCategoryEnum{
	AllergyIntoleranceCategoryEnumFood,
	AllergyIntoleranceCategoryEnumMedication,
	AllergyIntoleranceCategoryEnumEnvironment,
	AllergyIntoleranceCategoryEnumBiologic,
}

// IsValid ...
func (e AllergyIntoleranceCategoryEnum) IsValid() bool {
	switch e {
	case AllergyIntoleranceCategoryEnumFood, AllergyIntoleranceCategoryEnumMedication, AllergyIntoleranceCategoryEnumEnvironment, AllergyIntoleranceCategoryEnumBiologic:
		return true
	}

	return false
}

// String ...
func (e AllergyIntoleranceCategoryEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *AllergyIntoleranceCategoryEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AllergyIntoleranceCategoryEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AllergyIntoleranceCategoryEnum", str)
	}

	return nil
}

// MarshalGQL writes the allergy intolerance category to the supplied writer as a quoted string
func (e AllergyIntoleranceCategoryEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// AllergyIntoleranceCriticalityEnum is a FHIR enum
type AllergyIntoleranceCriticalityEnum string

const (
	// AllergyIntoleranceCriticalityEnumLow ...
	AllergyIntoleranceCriticalityEnumLow AllergyIntoleranceCriticalityEnum = "low"
	// AllergyIntoleranceCriticalityEnumHigh ...
	AllergyIntoleranceCriticalityEnumHigh AllergyIntoleranceCriticalityEnum = "high"
	// AllergyIntoleranceCriticalityEnumUnableToAssess ...
	AllergyIntoleranceCriticalityEnumUnableToAssess AllergyIntoleranceCriticalityEnum = "unable_to_assess"
)

// AllAllergyIntoleranceCriticalityEnum ...
var AllAllergyIntoleranceCriticalityEnum = []AllergyIntoleranceCriticalityEnum{
	AllergyIntoleranceCriticalityEnumLow,
	AllergyIntoleranceCriticalityEnumHigh,
	AllergyIntoleranceCriticalityEnumUnableToAssess,
}

// IsValid ...
func (e AllergyIntoleranceCriticalityEnum) IsValid() bool {
	switch e {
	case AllergyIntoleranceCriticalityEnumLow, AllergyIntoleranceCriticalityEnumHigh, AllergyIntoleranceCriticalityEnumUnableToAssess:
		return true
	}

	return false
}

// String ...
func (e AllergyIntoleranceCriticalityEnum) String() string {
	if e == AllergyIntoleranceCriticalityEnumUnableToAssess {
		return "unable-to-assess"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *AllergyIntoleranceCriticalityEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AllergyIntoleranceCriticalityEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AllergyIntoleranceCriticalityEnum", str)
	}

	return nil
}

// MarshalGQL writes the allergy intolerance criticality to the supplied writer as a quoted string
func (e AllergyIntoleranceCriticalityEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// AllergyIntoleranceTypeEnum is a FHIR enum
type AllergyIntoleranceTypeEnum string

const (
	// AllergyIntoleranceTypeEnumAllergy ...
	AllergyIntoleranceTypeEnumAllergy AllergyIntoleranceTypeEnum = "allergy"
	// AllergyIntoleranceTypeEnumIntolerance ...
	AllergyIntoleranceTypeEnumIntolerance AllergyIntoleranceTypeEnum = "intolerance"
)

// AllAllergyIntoleranceTypeEnum ...
var AllAllergyIntoleranceTypeEnum = []AllergyIntoleranceTypeEnum{
	AllergyIntoleranceTypeEnumAllergy,
	AllergyIntoleranceTypeEnumIntolerance,
}

// IsValid ...
func (e AllergyIntoleranceTypeEnum) IsValid() bool {
	switch e {
	case AllergyIntoleranceTypeEnumAllergy, AllergyIntoleranceTypeEnumIntolerance:
		return true
	}

	return false
}

// String ...
func (e AllergyIntoleranceTypeEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *AllergyIntoleranceTypeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AllergyIntoleranceTypeEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AllergyIntoleranceTypeEnum", str)
	}

	return nil
}

// MarshalGQL writes the allergy intolerance type to the supplied writer as a quoted string
func (e AllergyIntoleranceTypeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// AllergyIntoleranceReactionSeverityEnum is a FHIR enum
type AllergyIntoleranceReactionSeverityEnum string

const (
	// AllergyIntoleranceReactionSeverityEnumMild ...
	AllergyIntoleranceReactionSeverityEnumMild AllergyIntoleranceReactionSeverityEnum = "mild"
	// AllergyIntoleranceReactionSeverityEnumModerate ...
	AllergyIntoleranceReactionSeverityEnumModerate AllergyIntoleranceReactionSeverityEnum = "moderate"
	// AllergyIntoleranceReactionSeverityEnumSevere ...
	AllergyIntoleranceReactionSeverityEnumSevere AllergyIntoleranceReactionSeverityEnum = "severe"
)

// AllAllergyIntoleranceReactionSeverityEnum ...
var AllAllergyIntoleranceReactionSeverityEnum = []AllergyIntoleranceReactionSeverityEnum{
	AllergyIntoleranceReactionSeverityEnumMild,
	AllergyIntoleranceReactionSeverityEnumModerate,
	AllergyIntoleranceReactionSeverityEnumSevere,
}

// IsValid ...
func (e AllergyIntoleranceReactionSeverityEnum) IsValid() bool {
	switch e {
	case AllergyIntoleranceReactionSeverityEnumMild, AllergyIntoleranceReactionSeverityEnumModerate, AllergyIntoleranceReactionSeverityEnumSevere:
		return true
	}

	return false
}

// String ...
func (e AllergyIntoleranceReactionSeverityEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *AllergyIntoleranceReactionSeverityEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AllergyIntoleranceReactionSeverityEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AllergyIntolerance_ReactionSeverityEnum", str)
	}

	return nil
}

// MarshalGQL writes the allergy intolerance reaction severity to the supplied writer as a quoted string
func (e AllergyIntoleranceReactionSeverityEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// CompositionStatusEnum is a FHIR enum
type CompositionStatusEnum string

const (
	// CompositionStatusEnumPreliminary ...
	CompositionStatusEnumPreliminary CompositionStatusEnum = "preliminary"
	CompositionStatusEnumRegistered  CompositionStatusEnum = "registered"
	// CompositionStatusEnumFinal ...
	CompositionStatusEnumFinal CompositionStatusEnum = "final"
	// CompositionStatusEnumAmended ...
	CompositionStatusEnumAmended CompositionStatusEnum = "amended"
	// CompositionStatusEnumEnteredInError ...
	CompositionStatusEnumEnteredInError CompositionStatusEnum = "entered_in_error"
)

// AllCompositionStatusEnum ...
var AllCompositionStatusEnum = []CompositionStatusEnum{
	CompositionStatusEnumPreliminary,
	CompositionStatusEnumFinal,
	CompositionStatusEnumAmended,
	CompositionStatusEnumEnteredInError,
}

// IsValid ...
func (e CompositionStatusEnum) IsValid() bool {
	switch e {
	case CompositionStatusEnumPreliminary, CompositionStatusEnumFinal, CompositionStatusEnumAmended, CompositionStatusEnumEnteredInError:
		return true
	}

	return false
}

// String ...
func (e CompositionStatusEnum) String() string {
	if e == CompositionStatusEnumEnteredInError {
		return "entered-in-error"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *CompositionStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CompositionStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CompositionStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the composition status to the supplied writer as a quoted string
func (e CompositionStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// CompositionAttesterModeEnum is a FHIR enum
type CompositionAttesterModeEnum string

const (
	// CompositionAttesterModeEnumPersonal ...
	CompositionAttesterModeEnumPersonal CompositionAttesterModeEnum = "personal"
	// CompositionAttesterModeEnumProfessional ...
	CompositionAttesterModeEnumProfessional CompositionAttesterModeEnum = "professional"
	// CompositionAttesterModeEnumLegal ...
	CompositionAttesterModeEnumLegal CompositionAttesterModeEnum = "legal"
	// CompositionAttesterModeEnumOfficial ...
	CompositionAttesterModeEnumOfficial CompositionAttesterModeEnum = "official"
)

// AllCompositionAttesterModeEnum ...
var AllCompositionAttesterModeEnum = []CompositionAttesterModeEnum{
	CompositionAttesterModeEnumPersonal,
	CompositionAttesterModeEnumProfessional,
	CompositionAttesterModeEnumLegal,
	CompositionAttesterModeEnumOfficial,
}

// IsValid ...
func (e CompositionAttesterModeEnum) IsValid() bool {
	switch e {
	case CompositionAttesterModeEnumPersonal, CompositionAttesterModeEnumProfessional, CompositionAttesterModeEnumLegal, CompositionAttesterModeEnumOfficial:
		return true
	}

	return false
}

// String ...
func (e CompositionAttesterModeEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *CompositionAttesterModeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CompositionAttesterModeEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Composition_AttesterModeEnum", str)
	}

	return nil
}

// MarshalGQL writes the composition attester mode to the supplied writer as a quoted string
func (e CompositionAttesterModeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// EncounterStatusEnum is a FHIR enum
type EncounterStatusEnum string

const (
	// EncounterStatusEnumPlanned: The Encounter has not yet started.
	EncounterStatusEnumPlanned EncounterStatusEnum = "planned"
	// EncounterStatusEnumOnHold: The encounter has begun but it is currently on hold.
	EncounterStatusEnumOnHold EncounterStatusEnum = "on_hold"
	// EncounterStatusEnumDischarged: The encounter has been clinically completed and the patient has been discharged from the facility.
	EncounterStatusEnumDischarged EncounterStatusEnum = "discharged"
	// EncounterStatusEnumInProgress: The encounter has begun, the patient and the practitioner are meeting.
	EncounterStatusEnumInProgress EncounterStatusEnum = "in_progress"
	// EncounterStatusEnumDiscontinued: The encounter has started but it was not not able to be completed.
	// Further action man need to be performed related to this encounter.
	EncounterStatusEnumDiscontinued EncounterStatusEnum = "discontinued"
	// EncounterStatusEnumCompleted: The encounter has ended.
	EncounterStatusEnumCompleted EncounterStatusEnum = "completed"
	// EncounterStatusEnumCancelled: The encounter has been cancelled before it has begun.
	EncounterStatusEnumCancelled EncounterStatusEnum = "cancelled"
	// EncounterStatusEnumEnteredInError: This instance should not have been part of this patient's medical record.
	EncounterStatusEnumEnteredInError EncounterStatusEnum = "entered_in_error"
	// EncounterStatusEnumUnknown: The encounter status is unknown.
	EncounterStatusEnumUnknown EncounterStatusEnum = "unknown"
)

// AllEncounterStatusEnum ...
var AllEncounterStatusEnum = []EncounterStatusEnum{
	EncounterStatusEnumPlanned,
	EncounterStatusEnumOnHold,
	EncounterStatusEnumDischarged,
	EncounterStatusEnumInProgress,
	EncounterStatusEnumCompleted,
	EncounterStatusEnumDiscontinued,
	EncounterStatusEnumCancelled,
	EncounterStatusEnumEnteredInError,
	EncounterStatusEnumUnknown,
}

// IsValid ...
func (e EncounterStatusEnum) IsValid() bool {
	switch e {
	case EncounterStatusEnumPlanned, EncounterStatusEnumDiscontinued, EncounterStatusEnumCompleted, EncounterStatusEnumInProgress, EncounterStatusEnumDischarged, EncounterStatusEnumCancelled, EncounterStatusEnumEnteredInError, EncounterStatusEnumUnknown:
		return true
	}

	return false
}

// IsFinal ...
func (e EncounterStatusEnum) IsFinal() bool {
	switch e {
	case EncounterStatusEnumCompleted, EncounterStatusEnumCancelled, EncounterStatusEnumEnteredInError:
		return true
	}

	return false
}

// String ...
func (e EncounterStatusEnum) String() string {
	switch e {
	case EncounterStatusEnumInProgress:
		return "in-progress"
	case EncounterStatusEnumEnteredInError:
		return "entered-in-error"
	case EncounterStatusEnumOnHold:
		return "on-hold"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *EncounterStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EncounterStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EncounterStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the encounter status to the supplied writer as a quoted string
func (e EncounterStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// EncounterLocationStatusEnum is a FHIR enum
type EncounterLocationStatusEnum string

const (
	// EncounterLocationStatusEnumPlanned ...
	EncounterLocationStatusEnumPlanned EncounterLocationStatusEnum = "planned"
	// EncounterLocationStatusEnumActive ...
	EncounterLocationStatusEnumActive EncounterLocationStatusEnum = "active"
	// EncounterLocationStatusEnumReserved ...
	EncounterLocationStatusEnumReserved EncounterLocationStatusEnum = "reserved"
	// EncounterLocationStatusEnumCompleted ...
	EncounterLocationStatusEnumCompleted EncounterLocationStatusEnum = "completed"
)

// AllEncounterLocationStatusEnum ...
var AllEncounterLocationStatusEnum = []EncounterLocationStatusEnum{
	EncounterLocationStatusEnumPlanned,
	EncounterLocationStatusEnumActive,
	EncounterLocationStatusEnumReserved,
	EncounterLocationStatusEnumCompleted,
}

// IsValid ...
func (e EncounterLocationStatusEnum) IsValid() bool {
	switch e {
	case EncounterLocationStatusEnumPlanned, EncounterLocationStatusEnumActive, EncounterLocationStatusEnumReserved, EncounterLocationStatusEnumCompleted:
		return true
	}

	return false
}

// String ...
func (e EncounterLocationStatusEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *EncounterLocationStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EncounterLocationStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Encounter_LocationStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the encounter location status to the supplied writer as a quoted string
func (e EncounterLocationStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// EncounterStatusHistoryStatusEnum is a FHIR enum
type EncounterStatusHistoryStatusEnum string

const (
	// EncounterStatusHistoryStatusEnumPlanned ...
	EncounterStatusHistoryStatusEnumPlanned EncounterStatusHistoryStatusEnum = "planned"
	// EncounterStatusHistoryStatusEnumArrived ...
	EncounterStatusHistoryStatusEnumArrived EncounterStatusHistoryStatusEnum = "arrived"
	// EncounterStatusHistoryStatusEnumTriaged ...
	EncounterStatusHistoryStatusEnumTriaged EncounterStatusHistoryStatusEnum = "triaged"
	// EncounterStatusHistoryStatusEnumInProgress ...
	EncounterStatusHistoryStatusEnumInProgress EncounterStatusHistoryStatusEnum = "in_progress"
	// EncounterStatusHistoryStatusEnumOnleave ...
	EncounterStatusHistoryStatusEnumOnleave EncounterStatusHistoryStatusEnum = "onleave"
	// EncounterStatusHistoryStatusEnumFinished ...
	EncounterStatusHistoryStatusEnumFinished EncounterStatusHistoryStatusEnum = "finished"
	// EncounterStatusHistoryStatusEnumCancelled ...
	EncounterStatusHistoryStatusEnumCancelled EncounterStatusHistoryStatusEnum = "cancelled"
	// EncounterStatusHistoryStatusEnumEnteredInError ...
	EncounterStatusHistoryStatusEnumEnteredInError EncounterStatusHistoryStatusEnum = "entered_in_error"
	// EncounterStatusHistoryStatusEnumUnknown ...
	EncounterStatusHistoryStatusEnumUnknown EncounterStatusHistoryStatusEnum = "unknown"
)

// AllEncounterStatusHistoryStatusEnum ...
var AllEncounterStatusHistoryStatusEnum = []EncounterStatusHistoryStatusEnum{
	EncounterStatusHistoryStatusEnumPlanned,
	EncounterStatusHistoryStatusEnumArrived,
	EncounterStatusHistoryStatusEnumTriaged,
	EncounterStatusHistoryStatusEnumInProgress,
	EncounterStatusHistoryStatusEnumOnleave,
	EncounterStatusHistoryStatusEnumFinished,
	EncounterStatusHistoryStatusEnumCancelled,
	EncounterStatusHistoryStatusEnumEnteredInError,
	EncounterStatusHistoryStatusEnumUnknown,
}

// IsValid ...
func (e EncounterStatusHistoryStatusEnum) IsValid() bool {
	switch e {
	case EncounterStatusHistoryStatusEnumPlanned, EncounterStatusHistoryStatusEnumArrived, EncounterStatusHistoryStatusEnumTriaged, EncounterStatusHistoryStatusEnumInProgress, EncounterStatusHistoryStatusEnumOnleave, EncounterStatusHistoryStatusEnumFinished, EncounterStatusHistoryStatusEnumCancelled, EncounterStatusHistoryStatusEnumEnteredInError, EncounterStatusHistoryStatusEnumUnknown:
		return true
	}

	return false
}

// String ...
func (e EncounterStatusHistoryStatusEnum) String() string {
	switch e {
	case EncounterStatusHistoryStatusEnumInProgress:
		return "in-progress"
	case EncounterStatusHistoryStatusEnumEnteredInError:
		return "entered-in-error"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *EncounterStatusHistoryStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EncounterStatusHistoryStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Encounter_StatusHistoryStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the given enum  to the supplied writer as a quoted string
func (e EncounterStatusHistoryStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// EpisodeOfCareStatusEnum is a FHIR enum
type EpisodeOfCareStatusEnum string

const (
	// EpisodeOfCareStatusEnumPlanned ...
	EpisodeOfCareStatusEnumPlanned EpisodeOfCareStatusEnum = "planned"
	// EpisodeOfCareStatusEnumWaitlist ...
	EpisodeOfCareStatusEnumWaitlist EpisodeOfCareStatusEnum = "waitlist"
	// EpisodeOfCareStatusEnumActive ...
	EpisodeOfCareStatusEnumActive EpisodeOfCareStatusEnum = "active"
	// EpisodeOfCareStatusEnumOnhold ...
	EpisodeOfCareStatusEnumOnhold EpisodeOfCareStatusEnum = "onhold"
	// EpisodeOfCareStatusEnumFinished ...
	EpisodeOfCareStatusEnumFinished EpisodeOfCareStatusEnum = "finished"
	// EpisodeOfCareStatusEnumCancelled ...
	EpisodeOfCareStatusEnumCancelled EpisodeOfCareStatusEnum = "cancelled"
	// EpisodeOfCareStatusEnumEnteredInError ...
	EpisodeOfCareStatusEnumEnteredInError EpisodeOfCareStatusEnum = "entered_in_error"
)

// AllEpisodeOfCareStatusEnum ...
var AllEpisodeOfCareStatusEnum = []EpisodeOfCareStatusEnum{
	EpisodeOfCareStatusEnumPlanned,
	EpisodeOfCareStatusEnumWaitlist,
	EpisodeOfCareStatusEnumActive,
	EpisodeOfCareStatusEnumOnhold,
	EpisodeOfCareStatusEnumFinished,
	EpisodeOfCareStatusEnumCancelled,
	EpisodeOfCareStatusEnumEnteredInError,
}

// IsValid ...
func (e EpisodeOfCareStatusEnum) IsValid() bool {
	switch e {
	case EpisodeOfCareStatusEnumPlanned, EpisodeOfCareStatusEnumWaitlist, EpisodeOfCareStatusEnumActive, EpisodeOfCareStatusEnumOnhold, EpisodeOfCareStatusEnumFinished, EpisodeOfCareStatusEnumCancelled, EpisodeOfCareStatusEnumEnteredInError:
		return true
	}

	return false
}

// IsFinal ...
func (e EpisodeOfCareStatusEnum) IsFinal() bool {
	switch e {
	case EpisodeOfCareStatusEnumFinished, EpisodeOfCareStatusEnumCancelled, EpisodeOfCareStatusEnumEnteredInError:
		return true
	}

	return false
}

// String ...
func (e EpisodeOfCareStatusEnum) String() string {
	if e == EpisodeOfCareStatusEnumEnteredInError {
		return "entered-in-error"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *EpisodeOfCareStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EpisodeOfCareStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EpisodeOfCareStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the episode of care status to the supplied writer as a quoted string
func (e EpisodeOfCareStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// EpisodeOfCareStatusHistoryStatusEnum is a FHIR enum
type EpisodeOfCareStatusHistoryStatusEnum string

const (
	// EpisodeOfCareStatusHistoryStatusEnumPlanned ...
	EpisodeOfCareStatusHistoryStatusEnumPlanned EpisodeOfCareStatusHistoryStatusEnum = "planned"
	// EpisodeOfCareStatusHistoryStatusEnumWaitlist ...
	EpisodeOfCareStatusHistoryStatusEnumWaitlist EpisodeOfCareStatusHistoryStatusEnum = "waitlist"
	// EpisodeOfCareStatusHistoryStatusEnumActive ...
	EpisodeOfCareStatusHistoryStatusEnumActive EpisodeOfCareStatusHistoryStatusEnum = "active"
	// EpisodeOfCareStatusHistoryStatusEnumOnhold ...
	EpisodeOfCareStatusHistoryStatusEnumOnhold EpisodeOfCareStatusHistoryStatusEnum = "onhold"
	// EpisodeOfCareStatusHistoryStatusEnumFinished ...
	EpisodeOfCareStatusHistoryStatusEnumFinished EpisodeOfCareStatusHistoryStatusEnum = "finished"
	// EpisodeOfCareStatusHistoryStatusEnumCancelled ...
	EpisodeOfCareStatusHistoryStatusEnumCancelled EpisodeOfCareStatusHistoryStatusEnum = "cancelled"
	// EpisodeOfCareStatusHistoryStatusEnumEnteredInError ...
	EpisodeOfCareStatusHistoryStatusEnumEnteredInError EpisodeOfCareStatusHistoryStatusEnum = "entered_in_error"
)

// AllEpisodeOfCareStatusHistoryStatusEnum ...
var AllEpisodeOfCareStatusHistoryStatusEnum = []EpisodeOfCareStatusHistoryStatusEnum{
	EpisodeOfCareStatusHistoryStatusEnumPlanned,
	EpisodeOfCareStatusHistoryStatusEnumWaitlist,
	EpisodeOfCareStatusHistoryStatusEnumActive,
	EpisodeOfCareStatusHistoryStatusEnumOnhold,
	EpisodeOfCareStatusHistoryStatusEnumFinished,
	EpisodeOfCareStatusHistoryStatusEnumCancelled,
	EpisodeOfCareStatusHistoryStatusEnumEnteredInError,
}

// IsValid ...
func (e EpisodeOfCareStatusHistoryStatusEnum) IsValid() bool {
	switch e {
	case EpisodeOfCareStatusHistoryStatusEnumPlanned, EpisodeOfCareStatusHistoryStatusEnumWaitlist, EpisodeOfCareStatusHistoryStatusEnumActive, EpisodeOfCareStatusHistoryStatusEnumOnhold, EpisodeOfCareStatusHistoryStatusEnumFinished, EpisodeOfCareStatusHistoryStatusEnumCancelled, EpisodeOfCareStatusHistoryStatusEnumEnteredInError:
		return true
	}

	return false
}

// String ...
func (e EpisodeOfCareStatusHistoryStatusEnum) String() string {
	if e == EpisodeOfCareStatusHistoryStatusEnumEnteredInError {
		return "entered-in-error"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *EpisodeOfCareStatusHistoryStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EpisodeOfCareStatusHistoryStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EpisodeOfCare_StatusHistoryStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the status of the episode of care status history to the supplied writer as a quoted string
func (e EpisodeOfCareStatusHistoryStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// ObservationStatusEnum is a FHIR enum
type ObservationStatusEnum string

const (
	// ObservationStatusEnumRegistered ...
	ObservationStatusEnumRegistered ObservationStatusEnum = "registered"
	// ObservationStatusEnumPreliminary ...
	ObservationStatusEnumPreliminary ObservationStatusEnum = "preliminary"
	// ObservationStatusEnumFinal ...
	ObservationStatusEnumFinal ObservationStatusEnum = "final"
	// ObservationStatusEnumAmended ...
	ObservationStatusEnumAmended ObservationStatusEnum = "amended"
	// ObservationStatusEnumCorrected ...
	ObservationStatusEnumCorrected ObservationStatusEnum = "corrected"
	// ObservationStatusEnumCancelled ...
	ObservationStatusEnumCancelled ObservationStatusEnum = "cancelled"
	// ObservationStatusEnumEnteredInError ...
	ObservationStatusEnumEnteredInError ObservationStatusEnum = "entered_in_error"
	// ObservationStatusEnumUnknown ...
	ObservationStatusEnumUnknown ObservationStatusEnum = "unknown"
)

// AllObservationStatusEnum ...
var AllObservationStatusEnum = []ObservationStatusEnum{
	ObservationStatusEnumRegistered,
	ObservationStatusEnumPreliminary,
	ObservationStatusEnumFinal,
	ObservationStatusEnumAmended,
	ObservationStatusEnumCorrected,
	ObservationStatusEnumCancelled,
	ObservationStatusEnumEnteredInError,
	ObservationStatusEnumUnknown,
}

// IsValid ...
func (e ObservationStatusEnum) IsValid() bool {
	switch e {
	case ObservationStatusEnumRegistered, ObservationStatusEnumPreliminary, ObservationStatusEnumFinal, ObservationStatusEnumAmended, ObservationStatusEnumCorrected, ObservationStatusEnumCancelled, ObservationStatusEnumEnteredInError, ObservationStatusEnumUnknown:
		return true
	}

	return false
}

// String ...
func (e ObservationStatusEnum) String() string {
	if e == ObservationStatusEnumEnteredInError {
		return "entered-in-error"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *ObservationStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ObservationStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ObservationStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the observation status to the supplied writer as a quoted string
func (e ObservationStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// PatientGenderEnum is a FHIR enum
type PatientGenderEnum string

const (
	// PatientGenderEnumMale ...
	PatientGenderEnumMale PatientGenderEnum = "male"
	// PatientGenderEnumFemale ...
	PatientGenderEnumFemale PatientGenderEnum = "female"
	// PatientGenderEnumOther ...
	PatientGenderEnumOther PatientGenderEnum = "other"
	// PatientGenderEnumUnknown ...
	PatientGenderEnumUnknown PatientGenderEnum = "unknown"
)

// AllPatientGenderEnum ...
var AllPatientGenderEnum = []PatientGenderEnum{
	PatientGenderEnumMale,
	PatientGenderEnumFemale,
	PatientGenderEnumOther,
	PatientGenderEnumUnknown,
}

// IsValid ...
func (e PatientGenderEnum) IsValid() bool {
	switch e {
	case PatientGenderEnumMale, PatientGenderEnumFemale, PatientGenderEnumOther, PatientGenderEnumUnknown:
		return true
	}

	return false
}

// String ...
func (e PatientGenderEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *PatientGenderEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PatientGenderEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PatientGenderEnum", str)
	}

	return nil
}

// MarshalGQL writes the patient gender to the supplied writer as a quoted string
func (e PatientGenderEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// PatientContactGenderEnum is a FHIR enum
type PatientContactGenderEnum string

const (
	// PatientContactGenderEnumMale ...
	PatientContactGenderEnumMale PatientContactGenderEnum = "male"
	// PatientContactGenderEnumFemale ...
	PatientContactGenderEnumFemale PatientContactGenderEnum = "female"
	// PatientContactGenderEnumOther ...
	PatientContactGenderEnumOther PatientContactGenderEnum = "other"
	// PatientContactGenderEnumUnknown ...
	PatientContactGenderEnumUnknown PatientContactGenderEnum = "unknown"
)

// AllPatientContactGenderEnum ...
var AllPatientContactGenderEnum = []PatientContactGenderEnum{
	PatientContactGenderEnumMale,
	PatientContactGenderEnumFemale,
	PatientContactGenderEnumOther,
	PatientContactGenderEnumUnknown,
}

// IsValid ...
func (e PatientContactGenderEnum) IsValid() bool {
	switch e {
	case PatientContactGenderEnumMale, PatientContactGenderEnumFemale, PatientContactGenderEnumOther, PatientContactGenderEnumUnknown:
		return true
	}

	return false
}

// String ...
func (e PatientContactGenderEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *PatientContactGenderEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PatientContactGenderEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Patient_ContactGenderEnum", str)
	}

	return nil
}

// MarshalGQL writes the patient contact gender to the supplied writer as a quoted string
func (e PatientContactGenderEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// PatientLinkTypeEnum is a FHIR enum
type PatientLinkTypeEnum string

const (
	// PatientLinkTypeEnumReplacedBy ...
	PatientLinkTypeEnumReplacedBy PatientLinkTypeEnum = "replaced_by"
	// PatientLinkTypeEnumReplaces ...
	PatientLinkTypeEnumReplaces PatientLinkTypeEnum = "replaces"
	// PatientLinkTypeEnumRefer ...
	PatientLinkTypeEnumRefer PatientLinkTypeEnum = "refer"
	// PatientLinkTypeEnumSeealso ...
	PatientLinkTypeEnumSeealso PatientLinkTypeEnum = "seealso"
)

// AllPatientLinkTypeEnum ...
var AllPatientLinkTypeEnum = []PatientLinkTypeEnum{
	PatientLinkTypeEnumReplacedBy,
	PatientLinkTypeEnumReplaces,
	PatientLinkTypeEnumRefer,
	PatientLinkTypeEnumSeealso,
}

// IsValid ...
func (e PatientLinkTypeEnum) IsValid() bool {
	switch e {
	case PatientLinkTypeEnumReplacedBy, PatientLinkTypeEnumReplaces, PatientLinkTypeEnumRefer, PatientLinkTypeEnumSeealso:
		return true
	}

	return false
}

// String ...
func (e PatientLinkTypeEnum) String() string {
	if e == PatientLinkTypeEnumReplacedBy {
		return "replaced-by"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *PatientLinkTypeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PatientLinkTypeEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Patient_LinkTypeEnum", str)
	}

	return nil
}

// MarshalGQL writes the patient link type to the supplied writer as a quoted string
func (e PatientLinkTypeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// MedicationStatementStatusEnum indicates the status of a medication statement
type MedicationStatementStatusEnum string

const (
	// MedicationStatementStatusEnumActive is The medication is still being taken.
	MedicationStatementStatusEnumActive MedicationStatementStatusEnum = "active"

	// MedicationStatementStatusEnumInActive is The medication is no longer being taken.
	MedicationStatementStatusEnumInActive MedicationStatementStatusEnum = "inactive"

	// MedicationStatementStatusEnumEnteredInError is Some of the actions that are implied by the medication statement may have occurred.
	MedicationStatementStatusEnumEnteredInError MedicationStatementStatusEnum = "entered-in-error"

	// MedicationStatementStatusEnumIntended is The medication may be taken at some time in the future.
	MedicationStatementStatusEnumIntended MedicationStatementStatusEnum = "intended"

	// MedicationStatementStatusEnumStopped is Actions implied by the statement have been permanently halted, before all of them occurred
	MedicationStatementStatusEnumStopped MedicationStatementStatusEnum = "stopped"

	// MedicationStatementStatusEnumOnHold is Actions implied by the statement have been temporarily halted, but are expected to continue later.
	MedicationStatementStatusEnumOnHold MedicationStatementStatusEnum = "on-hold"

	// MedicationStatementStatusEnumUnknown is The state of the medication use is not currently known.
	MedicationStatementStatusEnumUnknown MedicationStatementStatusEnum = "unknown"

	// MedicationStatementStatusEnumNotTaken is The medication was not consumed by the patient
	MedicationStatementStatusEnumNotTaken MedicationStatementStatusEnum = "not-taken"
)

// AllMedicationStatementStatusEnum is a list of all possible medication statements status
var AllMedicationStatementStatusEnum = []MedicationStatementStatusEnum{
	MedicationStatementStatusEnumActive,
	MedicationStatementStatusEnumInActive,
	MedicationStatementStatusEnumEnteredInError,
	MedicationStatementStatusEnumIntended,
	MedicationStatementStatusEnumStopped,
	MedicationStatementStatusEnumOnHold,
	MedicationStatementStatusEnumUnknown,
	MedicationStatementStatusEnumNotTaken,
}

// IsValid ...
func (e MedicationStatementStatusEnum) IsValid() bool {
	switch e {
	case MedicationStatementStatusEnumActive,
		MedicationStatementStatusEnumInActive,
		MedicationStatementStatusEnumEnteredInError,
		MedicationStatementStatusEnumIntended,
		MedicationStatementStatusEnumStopped,
		MedicationStatementStatusEnumOnHold,
		MedicationStatementStatusEnumUnknown,
		MedicationStatementStatusEnumNotTaken:
		return true
	}

	return false
}

// String ...
func (e MedicationStatementStatusEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *MedicationStatementStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MedicationStatementStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Patient_LinkTypeEnum", str)
	}

	return nil
}

// MarshalGQL writes the patient link type to the supplied writer as a quoted string
func (e MedicationStatementStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// MedicationStatusEnum indicates the medication status
type MedicationStatusEnum string

const (
	// MedicationStatusEnumActive is The medication is available for use
	MedicationStatusEnumActive MedicationStatusEnum = "active"

	// MedicationStatusEnumInActive is The medication is not available for use.
	MedicationStatusEnumInActive MedicationStatusEnum = "inactive"

	// MedicationStatusEnumEnteredInError is The medication was entered in error.
	MedicationStatusEnumEnteredInError MedicationStatusEnum = "entered-in-error"
)

// AllMedicationStatusEnum is a list of all possible medication status
var AllMedicationStatusEnum = []MedicationStatusEnum{
	MedicationStatusEnumActive,
	MedicationStatusEnumInActive,
	MedicationStatusEnumEnteredInError,
}

// IsValid ...
func (e MedicationStatusEnum) IsValid() bool {
	switch e {
	case MedicationStatusEnumActive,
		MedicationStatusEnumInActive,
		MedicationStatusEnumEnteredInError:
		return true
	}

	return false
}

// String ...
func (e MedicationStatusEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *MedicationStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MedicationStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Patient_LinkTypeEnum", str)
	}

	return nil
}

// MarshalGQL writes the patient link type to the supplied writer as a quoted string
func (e MedicationStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// MediaStatusEnum is a FHIR enum
type MediaStatusEnum string

const (
	// MediaStatusInProgress ...
	MediaStatusInProgress MediaStatusEnum = "in-progress"

	// MediaStatusNotDone ...
	MediaStatusNotDone MediaStatusEnum = "not-done"

	// MediaStatusEnumOnHold ...
	MediaStatusOnHold MediaStatusEnum = "on-hold"

	// MediaStatusEnumOfficial ...
	MediaStatusStopped MediaStatusEnum = "stopped"

	// MediaStatusCompleted ....
	MediaStatusCompleted MediaStatusEnum = "completed"

	// MediaStatusEnteredInError ...
	MediaStatusEnteredInError MediaStatusEnum = "entered-in-error"

	// MediaStatusUnknown ....
	MediaStatusUnknown MediaStatusEnum = "unknown"
)

// AllMediaStatusEnum ...
var AllMediaStatusEnum = []MediaStatusEnum{
	MediaStatusInProgress,
	MediaStatusNotDone,
	MediaStatusOnHold,
	MediaStatusStopped,
	MediaStatusCompleted,
	MediaStatusEnteredInError,
	MediaStatusUnknown,
}

// IsValid ...
func (e MediaStatusEnum) IsValid() bool {
	switch e {
	case MediaStatusInProgress, MediaStatusNotDone, MediaStatusOnHold,
		MediaStatusStopped, MediaStatusCompleted, MediaStatusEnteredInError, MediaStatusUnknown:
		return true
	}

	return false
}

// String ...
func (e MediaStatusEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *MediaStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MediaStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Composition_AttesterModeEnum", str)
	}

	return nil
}

// MarshalGQL writes the composition attester mode to the supplied writer as a quoted string
func (e MediaStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// DiagnosticReportStatusEnum is documented here http://hl7.org/fhir/ValueSet/diagnostic-report-status
type DiagnosticReportStatusEnum string

const (
	DiagnosticReportStatusRegistered     = "registered"
	DiagnosticReportStatusPartial        = "partial"
	DiagnosticReportStatusPreliminary    = "preliminary"
	DiagnosticReportStatusFinal          = "final"
	DiagnosticReportStatusAmended        = "amended"
	DiagnosticReportStatusCorrected      = "corrected"
	DiagnosticReportStatusAppended       = "appended"
	DiagnosticReportStatusCancelled      = "cancelled"
	DiagnosticReportStatusEnteredInError = "entered-in-error"
	DiagnosticReportStatusUnknown        = "unknown"
)

// AllDiagnosticReportStatusEnum ...
var AllDiagnosticReportStatusEnum = []DiagnosticReportStatusEnum{
	DiagnosticReportStatusRegistered,
	DiagnosticReportStatusPartial,
	DiagnosticReportStatusPreliminary,
	DiagnosticReportStatusFinal,
	DiagnosticReportStatusAmended,
	DiagnosticReportStatusCorrected,
	DiagnosticReportStatusAppended,
	DiagnosticReportStatusCancelled,
	DiagnosticReportStatusEnteredInError,
	DiagnosticReportStatusUnknown,
}

// IsValid ...
func (e DiagnosticReportStatusEnum) IsValid() bool {
	switch e {
	case DiagnosticReportStatusRegistered, DiagnosticReportStatusPartial, DiagnosticReportStatusPreliminary,
		DiagnosticReportStatusFinal, DiagnosticReportStatusAmended, DiagnosticReportStatusCorrected, DiagnosticReportStatusAppended,
		DiagnosticReportStatusCancelled, DiagnosticReportStatusEnteredInError, DiagnosticReportStatusUnknown:
		return true
	}

	return false
}

// String ...
func (e DiagnosticReportStatusEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *DiagnosticReportStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DiagnosticReportStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DiagnosticReportStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the composition attester mode to the supplied writer as a quoted string
func (e DiagnosticReportStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// ServiceRequestStatusEnum represents the possible statuses of a ServiceRequest
type ServiceRequestStatusEnum string

const (
	ServiceRequestStatusDraft          ServiceRequestStatusEnum = "draft"
	ServiceRequestStatusActive         ServiceRequestStatusEnum = "active"
	ServiceRequestStatusOnHold         ServiceRequestStatusEnum = "on-hold"
	ServiceRequestStatusRevoked        ServiceRequestStatusEnum = "revoked"
	ServiceRequestStatusCompleted      ServiceRequestStatusEnum = "completed"
	ServiceRequestStatusEnteredInError ServiceRequestStatusEnum = "entered-in-error"
	ServiceRequestStatusUnknown        ServiceRequestStatusEnum = "unknown"
)

// ServiceRequestIntentEnum represents the possible intents of a ServiceRequest
type ServiceRequestIntentEnum string

const (
	ServiceRequestIntentProposal      ServiceRequestIntentEnum = "proposal"
	ServiceRequestIntentPlan          ServiceRequestIntentEnum = "plan"
	ServiceRequestIntentDirective     ServiceRequestIntentEnum = "directive"
	ServiceRequestIntentOrder         ServiceRequestIntentEnum = "order"
	ServiceRequestIntentOriginalOrder ServiceRequestIntentEnum = "original-order"
	ServiceRequestIntentReflexOrder   ServiceRequestIntentEnum = "reflex-order"
	ServiceRequestIntentFillerOrder   ServiceRequestIntentEnum = "filler-order"
	ServiceRequestIntentInstanceOrder ServiceRequestIntentEnum = "instance-order"
	ServiceRequestIntentOption        ServiceRequestIntentEnum = "option"
)

// ServiceRequestPriorityEnum represents the possible priorities of a ServiceRequest
type ServiceRequestPriorityEnum string

const (
	ServiceRequestPriorityRoutine ServiceRequestPriorityEnum = "routine"
	ServiceRequestPriorityUrgent  ServiceRequestPriorityEnum = "urgent"
	ServiceRequestPriorityAsap    ServiceRequestPriorityEnum = "asap"
	ServiceRequestPriorityStat    ServiceRequestPriorityEnum = "stat"
)

func (s ServiceRequestPriorityEnum) IsValid() bool {
	switch s {
	case ServiceRequestPriorityRoutine, ServiceRequestPriorityUrgent, ServiceRequestPriorityAsap, ServiceRequestPriorityStat:
		return true
	}

	return false
}

// SubscriptionType defines the types of methods used to execute a subscription
type SubscriptionTypeEnum string

const (
	SubscriptionTypeRestHook  SubscriptionTypeEnum = "rest-hook"
	SubscriptionTypeWebSocket SubscriptionTypeEnum = "websocket"
	SubscriptionTypeEmail     SubscriptionTypeEnum = "email"
	SubscriptionTypeSMS       SubscriptionTypeEnum = "sms"
	SubscriptionTypeMessage   SubscriptionTypeEnum = "message"
)

// DocumentReferenceStatusEnum is a FHIR enum for document reference statuses
type DocumentReferenceStatusEnum string

const (
	DocumentReferenceStatusEnumCurrent        DocumentReferenceStatusEnum = "current"
	DocumentReferenceStatusEnumSuperseded     DocumentReferenceStatusEnum = "superseded"
	DocumentReferenceStatusEnumEnteredInError DocumentReferenceStatusEnum = "entered-in-error"
)

// AllDocumentReferenceStatusEnum lists all document reference statuses
var AllDocumentReferenceStatusEnum = []DocumentReferenceStatusEnum{
	DocumentReferenceStatusEnumCurrent,
	DocumentReferenceStatusEnumSuperseded,
	DocumentReferenceStatusEnumEnteredInError,
}

// IsValid checks if the enum value is valid
func (e DocumentReferenceStatusEnum) IsValid() bool {
	switch e {
	case DocumentReferenceStatusEnumCurrent, DocumentReferenceStatusEnumSuperseded, DocumentReferenceStatusEnumEnteredInError:
		return true
	}

	return false
}

// String converts the enum to its string representation
func (e DocumentReferenceStatusEnum) String() string {
	return string(e)
}

// DocumentRelationshipTypeEnum is a FHIR enum for document relationship types
type DocumentRelationshipTypeEnum string

const (
	DocumentRelationshipTypeEnumReplaces   DocumentRelationshipTypeEnum = "replaces"
	DocumentRelationshipTypeEnumTransforms DocumentRelationshipTypeEnum = "transforms"
	DocumentRelationshipTypeEnumSigns      DocumentRelationshipTypeEnum = "signs"
	DocumentRelationshipTypeEnumAppends    DocumentRelationshipTypeEnum = "appends"
)

// AllDocumentRelationshipTypeEnum lists all document relationship types
var AllDocumentRelationshipTypeEnum = []DocumentRelationshipTypeEnum{
	DocumentRelationshipTypeEnumReplaces,
	DocumentRelationshipTypeEnumTransforms,
	DocumentRelationshipTypeEnumSigns,
	DocumentRelationshipTypeEnumAppends,
}

// IsValid checks if the enum value is valid
func (e DocumentRelationshipTypeEnum) IsValid() bool {
	switch e {
	case DocumentRelationshipTypeEnumReplaces, DocumentRelationshipTypeEnumTransforms, DocumentRelationshipTypeEnumSigns, DocumentRelationshipTypeEnumAppends:
		return true
	}

	return false
}

// String converts the enum to its string representation
func (e DocumentRelationshipTypeEnum) String() string {
	return string(e)
}

// ConsentStatusEnum a type enum that represents a Consent Status field of consent resource
type ConsentStatusEnum string

const (
	ConsentStatusActive   ConsentStatusEnum = "active"
	ConsentStatusInactive ConsentStatusEnum = "inactive"
)

// IsValid ...
func (c ConsentStatusEnum) IsValid() bool {
	switch c {
	case ConsentStatusActive, ConsentStatusInactive:
		return true
	}

	return false
}

// String converts status to string
func (c ConsentStatusEnum) String() string {
	return string(c)
}

// MarshalGQL writes the consent status as a quoted string
func (c ConsentStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a consent status enum
func (c *ConsentStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*c = ConsentStatusEnum(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid ConsentStatus Enum", str)
	}

	return nil
}

// ConsentProvisionTypeEnum a type enum tha represents a Consent Provision field of consent resource
type ConsentDecisionEnum string

const (
	ConsentDecisionDeny   ConsentDecisionEnum = "deny"
	ConsentDecisionPermit ConsentDecisionEnum = "permit"
)

// IsValid ...
func (c ConsentDecisionEnum) IsValid() bool {
	switch c {
	case ConsentDecisionDeny, ConsentDecisionPermit:
		return true
	}

	return false
}

// String converts consent provision type to string
func (c ConsentDecisionEnum) String() string {
	return string(c)
}

// MarshalGQL writes the consent provision type as a quoted string
func (c ConsentDecisionEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a consent provision type enum
func (c *ConsentDecisionEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*c = ConsentDecisionEnum(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid ConsentProvisionTypeEnum Enum", str)
	}

	return nil
}

// ConsentDataMeaningEnum represents the meaning of consent data
type ConsentDataMeaningEnum string

const (
	ConsentDataMeaningInstance   ConsentDataMeaningEnum = "instance"
	ConsentDataMeaningRelated    ConsentDataMeaningEnum = "related"
	ConsentDataMeaningDependents ConsentDataMeaningEnum = "dependents"
	ConsentDataMeaningAuthoredBy ConsentDataMeaningEnum = "authoredby"
)

// IsValid checks if the consent data meaning is valid
func (c ConsentDataMeaningEnum) IsValid() bool {
	switch c {
	case ConsentDataMeaningInstance, ConsentDataMeaningRelated, ConsentDataMeaningDependents, ConsentDataMeaningAuthoredBy:
		return true
	}

	return false
}

// String converts the consent data meaning to string
func (c ConsentDataMeaningEnum) String() string {
	return string(c)
}

// MarshalGQL writes the consent data meaning as a quoted string
func (c ConsentDataMeaningEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a JSON and converts it to a consent data meaning enum
func (c *ConsentDataMeaningEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*c = ConsentDataMeaningEnum(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid ConsentDataMeaningEnum", str)
	}

	return nil
}

// QuestionnaireResponseStatusEnum a type enum tha represents a questionnaire response status field of questionnaire response
type QuestionnaireResponseStatusEnum string

const (
	QuestionnaireResponseStatusEnumCompleted QuestionnaireResponseStatusEnum = "completed"
)

// IsValid ...
func (c QuestionnaireResponseStatusEnum) IsValid() bool {
	return c == QuestionnaireResponseStatusEnumCompleted
}

// String converts questionnaire response status type to string
func (c QuestionnaireResponseStatusEnum) String() string {
	return string(c)
}

// MarshalGQL writes the questionnaire response status type as a quoted string
func (c QuestionnaireResponseStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a questionnaire response status type enum
func (c *QuestionnaireResponseStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*c = QuestionnaireResponseStatusEnum(strings.ReplaceAll(str, "_", "-"))
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid QuestionnaireResponseStatus Enum", str)
	}

	return nil
}

// MedicationRequestStatus is used to define the allowed statuses of a medicaton request
type MedicationRequestStatus string

// NOTE: This list is not exhaustive. More status can be found here: https://hl7.org/fhir/r4/medicationrequest-definitions.html#MedicationRequest.status
const (
	ActiveMedicationStatus    MedicationRequestStatus = "active"
	CompletedMedicationStatus MedicationRequestStatus = "completed"
	StoppedMedicationStatus   MedicationRequestStatus = "stopped"
	CancelledMedicationStatus MedicationRequestStatus = "cancelled"
	OnHoldMedicationStatus    MedicationRequestStatus = "on_hold"
)

// AllMedicationStatus is a list of all possible medication request statuses
var AllMedicationStatus = []MedicationRequestStatus{
	CompletedMedicationStatus,
	OnHoldMedicationStatus,
	CancelledMedicationStatus,
	StoppedMedicationStatus,
	ActiveMedicationStatus,
}

// IsValid ...
func (t MedicationRequestStatus) IsValid() bool {
	switch t {
	case CompletedMedicationStatus, OnHoldMedicationStatus,
		CancelledMedicationStatus, StoppedMedicationStatus, ActiveMedicationStatus:
		return true
	}

	return false
}

// String ...
func (t MedicationRequestStatus) String() string {
	if t == OnHoldMedicationStatus {
		return "on-hold"
	}

	return string(t)
}

// UnmarshalGQL ...
func (t *MedicationRequestStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*t = MedicationRequestStatus(str)
	if !t.IsValid() {
		return fmt.Errorf("%s is not a valid task status", str)
	}

	return nil
}

// MarshalGQL writes the task status to the supplied writer as a quoted string
func (t MedicationRequestStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}

// MedicationRequestIntent is used to define the allowed intents of a medication request
type MedicationRequestIntent string

const (
	ProposalMedicationIntent      MedicationRequestIntent = "proposal"
	PlanMedicationIntent          MedicationRequestIntent = "plan"
	OrderMedicationIntent         MedicationRequestIntent = "order"
	OriginalOrderMedicationIntent MedicationRequestIntent = "original-order"
	ReflexOrderMedicationIntent   MedicationRequestIntent = "reflex-order"
	FillerOrderMedicationIntent   MedicationRequestIntent = "filler-order"
	InstanceOrderMedicationIntent MedicationRequestIntent = "instance-order"
	OptionMedicationIntent        MedicationRequestIntent = "option"
)

// AllMedicationIntents is a list of all possible medication request intents
var AllMedicationIntents = []MedicationRequestIntent{
	ProposalMedicationIntent,
	PlanMedicationIntent,
	OrderMedicationIntent,
	OriginalOrderMedicationIntent,
	ReflexOrderMedicationIntent,
	FillerOrderMedicationIntent,
	InstanceOrderMedicationIntent,
	OptionMedicationIntent,
}

// IsValid checks if the medication request intent is valid
func (t MedicationRequestIntent) IsValid() bool {
	switch t {
	case ProposalMedicationIntent, PlanMedicationIntent, OrderMedicationIntent,
		OriginalOrderMedicationIntent, ReflexOrderMedicationIntent,
		FillerOrderMedicationIntent, InstanceOrderMedicationIntent, OptionMedicationIntent:
		return true
	}

	return false
}

// String returns the string representation of the medication request intent
func (t MedicationRequestIntent) String() string {
	return string(t)
}

// UnmarshalGQL unmarshals a GraphQL input into a MedicationRequestIntent
func (t *MedicationRequestIntent) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*t = MedicationRequestIntent(str)
	if !t.IsValid() {
		return fmt.Errorf("%s is not a valid medication request intent", str)
	}

	return nil
}

// MarshalGQL writes the medication request intent to the supplied writer as a quoted string
func (t MedicationRequestIntent) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}

// IssueType is documented here http://hl7.org/fhir/ValueSet/issue-type
// IssueType is a FHIR enum
type IssueType string

const (
	IssueTypeInvalid         IssueType = "invalid"
	IssueTypeStructure       IssueType = "structure"
	IssueTypeRequired        IssueType = "required"
	IssueTypeValue           IssueType = "value"
	IssueTypeInvariant       IssueType = "invariant"
	IssueTypeSecurity        IssueType = "security"
	IssueTypeLogin           IssueType = "login"
	IssueTypeUnknown         IssueType = "unknown"
	IssueTypeExpired         IssueType = "expired"
	IssueTypeForbidden       IssueType = "forbidden"
	IssueTypeSuppressed      IssueType = "suppressed"
	IssueTypeProcessing      IssueType = "processing"
	IssueTypeNotSupported    IssueType = "not-supported"
	IssueTypeDuplicate       IssueType = "duplicate"
	IssueTypeMultipleMatches IssueType = "multiple-matches"
	IssueTypeNotFound        IssueType = "not-found"
	IssueTypeDeleted         IssueType = "deleted"
	IssueTypeTooLong         IssueType = "too-long"
	IssueTypeCodeInvalid     IssueType = "code-invalid"
	IssueTypeExtension       IssueType = "extension"
	IssueTypeTooCostly       IssueType = "too-costly"
	IssueTypeBusinessRule    IssueType = "business-rule"
	IssueTypeConflict        IssueType = "conflict"
	IssueTypeTransient       IssueType = "transient"
	IssueTypeLockError       IssueType = "lock-error"
	IssueTypeNoStore         IssueType = "no-store"
	IssueTypeException       IssueType = "exception"
	IssueTypeTimeout         IssueType = "timeout"
	IssueTypeIncomplete      IssueType = "incomplete"
	IssueTypeThrottled       IssueType = "throttled"
	IssueTypeInformational   IssueType = "informational"
)

// AllIssueTypes ...
var AllIssueTypes = []IssueType{
	IssueTypeInvalid,
	IssueTypeStructure,
	IssueTypeRequired,
	IssueTypeValue,
	IssueTypeInvariant,
	IssueTypeSecurity,
	IssueTypeLogin,
	IssueTypeUnknown,
	IssueTypeExpired,
	IssueTypeForbidden,
	IssueTypeSuppressed,
	IssueTypeProcessing,
	IssueTypeNotSupported,
	IssueTypeDuplicate,
	IssueTypeMultipleMatches,
	IssueTypeNotFound,
	IssueTypeDeleted,
	IssueTypeTooLong,
	IssueTypeCodeInvalid,
	IssueTypeExtension,
	IssueTypeTooCostly,
	IssueTypeBusinessRule,
	IssueTypeConflict,
	IssueTypeTransient,
	IssueTypeLockError,
	IssueTypeNoStore,
	IssueTypeException,
	IssueTypeTimeout,
	IssueTypeIncomplete,
	IssueTypeThrottled,
	IssueTypeInformational,
}

// IsValid ...
func (e IssueType) IsValid() bool {
	switch e {
	case IssueTypeInvalid, IssueTypeStructure, IssueTypeRequired, IssueTypeValue, IssueTypeInvariant,
		IssueTypeSecurity, IssueTypeLogin, IssueTypeUnknown, IssueTypeExpired, IssueTypeForbidden,
		IssueTypeSuppressed, IssueTypeProcessing, IssueTypeNotSupported, IssueTypeDuplicate,
		IssueTypeMultipleMatches, IssueTypeNotFound, IssueTypeDeleted, IssueTypeTooLong,
		IssueTypeCodeInvalid, IssueTypeExtension, IssueTypeTooCostly, IssueTypeBusinessRule,
		IssueTypeConflict, IssueTypeTransient, IssueTypeLockError, IssueTypeNoStore,
		IssueTypeException, IssueTypeTimeout, IssueTypeIncomplete, IssueTypeThrottled,
		IssueTypeInformational:
		return true
	}

	return false
}

// String ...
func (e IssueType) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *IssueType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IssueType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IssueType", str)
	}

	return nil
}

// MarshalGQL writes the issue type to the supplied writer as a quoted string
func (e IssueType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// IssueSeverity is a FHIR enum
type IssueSeverity string

const (
	IssueSeverityFatal       IssueSeverity = "fatal"
	IssueSeverityError       IssueSeverity = "error"
	IssueSeverityWarning     IssueSeverity = "warning"
	IssueSeverityInformation IssueSeverity = "information"
)

// AllIssueSeverities ...
var AllIssueSeverities = []IssueSeverity{
	IssueSeverityFatal,
	IssueSeverityError,
	IssueSeverityWarning,
	IssueSeverityInformation,
}

// IsValid ...
func (e IssueSeverity) IsValid() bool {
	switch e {
	case IssueSeverityFatal, IssueSeverityError, IssueSeverityWarning, IssueSeverityInformation:
		return true
	}

	return false
}

// String ...
func (e IssueSeverity) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *IssueSeverity) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IssueSeverity(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IssueSeverity", str)
	}

	return nil
}

// MarshalGQL writes the issue severity to the supplied writer as a quoted string
func (e IssueSeverity) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Proceduer enums
type ProcedureStatusEnum string

const (
	ProcedureStatusPreparation    ProcedureStatusEnum = "preparation"
	ProcedureStatusInProgress     ProcedureStatusEnum = "in-progress"
	ProcedureStatusNotDone        ProcedureStatusEnum = "not-done"
	ProcedureStatusOnHold         ProcedureStatusEnum = "on-hold"
	ProcedureStatusStopped        ProcedureStatusEnum = "stopped"
	ProcedureStatusCompleted      ProcedureStatusEnum = "completed"
	ProcedureStatusEnteredInError ProcedureStatusEnum = "entered-in-error"
	ProcedureStatusUnknown        ProcedureStatusEnum = "unknown"
)

func (p ProcedureStatusEnum) IsValid() bool {
	switch p {
	case ProcedureStatusCompleted, ProcedureStatusEnteredInError, ProcedureStatusInProgress,
		ProcedureStatusNotDone, ProcedureStatusOnHold, ProcedureStatusPreparation, ProcedureStatusStopped,
		ProcedureStatusUnknown:
		return true
	}

	return false
}

func (p ProcedureStatusEnum) Code() string {
	switch p {
	case ProcedureStatusPreparation:
		return "preparation"
	case ProcedureStatusInProgress:
		return "in-progress"
	case ProcedureStatusNotDone:
		return "not-done"
	case ProcedureStatusOnHold:
		return "on-hold"
	case ProcedureStatusStopped:
		return "stopped"
	case ProcedureStatusCompleted:
		return "completed"
	case ProcedureStatusEnteredInError:
		return "entered-in-error"
	case ProcedureStatusUnknown:
		return "unknown"
	}

	return "<unknown>"
}

func (p ProcedureStatusEnum) Display() string {
	switch p {
	case ProcedureStatusPreparation:
		return "Preparation"
	case ProcedureStatusInProgress:
		return "In Progress"
	case ProcedureStatusNotDone:
		return "Not Done"
	case ProcedureStatusOnHold:
		return "On Hold"
	case ProcedureStatusStopped:
		return "Stopped"
	case ProcedureStatusCompleted:
		return "Completed"
	case ProcedureStatusEnteredInError:
		return "Entered In Error"
	case ProcedureStatusUnknown:
		return "Unknown"
	}

	return "<unknown>"
}

// MedicationDispenseStatusCodes is documented here http://hl7.org/fhir/ValueSet/medicationdispense-status
type MedicationDispenseStatusCodes string

const (
	MedicationDispenseStatusCodesPreparation    MedicationDispenseStatusCodes = "preparation"
	MedicationDispenseStatusCodesInProgress     MedicationDispenseStatusCodes = "in-progress"
	MedicationDispenseStatusCodesCancelled      MedicationDispenseStatusCodes = "cancelled"
	MedicationDispenseStatusCodesOnHold         MedicationDispenseStatusCodes = "on-hold"
	MedicationDispenseStatusCodesCompleted      MedicationDispenseStatusCodes = "completed"
	MedicationDispenseStatusCodesEnteredInError MedicationDispenseStatusCodes = "entered-in-error"
	MedicationDispenseStatusCodesStopped        MedicationDispenseStatusCodes = "stopped"
	MedicationDispenseStatusCodesDeclined       MedicationDispenseStatusCodes = "declined"
	MedicationDispenseStatusCodesUnknown        MedicationDispenseStatusCodes = "unknown"
)

func (code MedicationDispenseStatusCodes) Code() string {
	switch code {
	case MedicationDispenseStatusCodesPreparation:
		return "preparation"
	case MedicationDispenseStatusCodesInProgress:
		return "in-progress"
	case MedicationDispenseStatusCodesCancelled:
		return "cancelled"
	case MedicationDispenseStatusCodesOnHold:
		return "on-hold"
	case MedicationDispenseStatusCodesCompleted:
		return "completed"
	case MedicationDispenseStatusCodesEnteredInError:
		return "entered-in-error"
	case MedicationDispenseStatusCodesStopped:
		return "stopped"
	case MedicationDispenseStatusCodesDeclined:
		return "declined"
	case MedicationDispenseStatusCodesUnknown:
		return "unknown"
	}

	return "unknown"
}

// ParticipantStatusEnum represents the appointment participant status
type ParticipantStatusEnum string

const (
	Accepted    ParticipantStatusEnum = "accepted"
	Declined    ParticipantStatusEnum = "declined"
	Tenative    ParticipantStatusEnum = "tenative"
	NeedsAction ParticipantStatusEnum = "needs-action"
)

func (p ParticipantStatusEnum) IsValid() bool {
	switch p {
	case Accepted, Declined, Tenative, NeedsAction:
		return true
	}

	return false
}

func (p ParticipantStatusEnum) Code() string {
	switch p {
	case Accepted:
		return "accepted"
	case Declined:
		return "declined"
	case Tenative:
		return "tenative"
	case NeedsAction:
		return "needs-action"
	}

	return "unknown"
}

func (code MedicationDispenseStatusCodes) Display() string {
	switch code {
	case MedicationDispenseStatusCodesPreparation:
		return "Preparation"
	case MedicationDispenseStatusCodesInProgress:
		return "In Progress"
	case MedicationDispenseStatusCodesCancelled:
		return "Cancelled"
	case MedicationDispenseStatusCodesOnHold:
		return "On Hold"
	case MedicationDispenseStatusCodesCompleted:
		return "Completed"
	case MedicationDispenseStatusCodesEnteredInError:
		return "Entered in Error"
	case MedicationDispenseStatusCodesStopped:
		return "Stopped"
	case MedicationDispenseStatusCodesDeclined:
		return "Declined"
	case MedicationDispenseStatusCodesUnknown:
		return "Unknown"
	}

	return "unknown"
}

func (p ParticipantStatusEnum) Display() string {
	switch p {
	case Accepted:
		return "Accepted"
	case Declined:
		return "Declined"
	case Tenative:
		return "Tenative"
	case NeedsAction:
		return "Needs Action"
	}

	return "Unknown"
}

type PriorityEnum string

const (
	ASAP      PriorityEnum = "A"
	Emergency PriorityEnum = "EM"
	Routine   PriorityEnum = "R"
	Urgent    PriorityEnum = "UR"
)

func (p PriorityEnum) IsValid() bool {
	switch p {
	case ASAP, Emergency, Routine, Urgent:
		return true
	}

	return false
}

func (p PriorityEnum) Code() string {
	switch p {
	case ASAP:
		return "A"
	case Emergency:
		return "EM"
	case Routine:
		return "R"
	case Urgent:
		return "UR"
	}

	return "unknown"
}

func (p PriorityEnum) Display() string {
	switch p {
	case ASAP:
		return "ASAP"
	case Emergency:
		return "emergency"
	case Routine:
		return "routine"
	case Urgent:
		return "urgent"
	}

	return "Unknown"
}

// AppointmentTypeEnum
type AppointmentTypeEnum string

const (
	RoutineAppointment   AppointmentTypeEnum = "ROUTINE"
	CheckUpAppointment   AppointmentTypeEnum = "CHECKUP"
	FollowUpAppointment  AppointmentTypeEnum = "FOLLOWUP"
	EmergencyAppointment AppointmentTypeEnum = "EMERGENCY"
	WalkInAppointment    AppointmentTypeEnum = "WALKIN"
)

func (a AppointmentTypeEnum) IsValid() bool {
	switch a {
	case RoutineAppointment, CheckUpAppointment, FollowUpAppointment, EmergencyAppointment, WalkInAppointment:
		return true
	}

	return false
}

func (a AppointmentTypeEnum) Code() string {
	switch a {
	case RoutineAppointment:
		return "ROUTINE"
	case CheckUpAppointment:
		return "CHECKUP"
	case EmergencyAppointment:
		return "EMERGENCY"
	case WalkInAppointment:
		return "WALKIN"
	}

	return "Unknown"
}

func (a AppointmentTypeEnum) Display() string {
	switch a {
	case RoutineAppointment:
		return "Routine"
	case CheckUpAppointment:
		return "Check up"
	case EmergencyAppointment:
		return "Emergency"
	case WalkInAppointment:
		return "Walk In"
	}

	return "Unknown"
}

// RequestStatus is documented here http://hl7.org/fhir/ValueSet/request-status
type RequestStatus string

const (
	RequestStatusDraft          RequestStatus = "draft"
	RequestStatusActive         RequestStatus = "active"
	RequestStatusOnHold         RequestStatus = "on-hold"
	RequestStatusRevoked        RequestStatus = "revoked"
	RequestStatusCompleted      RequestStatus = "completed"
	RequestStatusEnteredInError RequestStatus = "entered-in-error"
	RequestStatusUnknown        RequestStatus = "unknown"
)

type CarePlanIntent string

const (
	CarePlanIntentProposal      CarePlanIntent = "proposal"
	CarePlanIntentPlan          CarePlanIntent = "intent"
	CarePlanIntentDirective     CarePlanIntent = "directive"
	CarePlanIntentOrder         CarePlanIntent = "intent"
	CarePlanIntentOriginalOrder CarePlanIntent = "original-order"
	CarePlanIntentReflexOrder   CarePlanIntent = "reflex-order"
	CarePlanIntentFillerOrder   CarePlanIntent = "filler-order"
	CarePlanIntentInstanceOrder CarePlanIntent = "instance-order"
	CarePlanIntentOption        CarePlanIntent = "option"
)

func (c FHIRCodeableConcept) ConditionCode() *FHIRCodeableReferenceInput {
	for _, code := range c.Coding {
		if code.Code == nil {
			continue
		}

		codeStr := *code.Code

		return &FHIRCodeableReferenceInput{
			Concept: &FHIRCodeableConceptInput{
				Coding: []*FHIRCodingInput{
					{
						Code:    codeStr,
						System:  code.System,
						Display: code.Display,
						ID:      *code.ID,
					},
				},
			},
		}
	}

	return nil
}
