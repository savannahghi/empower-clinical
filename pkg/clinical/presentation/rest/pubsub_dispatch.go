package rest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// HandlePubSubEvent routes a published event to the usecase that consumes it.
//
// It carries no transport: the Cloud Pub/Sub push endpoint and the NATS
// subscriber each decode their own envelope, then hand the topic and message
// body here. One switch keeps the two from drifting apart.
func (p PresentationHandlersImpl) HandlePubSubEvent(ctx context.Context, topicID string, payload []byte) error {
	// decode unmarshals the message body into the type a topic expects.
	decode := func(into interface{}) error {
		if err := json.Unmarshal(payload, into); err != nil {
			return fmt.Errorf("unable to unmarshal payload for topic %q: %w", topicID, err)
		}

		return nil
	}

	switch topicID {
	case utils.AddPubSubNamespace(common.CreatePatientTopic, common.ClinicalServiceName):
		var data dto.PatientPubSubMessage
		if err := decode(&data); err != nil {
			return err
		}

		return p.CreatePubsubPatient(ctx, data)

	case utils.AddPubSubNamespace(common.OrganizationTopicName, common.ClinicalServiceName):
		var data dto.FacilityPubSubMessage
		if err := decode(&data); err != nil {
			return err
		}

		return p.CreatePubsubOrganization(ctx, data)

	case utils.AddPubSubNamespace(common.VitalsTopicName, common.ClinicalServiceName):
		var data dto.VitalSignPubSubMessage
		if err := decode(&data); err != nil {
			return err
		}

		return p.CreatePubsubVitals(ctx, data)

	case utils.AddPubSubNamespace(common.AllergyTopicName, common.ClinicalServiceName):
		var data dto.PatientAllergyPubSubMessage
		if err := decode(&data); err != nil {
			return err
		}

		return p.CreatePubsubAllergyIntolerance(ctx, data)

	case utils.AddPubSubNamespace(common.MedicationTopicName, common.ClinicalServiceName):
		var data dto.MedicationPubSubMessage
		if err := decode(&data); err != nil {
			return err
		}

		return p.CreatePubsubMedicationStatement(ctx, data)

	case utils.AddPubSubNamespace(common.TenantTopicName, common.ClinicalServiceName):
		var data dto.OrganizationInput
		if err := decode(&data); err != nil {
			return err
		}

		return p.CreatePubsubTenant(ctx, data)

	case utils.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName):
		var data dto.PatientTestResultPubSubMessage
		if err := decode(&data); err != nil {
			return err
		}

		return p.CreatePubsubTestResult(ctx, data)

	case utils.AddPubSubNamespace(common.SegmentationTopicName, common.ClinicalServiceName):
		var data dto.SegmentationPayload
		if err := decode(&data); err != nil {
			return err
		}

		return p.SegmentPatient(ctx, data)

	case utils.AddPubSubNamespace(common.ReferralTopicName, common.ClinicalServiceName):
		var data dto.PatientReferralTaskPayload
		if err := decode(&data); err != nil {
			return err
		}

		_, err := p.CreateReferralTask(ctx, data.Meta, data.Referral)

		return err

	case utils.AddPubSubNamespace(common.ReferralReportNotificationTopic, common.ClinicalServiceName):
		// No consumer behaviour: sending the patient SMS and the facility email
		// was never implemented. Decoded so a malformed message still errors.
		var data dto.ReferralReportNotification

		return decode(&data)

	case utils.AddPubSubNamespace(common.CreateCarePlanTopic, common.ClinicalServiceName):
		var data dto.CarePlanPayload
		if err := decode(&data); err != nil {
			return err
		}

		_, err := p.PatientCarePlan(ctx, &data)

		return err

	case utils.AddPubSubNamespace(common.FollowUpTaskTopic, common.ClinicalServiceName):
		var data domain.FHIRTaskInput
		if err := decode(&data); err != nil {
			return err
		}

		_, err := p.CreateTask(ctx, &data)

		return err

	default:
		return fmt.Errorf("%w: %v", common.ErrUnknownTopic, topicID)
	}
}
