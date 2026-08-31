package domain

import (
	"github.com/savannahghi/firebasetools"
)

// FHIRMedication definition:
type FHIRMedication struct {
	ID                           *string                    `json:"id,omitempty"`
	Meta                         *FHIRMeta                  `json:"meta,omitempty"`
	ImplicitRules                *string                    `json:"implicitRules,omitempty"`
	Language                     *string                    `json:"language,omitempty"`
	Text                         *FHIRNarrative             `json:"text,omitempty"`
	Extension                    []Extension                `json:"extension,omitempty"`
	ModifierExtension            []FHIRExtension            `json:"modifierExtension,omitempty"`
	Identifier                   []FHIRIdentifier           `json:"identifier,omitempty"`
	Code                         *FHIRCodeableConcept       `json:"code,omitempty"`
	Status                       *MedicationStatusEnum      `json:"status,omitempty"`
	MarketingAuthorizationHolder *FHIRReference             `json:"marketingAuthorizationHolder,omitempty"`
	DoseForm                     *FHIRCodeableConcept       `json:"doseForm,omitempty"`
	TotalVolume                  *FHIRQuantity              `json:"totalVolume,omitempty"`
	Ingredient                   []FHIRMedicationIngredient `json:"ingredient,omitempty"`
	Batch                        *FHIRMedicationBatch       `json:"batch,omitempty"`
	Definition                   *FHIRReference             `json:"definition,omitempty"`
}

func (m *FHIRMedication) GetText() string {
	if len(m.Extension) > 0 {
		for _, ext := range m.Extension {
			return ext.ValueString
		}
	}

	return ""
}

func (m *FHIRMedication) GetDoseForm() (string, string) {
	if m.DoseForm != nil && len(m.DoseForm.Coding) > 0 {
		for _, df := range m.DoseForm.Coding {
			return string(*df.Code), df.Display
		}
	}

	return "", ""
}

type FHIRMedicationIngredient struct {
	ID                      *string                `json:"id,omitempty"`
	Extension               []FHIRExtension        `json:"extension,omitempty"`
	ModifierExtension       []FHIRExtension        `json:"modifierExtension,omitempty"`
	Item                    *FHIRCodeableReference `json:"item"`
	IsActive                *bool                  `json:"isActive,omitempty"`
	StrengthRatio           *FHIRRatio             `json:"strengthRatio,omitempty"`
	StrengthCodeableConcept *FHIRCodeableConcept   `json:"strengthCodeableConcept,omitempty"`
	StrengthQuantity        *FHIRQuantity          `json:"strengthQuantity,omitempty"`
}

type FHIRMedicationBatch struct {
	ID                *string         `json:"id,omitempty"`
	Extension         []FHIRExtension `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension `json:"modifierExtension,omitempty"`
	LotNumber         *string         `json:"lotNumber,omitempty"`
	ExpirationDate    *string         `json:"expirationDate,omitempty"`
}

// FHIRMedicationInput ...
type FHIRMedicationInput struct {
	ID *string `json:"id,omitempty"`

	Text *FHIRNarrativeInput `json:"text,omitempty"`

	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	Code *FHIRCodeableConceptInput `json:"code,omitempty"`

	Status MedicationStatusEnum `json:"status,omitempty"`

	Manufacturer *FHIROrganizationInput `json:"manufacturer,omitempty"`

	DoseForm *FHIRCodeableConceptInput `json:"doseForm,omitempty"`

	Amount *FHIRRatioInput `json:"amount,omitempty"`

	Ingredient []*FHIRMedicationIngredient `json:"ingredient,omitempty"`

	Batch *FHIRMedicationBatch `json:"batch,omitempty"`

	// Meta stores more information about the resource
	Meta FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*Extension `json:"extension,omitempty"`
}

// FHIRMedicationRelayConnection is a Relay connection for Medication
type FHIRMedicationRelayConnection struct {
	Edges []*FHIRMedicationRelayEdge `json:"edges,omitempty"`

	PageInfo *firebasetools.PageInfo `json:"pageInfo,omitempty"`
}

// FHIRMedicationRelayEdge is a Relay edge for Medication
type FHIRMedicationRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIRMedication `json:"node,omitempty"`
}

// FHIRMedicationRelayPayload is used to return single instances of Medication
type FHIRMedicationRelayPayload struct {
	Resource *FHIRMedication `json:"resource,omitempty"`
}
