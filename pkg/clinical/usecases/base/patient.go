package base

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/foundation"
)

type BaseImpl struct {
	infrastructure.Infrastructure
	clinical.ClinicalImpl
	foundation.FoundationImpl
	*slog.Logger
}

// NewBaseImpl initializes new Clinical/Patient implementation
func NewBaseImpl(infra infrastructure.Infrastructure, clinical clinical.ClinicalImpl, foundation foundation.FoundationImpl, log *slog.Logger) *BaseImpl {
	return &BaseImpl{
		infra,
		clinical,
		foundation,
		log,
	}
}

// GetMedicalData returns a limited subset of specific medical data that for a specific patient
// These include: Allergies, Viral Load, Body Mass Index, Weight, CD4 Count using their respective OCL CIEL Terminology
// For each category the latest three records are fetched
func (c *BaseImpl) GetMedicalData(ctx context.Context, patientID string) (*dto.MedicalData, error) {
	data := &dto.MedicalData{}
	filterParams := map[string]interface{}{
		"patient": fmt.Sprintf("Patient/%v", patientID),
		"_count":  common.MedicalDataCount,
		"_sort":   "_lastUpdated",
		"_total":  "accurate",
	}

	fields := []string{
		"Regimen",
		"AllergyIntolerance",
		"Weight",
		"BMI",
		"ViralLoad",
		"CD4Count",
	}

	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	for _, field := range fields {
		switch field {
		case "Regimen":
			conn, err := c.FHIR.SearchFHIRMedicationStatement(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				if edge.Node.ID == nil {
					continue
				}

				if edge.Node.Status == nil {
					continue
				}

				if edge.Node.MedicationCodeableConcept.Coding == nil {
					continue
				}

				if len(edge.Node.MedicationCodeableConcept.Coding) < 1 {
					continue
				}

				if edge.Node.Subject.ID == nil {
					continue
				}

				data.Regimen = append(data.Regimen, mapFHIRMedicationStatementToMedicationStatementDTO(edge.Node))
			}
		case "AllergyIntolerance":
			conn, err := c.FHIR.SearchFHIRAllergyIntolerance(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Allergies {
				if edge.ID == nil {
					continue
				}

				if edge.Patient == nil {
					continue
				}

				if edge.Patient.ID == nil {
					continue
				}

				if edge.Code == nil {
					continue
				}

				if edge.Code.Coding == nil {
					continue
				}

				if len(edge.Code.Coding) < 1 {
					continue
				}

				if edge.Code.Coding[0].ID == nil {
					continue
				}

				data.Allergies = append(data.Allergies, foundation.MapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(edge))
			}

		case "Weight":
			filterParams["code"] = common.WeightLOINCTerminologyCode

			conn, err := c.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, observation := range conn.Observations {
				if !hasNilInObservation(observation) {
					data.Weight = append(data.Weight, foundation.MapFHIRObservationToObservationDTO(observation))
				}
			}

		case "BMI":
			filterParams["code"] = common.BMILOINCTerminologyCode

			conn, err := c.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, observation := range conn.Observations {
				if !hasNilInObservation(observation) {
					data.BMI = append(data.BMI, foundation.MapFHIRObservationToObservationDTO(observation))
				}
			}

		case "ViralLoad":
			filterParams["code"] = common.ViralLoadLOINCTerminologyCode

			conn, err := c.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, observation := range conn.Observations {
				if !hasNilInObservation(observation) {
					data.ViralLoad = append(data.ViralLoad, foundation.MapFHIRObservationToObservationDTO(observation))
				}
			}

		case "CD4Count":
			filterParams["code"] = common.CD4CountLOINCTerminologyCode

			conn, err := c.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, observation := range conn.Observations {
				if !hasNilInObservation(observation) {
					data.CD4Count = append(data.CD4Count, foundation.MapFHIRObservationToObservationDTO(observation))
				}
			}
		}
	}

	return data, nil
}

func hasNilInObservation(observation domain.FHIRObservation) bool {
	if observation.ID == nil {
		return true
	}

	if observation.Code == nil {
		return true
	}

	if observation.Code.Coding == nil {
		return true
	}

	if len(observation.Code.Coding) < 1 {
		return true
	}

	if observation.Subject == nil {
		return true
	}

	if observation.Subject.ID == nil {
		return true
	}

	return false
}

func mapFHIRMedicationStatementToMedicationStatementDTO(fhirAllergyIntolerance *domain.FHIRMedicationStatement) *dto.MedicationStatement {
	return &dto.MedicationStatement{
		ID:     *fhirAllergyIntolerance.ID,
		Status: dto.MedicationStatementStatusEnum(*fhirAllergyIntolerance.Status),
		Medication: dto.Medication{
			Name: fhirAllergyIntolerance.MedicationCodeableConcept.Coding[0].Display,
			Code: foundation.FhirAllergyIntoleranceClinicalStatusURL,
		},
		PatientID: fhirAllergyIntolerance.Subject.ResourceID(),
	}
}

// CreatePatient creates a new patient
func (c *BaseImpl) CreatePatient(ctx context.Context, input dto.PatientInput) (*dto.Patient, error) {
	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	facility, err := c.FHIR.GetFHIROrganization(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	nameInput := &domain.NameInput{
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		OtherNames: input.OtherNames,
	}

	phoneNumbers := []*domain.PhoneNumberInput{}

	for _, contact := range input.Contacts {
		if contact.Type == dto.ContactTypePhoneNumber {
			number := &domain.PhoneNumberInput{
				Msisdn: contact.Value,
			}
			phoneNumbers = append(phoneNumbers, number)
		}
	}

	registrationInput := domain.SimplePatientRegistrationInput{
		Names:        []*domain.NameInput{nameInput},
		BirthDate:    input.BirthDate,
		PhoneNumbers: phoneNumbers,
		Gender:       string(input.Gender),
		Active:       true,
	}

	if len(input.Identifiers) > 0 {
		documents := make([]*domain.IdentificationDocument, 0, len(input.Identifiers))

		for _, identifier := range input.Identifiers {
			if identifier.Value != "" {
				doc := &domain.IdentificationDocument{
					DocumentType:   domain.IDDocumentType(identifier.Type),
					DocumentNumber: identifier.Value,
				}

				documents = append(documents, doc)
			}
		}

		registrationInput.IdentificationDocuments = documents
	}

	patientInput, err := c.SimplePatientRegistrationInputToPatientInput(ctx, registrationInput, *facility.Resource.ID)
	if err != nil {
		return nil, err
	}

	orgRef := fmt.Sprintf("Organization/%s", *facility.Resource.ID)
	orgType := scalarutils.URI("Organization")

	patientInput.ManagingOrganization = &domain.FHIRReferenceInput{
		ID:        facility.Resource.ID,
		Reference: &orgRef,
		Display:   *facility.Resource.Name,
		Type:      &orgType,
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	patientInput.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	patient, err := c.FHIR.CreateFHIRPatient(ctx, *patientInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRPatientToPatientDTO(patient.PatientRecord), nil
}

func (c *BaseImpl) PatchPatient(ctx context.Context, id string, input dto.PatientInput) (*dto.Patient, error) {
	if id == "" {
		return nil, fmt.Errorf("a patient ID is required")
	}

	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	registrationInput := domain.SimplePatientRegistrationInput{
		Gender:    string(input.Gender),
		BirthDate: input.BirthDate,
	}

	if input.FirstName != "" {
		nameInput := &domain.NameInput{
			FirstName:  input.FirstName,
			LastName:   input.LastName,
			OtherNames: input.OtherNames,
		}
		registrationInput.Names = append(registrationInput.Names, nameInput)
	}

	for _, contact := range input.Contacts {
		if contact.Type == dto.ContactTypePhoneNumber {
			number := &domain.PhoneNumberInput{
				Msisdn: contact.Value,
			}
			registrationInput.PhoneNumbers = append(
				registrationInput.PhoneNumbers,
				number,
			)
		}
	}

	for _, identifier := range input.Identifiers {
		doc := &domain.IdentificationDocument{
			DocumentType:   domain.IDDocumentType(identifier.Type),
			DocumentNumber: identifier.Value,
		}

		registrationInput.IdentificationDocuments = append(
			registrationInput.IdentificationDocuments,
			doc,
		)
	}

	patientInput, err := c.SimplePatientRegistrationInputToPatientInput(ctx, registrationInput, facilityID)
	if err != nil {
		return nil, err
	}

	patient, err := c.FHIR.PatchFHIRPatient(ctx, id, *patientInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRPatientToPatientDTO(patient), nil
}

func (c *BaseImpl) DeletePatient(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("a patient ID is required")
	}

	ok, err := c.FHIR.DeleteFHIRPatient(ctx, id)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func mapFHIRPatientToPatientDTO(patient *domain.FHIRPatient) *dto.Patient {
	numbers := []string{}

	for _, phone := range patient.Telecom {
		if *phone.System == domain.ContactPointSystemEnumPhone {
			numbers = append(numbers, *phone.Value)
		}
	}

	return &dto.Patient{
		ID:          *patient.ID,
		Active:      *patient.Active,
		Name:        patient.Name[0].Text,
		PhoneNumber: numbers,
		Gender:      dto.Gender(patient.Gender.String()),
		BirthDate:   patient.BirthDate,
	}
}

// GetPatientEverything retrieves all resources related to a patient.
// On success, the response body contains a JSON-encoded representation of a Bundle resource of type
// searchset, containing the results of the operation. Errors generated by the FHIR store contain a
// JSON-encoded OperationOutcome resource describing the reason for the error. If the request cannot
// be mapped to a valid API method on a FHIR store, a generic GCP error might be returned instead.
//
// The API supports various query parameters to retrieve patient compartments. These parameters include:
//
// (1). `_count" sets the optional parameter "_count": Maximum number of
// resources in a page. If not specified, 100 is used. May not be larger
// than 1000.
//
// (2). `_page_token` sets the optional parameter "_page_token": Used to retrieve
// the next or previous page of results when using pagination. Set
// `_page_token` to the value of _page_token set in next or previous
// page links' url. Next and previous page are returned in the response
// bundle's links field, where `link.relation` is "previous" or "next".
// Omit `_page_token` if no previous request has been made.
//
// (3). `since` sets the optional parameter "_since": If provided, only
// resources updated after this time are returned. The time uses the
// format YYYY-MM-DDThh:mm:ss.sss+zz:zz. For example,
// `2015-02-07T13:28:17.239+02:00` or `2017-01-01T00:00:00Z`. The time
// must be specified to the second and include a time zone.
//
// (4). `_type` sets the optional parameter "_type": String of comma-delimited
// FHIR resource types. If provided, only resources of the specified
// resource type(s) are returned. Specifying multiple `_type` parameters
// isn't supported. For example, the result of
// `_type=Observation&_type=Encounter` is undefined. Use
// `_type=Observation,Encounter` instead.
//
// (5). `end` sets the optional parameter "end": The response includes records
// prior to the end date. The date uses the format YYYY-MM-DD. If no end
// date is provided, all records subsequent to the start date are in
// scope.
//
// (6). `start` sets the optional parameter "start": The response includes
// records subsequent to the start date. The date uses the format
// YYYY-MM-DD. If no start date is provided, all records prior to the
// end date are in scope.
//
// (7). `fields` allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *BaseImpl) GetPatientEverything(ctx context.Context, patientID string, params *dto.PatientEverythingFilterParams) (*dto.PatientEverythingConnection, error) {
	_, err := uuid.Parse(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient id: %s", patientID)
	}

	patient, err := c.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	var searchParams map[string]interface{}

	if params != nil {
		sp := make(map[string]interface{}, 8)

		if params.Count != "" {
			sp["_count"] = params.Count
		} else {
			sp["_count"] = "100"
		}

		if params.Type != "" {
			sp["_type"] = params.Type
		}

		if params.GetPages != "" {
			sp["_getpages"] = params.GetPages
		}

		if params.GetPagesOffset != "" {
			sp["_getpagesoffset"] = params.GetPagesOffset
		}

		if params.PageToken != "" && params.GetPages == "" {
			sp["_page_token"] = params.PageToken
		}

		if params.Start != "" {
			sp["start"] = params.Start
		}

		if params.End != "" {
			sp["end"] = params.End
		}

		if params.Since != "" {
			sp["_since"] = params.Since
		}

		if len(sp) > 0 {
			searchParams = sp
		}
	}

	searchParams["_total"] = "accurate"

	patientData, err := c.FHIR.GetFHIRPatientEverything(ctx, *patient.Resource.ID, searchParams)
	if err != nil {
		return nil, err
	}

	return &dto.PatientEverythingConnection{
		TotalCount: patientData.TotalCount,
		Edges:      patientData.Resources,
		PageInfo: dto.PageInfo{
			HasNextPage:     patientData.HasNextPage,
			EndCursor:       &patientData.NextCursor,
			HasPreviousPage: patientData.HasPreviousPage,
			StartCursor:     &patientData.PreviousCursor,
			BundleID:        &patientData.BundleID,
		},
	}, nil
}
