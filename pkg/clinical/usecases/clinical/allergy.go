package clinical

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/foundation"
)

type ClinicalImpl struct {
	infrastructure.Infrastructure
	foundation.FoundationImpl
	*slog.Logger
}

// NewClinicalImpl initializes new Clinical/Patient implementation
func NewClinicalImpl(infra infrastructure.Infrastructure, foundation foundation.FoundationImpl, log *slog.Logger) *ClinicalImpl {
	return &ClinicalImpl{
		infra,
		foundation,
		log,
	}
}

// CreateAllergyIntolerance creates a new allergy intolerance
func (c *ClinicalImpl) CreateAllergyIntolerance(ctx context.Context, input dto.AllergyInput) (*dto.Allergy, error) {
	err := helpers.Validate(input)
	if err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot record an allergy in a finished encounter")
	}

	encounterRef := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	allergyConcept, err := c.GetConcept(ctx, domain.TerminologySource(input.TerminologySource.Underscore()), input.Code)
	if err != nil {
		return nil, err
	}

	clinicalStatusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical")
	verificationSystem := "http://terminology.hl7.org/CodeSystem/allergyintolerance-verification"

	allergyIntoleranceTypeAllergy := domain.AllergyIntoleranceTypeEnumAllergy

	clinicalStatusCodeActive := "active"
	verificationDisplay := "confirmed"

	currentTime := time.Now().Format(time.DateOnly)

	allergyIntoleranceInput := domain.FHIRAllergyIntoleranceInput{
		ClinicalStatus: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{{
				System:  &clinicalStatusSystem,
				Code:    scalarutils.Code(clinicalStatusCodeActive),
				Display: clinicalStatusCodeActive,
			}},
		},
		Code: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&common.SystemURLICHI),
					Code:    scalarutils.Code(allergyConcept.ID),
					Display: allergyConcept.GetConceptDisplay(),
				},
			},
			Text: allergyConcept.GetConceptDisplay(),
		},
		Patient: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.Subject.ID,
			Reference: encounter.Resource.Subject.Reference,
			Display:   encounter.Resource.Subject.Display,
		},

		Encounter: &domain.FHIRReferenceInput{
			Reference: &encounterRef,
			Display:   *encounter.Resource.ID,
		},
		RecordedDate: (*scalarutils.DateTime)(&currentTime),
		Type: &domain.FHIRCodeableConceptInput{
			Text: allergyIntoleranceTypeAllergy.String(),
		},
		VerificationStatus: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&verificationSystem),
					Code:    scalarutils.Code(verificationDisplay),
					Display: verificationDisplay,
				},
			},
			Text: verificationDisplay,
		},
		OnsetDateTime: (*scalarutils.DateTime)(&currentTime),
		Participant: []*domain.AllergyIntoleranceParticipantInput{
			{
				Actor: &domain.FHIRReferenceInput{
					Reference: encounter.Resource.Subject.Reference,
					Display:   encounter.Resource.Subject.Display,
				},
			},
		},
	}

	if input.Reaction != nil {
		manifestationConcept, err := c.GetConcept(ctx, domain.TerminologySourceICD11WHO, input.Reaction.Code)
		if err != nil {
			return nil, err
		}

		severity := strings.ToLower(string(input.Reaction.Severity))

		allergyIntoleranceInput.Reaction = []*domain.FHIRAllergyintoleranceReactionInput{{
			Description: &severity,
			Manifestation: []*domain.FHIRCodeableReferenceInput{{
				Concept: &domain.FHIRCodeableConceptInput{
					Coding: []*domain.FHIRCodingInput{
						{
							System:  (*scalarutils.URI)(&common.SystemURLICHI),
							Code:    scalarutils.Code(manifestationConcept.ID),
							Display: manifestationConcept.GetConceptDisplay(),
						},
					},
					Text: manifestationConcept.GetConceptDisplay(),
				},
			}},
			Severity: (*domain.AllergyIntoleranceReactionSeverityEnum)(&severity),
		}}
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	allergyIntoleranceInput.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	allergyIntolerance, err := c.FHIR.CreateFHIRAllergyIntolerance(ctx, allergyIntoleranceInput)
	if err != nil {
		return nil, err
	}

	allergyIntoleranceObj := foundation.MapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(*allergyIntolerance.Resource)
	allergyIntoleranceObj.TerminologySource = input.TerminologySource

	return allergyIntoleranceObj, nil
}

// SearchAllergy is used to retrieve allergy from OCL
func (c *ClinicalImpl) SearchAllergy(ctx context.Context, name string, pagination dto.Pagination) (*dto.TerminologyConnection, error) {
	err := pagination.Validate()
	if err != nil {
		return nil, err
	}

	org := []string{string(dto.OrganisationSourceWHO)}
	source := []string{string(domain.TerminologySourceICD11WHO.Hyphenated())}

	conceptPage, err := c.OpenConceptLab.
		ListConcepts(ctx, org, source, true, &name, nil, nil, nil, nil, nil, nil, nil, nil, &pagination)
	if err != nil {
		return nil, err
	}

	terminologyPage := &dto.TerminologyConnection{
		TotalCount: conceptPage.Count,
		Edges:      []dto.TerminologyEdge{},
		PageInfo:   dto.PageInfo{},
	}

	if conceptPage.Next != nil {
		terminologyPage.PageInfo.HasNextPage = true
		terminologyPage.PageInfo.StartCursor = conceptPage.Next
	}

	if conceptPage.Previous != nil {
		terminologyPage.PageInfo.HasPreviousPage = true
		terminologyPage.PageInfo.EndCursor = conceptPage.Previous
	}

	for _, concept := range conceptPage.Results {
		terminologyPage.Edges = append(terminologyPage.Edges, dto.TerminologyEdge{
			Node: dto.Terminology{
				Code:   concept.ID,
				System: domain.TerminologySource(concept.Source),
				Name:   concept.GetConceptDisplay(),
			},
		})
	}

	return terminologyPage, nil
}

// GetAllergyIntolerance fetches all the allergy intolerance from FHIR by allergy intolerance ID
func (c *ClinicalImpl) GetAllergyIntolerance(ctx context.Context, id string) (*dto.Allergy, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid allergy intolerance id: %s", id)
	}

	allergyIntolerance, err := c.FHIR.GetFHIRAllergyIntolerance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to search for allergy intolerance: %w", err)
	}

	intolerance := foundation.MapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(*allergyIntolerance.Resource)

	return intolerance, nil
}

// ListPatientAllergies is used to list all allergies associated with a specific patient
func (c *ClinicalImpl) ListPatientAllergies(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.AllergyConnection, error) {
	_, err := uuid.Parse(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient id: %s", patientID)
	}

	err = pagination.Validate()
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", patientID)

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	allergyResponses, err := c.FHIR.SearchPatientAllergyIntolerance(ctx, patientReference, *identifiers, pagination)
	if err != nil {
		return nil, err
	}

	patientAllergyIntolerances := []*dto.Allergy{}

	for _, allergyResponse := range allergyResponses.Allergies {
		patientAllergyIntolerances = append(patientAllergyIntolerances, foundation.MapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(allergyResponse))
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     allergyResponses.HasNextPage,
		EndCursor:       &allergyResponses.NextCursor,
		HasPreviousPage: allergyResponses.HasPreviousPage,
		StartCursor:     &allergyResponses.PreviousCursor,
	}

	connection := dto.CreateAllergyConnection(patientAllergyIntolerances, pageInfo, allergyResponses.TotalCount)

	return &connection, nil
}
