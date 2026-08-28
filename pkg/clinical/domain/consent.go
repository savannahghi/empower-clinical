package domain

import (
	"github.com/savannahghi/scalarutils"
)

// Consent models a fhir consent resource.
type FHIRConsent struct {
	ID               *string                   `json:"id,omitempty"`
	Identifier       []*FHIRIdentifier         `json:"identifier,omitempty"`
	Status           *ConsentStatusEnum        `json:"status,omitempty"`
	Category         []*FHIRCodeableConcept    `json:"category,omitempty"`
	Provision        []*FHIRConsentProvision   `json:"provision,omitempty"`
	Subject          *FHIRReference            `json:"subject,omitempty"`
	Grantor          []*FHIRReference          `json:"grantor,omitempty"`    // This is the patient that grants permision
	Grantee          []*FHIRReference          `json:"grantee,omitempty"`    // This is the entity that is granted access
	Manager          []*FHIRReference          `json:"manager,omitempty"`    // This is the entity that manages consent
	Controller       []*FHIRReference          `json:"controller,omitempty"` // Consent Enforcer
	SourceAttachment []FHIRAttachment          `json:"sourceAttachment,omitempty"`
	SourceReference  []FHIRReference           `json:"sourceReference,omitempty"`
	VerificationType []*FHIRCodeableConcept    `json:"verificaitonType,omitempty"`
	Decision         *ConsentDecisionEnum      `json:"decision,omitempty"`
	Meta             *FHIRMetaInput            `json:"meta,omitempty"`
	Extension        []Extension               `json:"extension,omitempty"`
	Date             *scalarutils.Date         `json:"date,omitempty"`
	Period           *FHIRPeriod               `json:"period,omitempty"`
	Verifcation      []FHIRConsentVerificaiton `json:"verification,omitempty"`
	ImplicitRules    *string                   `json:"implicitRules,omitempty"`
	Text             *FHIRNarrative            `json:"text,omitempty"`
	RegulatoryBasis  []FHIRCodeableConcept     `json:"regulatoryBasis,omitempty"`
	PolicyText       []FHIRReference           `json:"policyText,omitempty"`
	PolicyBasis      *FRHIConsentPolicyBasis   `json:"policyBasis,omitempty"`
	Language         *string                   `json:"language,omitempty"`
}

// FHIRConsentProvision models a fhir consent provision
type FHIRConsentProvision struct {
	ID   *string                    `json:"id,omitempty"`
	Data []FHIRConsentProvisionData `json:"data,omitempty"`
}

// FHIRConsentProvisionData models a consent provision data
type FHIRConsentProvisionData struct {
	ID                *string                `json:"id,omitempty"`
	Extension         []Extension            `json:"extension,omitempty"`
	ModifierExtension []Extension            `json:"modifierExtension,omitempty"`
	Meaning           ConsentDataMeaningEnum `json:"meaning,omitempty"`
	Reference         *FHIRReference         `json:"reference,omitempty"`
}

// FHIRConsentVerificaiton models a consent verification
type FHIRConsentVerificaiton struct {
	ID                *string              `json:"id,omitempty"`
	Extension         []Extension          `json:"extension,omitempty"`
	ModifierExtension []Extension          `json:"modifierExtension,omitempty"`
	Verfied           *bool                `json:"verified"`
	VerificaitonType  *FHIRCodeableConcept `json:"verificationType,omitempty"`
	VerifiedBy        *FHIRReference       `json:"verifiedBy"`
	VerifiedWith      *FHIRReference       `json:"verifiedWith"`
	VerificationDate  []*scalarutils.Date  `json:"verificationDate"`
}

// ConsentPolicy models the basis of a policy
type FRHIConsentPolicyBasis struct {
	ID                *string     `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []Extension `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []Extension `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Reference         *Reference  `bson:"reference,omitempty" json:"reference,omitempty"`
	URL               *string     `bson:"url,omitempty" json:"url,omitempty"`
}
