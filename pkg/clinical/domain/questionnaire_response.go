package domain

import (
	"strconv"
	"strings"
)

// FHIRQuestionnaireResponse models questionnaire response resource
type FHIRQuestionnaireResponse struct {
	ID                *string                         `json:"id,omitempty"`
	Meta              *FHIRMetaInput                  `json:"meta,omitempty"`
	ImplicitRules     *string                         `json:"implicitRules,omitempty"`
	Language          *string                         `json:"language,omitempty"`
	Text              *FHIRNarrative                  `json:"text,omitempty"`
	Extension         []*FHIRExtension                `json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension                `json:"modifierExtension,omitempty"`
	Identifier        []*FHIRIdentifier               `json:"identifier,omitempty"`
	BasedOn           []*FHIRReference                `json:"basedOn,omitempty"`
	PartOf            []*FHIRReference                `json:"partOf,omitempty"`
	Questionnaire     *string                         `json:"questionnaire"`
	Status            QuestionnaireResponseStatusEnum `json:"status"`
	Subject           *FHIRReference                  `json:"subject,omitempty"`
	Encounter         *FHIRReference                  `json:"encounter,omitempty"`
	Authored          *string                         `json:"authored,omitempty"`
	Author            *FHIRReference                  `json:"author,omitempty"`
	Source            *FHIRReference                  `json:"source,omitempty"`
	Item              []FHIRQuestionnaireResponseItem `json:"item,omitempty"`
}

type FHIRQuestionnaireResponseItem struct {
	ID                *string                               `json:"id,omitempty"`
	Extension         []*FHIRExtension                      `json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension                      `json:"modifierExtension,omitempty"`
	LinkID            string                                `json:"linkId"`
	Definition        *string                               `json:"definition,omitempty"`
	Text              *string                               `json:"text,omitempty"`
	Answer            []FHIRQuestionnaireResponseItemAnswer `json:"answer,omitempty"`
	Item              []FHIRQuestionnaireResponseItem       `json:"item,omitempty"`
}

// FHIRQuestionnaireResponseItemAnswer models item answer object of questionnaire response resource
type FHIRQuestionnaireResponseItemAnswer struct {
	ID                *string                         `json:"id,omitempty"`
	Extension         []*FHIRExtension                `json:"extension,omitempty"`
	ModifierExtension []*FHIRExtension                `json:"modifierExtension,omitempty"`
	ValueBoolean      *bool                           `json:"valueBoolean"`
	ValueDecimal      *float64                        `json:"valueDecimal"`
	ValueInteger      *int                            `json:"valueInteger"`
	ValueDate         *string                         `json:"valueDate"`
	ValueDateTime     *string                         `json:"valueDateTime"`
	ValueTime         *string                         `json:"valueTime"`
	ValueString       *string                         `json:"valueString"`
	ValueURI          *string                         `json:"valueUri"`
	ValueAttachment   *FHIRAttachment                 `json:"valueAttachment"`
	ValueCoding       *FHIRCoding                     `json:"valueCoding"`
	ValueQuantity     *FHIRQuantity                   `json:"valueQuantity"`
	ValueReference    *FHIRReference                  `json:"valueReference"`
	Item              []FHIRQuestionnaireResponseItem `json:"item,omitempty"`
}

// SimpleQuestionnaireResponse is used to model a simple response of the questionnaire. We return a 'just enough'
// set of fields to make it easier to render on the UI
type SimpleQuestionnaireResponse struct {
	Group     *string     `json:"group,omitempty"`
	Questions []Questions `json:"questions,omitempty"`
}

type Questions struct {
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`

	// This will be used to hold nested questions
	ChildQuestions []Questions `json:"childQuestions,omitempty"`
}

// GetFHIRQuestionnaireResponse returns all the responses in a Question:Answer format
func (f *FHIRQuestionnaireResponse) GetFHIRQuestionnaireResponse() []*SimpleQuestionnaireResponse {
	responses := []*SimpleQuestionnaireResponse{}

	for _, group := range f.Item {
		groupAnswers := SimpleQuestionnaireResponse{}

		if group.Text != nil {
			groupAnswers.Group = group.Text

			for _, questionItem := range group.Item {
				if questionItem.Text != nil && strings.Contains(strings.ToLower(*questionItem.Text), "score") {
					continue
				}

				questions := Questions{
					Question: *questionItem.Text,
				}

				for _, answer := range questionItem.Answer {
					questions.Answer = answer.ToString()

					// Check if there are nested questions
					// Presence of the 'item' key indicates that there is a nested question
					if answer.Item != nil {
						questions.ChildQuestions = extractChildQuestions(answer.Item)
					}
				}

				groupAnswers.Questions = append(groupAnswers.Questions, questions)
			}

			responses = append(responses, &groupAnswers)
		}
	}

	return responses
}

// extractChildQuestions extracts child questions from nested items.
func extractChildQuestions(items []FHIRQuestionnaireResponseItem) []Questions {
	var childQuestions []Questions

	for _, item := range items {
		childQuestion := Questions{
			Question: *item.Text,
		}

		for _, answer := range item.Answer {
			childQuestion.Answer = answer.ToString()

			// Recursively extract further nested questions
			if answer.Item != nil {
				childQuestion.ChildQuestions = extractChildQuestions(answer.Item)
			}
		}

		childQuestions = append(childQuestions, childQuestion)
	}

	return childQuestions
}

// ToString converts the existing item answer to a string
func (q *FHIRQuestionnaireResponseItemAnswer) ToString() string {
	switch {
	case q.ValueBoolean != nil:
		return strconv.FormatBool(*q.ValueBoolean)
	case q.ValueDecimal != nil:
		return strconv.FormatFloat(*q.ValueDecimal, 'f', -1, 64)
	case q.ValueInteger != nil:
		return strconv.Itoa(*q.ValueInteger)
	case q.ValueDate != nil:
		return *q.ValueDate
	case q.ValueDateTime != nil:
		return *q.ValueDateTime
	case q.ValueTime != nil:
		return *q.ValueTime
	case q.ValueString != nil:
		return *q.ValueString
	case q.ValueURI != nil:
		return *q.ValueURI
	case q.ValueCoding != nil:
		return q.ValueCoding.ToString()

	default:
		return ""
	}
}

// FHIRQuestionnaireResponseRelayPayload is used to return a single instance of Questionnaire response
type FHIRQuestionnaireResponseRelayPayload struct {
	Resource *FHIRQuestionnaireResponse `json:"resource,omitempty"`
}
