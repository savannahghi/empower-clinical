package foundation

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// RecordConsent records a user consent
func (c *FoundationImpl) RecordConsent(ctx context.Context, input dto.ConsentInput) (*dto.ConsentOutput, error) {
	year, month, day := time.Now().UTC().Date()
	currentDate := scalarutils.Date{
		Year:  year,
		Month: int(month),
		Day:   day,
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumCompleted {
		return nil, fmt.Errorf("cannot create a consent in a finished encounter")
	}

	encounterRef := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	encounterReference := &domain.FHIRReference{
		ID:        encounter.Resource.ID,
		Reference: &encounterRef,
		Display:   *encounter.Resource.ID,
	}

	subjectReference := &domain.FHIRReference{
		Reference: encounter.Resource.Subject.Reference,
		Display:   *encounter.Resource.ID,
	}

	var system scalarutils.URI = "http://terminology.hl7.org/CodeSystem/consentcategorycodes"

	code := scalarutils.Code("acd")

	coding := &domain.FHIRCoding{
		System:  &system,
		Code:    &code,
		Display: "Advance Directive",
	}

	policyRule := domain.FHIRCodeableConcept{
		Text:   "cric",
		Coding: []*domain.FHIRCoding{coding},
	}

	consentProvision := []*domain.FHIRConsentProvision{
		{
			Data: []domain.FHIRConsentProvisionData{
				{
					Meaning:   domain.ConsentDataMeaningRelated,
					Reference: encounterReference,
				},
			},
		},
	}

	category := &domain.FHIRCodeableConcept{
		Text:   "Advance Directive",
		Coding: []*domain.FHIRCoding{coding},
	}

	organizationID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	organizationRef := fmt.Sprintf("Organization/%s", organizationID)
	organizationReference := &domain.FHIRReference{
		Reference: &organizationRef,
		Display:   organizationID,
	}
	verified := true
	verification := []domain.FHIRConsentVerificaiton{
		{
			VerifiedWith: subjectReference,
			VerifiedBy:   organizationReference,
			Verfied:      &verified,
			VerificationDate: []*scalarutils.Date{
				{
					Year:  currentDate.Year,
					Month: currentDate.Month,
					Day:   currentDate.Day,
				},
			},
		},
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	consentMeta := domain.FHIRMetaInput{
		Tag: tags,
	}

	status := domain.ConsentStatusActive
	decision := domain.ConsentDecisionEnum(input.Decision)
	consent := domain.FHIRConsent{
		Status:          &status,
		Decision:        &decision,
		Subject:         subjectReference,
		Grantor:         []*domain.FHIRReference{subjectReference},
		Grantee:         []*domain.FHIRReference{organizationReference},
		Manager:         []*domain.FHIRReference{organizationReference},
		Controller:      []*domain.FHIRReference{organizationReference},
		Category:        []*domain.FHIRCodeableConcept{category},
		RegulatoryBasis: []domain.FHIRCodeableConcept{policyRule},
		Meta:            &consentMeta,
		Date:            &currentDate,
		Verifcation:     verification,
		Provision:       consentProvision,
	}

	additionalNarrativeStatus := domain.NarrativeStatusEnumAdditional.String()

	consent.Text = utils.NarrativeGenerator(input.ScreeningType.String(), &additionalNarrativeStatus)

	resp, err := c.FHIR.CreateFHIRConsent(ctx, consent)
	if err != nil {
		return nil, err
	}

	output := &dto.ConsentOutput{
		Status: resp.Status,
	}

	return output, nil
}
