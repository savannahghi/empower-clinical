package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FHIRSubstance is documented here http://hl7.org/fhir/StructureDefinition/Substance
type FHIRSubstance struct {
	ID                *string                    `bson:"id,omitempty" json:"id,omitempty"`
	Meta              *FHIRMeta                  `bson:"meta,omitempty" json:"meta,omitempty"`
	ImplicitRules     *string                    `bson:"implicitRules,omitempty" json:"implicitRules,omitempty"`
	Language          *string                    `bson:"language,omitempty" json:"language,omitempty"`
	Text              *FHIRNarrative             `bson:"text,omitempty" json:"text,omitempty"`
	Extension         []*FHIRExtension           `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension           `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Identifier        []*FHIRIdentifier          `bson:"identifier,omitempty" json:"identifier,omitempty"`
	Instance          bool                       `bson:"instance" json:"instance"`
	Status            *FHIRSubstanceStatus       `bson:"status,omitempty" json:"status,omitempty"`
	Category          []*FHIRCodeableConcept     `bson:"category,omitempty" json:"category,omitempty"`
	Code              FHIRCodeableReference      `bson:"code" json:"code"`
	Description       *string                    `bson:"description,omitempty" json:"description,omitempty"`
	Expiry            *string                    `bson:"expiry,omitempty" json:"expiry,omitempty"`
	Quantity          *FHIRQuantity              `bson:"quantity,omitempty" json:"quantity,omitempty"`
	Ingredient        []*FHIRSubstanceIngredient `bson:"ingredient,omitempty" json:"ingredient,omitempty"`
}

type FHIRSubstanceIngredient struct {
	ID                       *string             `bson:"id,omitempty" json:"id,omitempty"`
	Extension                []*FHIRExtension    `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension        []*FHIRExtension    `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Quantity                 *FHIRRatio          `bson:"quantity,omitempty" json:"quantity,omitempty"`
	SubstanceCodeableConcept FHIRCodeableConcept `bson:"substanceCodeableConcept" json:"substanceCodeableConcept"`
	SubstanceReference       FHIRReference       `bson:"substanceReference" json:"substanceReference"`
}

type PagedFHIRSubstance struct {
	Substance       []*FHIRSubstance
	HasNextPage     bool
	NextCursor      string
	HasPreviousPage bool
	PreviousCursor  string
	TotalCount      int
}

type FHIRSubstanceRelayPayload struct {
	Resource *FHIRSubstance `json:"resource,omitempty"`
}

// FHIRSubstanceStatus is documented here http://hl7.org/fhir/ValueSet/substance-status
type FHIRSubstanceStatus string

const (
	FHIRSubstanceStatusActive         FHIRSubstanceStatus = "active"
	FHIRSubstanceStatusInactive       FHIRSubstanceStatus = "inactive"
	FHIRSubstanceStatusEnteredInError FHIRSubstanceStatus = "entered-in-error"
)

func (code FHIRSubstanceStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(code.Code())
}

func (code *FHIRSubstanceStatus) UnmarshalJSON(json []byte) error {
	s := strings.Trim(string(json), "\"")
	switch s {
	case "active":
		*code = FHIRSubstanceStatusActive
	case "inactive":
		*code = FHIRSubstanceStatusInactive
	case "entered-in-error":
		*code = FHIRSubstanceStatusEnteredInError
	default:
		return fmt.Errorf("unknown FHIRSubstanceStatus code `%s`", s)
	}

	return nil
}

func (code FHIRSubstanceStatus) String() string {
	return code.Code()
}

// FHIRSubstanceStatus represents the status of a FHIR Substance
func (code FHIRSubstanceStatus) Code() string {
	switch code {
	case FHIRSubstanceStatusActive:
		return "active"
	case FHIRSubstanceStatusInactive:
		return "inactive"
	case FHIRSubstanceStatusEnteredInError:
		return "entered-in-error"
	}

	return "<unknown>"
}

// Display returns the display string for the given FHIRSubstanceStatus code.
func (code FHIRSubstanceStatus) Display() string {
	switch code {
	case FHIRSubstanceStatusActive:
		return "Active"
	case FHIRSubstanceStatusInactive:
		return "Inactive"
	case FHIRSubstanceStatusEnteredInError:
		return "Entered in Error"
	}

	return "<unknown>"
}
