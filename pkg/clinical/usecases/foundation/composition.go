package foundation

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure"
)

type FoundationImpl struct {
	infrastructure.Infrastructure
	*slog.Logger
}

func NewFoundationImpl(infra infrastructure.Infrastructure,
	log *slog.Logger,
) *FoundationImpl {
	return &FoundationImpl{
		infra,
		log,
	}
}

// CreateComposition creates a new composition
func (c *FoundationImpl) CreateComposition(ctx context.Context, input dto.CompositionInput) (*dto.Composition, error) {
	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot record a composition in a completed encounter")
	}

	patient, err := c.FHIR.GetFHIRPatient(ctx, encounter.Resource.Subject.ResourceID())
	if err != nil {
		return nil, err
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	patientRef := fmt.Sprintf("Patient/%s", *patient.Resource.ID)
	patientType := "Patient"
	compositionTitle := fmt.Sprintf("%s's assessment note", patient.Resource.Name[0].Text)

	encounterRef := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	encounterType := scalarutils.URI("Encounter")
	organizationRef := fmt.Sprintf("Organization/%s", identifiers.OrganizationID)

	today := time.Now()

	date, err := scalarutils.NewDate(today.Day(), int(today.Month()), today.Year())
	if err != nil {
		return nil, err
	}

	compositionCategoryCode, err := c.mapCategoryEnumToCode(input.Category)
	if err != nil {
		return nil, err
	}

	compositionConcept, err := c.mapCompositionConcepts(ctx, compositionCategoryCode, common.LOINCProgressNoteCode)
	if err != nil {
		return nil, err
	}

	status := input.Status.ToCode()
	generatedTextStatus := "generated"
	title := compositionConcept.CompositionCategoryConcept.GetConceptDisplay()

	compositionInput := domain.FHIRCompositionInput{
		Status: (*domain.CompositionStatusEnum)(&status),
		Type: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&common.LoincSystemURL),
					Code:    scalarutils.Code(compositionConcept.CompositionTypeConcept.ID),
					Display: compositionConcept.CompositionTypeConcept.GetConceptDisplay(),
				},
			},
			Text: compositionConcept.CompositionTypeConcept.GetConceptDisplay(),
		},
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&common.LoincSystemURL),
						Code:    scalarutils.Code(compositionConcept.CompositionCategoryConcept.ID),
						Display: title,
					},
				},
				Text: title,
			},
		},
		Name: &compositionTitle,
		Subject: []*domain.FHIRReferenceInput{
			{
				Reference: &patientRef,
				Type:      (*scalarutils.URI)(&patientType),
				Display:   *patient.Resource.ID,
			},
		},
		Encounter: &domain.FHIRReferenceInput{
			Reference: &encounterRef,
			Display:   *encounter.Resource.ID,
			Type:      &encounterType,
		},
		Date: date,
		Author: []*domain.FHIRReferenceInput{
			{
				Reference: &organizationRef,
				Display:   identifiers.OrganizationID,
			},
		},
		Title: &compositionTitle,
		Section: []*domain.FHIRCompositionSectionInput{
			{
				Title: &title,
				Code: &domain.FHIRCodeableConceptInput{
					Coding: []*domain.FHIRCodingInput{
						{
							System:  (*scalarutils.URI)(&common.LoincSystemURL),
							Code:    scalarutils.Code(compositionConcept.CompositionCategoryConcept.ID),
							Display: title,
						},
					},
					Text: compositionConcept.CompositionTypeConcept.GetConceptDisplay(),
				},
				Author: []*domain.FHIRReferenceInput{
					{
						Reference: &organizationRef,
						Display:   identifiers.OrganizationID,
					},
				},
				Text: &domain.FHIRNarrativeInput{
					Status: (*domain.NarrativeStatusEnum)(&generatedTextStatus),
					Div:    scalarutils.XHTML("<div xmlns=\"http://www.w3.org/1999/xhtml\">Generated text.</div>"),
				},
			},
		},
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	compositionInput.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	composition, err := c.FHIR.CreateFHIRComposition(ctx, compositionInput)
	if err != nil {
		return nil, err
	}

	result := mapFHIRCompositionToCompositionDTO(*composition.Resource)

	return &result.Edges[0].Node, nil
}

func mapFHIRCompositionToCompositionDTO(composition domain.FHIRComposition) *dto.CompositionConnection {
	var compositionSection []*dto.Section

	for _, item := range composition.Section {
		var itemSubSections []*dto.Section

		if len(item.Section) > 0 {
			for _, section := range item.Section {
				itemSubSections = append(itemSubSections, &dto.Section{
					ID:     section.ID,
					Title:  section.Title,
					Code:   section.Code.ID,
					Author: section.Author[0].Reference,
					Text:   helpers.ExtractTextFromHTML(string(section.Text.Div)),
				})
			}
		}

		compositionSection = append(compositionSection, &dto.Section{
			ID:      item.ID,
			Title:   item.Title,
			Code:    &item.Code.Coding[0].Display,
			Author:  item.Author[0].Reference,
			Text:    helpers.ExtractTextFromHTML(string(item.Text.Div)),
			Section: itemSubSections,
		})
	}

	output := dto.Composition{
		ID:          *composition.ID,
		Type:        dto.CompositionType(composition.Type.Text),
		Status:      dto.CompositionStatusEnum(*composition.Status),
		PatientID:   composition.Subject[0].Display,
		EncounterID: composition.Encounter.Display,
		Date:        composition.Date,
		Section:     compositionSection,
	}

	if len(composition.Section) != 0 {
		output.Text = helpers.ExtractTextFromHTML(string(composition.Section[0].Text.Div))
	}

	if len(composition.Category) != 0 {
		output.Category = dto.CompositionCategory(composition.Category[0].Text)
	}

	return &dto.CompositionConnection{
		TotalCount: 0,
		Edges: []dto.CompositionEdge{
			{
				Node:   output,
				Cursor: "",
			},
		},
		PageInfo: dto.PageInfo{},
	}
}

// ListPatientCompositions lists a patient's compositions
func (c FoundationImpl) ListPatientCompositions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.CompositionConnection, error) {
	_, err := uuid.Parse(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient id: %s", patientID)
	}

	err = pagination.Validate()
	if err != nil {
		return nil, err
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	patient, err := c.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	patientRef := fmt.Sprintf("Patient/%s", *patient.Resource.ID)
	params := map[string]interface{}{
		"subject": patientRef,
		"_sort":   "date",
		"_total":  "accurate",
	}

	if encounterID != nil {
		encounterReference := fmt.Sprintf("Encounter/%s", *encounterID)
		params["encounter"] = encounterReference
	}

	if date != nil {
		params["date"] = date.AsTime().Format(time.DateOnly)
	}

	compositionsResponse, err := c.FHIR.SearchFHIRComposition(ctx, params, *identifiers, pagination)
	if err != nil {
		return nil, err
	}

	compositions := []dto.Composition{}

	for _, resource := range compositionsResponse.Compositions {
		composition := mapFHIRCompositionToCompositionDTO(resource)
		compositions = append(compositions, composition.Edges[0].Node)
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     compositionsResponse.HasNextPage,
		EndCursor:       &compositionsResponse.NextCursor,
		HasPreviousPage: compositionsResponse.HasPreviousPage,
		StartCursor:     &compositionsResponse.PreviousCursor,
	}

	connection := dto.CreateCompositionConnection(compositions, pageInfo, compositionsResponse.TotalCount)

	return &connection, nil
}

// AppendNoteToComposition appends a note to the patient's composition information such as section
func (c *FoundationImpl) AppendNoteToComposition(ctx context.Context, id string, input dto.PatchCompositionInput) (*dto.Composition, error) {
	if id == "" {
		return nil, fmt.Errorf("a composition id is required")
	}

	composition, err := c.FHIR.GetFHIRComposition(ctx, id)
	if err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, composition.Resource.Encounter.Display)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot record a composition in a completed encounter")
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	organizationRef := fmt.Sprintf("Organization/%s", identifiers.OrganizationID)

	compositionCategoryCode, err := c.mapCategoryEnumToCode(input.Category)
	if err != nil {
		return nil, err
	}

	compositionCategoryConcept, err := c.GetConcept(ctx, domain.TerminologySourceLOINC, compositionCategoryCode)
	if err != nil {
		return nil, err
	}

	compositionSectionTextStatus := "generated"
	title := compositionCategoryConcept.GetConceptDisplay()

	compositionSection := &domain.FHIRCompositionSection{
		Title: &title,
		Code: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&compositionCategoryConcept.URL),
					Code:    scalarutils.Code(compositionCategoryConcept.ID),
					Display: compositionCategoryConcept.GetConceptDisplay(),
				},
			},
			Text: compositionCategoryConcept.GetConceptDisplay(),
		},
		Author: []*domain.FHIRReference{
			{
				Reference: &organizationRef,
			},
		},
		Focus: &domain.FHIRReference{},
		Text: &domain.FHIRNarrative{
			Status: (*domain.NarrativeStatusEnum)(&compositionSectionTextStatus),
			Div:    scalarutils.XHTML(input.Note),
		},
	}

	composition.Resource.Section = append(composition.Resource.Section, compositionSection)

	var sectionInput []*domain.FHIRCompositionSectionInput

	for _, s := range composition.Resource.Section {
		sectionInput = append(sectionInput, &domain.FHIRCompositionSectionInput{
			ID:    s.ID,
			Title: s.Title,
			Code: &domain.FHIRCodeableConceptInput{
				ID:     s.Code.ID,
				Coding: s.Code.Coding,
				Text:   compositionSectionTextStatus,
			},
			Author: []*domain.FHIRReferenceInput{
				{
					Reference: &organizationRef,
				},
			},
			Text: &domain.FHIRNarrativeInput{
				ID:     s.Text.ID,
				Status: s.Text.Status,
				Div:    s.Text.Div,
			},
		})

		if len(s.Section) > 0 {
			var nestedsectionInput []*domain.FHIRCompositionSectionInput

			for _, r := range s.Section {
				nestedsectionInput = append(nestedsectionInput, &domain.FHIRCompositionSectionInput{
					ID:    r.ID,
					Title: r.Title,
					Code: &domain.FHIRCodeableConceptInput{
						ID:     r.Code.ID,
						Coding: r.Code.Coding,
						Text:   compositionSectionTextStatus,
					},
					Author: []*domain.FHIRReferenceInput{
						{
							Reference: &organizationRef,
						},
					},
					Text: &domain.FHIRNarrativeInput{
						ID:     r.Text.ID,
						Status: r.Text.Status,
						Div:    r.Text.Div,
					},
				})
			}

			for _, m := range sectionInput {
				m.Section = append(m.Section, nestedsectionInput...)
			}
		}
	}

	compositionInput := &domain.FHIRCompositionInput{
		Section: sectionInput,
	}

	output, err := c.FHIR.PatchFHIRComposition(ctx, id, *compositionInput)
	if err != nil {
		return nil, err
	}

	result := mapFHIRCompositionToCompositionDTO(*output)

	return &result.Edges[0].Node, nil
}
