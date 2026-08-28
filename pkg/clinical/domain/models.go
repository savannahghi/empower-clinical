package domain

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
)

// Dummy ..
type Dummy struct {
	ID string `json:"id"`
}

// IsEntity ...
func (d Dummy) IsEntity() {}

// IsNode ..
func (d *Dummy) IsNode() {}

// SetID sets the trace's ID
func (d *Dummy) SetID(id string) {
	d.ID = id
}

// NameInput is used to input patient names.
type NameInput struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	OtherNames *string `json:"otherNames"`
}

// IdentificationDocument is used to input e.g National ID or passport document
// numbers at patient registration.
type IdentificationDocument struct {
	DocumentType     IDDocumentType         `json:"documentType"`
	DocumentNumber   string                 `json:"documentNumber"`
	Title            *string                `json:"title,omitempty"`
	ImageContentType *enumutils.ContentType `json:"imageContentType,omitempty"`
	ImageBase64      *string                `json:"imageBase64,omitempty"`
}

// PhoneNumberInput is used to input phone numbers.
type PhoneNumberInput struct {
	Msisdn             string `json:"msisdn"`
	VerificationCode   string `json:"verificationCode"`
	IsUssd             bool   `json:"isUSSD"`
	CommunicationOptIn bool   `json:"communicationOptIn"`
}

// PhotoInput is used to upload patient photos.
type PhotoInput struct {
	PhotoContentType enumutils.ContentType `json:"photoContentType"`
	PhotoBase64data  string                `json:"photoBase64data"`
	PhotoFilename    string                `json:"photoFilename"`
}

// EmailInput is used to register patient emails.
type EmailInput struct {
	Email              *string `json:"email"`
	CommunicationOptIn bool    `json:"communicationOptIn"`
}

// PhysicalAddress is used to record a precise physical address.
type PhysicalAddress struct {
	MapsCode        string  `json:"mapsCode"`
	PhysicalAddress *string `json:"physicalAddress"`
}

// PostalAddress is used to record patient's postal addresses
type PostalAddress struct {
	PostalAddress *string `json:"postalAddress"`
	PostalCode    string  `json:"postalCode"`
}

// SimplePatientRegistrationInput provides a simplified API to support registration
// of patients.
type SimplePatientRegistrationInput struct {
	ID                      string                    `json:"id,omitempty"`
	Names                   []*NameInput              `json:"names,omitempty"`
	IdentificationDocuments []*IdentificationDocument `json:"identificationDocuments,omitempty"`
	BirthDate               *scalarutils.Date         `json:"birthDate,omitempty"`
	PhoneNumbers            []*PhoneNumberInput       `json:"phoneNumbers,omitempty"`
	Photos                  []*PhotoInput             `json:"photos,omitempty"`
	Emails                  []*EmailInput             `json:"emails,omitempty"`
	PhysicalAddresses       []*PhysicalAddress        `json:"physicalAddresses,omitempty"`
	PostalAddresses         []*PostalAddress          `json:"postalAddresses,omitempty"`
	Gender                  string                    `json:"gender,omitempty"`
	Active                  bool                      `json:"active,omitempty"`
	MaritalStatus           MaritalStatus             `json:"maritalStatus,omitempty"`
	Languages               []enumutils.Language      `json:"languages,omitempty"`
	ReplicateUSSD           bool                      `json:"replicate_ussd,omitempty"`
}

// BreakGlassEpisodeCreationInput is used to start emergency episodes via a
// break glass protocol
type BreakGlassEpisodeCreationInput struct {
	PatientID       string `json:"patientID" firestore:"patientID"`
	ProviderCode    string `json:"providerCode" firestore:"providerCode"`
	PractitionerUID string `json:"practitionerUID" firestore:"practitionerUID"`
	// ProviderPhone is the provider phone number
	ProviderPhone string `json:"providerPhone" firestore:"providerPhone"`
	Otp           string `json:"otp" firestore:"otp"`
	FullAccess    bool   `json:"fullAccess" firestore:"fullAccess"`
	// PatientPhone is the patient phone number used to send alert to patient
	PatientPhone string `json:"patient_phone" firestore:"patient_phone"`
}

// OTPEpisodeCreationInput is used to start patient visits via OTP
type OTPEpisodeCreationInput struct {
	PatientID    string `json:"patientID"`
	ProviderCode string `json:"providerCode"`
	Msisdn       string `json:"msisdn"`
	Otp          string `json:"otp"`
	FullAccess   bool   `json:"fullAccess"`
}

// OTPEpisodeUpgradeInput is used to upgrade existing open episodes
type OTPEpisodeUpgradeInput struct {
	EpisodeID string `json:"episodeID"`
	Msisdn    string `json:"msisdn"`
	Otp       string `json:"otp"`
}

// SimpleNHIFInput adds NHIF membership details as an extra identifier.
type SimpleNHIFInput struct {
	PatientID             string                 `json:"patientID"`
	MembershipNumber      string                 `json:"membershipNumber"`
	FrontImageBase64      *string                `json:"frontImageBase64"`
	FrontImageContentType *enumutils.ContentType `json:"frontImageContentType"`
	RearImageBase64       *string                `json:"rearImageBase64"`
	RearImageContentType  *enumutils.ContentType `json:"rearImageContentType"`
}

// SimpleNextOfKinInput is used to add next of kin to a patient.
type SimpleNextOfKinInput struct {
	PatientID         string              `json:"patientID"`
	Names             []*NameInput        `json:"names"`
	PhoneNumbers      []*PhoneNumberInput `json:"phoneNumbers"`
	Emails            []*EmailInput       `json:"emails"`
	PhysicalAddresses []*PhysicalAddress  `json:"physicalAddresses"`
	PostalAddresses   []*PostalAddress    `json:"postalAddresses"`
	Gender            string              `json:"gender"`
	Relationship      RelationshipType    `json:"relationship"`
	Active            bool                `json:"active"`
	BirthDate         scalarutils.Date    `json:"birthDate"`
}

// Reference defines references to other FHIR resources.
type Reference struct {
	ID         *string        `json:"id,omitempty"`
	Reference  string         `json:"reference,omitempty"`
	Type       string         `json:"type,omitempty"`
	Identifier FHIRIdentifier `json:"identifier,omitempty"`
	Display    *string        `json:"display,omitempty"`
}

// PatientPayload is used to return patient records and ancillary data after
// mutations.
type PatientPayload struct {
	PatientRecord *FHIRPatient `json:"patientRecord,omitempty"`
}

// RetirePatientInput is used to retire patient records.
type RetirePatientInput struct {
	ID string `json:"id"`
}

// PatientExtraInformationInput is used to update patient records metadata.
type PatientExtraInformationInput struct {
	PatientID     string                `json:"patientID"`
	MaritalStatus *MaritalStatus        `json:"maritalStatus"`
	Languages     []*enumutils.Language `json:"languages"`
	Emails        []*EmailInput         `json:"emails"`
}

// USSDNextOfKinCreationInput is used to register next of kin via USSD.
type USSDNextOfKinCreationInput struct {
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	OtherNames string    `json:"otherNames"`
	BirthDate  time.Time `json:"birthDate"`
	Gender     string    `json:"gender"`
	Active     bool      `json:"active"`
	ParentID   string    `json:"parentID"`
}

// IDDocumentType is an internal code system for identification document types.
type IDDocumentType string

// ID type constants
const (
	// IDDocumentTypeNationalID represents a National ID
	IDDocumentTypeNationalID IDDocumentType = "national_id"
	// IDDocumentTypePassport represents a Passport ID
	IDDocumentTypePassport IDDocumentType = "passport"
	// IDDocumentTypeAlienID represents an Alien ID
	IDDocumentTypeAlienID IDDocumentType = "alien_id"
	// IDDocumentTypeHealthID represents a Health ID
	IDDocumentTypeHealthID IDDocumentType = "slade_health_id"
	// IDDocumentTypePatientNumber represents a Patient Number
	IDDocumentTypePatientNumber IDDocumentType = "patient_number"
	// IDDocumentTypeMilitaryID represents a Military ID
	IDDocumentTypeMilitaryID IDDocumentType = "military_id"
	// IDDocumentTypeNationalHospitalInsuranceFundID represents a NHIF ID
	IDDocumentTypeNationalHospitalInsuranceFundID IDDocumentType = "nhif_id"
	// IDDocumentTypeSmartMemberNumber represents a Smart Member Number
	IDDocumentTypeSmartMemberNumber IDDocumentType = "smart_member_number"
	// IDDocumentTypeDrChronoChartID represents a Dr Chrono Chart ID
	IDDocumentTypeDrChronoChartID IDDocumentType = "dr_chrono_chart_id"
	// IDDocumentTypeFHIRPatientID represents a FHIR Patient ID
	IDDocumentTypeFHIRPatientID IDDocumentType = "fhir_patient_id"
	// IDDocumentTypeERPCustomerID represents an ERP Customer ID
	IDDocumentTypeERPCustomerID IDDocumentType = "erp_customer_id"
	// IDDocumentTypeComprehensiveCareClinicNumber represents a CCC Number
	IDDocumentTypeComprehensiveCareClinicNumber IDDocumentType = "ccc_number"
	// IDDocumentTypeRefugeeID represents a Refugee ID
	IDDocumentTypeRefugeeID IDDocumentType = "refugee_id"
	// IDDocumentTypeBirthCertificateNumber represents a Birth Certificate Number
	IDDocumentTypeBirthCertificateNumber IDDocumentType = "birth_certificate_number"
	// IDDocumentTypeMandateNumber represents a Mandate Number
	IDDocumentTypeMandateNumber IDDocumentType = "mandate_number"
	// IDDocumentTypeClientRegistryNumber represents a Client Registry Number
	IDDocumentTypeClientRegistryNumber IDDocumentType = "client_registry_number"
)

// AllIDDocumentType is a list of known ID types
var AllIDDocumentType = []IDDocumentType{
	IDDocumentTypeNationalID,
	IDDocumentTypePassport,
	IDDocumentTypeAlienID,
	IDDocumentTypeHealthID,
	IDDocumentTypePatientNumber,
	IDDocumentTypeMilitaryID,
	IDDocumentTypeNationalHospitalInsuranceFundID,
	IDDocumentTypeSmartMemberNumber,
	IDDocumentTypeDrChronoChartID,
	IDDocumentTypeFHIRPatientID,
	IDDocumentTypeERPCustomerID,
	IDDocumentTypeComprehensiveCareClinicNumber,
	IDDocumentTypeRefugeeID,
	IDDocumentTypeBirthCertificateNumber,
	IDDocumentTypeMandateNumber,
	IDDocumentTypeClientRegistryNumber,
}

// IsValid checks if the given ID type is valid
func (e IDDocumentType) IsValid() bool {
	switch e {
	case IDDocumentTypeNationalID,
		IDDocumentTypePassport,
		IDDocumentTypeAlienID,
		IDDocumentTypeHealthID,
		IDDocumentTypePatientNumber,
		IDDocumentTypeMilitaryID,
		IDDocumentTypeNationalHospitalInsuranceFundID,
		IDDocumentTypeSmartMemberNumber,
		IDDocumentTypeDrChronoChartID,
		IDDocumentTypeFHIRPatientID,
		IDDocumentTypeERPCustomerID,
		IDDocumentTypeComprehensiveCareClinicNumber,
		IDDocumentTypeRefugeeID,
		IDDocumentTypeBirthCertificateNumber,
		IDDocumentTypeMandateNumber,
		IDDocumentTypeClientRegistryNumber:
		return true
	}

	return false
}

// String ...
func (e IDDocumentType) String() string {
	return strings.ReplaceAll(string(e), "_", "-")
}

// GetDisplayName returns the human-readable display name for an ID type.
func (e IDDocumentType) GetDisplayName() string {
	normalized := IDDocumentType(e.String())

	mapping := map[IDDocumentType]string{
		IDDocumentType(IDDocumentTypeNationalID.String()):                      "National ID",
		IDDocumentType(IDDocumentTypePassport.String()):                        "Passport Number",
		IDDocumentType(IDDocumentTypeAlienID.String()):                         "Alien ID",
		IDDocumentType(IDDocumentTypeHealthID.String()):                        "Slade Health ID",
		IDDocumentType(IDDocumentTypePatientNumber.String()):                   "Patient Number",
		IDDocumentType(IDDocumentTypeMilitaryID.String()):                      "Military ID",
		IDDocumentType(IDDocumentTypeNationalHospitalInsuranceFundID.String()): "National Hospital Insurance Fund ID",
		IDDocumentType(IDDocumentTypeSmartMemberNumber.String()):               "Smart Member Number",
		IDDocumentType(IDDocumentTypeDrChronoChartID.String()):                 "Dr Chrono Chart ID",
		IDDocumentType(IDDocumentTypeFHIRPatientID.String()):                   "FHIR Patient ID",
		IDDocumentType(IDDocumentTypeERPCustomerID.String()):                   "ERP Customer ID",
		IDDocumentType(IDDocumentTypeComprehensiveCareClinicNumber.String()):   "Comprehensive Care Clinic Number",
		IDDocumentType(IDDocumentTypeRefugeeID.String()):                       "Refugee ID",
		IDDocumentType(IDDocumentTypeBirthCertificateNumber.String()):          "Birth Certificate Number",
		IDDocumentType(IDDocumentTypeMandateNumber.String()):                   "Mandate Number",
		IDDocumentType(IDDocumentTypeClientRegistryNumber.String()):            "Client Registry Number",
	}

	if display, exists := mapping[normalized]; exists {
		return display
	}

	return "Unknown ID Type"
}

// UnmarshalGQL translates the input value to an ID type
func (e *IDDocumentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IDDocumentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IDDocumentType", str)
	}

	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e IDDocumentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// MaritalStatus is used to code individuals' marital statuses.
//
// See: https://www.hl7.org/fhir/valueset-marital-status.html
type MaritalStatus string

// known marital statuses
const (
	// MaritalStatusA ...
	MaritalStatusA MaritalStatus = "A"
	// MaritalStatusD ...
	MaritalStatusD MaritalStatus = "D"
	// MaritalStatusI ...
	MaritalStatusI MaritalStatus = "I"
	// MaritalStatusL ...
	MaritalStatusL MaritalStatus = "L"
	// MaritalStatusM ...
	MaritalStatusM MaritalStatus = "M"
	// MaritalStatusP ...
	MaritalStatusP MaritalStatus = "P"
	// MaritalStatusS ...
	MaritalStatusS MaritalStatus = "S"
	// MaritalStatusT ...
	MaritalStatusT MaritalStatus = "T"
	// MaritalStatusU ...
	MaritalStatusU MaritalStatus = "U"
	// MaritalStatusW ...
	MaritalStatusW MaritalStatus = "W"
	// MaritalStatusUnk ...
	MaritalStatusUnk MaritalStatus = "UNK"
)

// AllMaritalStatus is a list of known marital statuses
var AllMaritalStatus = []MaritalStatus{
	MaritalStatusA,
	MaritalStatusD,
	MaritalStatusI,
	MaritalStatusL,
	MaritalStatusM,
	MaritalStatusP,
	MaritalStatusS,
	MaritalStatusT,
	MaritalStatusU,
	MaritalStatusW,
	MaritalStatusUnk,
}

// IsValid checks that the marital status is valid
func (e MaritalStatus) IsValid() bool {
	switch e {
	case MaritalStatusA, MaritalStatusD, MaritalStatusI, MaritalStatusL, MaritalStatusM, MaritalStatusP, MaritalStatusS, MaritalStatusT, MaritalStatusU, MaritalStatusW, MaritalStatusUnk:
		return true
	}

	return false
}

// String ...
func (e MaritalStatus) String() string {
	return string(e)
}

// UnmarshalGQL turns the supplied input into a marital status enum value
func (e *MaritalStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MaritalStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MaritalStatus", str)
	}

	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e MaritalStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// RelationshipType defines relationship types for patients.
//
// See: https://www.hl7.org/fhir/valueset-relatedperson-relationshiptype.html
type RelationshipType string

// list of known relationship types
const (
	// RelationshipTypeC ...
	RelationshipTypeC RelationshipType = "C"
	// RelationshipTypeE ...
	RelationshipTypeE RelationshipType = "E"
	// RelationshipTypeF ...
	RelationshipTypeF RelationshipType = "F"
	// RelationshipTypeI ...
	RelationshipTypeI RelationshipType = "I"
	// RelationshipTypeN ...
	RelationshipTypeN RelationshipType = "N"
	// RelationshipTypeO ...
	RelationshipTypeO RelationshipType = "O"
	// RelationshipTypeS ...
	RelationshipTypeS RelationshipType = "S"
	// RelationshipTypeU ...
	RelationshipTypeU RelationshipType = "U"
)

// AllRelationshipType is a list of all known relationship types
var AllRelationshipType = []RelationshipType{
	RelationshipTypeC,
	RelationshipTypeE,
	RelationshipTypeF,
	RelationshipTypeI,
	RelationshipTypeN,
	RelationshipTypeO,
	RelationshipTypeS,
	RelationshipTypeU,
}

// IsValid ensures that the relationship type is valid
func (e RelationshipType) IsValid() bool {
	switch e {
	case RelationshipTypeC, RelationshipTypeE, RelationshipTypeF, RelationshipTypeI, RelationshipTypeN, RelationshipTypeO, RelationshipTypeS, RelationshipTypeU:
		return true
	}

	return false
}

// String ...
func (e RelationshipType) String() string {
	return string(e)
}

// UnmarshalGQL converts its input (if valid) into a relationship type
func (e *RelationshipType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RelationshipType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RelationshipType", str)
	}

	return nil
}

// MarshalGQL writes the relationship type to the supplied writer
func (e RelationshipType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// PhoneNumberPayload is a D.T.O that accepts a phone number
type PhoneNumberPayload struct {
	PhoneNumber string `json:"phoneNumber"`
}

// EmailOptIn is used to persist and manage email communication whitelists
type EmailOptIn struct {
	Email   string `json:"email" firestore:"email"`
	OptedIn bool   `json:"optedIn" firestore:"optedIn"`
}

// IsEntity ...
func (e EmailOptIn) IsEntity() {}

// RelationshipTypeDisplay computes friendly string for relationship types
func RelationshipTypeDisplay(val RelationshipType) string {
	switch val {
	case RelationshipTypeC:
		return "Emergency Contact"
	case RelationshipTypeE:
		return "Employer"
	case RelationshipTypeF:
		return "Federal Agency"
	case RelationshipTypeI:
		return "Insurance Company"
	case RelationshipTypeN:
		return "Next-of-Kin"
	case RelationshipTypeO:
		return "Other"
	case RelationshipTypeS:
		return "State Agency"
	case RelationshipTypeU:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// MaritalStatusDisplay calculates the text display for a marital status
// See: https://www.hl7.org/fhir/valueset-marital-status.html
func MaritalStatusDisplay(val MaritalStatus) string {
	switch val {
	case MaritalStatusA:
		return "Annulled"
	case MaritalStatusD:
		return "Divorced"
	case MaritalStatusI:
		return "Interlocutory"
	case MaritalStatusL:
		return "Legally Separated"
	case MaritalStatusM:
		return "Married"
	case MaritalStatusP:
		return "Polygamous"
	case MaritalStatusS:
		return "Never Married"
	case MaritalStatusT:
		return "Domestic Partner"
	case MaritalStatusU:
		return "unmarried"
	case MaritalStatusW:
		return "Widowed"
	case MaritalStatusUnk:
		return "unknown"
	default:
		return "unknown"
	}
}

// PagedFHIRResource is a universal model for fetching FHIR resources with PageInfo details
type PagedFHIRResource struct {
	Resources       []map[string]interface{}
	HasNextPage     bool
	NextCursor      string
	HasPreviousPage bool
	PreviousCursor  string
	TotalCount      int
	BundleID        string
}
