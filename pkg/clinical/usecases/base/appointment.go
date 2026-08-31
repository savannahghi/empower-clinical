package base

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// AppointmentPayload is used to model appointment data transfer object that is used to create a new appointment
type AppointmentPayload struct {
	Tags    *domain.FHIRMetaInput
	Subject *domain.FHIRReference
	Reason  string
}

func (c *BaseImpl) CreateAppointment(ctx context.Context, appointment *domain.FHIRAppointmentInput) (*domain.FHIRAppointment, error) {
	output, err := c.FHIR.CreateFHIRAppointment(ctx, appointment)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// CreateCheckIn creates a future date check-in for a patient in advantage EMR
func (c *BaseImpl) CreateCheckIn(
	ctx context.Context, patientID string, appointmentPayload *AppointmentPayload,
	date *scalarutils.Date, headers *dto.AdvantageHeaders) (*domain.FHIRAppointment, error) {
	schedule, err := c.AdvantageService.GetSchedules(ctx, headers)
	if err != nil {
		return nil, err
	}

	if len(schedule) < 1 {
		return nil, fmt.Errorf("no schedules found")
	}

	bookingDate := date.AsTime().Format(time.DateOnly)

	slots, err := c.AdvantageService.GetSlots(ctx, bookingDate, schedule[0].ID, headers)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	endDate := date.AsTime()
	endTime := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 0, 0, currentTime.Location()).Format("2006-01-02 15:04:05")

	checkin := &dto.Checkin{
		Slot:    slots[0].ID,
		Start:   bookingDate,
		End:     endTime,
		Patient: patientID,
	}

	err = c.AdvantageService.CreateCheckin(ctx, checkin, headers)
	if err != nil {
		return nil, err
	}

	bookedAppointment := dto.AppointmentStatusBooked
	patientRef := *appointmentPayload.Subject.Reference
	startInstant := date.AsTime().Format(timeFormatStr)
	endInstant := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 0, 0, currentTime.Location()).Format(timeFormatStr)

	concept, err := c.GetConcept(ctx, domain.TerminologySourceLOINC, common.ReferForMedicalConsultationLOINCTerminology)
	if err != nil {
		return nil, err
	}

	userSelected := false

	needsAction := domain.NeedsAction.Code()

	specialtyCodeSystem := helpers.CodeSystem(common.SpecialtyCodeSystem)
	actionPriorityCodeSystem := common.ActionPriorityCodeSystem
	appopintmentTypeCodeSystem := common.AppointmentTypeCodeSystem

	appointment := &domain.FHIRAppointmentInput{
		Meta:   appointmentPayload.Tags,
		Status: (*scalarutils.Code)(&bookedAppointment),
		AppointmentType: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					Code:    scalarutils.Code(domain.CheckUpAppointment),
					System:  (*scalarutils.URI)(&appopintmentTypeCodeSystem),
					Display: domain.CheckUpAppointment.Display(),
				},
			},
		},
		Reason: []*domain.FHIRCodeableReferenceInput{
			{
				Concept: &domain.FHIRCodeableConceptInput{
					Coding: []*domain.FHIRCodingInput{
						{
							System:       (*scalarutils.URI)(&concept.URL),
							Code:         scalarutils.Code(concept.ID),
							Display:      concept.GetConceptDisplay(),
							UserSelected: &userSelected,
						},
					},
					Text: concept.GetConceptDisplay(),
				},
			},
		},
		Start: (*scalarutils.Instant)(&startInstant),
		End:   (*scalarutils.Instant)(&endInstant),
		Participant: []*domain.FHIRAppointmentParticipant{
			{
				Actor: &domain.FHIRReference{
					ID:        &patientID,
					Reference: &patientRef,
					Display:   patientID,
				},
				Status: (*scalarutils.Code)(&needsAction),
			},
		},

		Note: []*domain.FHIRAnnotationInput{
			{
				Text: (*scalarutils.Markdown)(&appointmentPayload.Reason),
			},
		},
		Specialty: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  specialtyCodeSystem,
						Code:    scalarutils.Code(dto.Dentist),
						Display: dto.Dentist.Display(),
					},
				},
			},
		},
		Priority: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&actionPriorityCodeSystem),
					Code:    scalarutils.Code(domain.ASAP),
					Display: domain.ASAP.Display(),
				},
			},
		},
		Created: (*scalarutils.DateTime)(&bookingDate),
		Subject: &domain.FHIRReferenceInput{
			Reference: &patientRef,
			Display:   patientID,
		},
	}

	return c.CreateAppointment(ctx, appointment)
}

// ScheduleAppointment creates a check-in for a later date in advantage EMR
//
// ------------------ DEVELOPER NOTE! ------------------
//
// The field patientID in one of the methods parameter, i.e `input.patientID`, has to be the patient ID from advantage EMR.
// and not the regular FHIR patient ID from Google Cloud Healthcare.
//
// This is because when creating a patient check-in in advantage EMR, the EMR expects that the patient exists in their database.
//
// It is also to be noted that using the subject reference from encounter resource will NOT work due to the aforementioned reason.
func (c *BaseImpl) ScheduleAppointment(ctx context.Context, input *dto.ScheduleAppointmentInput, headers *dto.AdvantageHeaders) (bool, error) {
	if input.EncounterID == "" {
		return false, fmt.Errorf("encounter ID is required")
	}

	if input.Date == nil {
		return false, fmt.Errorf("appointment date is required")
	}

	if input.PatientID == "" {
		return false, fmt.Errorf("patient ID is required")
	}

	encounter, err := c.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return false, err
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return false, err
	}

	metaTags := &domain.FHIRMetaInput{
		Tag: tags,
	}

	appointmentPayload := &AppointmentPayload{
		Tags:    metaTags,
		Subject: encounter.Resource.Subject,
		Reason:  input.Reason,
	}

	appointmentOutput, err := c.CreateCheckIn(ctx, input.PatientID, appointmentPayload, input.Date, headers)
	if err != nil {
		return false, err
	}

	_, err = c.CreateAppointmentTask(ctx, *encounter.Resource.ID, metaTags, appointmentOutput)
	if err != nil {
		return false, err
	}

	return true, nil
}
