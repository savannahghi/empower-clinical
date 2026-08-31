package domain

import (
	"fmt"

	"github.com/savannahghi/firebasetools"
)

// FHIROrganization definition: The organization (facility) responsible for this organization
type FHIROrganization struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`
	// Whether the organization's record is still in active use
	Active *bool `json:"active,omitempty"`
	// Identifier(s) by which this organization is known.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`
	// Specific type of organization (e.g. Healthcare Provider, Hospital Department, Insurance Company).
	Type []*FHIRCodeableConcept `json:"type,omitempty"`
	// Name used for the organization
	Name *string `json:"name,omitempty"`
	// A list of alternate names that the organization is known as, or was known as in the past
	Alias []string `json:"alias,omitempty"`
	// A contact detail for the organization
	// ! Rule: The telecom of an organization can never be of use 'home'
	// An address for the organization.
	Address []*FHIRAddress `json:"address,omitempty"`

	Meta    *FHIRMeta                 `json:"meta,omitempty"`
	Text    *FHIRNarrative            `json:"text,omitempty"`
	Contact []FHIROrganizationContact `json:"contact,omitempty"`
	PartOf  *FHIRReference            `json:"partOf,omitempty"`
}

// FHIROrganizationInput definition: The organization (facility) responsible for this organization
type FHIROrganizationInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID           *string `json:"id,omitempty"`
	ResourceType string  `json:"resourceType,omitempty"`
	// Whether the organization's record is still in active use
	Active *bool `json:"active,omitempty"`
	// Identifier(s) by which this organization is known.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`
	// Specific type of organization (e.g. Healthcare Provider, Hospital Department, Insurance Company).
	Type []*FHIRCodeableConceptInput `json:"type,omitempty"`
	// Name used for the organization
	Name *string `json:"name,omitempty"`
	// A list of alternate names that the organization is known as, or was known as in the past
	Alias []string `json:"alias,omitempty"`

	// A contact detail for the organization
	// ! Rule: The telecom of an organization can never be of use 'home'
	Contact []FHIROrganizationContactInput `json:"contact,omitempty"`

	// An address for the organization.
	Address []*FHIRAddressInput `json:"address,omitempty"`

	Text *FHIRNarrative `json:"text,omitempty"`

	PartOf *FHIRReference `json:"partOf,omitempty"`
}

// FHIROrganizationRelayPayload is used to return single instances of Organization
type FHIROrganizationRelayPayload struct {
	Resource *FHIROrganization `json:"resource,omitempty"`
}

// FHIROrganizationRelayConnection is a Relay connection for Organization
type FHIROrganizationRelayConnection struct {
	Edges    []*FHIROrganizationRelayEdge `json:"edges,omitempty"`
	PageInfo *firebasetools.PageInfo      `json:"pageInfo,omitempty"`
}

// FHIROrganizationRelayEdge is a Relay edge for Organization
type FHIROrganizationRelayEdge struct {
	Cursor *string           `json:"cursor,omitempty"`
	Node   *FHIROrganization `json:"node,omitempty"`
}

type FHIROrganizationContact struct {
	ID      *string            `json:"id,omitempty"`
	Telecom []FHIRContactPoint `json:"telecom,omitempty"`
}

type FHIROrganizationContactInput struct {
	Telecom []FHIRContactPointInput `json:"telecom,omitempty"`
}

type ReceivingFacilityContactDetails struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
}

func (o *FHIROrganization) GetOrganizationContacts() []string {
	var organisationContacts []string

	for _, contacts := range o.Contact {
		if len(contacts.Telecom) > 0 {
			for _, contact := range contacts.Telecom {
				if contact.Value != nil {
					organisationContacts = append(organisationContacts, *contact.Value)
				}
			}
		}
	}

	return organisationContacts
}
func SafeStringPointerDeference(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func (o *FHIROrganization) GetReceivingFacilityContactDetails() (*ReceivingFacilityContactDetails, error) {
	var contactDetails ReceivingFacilityContactDetails

	if o == nil {
		return nil, fmt.Errorf("missing organization")
	}

	contactDetails.Name = SafeStringPointerDeference(o.Name)

	for _, contact := range o.Contact {
		for _, telecom := range contact.Telecom {
			switch telecom.System.String() {
			case "email":
				contactDetails.Email = SafeStringPointerDeference(telecom.Value)
			case "phone":
				contactDetails.Phone = SafeStringPointerDeference(telecom.Value)
			}

			if contactDetails.Email != "" || contactDetails.Phone != "" {
				return &contactDetails, nil
			}
		}
	}

	if contactDetails.Email == "" && contactDetails.Phone == "" {
		return nil, fmt.Errorf("receiving facility, %s, missing email and phone contact details", contactDetails.Name)
	}

	// TODO: Should NOT fail if no email is provided. However, a permanent change has to be made which requires all the 'client(s)'
	// to provide an email when creating and organization. So the below check shall then be uncommented
	// if contactDetails.Email == "" {
	// 	return nil, fmt.Errorf("receiving facility, %s, missing email contact", contactDetails.Name)
	// }

	return nil, fmt.Errorf("receiving facility, %s, missing phone contact", contactDetails.Name)
}
