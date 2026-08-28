package domain

import (
	"time"

	"github.com/savannahghi/scalarutils"
)

// PagedFHIRQuestionnaire is used to return paginated list of questionnaires
type PagedFHIRQuestionnaires struct {
	Questionnaires  []FHIRQuestionnaire `mapstructure:"questionnaires"`
	HasNextPage     bool                `mapstructure:"hasNextPage"`
	NextCursor      string              `mapstructure:"nextCursor"`
	HasPreviousPage bool                `mapstructure:"hasPreviousPage"`
	PreviousCursor  string              `mapstructure:"previousCursor"`
	TotalCount      int                 `mapstructure:"totalCount"`
}

// FHIRQuestionnaire models the FHIR questionnaire model as described in https://www.hl7.org/fhir/questionnaire.html
type FHIRQuestionnaire struct {
	ID                *string                  `json:"id,omitempty" mapstructure:"id"`
	ResourceType      string                   `json:"resourceType,omitempty" mapstructure:"resourceType"`
	Meta              *FHIRMetaInput           `json:"meta,omitempty" mapstructure:"meta"`
	ImplicitRules     *string                  `json:"implicitRules,omitempty" mapstructure:"implicitRules"`
	Language          *string                  `json:"language,omitempty" mapstructure:"language"`
	Text              *FHIRNarrative           `json:"text,omitempty" mapstructure:"text"`
	Extension         []*Extension             `json:"extension,omitempty" mapstructure:"extension"`
	ModifierExtension []*Extension             `json:"modifierExtension,omitempty" mapstructure:"modifierExtension"`
	URL               *scalarutils.URI         `json:"url,omitempty" mapstructure:"url"`
	Identifier        []*FHIRIdentifier        `json:"identifier,omitempty" mapstructure:"identifier"`
	Version           *string                  `json:"version,omitempty" mapstructure:"version"`
	Name              *string                  `json:"name,omitempty" mapstructure:"name"`
	Title             *string                  `json:"title,omitempty" mapstructure:"title"`
	DerivedFrom       []*string                `json:"derivedFrom,omitempty" mapstructure:"derivedFrom"`
	Status            *scalarutils.Code        `json:"status,omitempty" mapstructure:"status"`
	Experimental      *bool                    `json:"experimental,omitempty" mapstructure:"experimental"`
	Date              *scalarutils.DateTime    `json:"date,omitempty" mapstructure:"date"`
	Publisher         *string                  `json:"publisher,omitempty" mapstructure:"publisher"`
	Description       *string                  `json:"description,omitempty" mapstructure:"description"`
	UseContext        *FHIRUsageContext        `json:"useContext,omitempty" mapstructure:"useContext"`
	Jurisdiction      []*FHIRCodeableConcept   `json:"jurisdiction,omitempty" mapstructure:"jurisdiction"`
	Purpose           *string                  `json:"purpose,omitempty" mapstructure:"purpose"`
	EffectivePeriod   *FHIRPeriod              `json:"effectivePeriod,omitempty" mapstructure:"effectivePeriod"`
	Code              []*FHIRCoding            `json:"code,omitempty" mapstructure:"code"`
	Item              []*FHIRQuestionnaireItem `json:"item,omitempty" mapstructure:"item"`
}

// FHIRQuestionnaireItem represents the questions and sections within a FHIR questionnaire
type FHIRQuestionnaireItem struct {
	ID                *string                              `json:"id,omitempty" mapstructure:"id"`
	Meta              *FHIRMeta                            `json:"meta,omitempty" mapstructure:"meta"`
	Extension         []*Extension                         `json:"extension,omitempty" mapstructure:"extension"`
	ModifierExtension []*Extension                         `json:"modifierExtension,omitempty" mapstructure:"modifierExtension"`
	LinkID            *string                              `json:"linkId,omitempty" mapstructure:"linkId"`
	Definition        *scalarutils.URI                     `json:"definition,omitempty" mapstructure:"definition"`
	Code              []*FHIRCoding                        `json:"code,omitempty" mapstructure:"code"`
	Prefix            *string                              `json:"prefix,omitempty" mapstructure:"prefix"`
	Text              *string                              `json:"text,omitempty" mapstructure:"text"`
	Type              *scalarutils.Code                    `json:"type,omitempty" mapstructure:"type"`
	EnableWhen        []*FHIRQuestionnaireItemEnableWhen   `json:"enableWhen,omitempty" mapstructure:"enableWhen"`
	EnableBehavior    *scalarutils.Code                    `json:"enableBehavior,omitempty" mapstructure:"enableBehavior"`
	DisabledDisplay   *scalarutils.Code                    `json:"disabledDisplay,omitempty" mapstructure:"disabledDisplay"`
	Required          *bool                                `json:"required,omitempty" mapstructure:"required"`
	Repeats           *bool                                `json:"repeats,omitempty" mapstructure:"repeats"`
	ReadOnly          *bool                                `json:"readOnly,omitempty" mapstructure:"readOnly"`
	MaxLength         *int                                 `json:"maxLength,omitempty" mapstructure:"maxLength"`
	AnswerValueSet    *string                              `json:"answerValueSet,omitempty" mapstructure:"answerValueSet"`
	AnswerOption      []*FHIRQuestionnaireItemAnswerOption `json:"answerOption,omitempty" mapstructure:"answerOption"`
	Initial           []*FHIRQuestionnaireItemInitial      `json:"initial,omitempty" mapstructure:"initial"`
	Item              []*FHIRQuestionnaireItem             `json:"item,omitempty" mapstructure:"item"`
}

// FHIRQuestionnaireItemEnableWhen defines when to enable the FHIR Questionnaire item.
type FHIRQuestionnaireItemEnableWhen struct {
	ID                *string               `json:"id,omitempty" mapstructure:"id"`
	Extension         []*Extension          `json:"extension,omitempty" mapstructure:"extension"`
	ModifierExtension []*Extension          `json:"modifierExtension,omitempty" mapstructure:"modifierExtension"`
	Question          *string               `json:"question,omitempty" mapstructure:"question"`
	Operator          *scalarutils.Code     `json:"operator,omitempty" mapstructure:"operator"`
	AnswerBoolean     *bool                 `json:"answerBoolean,omitempty" mapstructure:"answerBoolean"`
	AnswerDecimal     *float64              `json:"answerDecimal,omitempty" mapstructure:"answerDecimal"`
	AnswerInteger     *int                  `json:"answerInteger,omitempty" mapstructure:"answerInteger"`
	AnswerDate        *scalarutils.Date     `json:"answerDate,omitempty" mapstructure:"answerDate"`
	AnswerDateTime    *scalarutils.DateTime `json:"answerDateTime,omitempty" mapstructure:"answerDateTime"`
	AnswerTime        *scalarutils.DateTime `json:"answerTime,omitempty" mapstructure:"answerTime"`
	AnswerString      *string               `json:"answerString,omitempty" mapstructure:"answerString"`
	AnswerCoding      *FHIRCoding           `json:"answerCoding,omitempty" mapstructure:"answerCoding"`
	AnswerQuantity    *FHIRQuantity         `json:"answerQuantity,omitempty" mapstructure:"answerQuantity"`
	AnswerReference   *FHIRReference        `json:"answerReference,omitempty" mapstructure:"answerReference"`
}

// FHIRQuestionnaireItemAnswerOption represents the permitted answers to a questionnaire.
// ! Rule: A question cannot have both answerOption and answerValueSet
// ! Rule: Only coding, decimal, integer, date, dateTime, time, string or quantity items can have answerOption or answerValueSet
// ! Rule: If one or more answerOption is present, initial cannot be present. Use answerOption.initialSelected instead
type FHIRQuestionnaireItemAnswerOption struct {
	ID                *string           `json:"id,omitempty" mapstructure:"id"`
	Extension         []*Extension      `json:"extension,omitempty" mapstructure:"extension"`
	ModifierExtension []*Extension      `json:"modifierExtension,omitempty" mapstructure:"modifierExtension"`
	ValueInteger      *int              `json:"valueInteger,omitempty" mapstructure:"valueInteger"`
	ValueDate         *scalarutils.Date `json:"valueDate,omitempty" mapstructure:"valueDate"`
	ValueTime         *time.Time        `json:"valueTime,omitempty" mapstructure:"valueTime"`
	ValueString       *string           `json:"valueString,omitempty" mapstructure:"valueString"`
	ValueCoding       *FHIRCoding       `json:"valueCoding,omitempty" mapstructure:"valueCoding"`
	ValueReference    *FHIRReference    `json:"valueReference,omitempty" mapstructure:"valueReference"`
	InitialSelected   *bool             `json:"initialSelected,omitempty" mapstructure:"initialSelected"`
}

// FHIRQuestionnaireItemInitial defines the initial value(s) when a questionnaire item is first rendered
type FHIRQuestionnaireItemInitial struct {
	ID                *string               `json:"id,omitempty" mapstructure:"id"`
	Extension         []*Extension          `json:"extension,omitempty" mapstructure:"extension"`
	ModifierExtension []*Extension          `json:"modifierExtension,omitempty" mapstructure:"modifierExtension"`
	ValueBoolean      *bool                 `json:"valueBoolean,omitempty" mapstructure:"valueBoolean"`
	ValueDecimal      *float64              `json:"valueDecimal,omitempty" mapstructure:"valueDecimal"`
	ValueInteger      *int                  `json:"valueInteger,omitempty" mapstructure:"valueInteger"`
	ValueDate         *scalarutils.Date     `json:"valueDate,omitempty" mapstructure:"valueDate"`
	ValueDateTime     *scalarutils.DateTime `json:"valueDateTime,omitempty" mapstructure:"valueDateTime"`
	ValueTime         *time.Time            `json:"valueTime,omitempty" mapstructure:"valueTime"`
	ValueString       string                `json:"valueString,omitempty" mapstructure:"valueString"`
	ValueURI          *scalarutils.URI      `json:"valueUri,omitempty" mapstructure:"valueURI"`
	ValueAttachment   *FHIRAttachment       `json:"valueAttachment,omitempty" mapstructure:"valueAttachment"`
	ValueCoding       *FHIRCoding           `json:"valueCoding,omitempty" mapstructure:"valueCoding"`
	ValueQuantity     *FHIRQuantity         `json:"valueQuantity,omitempty" mapstructure:"valueQuantity"`
	ValueReference    *FHIRReference        `json:"valueReference,omitempty" mapstructure:"valueReference"`
}

// FHIRUsageContext describes the context that the questionnaire content is intended to support
type FHIRUsageContext struct {
	ID                   *string              `json:"id,omitempty" mapstructure:"id"`
	Extension            []*Extension         `json:"extension,omitempty" mapstructure:"extension"`
	Code                 *FHIRCoding          `json:"code,omitempty" mapstructure:"code"`
	ValueCodeableConcept *FHIRCodeableConcept `json:"valueCodeableConcept,omitempty" mapstructure:"valueCodeableConcept"`
	ValueQuantity        *FHIRQuantity        `json:"valueQuantity,omitempty" mapstructure:"valueQuantity"`
	ValueRange           *FHIRRange           `json:"valueRange,omitempty" mapstructure:"valueRange"`
	ValueReference       *FHIRReference       `json:"valueReference,omitempty" mapstructure:"valueReference"`
}

// FHIRQuestionnaireRelayPayload is used to return a single instance of Questionnaire
type FHIRQuestionnaireRelayPayload struct {
	Resource *FHIRQuestionnaire `json:"resource,omitempty"`
}
