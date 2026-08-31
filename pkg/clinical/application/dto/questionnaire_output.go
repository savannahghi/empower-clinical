package dto

import (
	"time"

	"github.com/savannahghi/scalarutils"
)

// QuestionnaireEdge defines the structure of a Questionnaire edge
type QuestionnaireEdge struct {
	Node   Questionnaire
	Cursor string
}

// QuestionnaireConnection models the questionnaire connection data class
type QuestionnaireConnection struct {
	TotalCount int
	Edges      []QuestionnaireEdge
	PageInfo   PageInfo
}

// CreateQuestionnaireConnection creates a connection to map out questionnaire results following the GraphQl Cursor Connection Specification
func CreateQuestionnaireConnection(questionnaires []*Questionnaire, pageInfo PageInfo, total int) QuestionnaireConnection {
	connection := QuestionnaireConnection{
		TotalCount: total,
		Edges:      []QuestionnaireEdge{},
		PageInfo:   pageInfo,
	}

	for _, questionnaire := range questionnaires {
		edge := QuestionnaireEdge{
			Node:   *questionnaire,
			Cursor: questionnaire.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}

// Questionnaire models the dataclass to display questionnaires
type Questionnaire struct {
	ID                string                `json:"id,omitempty"`
	ResourceType      *string               `json:"resourceType,omitempty"`
	Meta              *Meta                 `json:"meta,omitempty"`
	ImplicitRules     *string               `json:"implicitRules,omitempty"`
	Language          *string               `json:"language,omitempty"`
	Text              *Narrative            `json:"text,omitempty"`
	Extension         []*Extension          `json:"extension,omitempty"`
	ModifierExtension []*Extension          `json:"modifierExtension,omitempty"`
	URL               *scalarutils.URI      `json:"url,omitempty"`
	Identifier        *[]Identifier         `json:"identifier,omitempty"`
	Version           *string               `json:"version,omitempty"`
	Name              *string               `json:"name,omitempty"`
	Title             *string               `json:"title,omitempty"`
	DerivedFrom       []*string             `json:"derivedFrom,omitempty"`
	Status            *scalarutils.Code     `json:"status,omitempty"`
	Experimental      *bool                 `json:"experimental,omitempty"`
	Date              *scalarutils.DateTime `json:"date,omitempty"`
	Publisher         *string               `json:"publisher,omitempty"`
	Description       *string               `json:"description,omitempty"`
	UseContext        *UsageContext         `json:"useContext,omitempty"`
	Jurisdiction      []*CodeableConcept    `json:"jurisdiction,omitempty"`
	Purpose           *string               `json:"purpose,omitempty"`
	EffectivePeriod   *Period               `json:"effectivePeriod,omitempty"`
	Code              []*Coding             `json:"code,omitempty"`
	Item              []*QuestionnaireItem  `json:"item,omitempty"`
}

// IsEntity marks a questionnaire as an entity
func (q Questionnaire) IsEntity() {}

// Narrative models the questionnaire narrative
type Narrative struct {
	ID     string             `json:"id,omitempty"`
	Status string             `json:"status,omitempty"`
	Div    *scalarutils.XHTML `json:"div,omitempty"`
}

// Extension models the various extensions of a questionnaire resource
type Extension struct {
	URL                  string           `json:"url,omitempty"`
	ValueBoolean         bool             `json:"valueBoolean,omitempty"`
	ValueInteger         int              `json:"valueInteger"`
	ValueDecimal         float64          `json:"valueDecimal"`
	ValueBase64Binary    string           `json:"valueBase64Binary,omitempty"`
	ValueInstant         string           `json:"valueInstant,omitempty"`
	ValueString          string           `json:"valueString,omitempty"`
	ValueURI             string           `json:"valueURI,omitempty"`
	ValueDate            string           `json:"valueDate,omitempty"`
	ValueDateTime        string           `json:"valueDateTime,omitempty"`
	ValueTime            string           `json:"valueTime,omitempty"`
	ValueCode            string           `json:"valueCode,omitempty"`
	ValueOid             string           `json:"valueOid,omitempty"`
	ValueUUID            string           `json:"valueUUID,omitempty"`
	ValueID              string           `json:"valueID,omitempty"`
	ValueUnsignedInt     int              `json:"valueUnsignedInt,omitempty"`
	ValuePositiveInt     int              `json:"valuePositiveInt,omitempty"`
	ValueMarkdown        string           `json:"valueMarkdown,omitempty"`
	ValueAnnotation      *Annotation      `json:"valueAnnotation,omitempty"`
	ValueAttachment      *Attachment      `json:"valueAttachment,omitempty"`
	ValueIdentifier      *Identifier      `json:"valueIdentifier,omitempty"`
	ValueCodeableConcept *CodeableConcept `json:"valueCodeableConcept,omitempty"`
	ValueCoding          *Coding          `json:"valueCoding,omitempty"`
	ValueQuantity        *Quantity        `json:"valueQuantity,omitempty"`
	ValueRange           *Range           `json:"valueRange,omitempty"`
	ValuePeriod          *Period          `json:"valuePeriod,omitempty"`
	ValueRatio           *Ratio           `json:"valueRatio,omitempty"`
	ValueReference       *Reference       `json:"valueReference,omitempty"`
	ValueExpression      *Expression      `json:"valueExpression,omitempty"`
}

// Annotation is a resources annotation
type Annotation struct {
	ID              *string               `json:"id,omitempty"`
	AuthorReference *Reference            `json:"authorReference,omitempty" swaggertype:"primitive,string"`
	AuthorString    *string               `json:"authorString,omitempty"`
	Time            *scalarutils.DateTime `json:"time,omitempty" swaggertype:"primitive,string"`
	Text            scalarutils.Markdown  `json:"text,omitempty" swaggertype:"primitive,string"`
}

// Range models the range of a resource
type Range struct {
	ID   string   `json:"id,omitempty"`
	Low  Quantity `json:"low,omitempty"`
	High Quantity `json:"high,omitempty"`
}

// Ratio defines the ratio of a resource
type Ratio struct {
	ID          string   `json:"id,omitempty"`
	Numerator   Quantity `json:"numerator,omitempty"`
	Denominator Quantity `json:"denominator,omitempty"`
}

// QuestionnaireItem contains the questionnaires questions
type QuestionnaireItem struct {
	ID                *string                          `json:"id,omitempty"`
	Meta              *Meta                            `json:"meta,omitempty"`
	Extension         []*Extension                     `json:"extension,omitempty"`
	ModifierExtension []*Extension                     `json:"modifierExtension,omitempty"`
	LinkID            *string                          `json:"linkId,omitempty"`
	Definition        *scalarutils.URI                 `json:"definition,omitempty"`
	Code              []*Coding                        `json:"code,omitempty"`
	Prefix            *string                          `json:"prefix,omitempty"`
	Text              string                           `json:"text,omitempty"`
	Type              *scalarutils.Code                `json:"type,omitempty"`
	EnableWhen        []*QuestionnaireItemEnableWhen   `json:"enableWhen,omitempty"`
	EnableBehavior    *scalarutils.Code                `json:"enableBehavior,omitempty"`
	DisabledDisplay   *scalarutils.Code                `json:"disabledDisplay,omitempty"`
	Required          *bool                            `json:"required,omitempty"`
	Repeats           *bool                            `json:"repeats,omitempty"`
	ReadOnly          *bool                            `json:"readOnly,omitempty"`
	MaxLength         *int                             `json:"maxLength,omitempty"`
	AnswerValueSet    *string                          `json:"answerValueSet,omitempty"`
	AnswerOption      []*QuestionnaireItemAnswerOption `json:"answerOption,omitempty"`
	Initial           []*QuestionnaireItemInitial      `json:"initial,omitempty"`
	Item              []*QuestionnaireItem             `json:"item,omitempty"`
}

// QuestionnaireItemEnableWhen models an output that represents the initial values of a questionnaire item
type QuestionnaireItemEnableWhen struct {
	ID                *string               `json:"id,omitempty"`
	Extension         []*Extension          `json:"extension,omitempty"`
	ModifierExtension []*Extension          `json:"modifierExtension,omitempty"`
	Question          *string               `json:"question,omitempty"`
	Operator          *scalarutils.Code     `json:"operator,omitempty"`
	AnswerBoolean     *bool                 `json:"answerBoolean,omitempty"`
	AnswerDecimal     *float64              `json:"answerDecimal,omitempty"`
	AnswerInteger     *int                  `json:"answerInteger,omitempty"`
	AnswerDate        *scalarutils.Date     `json:"answerDate,omitempty"`
	AnswerDateTime    *scalarutils.DateTime `json:"answerDateTime,omitempty"`
	AnswerTime        *scalarutils.DateTime `json:"answerTime,omitempty"`
	AnswerString      *string               `json:"answerString,omitempty"`
	AnswerCoding      *Coding               `json:"answerCoding,omitempty"`
	AnswerQuantity    *Quantity             `json:"answerQuantity,omitempty"`
	AnswerReference   *Reference            `json:"answerReference,omitempty"`
}

// QuestionnaireItemAnswerOption defines the answer option for a Questionnaire item
type QuestionnaireItemAnswerOption struct {
	ID                *string           `json:"id,omitempty"`
	Extension         []*Extension      `json:"extension,omitempty"`
	ModifierExtension []*Extension      `json:"modifierExtension,omitempty"`
	ValueInteger      int               `json:"valueInteger,omitempty"`
	ValueDate         *scalarutils.Date `json:"valueDate,omitempty"`
	ValueString       *string           `json:"valueString,omitempty"`
	ValueCoding       *Coding           `json:"valueCoding,omitempty"`
	ValueReference    *Reference        `json:"valueReference,omitempty"`
	InitialSelected   *bool             `json:"initialSelected,omitempty"`
}

// QuestionnaireItemInitial models the initial values of a Questionnaire item.
type QuestionnaireItemInitial struct {
	ID                *string               `json:"id,omitempty"`
	Extension         []*Extension          `json:"extension,omitempty"`
	ModifierExtension []*Extension          `json:"modifierExtension,omitempty"`
	ValueBoolean      *bool                 `json:"valueBoolean,omitempty"`
	ValueDecimal      *float64              `json:"valueDecimal,omitempty"`
	ValueInteger      *int                  `json:"valueInteger,omitempty"`
	ValueDate         *scalarutils.Date     `json:"valueDate,omitempty"`
	ValueDateTime     *scalarutils.DateTime `json:"valuescalarutils.DateTime,omitempty"`
	ValueString       *string               `json:"valueString,omitempty"`
	ValueURI          *scalarutils.URI      `json:"valueUri,omitempty"`
	ValueAttachment   *Attachment           `json:"valueAttachment,omitempty"`
	ValueCoding       *Coding               `json:"valueCoding,omitempty"`
	ValueQuantity     *Quantity             `json:"valueQuantity,omitempty"`
	ValueReference    *Reference            `json:"valueReference,omitempty"`
}

// UsageContext defines the questionnaires usage context
type UsageContext struct {
	ID                   *string          `json:"id,omitempty"`
	Extension            []*Extension     `json:"extension,omitempty"`
	Code                 *Coding          `json:"code,omitempty"`
	ValueCodeableConcept *CodeableConcept `json:"valueCodeableConcept,omitempty"`
	ValueQuantity        *Quantity        `json:"valueQuantity,omitempty"`
	ValueRange           *Range           `json:"valueRange,omitempty"`
	ValueReference       *Reference       `json:"valueReference,omitempty"`
}

// Identifier models tha data is associated with a resource
type Identifier struct {
	ID       *string          `json:"id,omitempty"`
	Use      *string          `json:"use,omitempty"`
	Type     *CodeableConcept `json:"type,omitempty"`
	System   *scalarutils.URI `json:"system,omitempty" swaggertype:"primitive,string"`
	Value    *string          `json:"value,omitempty"`
	Period   *Period          `json:"period,omitempty"`
	Assigner *Reference       `json:"assigner,omitempty"`
}

// CodeableConcept represents a codeable concept associated with a questionnaire item
type CodeableConcept struct {
	ID     *string   `json:"id,omitempty"`
	Coding []*Coding `json:"coding,omitempty"`
	Text   string    `json:"text,omitempty"`
}

// Period is the questionnaires period
type Period struct {
	ID    *string               `json:"id,omitempty"`
	Start *scalarutils.DateTime `json:"start,omitempty" swaggertype:"primitive,string"`
	End   *scalarutils.DateTime `json:"end,omitempty" swaggertype:"primitive,string"`
}

// Meta models questionnaire output metadata
type Meta struct {
	VersionID   string    `json:"versionId,omitempty"`
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	Source      string    `json:"source,omitempty"`
	Tag         []Coding  `json:"tag,omitempty"`
	Security    []Coding  `json:"security,omitempty"`
}

// IsEntity marks a meta as an entity
func (m Meta) IsEntity() {}
