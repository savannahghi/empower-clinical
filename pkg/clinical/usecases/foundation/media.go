package foundation

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/silurlshortener"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// UploadMedia uploads media to GCS and creates the resource in FHIR
func (c *FoundationImpl) UploadMedia(ctx context.Context, encounterID, serviceRequestID string, file io.Reader, contentType string) (*dto.Media, error) {
	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	facility, err := c.FHIR.GetFHIROrganization(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	patientID := encounter.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)
	mediaObjectName := fmt.Sprintf("%s@%s", patientReference, time.Now())

	patient, err := c.FHIR.GetFHIRPatient(ctx, *patientID)
	if err != nil {
		return nil, err
	}

	mediaUploadOutput, err := c.Upload.UploadMedia(ctx, mediaObjectName, file, contentType)
	if err != nil {
		return nil, err
	}

	shortener := silurlshortener.ShortenURLPayload{
		LongURL:         mediaUploadOutput.SignedURL,
		Title:           fmt.Sprintf("%s%s", encounter.Resource.Subject.Display, "file"),
		Domain:          serverutils.MustGetEnvVar("URL_SHORTENER_DOMAIN"),
		ShortCodeLength: 6,
	}

	url, err := c.URLShorteningService.Shorten(ctx, &shortener)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	serviceRequestRef := fmt.Sprintf("ServiceRequest/%s", serviceRequestID)

	orgRef := fmt.Sprintf("Organization/%s", *facility.Resource.ID)
	orgType := scalarutils.URI("Organization")

	title := fmt.Sprintf("%s's Test Result", patient.Resource.Names())

	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))
	preliminaryDocStatus := domain.CompositionStatusEnumPreliminary

	payload := &domain.FHIRDocumentReferenceInput{
		Subject: &domain.FHIRReferenceInput{
			ID:        patientID,
			Reference: &patientReference,
			Display:   patient.Resource.Names(),
		},
		Content: []domain.FHIRDocumentReferenceContent{
			{
				Attachment: domain.FHIRAttachment{
					ContentType: (*scalarutils.Code)(&mediaUploadOutput.ContentType),
					URL:         (*scalarutils.URL)(&url.ShortURL),
					Title:       &title,
				},
			},
		},
		Context: []*domain.FHIRReferenceInput{
			{
				Reference: &encounterReference,
				Display:   *encounter.Resource.ID,
			},
		},
		Custodian: &domain.FHIRReferenceInput{
			ID:        facility.Resource.ID,
			Reference: &orgRef,
			Display:   *facility.Resource.Name,
			Type:      &orgType,
		},
		Author: []*domain.FHIRReferenceInput{
			{
				Reference: &orgRef,
				Display:   facilityID,
			},
		},
		Date: &instant,
		Type: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&common.LoincSystemURL),
					Code:    common.LOINCMedicalRecordCode,
					Display: "Medical records",
				},
			},
			Text: "Medical records",
		},
		DocStatus: &preliminaryDocStatus,
		Status:    domain.DocumentReferenceStatusEnumCurrent,
	}

	if serviceRequestID != "" {
		payload.BasedOn = []*domain.FHIRReferenceInput{
			{
				Reference: &serviceRequestRef,
				Display:   serviceRequestID,
			},
		}
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	payload.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	_, err = c.FHIR.CreateFHIRDocumentReference(ctx, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	output := &dto.Media{
		ID:          mediaUploadOutput.ID,
		PatientID:   *patientID,
		PatientName: patient.Resource.Names(),
		MediaLink:   url.ShortURL,
		Name:        mediaUploadOutput.Name,
		ContentType: mediaUploadOutput.ContentType,
	}

	return output, nil
}

// ListPatientMedia list the patients media resources
func (c *FoundationImpl) ListPatientMedia(ctx context.Context, encounterID, serviceRequestID string, pagination dto.Pagination) (*dto.MediaConnection, error) {
	err := pagination.Validate()
	if err != nil {
		return nil, err
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	params := map[string]interface{}{
		"encounter": *encounter.Resource.ID,
		"related":   serviceRequestID,
		"_total":    "accurate",
	}

	mediaResources, err := c.FHIR.SearchFHIRDocumentReference(ctx, params, *identifiers, pagination)
	if err != nil {
		return nil, err
	}

	patientMediaList := []*dto.Media{}

	for _, documentRef := range mediaResources.DocumentReferences {
		media, err := c.mapFHIRDocumentReferenceToMediaDTO(ctx, documentRef)
		if err != nil {
			utils.ReportErrorToSentry(err)
			continue
		}

		patientMediaList = append(patientMediaList, media)
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     mediaResources.HasNextPage,
		EndCursor:       &mediaResources.NextCursor,
		HasPreviousPage: mediaResources.HasPreviousPage,
		StartCursor:     &mediaResources.PreviousCursor,
	}

	connection := dto.CreateMediaConnection(patientMediaList, pageInfo, mediaResources.TotalCount)

	return &connection, nil
}

func (c *FoundationImpl) mapFHIRDocumentReferenceToMediaDTO(ctx context.Context, fhirDocRef domain.FHIRDocumentReference) (*dto.Media, error) {
	attachment, err := fhirDocRef.GetDocumentAttachment()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	shortener := silurlshortener.ShortenURLPayload{
		LongURL:         string(*attachment.URL),
		Title:           fmt.Sprintf("%s%s", fhirDocRef.Subject.Display, "file"),
		Domain:          serverutils.MustGetEnvVar("URL_SHORTENER_DOMAIN"),
		ShortCodeLength: 6,
	}

	url, err := c.URLShorteningService.Shorten(ctx, &shortener)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	media := &dto.Media{
		ID:          fhirDocRef.ID,
		PatientID:   *fhirDocRef.Subject.ID,
		PatientName: fhirDocRef.Subject.Display,
		MediaLink:   url.ShortURL,
		Name:        *attachment.Title,
		ContentType: string(*attachment.ContentType),
		SignedURL:   url.ShortURL,
	}

	return media, nil
}
