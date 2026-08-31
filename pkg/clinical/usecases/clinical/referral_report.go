package clinical

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/urlshortener"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/silurlshortener"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// gotenbergURL is the base URL of the Gotenberg service used to render PDFs.
func gotenbergURL() string {
	if url := os.Getenv("GOTENBERG_URL"); url != "" {
		return strings.TrimSuffix(url, "/")
	}

	return "http://localhost:3000"
}

var referralFormTemplate = template.Must(template.New("ReferralFormTemplate").Parse(utils.ReferralFormTemplate))

// DocumentReferencePayload models data used to create caller specific reference document payload
type DocumentReferencePayload struct {
	Subject    *domain.FHIRReferenceInput
	Encounter  *domain.FHIRReference
	Attachment *domain.FHIRAttachment
	Related    *domain.FHIRReference
	// Procedure that caused this media to be created
	BasedOn           []*domain.FHIRReferenceInput
	TerminologySystem string
	Custodian         *domain.FHIRReference
}

// GenerateReferralReportPDF generates a PDF report for a given referral.
//
// The serviceRequestID is unique to each ServiceRequest resource, which,
// according to FHIR standards, is how referrals are represented. In FHIR, a referral
// is a specific type of ServiceRequest, which typically contains details such
// as the requester, the patient, the requested service, and other clinical information.
//
// By leveraging the serviceRequestID, this function retrieves the associated ServiceRequest
// from the FHIR server. It then extracts relevant data, including patient and encounter
// information, to construct a comprehensive referral report. The report is formatted
// as a PDF, making it suitable for clinical review, record-keeping, or sharing with
// other healthcare professionals.
func (c *ClinicalImpl) GenerateReferralReportPDF(ctx context.Context, serviceRequestID string) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "GenerateReferralReportPDF")

	defer func() { span.End() }()

	if serviceRequestID == "" {
		return nil, fmt.Errorf("service request ID cannot be empty")
	}

	serviceRequest, err := c.FHIR.GetFHIRServiceRequest(ctx, serviceRequestID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	var referralNote string
	if len(serviceRequest.Resource.Note) > 0 && serviceRequest.Resource.Note[0].Text != nil {
		referralNote = string(*serviceRequest.Resource.Note[0].Text)
	}

	var referralTest string

	if serviceRequest.Resource.Code != nil {
		for _, code := range serviceRequest.Resource.Code.Concept.Coding {
			if *code.Code == scalarutils.Code("TEST") {
				referralTest = code.Display
			}
		}
	}

	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch patient, facility, receiving facility, and observations in parallel
	// since they all only depend on the service request data.
	var (
		patient           *domain.FHIRPatientRelayPayload
		facility          *domain.FHIROrganizationRelayPayload
		receivingFacility *domain.FHIROrganizationRelayPayload
		observations      *dto.ObservationConnection
	)

	var wg sync.WaitGroup

	errs := make(chan error, 4)

	wg.Add(4)

	go func() {
		defer wg.Done()

		var fetchErr error

		patient, fetchErr = c.FHIR.GetFHIRPatient(ctx, *serviceRequest.Resource.Subject.ID)
		if fetchErr != nil {
			errs <- fetchErr
		}
	}()

	go func() {
		defer wg.Done()

		var fetchErr error

		facility, fetchErr = c.FHIR.GetFHIROrganization(ctx, facilityID)
		if fetchErr != nil {
			errs <- fetchErr
		}
	}()

	go func() {
		defer wg.Done()

		var fetchErr error

		receivingFacility, fetchErr = c.FHIR.GetFHIROrganization(ctx, serviceRequest.Resource.GetPerformerID())
		if fetchErr != nil {
			errs <- fetchErr
		}
	}()

	go func() {
		defer wg.Done()

		var fetchErr error

		observations, fetchErr = c.GetPatientObservations(ctx, &dto.FetchObservationPayload{
			PatientID:   *serviceRequest.Resource.Subject.ID,
			EncounterID: serviceRequest.Resource.Encounter.ID,
			Pagination:  &dto.Pagination{},
		})
		if fetchErr != nil {
			errs <- fetchErr
		}
	}()

	wg.Wait()
	close(errs)

	for fetchErr := range errs {
		utils.ReportErrorToSentry(fetchErr)
		return nil, fetchErr
	}

	patientIdentifier, err := patient.Resource.GetIDs()
	if err != nil {
		return nil, err
	}

	age := time.Since(patient.Resource.BirthDate.AsTime()).Hours() / 24 / 365

	patientData := domain.Patient{
		Name:        patient.Resource.Name[0].Text,
		EmpowerID:   patientIdentifier.HealthID,
		NationalID:  patientIdentifier.NationalID,
		PhoneNumber: *patient.Resource.Telecom[0].Value,
		DateOfBirth: patient.Resource.BirthDate.String(),
		Age:         int(age),
		Sex:         cases.Title(language.English).String(patient.Resource.Gender.String()),
	}

	var facilityContact string

	facilityContacts := facility.Resource.GetOrganizationContacts()
	if len(facilityContacts) > 0 {
		facilityContact = facilityContacts[0]
	}

	var tests []domain.Test

	for _, observation := range observations.Edges {
		date, err := utils.ConvertDateStringToDateScalar(time.RFC3339, observation.Node.TimeRecorded)
		if err != nil {
			return nil, err
		}

		tests = append(tests, domain.Test{
			Name:    observation.Node.Name,
			Results: observation.Node.Value,
			Date:    date.String(),
		})
	}

	receivingFacilityContacts, err := receivingFacility.Resource.GetReceivingFacilityContactDetails()
	if err != nil {
		return nil, err
	}

	data := domain.ReferralReport{
		Facility:  *facility.Resource.Name,
		Patient:   patientData,
		NextOfKin: domain.NextOfKin{},
		ReceivingFacility: domain.ReceivingFacility{
			FacilityName:    *receivingFacility.Resource.Name,
			FacilityCounty:  "Nairobi",
			FacilityContact: receivingFacilityContacts.Phone,
			FacilityEmail:   receivingFacilityContacts.Email,
		},
		ReferralReason: domain.ReferralReason{
			Reason: "",
			Test:   referralTest,
			Date:   time.Now().Format("Jan 2, 2006"),
			Note:   referralNote,
		},
		MedicalHistory: domain.MedicalHistory{Procedure: "Screening", Medication: "None", Tests: tests},
		Footer: domain.Footer{
			Phone: facilityContact,
		},
	}

	var htmlBuffer bytes.Buffer

	err = referralFormTemplate.Execute(&htmlBuffer, data)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	htmlContent := htmlBuffer.String()

	url := gotenbergURL() + "/forms/chromium/convert/html"

	var b bytes.Buffer

	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("files", "index.html")
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	_, err = io.Copy(fw, bytes.NewReader([]byte(htmlContent)))
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	w.Close()

	// nosemgrep: problem-based-packs.insecure-transport.go-stdlib.http-customized-request.http-customized-request
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	defer resp.Body.Close()

	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	currentTime := time.Now().Format("20060102T150405")
	filename := fmt.Sprintf("%s_%s.pdf", patientData.Name, currentTime)

	result, err := c.Upload.UploadMedia(ctx, filename, bytes.NewReader(pdfBytes), "application/pdf")
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	title := fmt.Sprintf("%s's Referral report", serviceRequest.Resource.Subject.Display)
	serviceRequestRef := fmt.Sprintf("ServiceRequest/%s", serviceRequestID)

	urlShortenerPayload := &silurlshortener.ShortenURLPayload{
		LongURL:         result.SignedURL,
		Tags:            []string{"clinical-service"},
		Domain:          os.Getenv(urlshortener.DomainEnvVarName),
		ShortCodeLength: 5,
	}

	signedURL, err := c.URLShorteningService.Shorten(ctx, urlShortenerPayload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	payload := &DocumentReferencePayload{
		Subject: &domain.FHIRReferenceInput{
			ID:        serviceRequest.Resource.Subject.ID,
			Reference: serviceRequest.Resource.Subject.Reference,
			Display:   serviceRequest.Resource.Subject.Display,
		},
		Attachment: &domain.FHIRAttachment{
			ContentType: (*scalarutils.Code)(&result.ContentType),
			URL:         (*scalarutils.URL)(&signedURL.ShortURL),
			Title:       &title,
		},
		BasedOn: []*domain.FHIRReferenceInput{
			{
				Reference: &serviceRequestRef,
				Display:   serviceRequestID,
			},
		},
		TerminologySystem: common.ReferralLOINCTerminologySystem,
		Encounter: &domain.FHIRReference{
			ID:        serviceRequest.Resource.Encounter.ID,
			Reference: serviceRequest.Resource.Encounter.Reference,
			Display:   serviceRequest.Resource.Subject.Display,
		},
	}

	err = c.CreateDocumentReference(ctx, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return pdfBytes, nil
}

// CreateDocumentReference is a helper method to abstract the creation of a document reference
func (c *ClinicalImpl) CreateDocumentReference(ctx context.Context, payload *DocumentReferencePayload) error {
	concept, err := c.GetConcept(ctx, domain.TerminologySourceLOINC, payload.TerminologySystem)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return err
	}

	finalDocStatus := domain.CompositionStatusEnumFinal
	status := domain.DocumentReferenceStatusEnumCurrent
	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	organizationID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return err
	}

	organizationRef := fmt.Sprintf("Organization/%s", organizationID)

	documentReference := &domain.FHIRDocumentReferenceInput{
		Author: []*domain.FHIRReferenceInput{
			{
				Reference: &organizationRef,
				Display:   organizationID,
			},
		},
		Status:    status,
		DocStatus: &finalDocStatus,
		BasedOn:   payload.BasedOn,
		Subject:   payload.Subject,
		Date:      &instant,
		Content: []domain.FHIRDocumentReferenceContent{
			{
				Attachment: *payload.Attachment,
			},
		},
		Context: []*domain.FHIRReferenceInput{
			{
				Reference: payload.Encounter.Reference,
				Display:   payload.Encounter.Display,
			},
		},
		Custodian: &domain.FHIRReferenceInput{
			Reference: &organizationRef,
			Display:   organizationID,
		},
	}

	if payload.TerminologySystem != "" {
		concept, err := c.GetConcept(ctx, domain.TerminologySourceLOINC, payload.TerminologySystem)
		if err != nil {
			utils.ReportErrorToSentry(err)
			return err
		}

		documentReference.Type = &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&concept.URL),
					Code:    scalarutils.Code(concept.ID),
					Display: concept.GetConceptDisplay(),
				},
			},
			Text: concept.GetConceptDisplay(),
		}
	}

	if concept != nil {
		documentReference.Type = &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  &common.LoincSystemURL,
					Code:    scalarutils.Code(concept.ID),
					Display: concept.GetConceptDisplay(),
				},
			},
			Text: concept.GetConceptDisplay(),
		}
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return err
	}

	documentReference.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	narrativeStatus := domain.NarrativeStatusEnumGenerated.String()
	narrative := utils.NarrativeGenerator("Generated text", &narrativeStatus)

	documentReference.Text = (*domain.FHIRNarrativeInput)(narrative)

	_, err = c.FHIR.CreateFHIRDocumentReference(ctx, documentReference)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return err
	}

	return nil
}

// ShareReferralForm is searches for a document reference using the service request ID associated and retrieves the document URL
// (which is the url of the referral form that we want to share) and sends it to the patient via SMS
func (c *ClinicalImpl) ShareReferralForm(ctx context.Context, input *dto.ShareReferralFormInput) (bool, error) {
	identifiers, err := c.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return false, err
	}

	params := map[string]interface{}{
		"relatesto": fmt.Sprintf("ServiceRequest/%s", input.ServiceRequestID),
		"_sort":     "_lastUpdated",
		"_count":    "1",
		"_total":    "accurate",
	}

	output, err := c.FHIR.SearchFHIRDocumentReference(ctx, params, *identifiers, dto.Pagination{})
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, err
	}

	if len(output.DocumentReferences) == 0 {
		return false, errors.New("no document reference found")
	}

	documentReference := output.DocumentReferences[0]

	patientID, err := domain.GetPatientID(&documentReference)
	if err != nil {
		return false, err
	}

	patient, err := c.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, err
	}

	attachment, err := documentReference.GetDocumentAttachment()
	if err != nil {
		return false, err
	}

	referralRequest, err := c.FHIR.GetFHIRServiceRequest(ctx, input.ServiceRequestID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, err
	}

	receivingFacility, err := referralRequest.Resource.GetReceivingFacilityDetails()
	if err != nil {
		return false, err
	}

	var wg sync.WaitGroup

	errorChan := make(chan error, 2)

	wg.Add(1)

	go func() {
		defer wg.Done()

		err := c.sendReferralSMS(ctx, patient, string(*attachment.URL), input.WorkstationID, input.BranchID)
		if err != nil {
			errorChan <- err
		}
	}()

	if receivingFacility.FacilityEmail != "" {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := c.sendReferralEmail(ctx, patient, string(*attachment.URL), receivingFacility)
			if err != nil {
				errorChan <- err
			}
		}()
	}

	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// sendReferralSMS sends an SMS to the patient with the link to the referral report
func (c *ClinicalImpl) sendReferralSMS(ctx context.Context, patient *domain.FHIRPatientRelayPayload, mediaURL, workstationID, branchID string) error {
	payload := &silurlshortener.ShortenURLPayload{
		LongURL:         mediaURL,
		Tags:            []string{"clinical-service"},
		Domain:          os.Getenv(urlshortener.DomainEnvVarName),
		ShortCodeLength: 5,
	}

	response, err := c.URLShorteningService.Shorten(ctx, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return err
	}

	patientContact, err := patient.Resource.GetPhoneNumbers()
	if err != nil {
		return err
	}

	textMessage := fmt.Sprintf("Hi %s. Click this link to download your referral form %s", *patient.Resource.Name[0].Given[0], response.ShortURL)

	smsPayload := &dto.SMSPayload{
		Intention:  "DIRECT_MESSAGE",
		Message:    textMessage,
		Recipients: patientContact,
	}

	err = c.AdvantageService.SendSMS(ctx, workstationID, branchID, *smsPayload)
	if err != nil {
		return err
	}

	return nil
}

// sendReferralEmail sends an email to the receiving facility regarding the patient being referred. If no facility email is provided, it does not send one
// Sending a text message to the facility can be an option
func (c *ClinicalImpl) sendReferralEmail(ctx context.Context, patient *domain.FHIRPatientRelayPayload, mediaURL string, receivingFacility *domain.ReceivingFacility) error {
	pdfBytes, err := fetchFileFromURL(mediaURL)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return err
	}

	patientPhone, err := patient.Resource.GetPhoneNumbers()
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("New Referral Request from %v", receivingFacility.FacilityName)

	tmpl, err := template.New("email").Parse(utils.ReferralEmailTemplate)
	if err != nil {
		return err
	}

	var emailHTMLBody bytes.Buffer

	err = tmpl.Execute(&emailHTMLBody, map[string]string{
		"PatientName":        patient.Resource.Names(),
		"PatientPhoneNumber": patientPhone[0],
	})
	if err != nil {
		return err
	}

	attachmentFileName := fmt.Sprintf("%v.pdf", patient.Resource.Names())

	destinations := []string{receivingFacility.FacilityEmail}

	rawEmail, err := createRawEmail(
		serverutils.MustGetEnvVar("DEFAULT_FROM_EMAIL"),
		destinations,
		subject,
		emailHTMLBody.String(),
		pdfBytes,
		attachmentFileName,
	)
	if err != nil {
		return err
	}

	_, err = c.Mail.SendRawEmail(ctx, &ses.SendRawEmailInput{
		RawMessage: &types.RawMessage{
			Data: rawEmail,
		},
		Destinations: destinations,
		Source:       aws.String(serverutils.MustGetEnvVar("DEFAULT_FROM_EMAIL")),
	})
	if err != nil {
		return err
	}

	return nil
}

// fetchFileFromURL retrieves the content of a file from the specified URL.
// In this context, it is used to fetch a referral report stored in GCS.
func fetchFileFromURL(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file: %s", resp.Status)
	}

	attachment, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

// createRawEmail constructs a raw MIME email message with the given sender, recipient, subject, body text,
// and a PDF attachment. It will be used to compose the email that will be passed to Amazon SES SendRawEmail api
func createRawEmail(sender string, recipients []string, subject, body string, attachment []byte, attachmentFilename string) ([]byte, error) {
	var emailRaw bytes.Buffer

	writer := multipart.NewWriter(&emailRaw)

	headers := make(map[string]string)
	headers["From"] = sender
	headers["To"] = strings.Join(recipients, ", ")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = `multipart/mixed; boundary="` + writer.Boundary() + `"`

	for k, v := range headers {
		emailRaw.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	emailRaw.WriteString("\r\n")

	bodyPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	if err != nil {
		return nil, err
	}

	_, err = bodyPart.Write([]byte(body))
	if err != nil {
		return nil, err
	}

	attachmentPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"application/pdf; name=\"" + attachmentFilename + "\""},
		"Content-Disposition":       {"attachment; filename=\"" + attachmentFilename + "\""},
		"Content-Transfer-Encoding": {"base64"},
	})
	if err != nil {
		return nil, err
	}

	base64Attachment := make([]byte, base64.StdEncoding.EncodedLen(len(attachment)))
	base64.StdEncoding.Encode(base64Attachment, attachment)

	_, err = attachmentPart.Write(base64Attachment)
	if err != nil {
		return nil, err
	}

	writer.Close()

	return emailRaw.Bytes(), nil
}
