package domain

import (
	"errors"
	"strings"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FHIRServiceRequest definition: a record of a request for service such as diagnostic investigations, treatments, or operations to be performed.
type FHIRServiceRequest struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// Identifiers assigned to this order instance by the orderer and/or the receiver and/or order fulfiller.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// The URL pointing to a FHIR-defined protocol, guideline, orderset or other definition that is adhered to in whole or in part by this ServiceRequest.
	InstantiatesCanonical *scalarutils.Canonical `json:"instantiatesCanonical,omitempty"`

	// The URL pointing to an externally maintained protocol, guideline, orderset or other definition that is adhered to in whole or in part by this ServiceRequest.
	InstantiatesURI *scalarutils.Instant `json:"instantiatesURI,omitempty"`

	// Plan/proposal/order fulfilled by this request.
	BasedOn []*FHIRReference `json:"basedOn,omitempty"`

	// The request takes the place of the referenced completed or terminated request(s).
	Replaces []*FHIRReference `json:"replaces,omitempty"`

	// A shared identifier common to all service requests that were authorized more or less simultaneously by a single author, representing the composite or group identifier.
	Requisition *FHIRIdentifier `json:"requisition,omitempty"`

	// The status of the order.
	Status ServiceRequestStatusEnum `json:"status,omitempty"`

	// Whether the request is a proposal, plan, an original order or a reflex order.
	Intent ServiceRequestIntentEnum `json:"intent,omitempty"`

	// A code that classifies the service for searching, sorting and display purposes (e.g. "Surgical Procedure").
	Category []*FHIRCodeableConcept `json:"category,omitempty"`

	// Indicates how quickly the ServiceRequest should be addressed with respect to other requests.
	Priority ServiceRequestPriorityEnum `json:"priority,omitempty"`

	// Set this to true if the record is saying that the service/procedure should NOT be performed.
	DoNotPerform *bool `json:"doNotPerform,omitempty"`

	// A code that identifies a particular service (i.e., procedure, diagnostic investigation, or panel of investigations) that have been requested.
	Code *FHIRCodeableReference `json:"code,omitempty"`

	// Additional details and instructions about the how the services are to be delivered.   For example, and order for a urinary catheter may have an order detail for an external or indwelling catheter, or an order for a bandage may require additional instructions specifying how the bandage should be applied.
	OrderDetail []*FHIRCodeableConcept `json:"orderDetail,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityQuantity *FHIRQuantity `json:"quantityQuantity,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityRatio *FHIRRatio `json:"quantityRatio,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityRange *FHIRRange `json:"quantityRange,omitempty"`

	// On whom or what the service is to be performed. This is usually a human patient, but can also be requested on animals, groups of humans or animals, devices such as dialysis machines, or even locations (typically for environmental scans).
	Subject *FHIRReference `json:"subject,omitempty"`

	// An encounter that provides additional information about the healthcare context in which this request is made.
	Encounter *FHIRReference `json:"encounter,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrenceDateTime *scalarutils.Date `json:"occurrenceDateTime,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrencePeriod *FHIRPeriod `json:"occurrencePeriod,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrenceTiming *FHIRTiming `json:"occurrenceTiming,omitempty"`

	// If a CodeableConcept is present, it indicates the pre-condition for performing the service.  For example "pain", "on flare-up", etc.
	AsNeededBoolean *bool `json:"asNeededBoolean,omitempty"`

	// If a CodeableConcept is present, it indicates the pre-condition for performing the service.  For example "pain", "on flare-up", etc.
	AsNeededCodeableConcept *scalarutils.Code `json:"asNeededCodeableConcept,omitempty"`

	// When the request transitioned to being actionable.
	AuthoredOn *scalarutils.DateTime `json:"authoredOn,omitempty"`

	// The individual who initiated the request and has responsibility for its activation.
	Requester *FHIRReference `json:"requester,omitempty"`

	// Desired type of performer for doing the requested service.
	PerformerType *FHIRCodeableConcept `json:"performerType,omitempty"`

	// The desired performer for doing the requested service.  For example, the surgeon, dermatopathologist, endoscopist, etc.
	Performer []*FHIRReference `json:"performer,omitempty"`

	// The preferred location(s) where the procedure should actually happen in coded or free text form. E.g. at home or nursing day care center.
	LocationCode *scalarutils.Code `json:"locationCode,omitempty"`

	// A reference to the the preferred location(s) where the procedure should actually happen. E.g. at home or nursing day care center.
	LocationReference []*FHIRReference `json:"locationReference,omitempty"`

	// An explanation or justification for why this service is being requested in coded or textual form.   This is often for billing purposes.  May relate to the resources referred to in `supportingInfo`.
	ReasonCode *scalarutils.Code `json:"reasonCode,omitempty"`

	// Indicates another resource that provides a justification for why this service is being requested.   May relate to the resources referred to in `supportingInfo`.
	ReasonReference []*FHIRReference `json:"reasonReference,omitempty"`

	// Insurance plans, coverage extensions, pre-authorizations and/or pre-determinations that may be needed for delivering the requested service.
	Insurance []*FHIRReference `json:"insurance,omitempty"`

	// Additional clinical information about the patient or specimen that may influence the services or their interpretations.     This information includes diagnosis, clinical findings and other observations.  In laboratory ordering these are typically referred to as "ask at order entry questions (AOEs)".  This includes observations explicitly requested by the producer (filler) to provide context or supporting information needed to complete the order. For example,  reporting the amount of inspired oxygen for blood gas measurements.
	SupportingInfo []*FHIRReference `json:"supportingInfo,omitempty"`

	// One or more specimens that the laboratory procedure will use.
	Specimen []*FHIRReference `json:"specimen,omitempty"`

	// Anatomic location where the procedure should be performed. This is the target site.
	BodySite []*FHIRCodeableConcept `json:"bodySite,omitempty"`

	// Any other notes and comments made about the service request. For example, internal billing notes.
	Note []*FHIRAnnotation `json:"note,omitempty"`

	// Instructions in terms that are understood by the patient or consumer.
	PatientInstruction *string `json:"patientInstruction,omitempty"`

	// Key events in the history of the request.
	RelevantHistory []*FHIRReference `json:"relevantHistory,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMeta `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`

	ImplicitRules     *string                  `json:"implicitRules,omitempty"`
	Language          *string                  `json:"language,omitempty"`
	ModifierExtension []Extension              `json:"modifierExtension,omitempty"`
	InstantiatesUri   []string                 `json:"instantiatesUri,omitempty"`
	Focus             []Reference              `json:"focus,omitempty"`
	Location          []*FHIRCodeableReference `json:"location,omitempty"`
	Reason            []*FHIRCodeableReference `json:"reason,omitempty"`
	BodyStructure     *Reference               `json:"bodyStructure,omitempty"`
}

func (s *FHIRServiceRequest) GetCoding() (string, string) {
	if s.Code != nil && s.Code.Concept != nil && len(s.Code.Concept.Coding) > 0 {
		for _, code := range s.Code.Concept.Coding {
			return string(*code.Code), code.Display
		}
	}

	return "", ""
}

// GetCategory returns the first category code recorded on the service request
// (e.g. "laboratory-procedure" for lab orders, "referral" for referrals). It returns an
// empty string when the resource carries no category coding.
func (s *FHIRServiceRequest) GetCategory() string {
	for _, category := range s.Category {
		if category == nil {
			continue
		}

		for _, coding := range category.Coding {
			if coding != nil && coding.Code != nil {
				return string(*coding.Code)
			}
		}
	}

	return ""
}

// FHIRServiceRequestInput is the input type for ServiceRequest
type FHIRServiceRequestInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrativeInput `json:"text,omitempty"`

	// Identifiers assigned to this order instance by the orderer and/or the receiver and/or order fulfiller.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	// The URL pointing to a FHIR-defined protocol, guideline, orderset or other definition that is adhered to in whole or in part by this ServiceRequest.
	InstantiatesCanonical *scalarutils.Canonical `json:"instantiatesCanonical,omitempty"`

	// The URL pointing to an externally maintained protocol, guideline, orderset or other definition that is adhered to in whole or in part by this ServiceRequest.
	InstantiatesURI *scalarutils.Instant `json:"instantiatesURI,omitempty"`

	// Plan/proposal/order fulfilled by this request.
	BasedOn []*FHIRReferenceInput `json:"basedOn,omitempty"`

	// The request takes the place of the referenced completed or terminated request(s).
	Replaces []*FHIRReferenceInput `json:"replaces,omitempty"`

	// A shared identifier common to all service requests that were authorized more or less simultaneously by a single author, representing the composite or group identifier.
	Requisition *FHIRIdentifierInput `json:"requisition,omitempty"`

	// The status of the order.
	Status ServiceRequestStatusEnum `json:"status,omitempty"`

	// Whether the request is a proposal, plan, an original order or a reflex order.
	Intent ServiceRequestIntentEnum `json:"intent,omitempty"`

	// A code that classifies the service for searching, sorting and display purposes (e.g. "Surgical Procedure").
	Category []*FHIRCodeableConceptInput `json:"category,omitempty"`

	// Indicates how quickly the ServiceRequest should be addressed with respect to other requests.
	Priority ServiceRequestPriorityEnum `json:"priority,omitempty"`

	// Set this to true if the record is saying that the service/procedure should NOT be performed.
	DoNotPerform *bool `json:"doNotPerform,omitempty"`

	// A code that identifies a particular service (i.e., procedure, diagnostic investigation, or panel of investigations) that have been requested.
	Code *FHIRCodeableReferenceInput `json:"code,omitempty"`

	// Additional details and instructions about the how the services are to be delivered.   For example, and order for a urinary catheter may have an order detail for an external or indwelling catheter, or an order for a bandage may require additional instructions specifying how the bandage should be applied.
	OrderDetail []*FHIRCodeableConceptInput `json:"orderDetail,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityQuantity *FHIRQuantityInput `json:"quantityQuantity,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityRatio *FHIRRatioInput `json:"quantityRatio,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityRange *FHIRRangeInput `json:"quantityRange,omitempty"`

	// On whom or what the service is to be performed. This is usually a human patient, but can also be requested on animals, groups of humans or animals, devices such as dialysis machines, or even locations (typically for environmental scans).
	Subject *FHIRReferenceInput `json:"subject,omitempty"`

	// An encounter that provides additional information about the healthcare context in which this request is made.
	Encounter *FHIRReferenceInput `json:"encounter,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrenceDateTime *scalarutils.Date `json:"occurrenceDateTime,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrencePeriod *FHIRPeriodInput `json:"occurrencePeriod,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrenceTiming *FHIRTimingInput `json:"occurrenceTiming,omitempty"`

	// If a CodeableConcept is present, it indicates the pre-condition for performing the service.  For example "pain", "on flare-up", etc.
	AsNeededBoolean *bool `json:"asNeededBoolean,omitempty"`

	// If a CodeableConcept is present, it indicates the pre-condition for performing the service.  For example "pain", "on flare-up", etc.
	AsNeededCodeableConcept *scalarutils.Code `json:"asNeededCodeableConcept,omitempty"`

	// When the request transitioned to being actionable.
	AuthoredOn *scalarutils.DateTime `json:"authoredOn,omitempty"`

	// The individual who initiated the request and has responsibility for its activation.
	Requester *FHIRReferenceInput `json:"requester,omitempty"`

	// Desired type of performer for doing the requested service.
	PerformerType *FHIRCodeableConceptInput `json:"performerType,omitempty"`

	// The desired performer for doing the requested service.  For example, the surgeon, dermatopathologist, endoscopist, etc.
	Performer []*FHIRReferenceInput `json:"performer,omitempty"`

	// The preferred location(s) where the procedure should actually happen in coded or free text form. E.g. at home or nursing day care center.
	LocationCode *scalarutils.Code `json:"locationCode,omitempty"`

	// A reference to the the preferred location(s) where the procedure should actually happen. E.g. at home or nursing day care center.
	LocationReference []*FHIRReferenceInput `json:"locationReference,omitempty"`

	// An explanation or justification for why this service is being requested in coded or textual form.   This is often for billing purposes.  May relate to the resources referred to in `supportingInfo`.
	ReasonCode *scalarutils.Code `json:"reasonCode,omitempty"`

	// Insurance plans, coverage extensions, pre-authorizations and/or pre-determinations that may be needed for delivering the requested service.
	Insurance []*FHIRReferenceInput `json:"insurance,omitempty"`

	// Additional clinical information about the patient or specimen that may influence the services or their interpretations.     This information includes diagnosis, clinical findings and other observations.  In laboratory ordering these are typically referred to as "ask at order entry questions (AOEs)".  This includes observations explicitly requested by the producer (filler) to provide context or supporting information needed to complete the order. For example,  reporting the amount of inspired oxygen for blood gas measurements.
	SupportingInfo []*FHIRReferenceInput `json:"supportingInfo,omitempty"`

	// One or more specimens that the laboratory procedure will use.
	Specimen []*FHIRReferenceInput `json:"specimen,omitempty"`

	// Anatomic location where the procedure should be performed. This is the target site.
	BodySite []*FHIRCodeableConceptInput `json:"bodySite,omitempty"`

	// Any other notes and comments made about the service request. For example, internal billing notes.
	Note []*FHIRAnnotationInput `json:"note,omitempty"`

	// Instructions in terms that are understood by the patient or consumer.
	PatientInstruction *string `json:"patientInstruction,omitempty"`

	// Key events in the history of the request.
	RelevantHistory []*FHIRReferenceInput `json:"relevantHistory,omitempty"`

	// Meta stores more information about the resource
	Meta FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`

	ImplicitRules     *string                       `json:"implicitRules,omitempty"`
	Language          *string                       `json:"language,omitempty"`
	ModifierExtension []Extension                   `json:"modifierExtension,omitempty"`
	InstantiatesUri   []string                      `json:"instantiatesUri,omitempty"`
	Focus             []Reference                   `json:"focus,omitempty"`
	Location          []*FHIRCodeableReferenceInput `json:"location,omitempty"`
	Reason            []*FHIRCodeableReferenceInput `json:"reason,omitempty"`
	BodyStructure     *Reference                    `json:"bodyStructure,omitempty"`
}

// FHIRServiceRequestRelayConnection is a Relay connection for ServiceRequest
type FHIRServiceRequestRelayConnection struct {
	Edges []*FHIRServiceRequestRelayEdge `json:"edges,omitempty"`

	PageInfo *firebasetools.PageInfo `json:"pageInfo,omitempty"`
}

// FHIRServiceRequestRelayEdge is a Relay edge for ServiceRequest
type FHIRServiceRequestRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIRServiceRequest `json:"node,omitempty"`
}

// FHIRServiceRequestRelayPayload is used to return single instances of ServiceRequest
type FHIRServiceRequestRelayPayload struct {
	Resource *FHIRServiceRequest `json:"resource,omitempty"`
}

// PagedFHIRServiceRequest is used to return paginated list of service requests
type PagedFHIRServiceRequest struct {
	ServiceRequests []FHIRServiceRequest `mapstructure:"serviceRequests"`
	HasNextPage     bool                 `mapstructure:"hasNextPage"`
	NextCursor      string               `mapstructure:"nextCursor"`
	HasPreviousPage bool                 `mapstructure:"hasPreviousPage"`
	PreviousCursor  string               `mapstructure:"previousCursor"`
	TotalCount      *int                 `mapstructure:"totalCount"`
}

// ReceivingFacility is the facility that is accepting the patient into its care
type ReceivingFacility struct {
	FacilityName    string `json:"facilityName,omitempty"`
	FacilityCounty  string `json:"facilityCounty,omitempty"`
	FacilityContact string `json:"facilityContact,omitempty"`
	FacilityEmail   string `json:"facilityEmail,omitempty"`
}

// GetReceivingFacilityDetails is used to fetch the details of the receiving facility
func (f *FHIRServiceRequest) GetReceivingFacilityDetails() (*ReceivingFacility, error) {
	if f.Extension == nil {
		return nil, errors.New("extension details is nil")
	}

	receivingFacility := &ReceivingFacility{}

	for _, extension := range f.Extension {
		if extension.URL == "http://savannahghi.org/fhir/StructureDefinition/referred-facility" {
			for _, ext := range extension.Extension {
				if ext.URL == "facilityName" {
					receivingFacility.FacilityName = ext.ValueString
				}

				if ext.URL == "facilityCounty" {
					receivingFacility.FacilityCounty = ext.ValueString
				}

				if ext.URL == "facilityContact" {
					receivingFacility.FacilityContact = ext.ValueString
				}

				if ext.URL == "facilityEmail" {
					receivingFacility.FacilityEmail = ext.ValueString
				}
			}
		}
	}

	return receivingFacility, nil
}

// GetSubject returns the patient details
func (f *FHIRServiceRequest) GetSubject() *FHIRReference {
	return f.Subject
}

func (f *FHIRServiceRequest) GetFacilityFromMeta() *FHIRMeta {
	return f.Meta
}

// GetFacilityName is used to get the name of the facility
func (f *FHIRServiceRequest) GetFacilityName() string {
	facilitySystem := scalarutils.URI("http://mycarehub/tenant-identification/facility")

	if f.Meta.Tag != nil {
		for _, meta := range f.Meta.Tag {
			if string(*meta.System) == string(facilitySystem) {
				return meta.Display
			}
		}
	}

	return ""
}

// GetPatientReferralReason fetches the reason why a patient has been referred from a FHIR service request.
// It specifically looks for the code matching the common.ReferralReasonCIELCode in the service request's coding.
// If a matching code is found, the display name of that code is returned as the name of the referred test.
// It defaults to "test" if a test name is not found
func (f *FHIRServiceRequest) GetPatientReferralReason() string {
	if f.Code.Concept != nil && f.Code.Concept.Text != "" {
		return f.Code.Concept.Text
	}

	return ""
}

// GetPatientReferralTest returns the test that the patient has been referred for
func (f *FHIRServiceRequest) GetPatientReferralTest() string {
	if f.Code != nil && len(f.Code.Concept.Coding) > 0 {
		for _, coding := range f.Code.Concept.Coding {
			if coding.Code != nil && *coding.Code == scalarutils.Code("TEST") && coding.Display != "" {
				return coding.Display
			}
		}
	}

	return ""
}

// GetRequestedServices retrieves What is being requested/ordered, in this case, it could be the tests that the patient
// has been referred for
func (f *FHIRServiceRequest) GetRequestedServices(serviceCIELCode string) []string {
	var services []string

	for _, coding := range f.Code.Concept.Coding {
		if coding.Code != nil && *coding.Code == scalarutils.Code(serviceCIELCode) {
			services = append(services, coding.Display)
		}
	}

	return services
}

// GetPractitionersNotes returns all practitioner's notes from the FHIR service request concatenated with a newline separator.
func (f *FHIRServiceRequest) GetPractitionersNotes() string {
	if len(f.Note) == 0 {
		return ""
	}

	var notes []string

	for _, note := range f.Note {
		if note.Text != nil {
			notes = append(notes, string(*note.Text))
		}
	}

	return strings.Join(notes, "\n")
}

// GetCode is a helper method that fetches a code from ServiceRequest
func (f *FHIRServiceRequest) GetCode() string {
	for _, code := range f.Code.Concept.Coding {
		if code.Code != nil {
			return (string)(*code.Code)
		}
	}

	return ""
}

// ServiceRequestType is a custom categorization of service requests.
//
// Service requests are used to create referrals, lab orders etc. These enums are just
// helpful in differentiating
type ServiceRequestType string

const (
	ReferralServiceRequestType ServiceRequestType = "REFERRAL_SERVICE_REQUEST"
	LabOrderServiceRequestType ServiceRequestType = "LAB_ORDER_SERVICE_REQUEST"
)

// String converts the service request data meaning to string
func (s ServiceRequestType) String() string {
	switch s {
	case ReferralServiceRequestType:
		return "Patient Referral (Service Request)"
	case LabOrderServiceRequestType:
		return "Laboratory Order (Service Request)"
	default:
		return string(s)
	}
}

// GetPerformerID fetches the performer ID
func (s *FHIRServiceRequest) GetPerformerID() string {
	for _, org := range s.Performer {
		if org.Display != "" {
			return org.Display
		}
	}

	return ""
}

type ServiceRequestCategoryType string

const (
	LaboratoryProcedureCategoryType ServiceRequestCategoryType = "laboratory-procedure"
	ImagingCategoryType             ServiceRequestCategoryType = "imaging"
	CounsellingCategoryType         ServiceRequestCategoryType = "counselling"
	EducationCategoryType           ServiceRequestCategoryType = "education"
	SurgicalProcedureCategoryType   ServiceRequestCategoryType = "surgical-procedure"
	ReferralCategoryType            ServiceRequestCategoryType = "referral"
)

func (s ServiceRequestCategoryType) Display() string {
	switch s {
	case LaboratoryProcedureCategoryType:
		return "Laboratory procedure"
	case ImagingCategoryType:
		return "Imaging"
	case EducationCategoryType:
		return "Education"
	case SurgicalProcedureCategoryType:
		return "Surgical procedure"
	case ReferralCategoryType:
		return "Referral"
	case CounsellingCategoryType:
		return "Counselling"
	}

	return "Unknown"
}
