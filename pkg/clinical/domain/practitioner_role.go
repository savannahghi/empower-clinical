package domain

// FHIRPractitionerRole is documented here http://hl7.org/fhir/StructureDefinition/PractitionerRole
type FHIRPractitionerRole struct {
	ID                *string                     `json:"id,omitempty"`
	Meta              *FHIRMeta                   `json:"meta,omitempty"`
	ImplicitRules     *string                     `json:"implicitRules,omitempty"`
	Language          *string                     `json:"language,omitempty"`
	Text              *FHIRNarrative              `json:"text,omitempty"`
	Extension         []FHIRExtension             `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension             `json:"modifierExtension,omitempty"`
	Identifier        []FHIRIdentifier            `json:"identifier,omitempty"`
	Active            *bool                       `json:"active,omitempty"`
	Period            *FHIRPeriod                 `json:"period,omitempty"`
	Practitioner      *FHIRReference              `json:"practitioner,omitempty"`
	Organization      *FHIRReference              `json:"organization,omitempty"`
	Code              []FHIRCodeableConcept       `json:"code,omitempty"`
	Specialty         []FHIRCodeableConcept       `json:"specialty,omitempty"`
	Location          []FHIRReference             `json:"location,omitempty"`
	HealthcareService []FHIRReference             `json:"healthcareService,omitempty"`
	Contact           []FHIRExtendedContactDetail `json:"contact,omitempty"`
	Characteristic    []FHIRCodeableConcept       `json:"characteristic,omitempty"`
	Communication     []FHIRCodeableConcept       `json:"communication,omitempty"`
	Availability      []FHIRAvailability          `json:"availability,omitempty"`
	Endpoint          []FHIRReference             `json:"endpoint,omitempty"`
}

// PagedFHIRPractitionerRole is an FHIR practitioner role's paginated model data class
type PagedFHIRPractitionerRole struct {
	PractitionerRoles []*FHIRPractitionerRole
	HasNextPage       bool
	NextCursor        string
	HasPreviousPage   bool
	PreviousCursor    string
	TotalCount        int
}

type FHIRPractitionerRoleRelayPayload struct {
	Resource *FHIRPractitionerRole `json:"resource,omitempty"`
}
