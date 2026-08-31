package pubsubmessaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

func (ps ServicePubSubMessaging) newPublish(
	ctx context.Context,
	data interface{},
	topic, serviceName string,
) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data received: %w", err)
	}

	return ps.PublishToPubsub(
		ctx,
		ps.AddPubSubNamespace(topic, serviceName),
		serviceName,
		payload,
	)
}

// NotifyPatientFHIRIDUpdate publishes to patient fhir id update topic. Mycarehub service will subscribe to this topic
// and update the patient's FHIR ID in it's database
func (ps ServicePubSubMessaging) NotifyPatientFHIRIDUpdate(ctx context.Context, data dto.UpdatePatientFHIRID) error {
	return ps.newPublish(ctx, data, common.AddFHIRIDToPatientProfile, MyCareHubServiceName)
}

// NotifyFacilityFHIRIDUpdate publishes to a topic. The idea is that after a mycarehub facility is created as an organization in FHIR,
// we should send back the ID to mycarehub and store in the database
func (ps ServicePubSubMessaging) NotifyFacilityFHIRIDUpdate(ctx context.Context, data dto.UpdateFacilityFHIRID) error {
	return ps.newPublish(ctx, data, common.AddFHIRIDToFacility, MyCareHubServiceName)
}

// NotifyProgramFHIRIDUpdate publishes to the program fhir id update topic
func (ps ServicePubSubMessaging) NotifyProgramFHIRIDUpdate(ctx context.Context, data dto.UpdateProgramFHIRID) error {
	return ps.newPublish(ctx, data, common.AddFHIRIDToProgram, MyCareHubServiceName)
}

// NotifySegmentation publishes the the data used to update the segmentation data in advantage
func (ps ServicePubSubMessaging) NotifySegmentation(ctx context.Context, data dto.SegmentationPayload) error {
	return ps.newPublish(ctx, data, common.SegmentationTopicName, common.ClinicalServiceName)
}

// NotifyCreatePatientReferralTask publishes the data used to create a patient referral task
func (ps ServicePubSubMessaging) NotifyCreatePatientReferralTask(ctx context.Context, data dto.PatientReferralTaskPayload) error {
	return ps.newPublish(ctx, data, common.ReferralTopicName, common.ClinicalServiceName)
}

// SendReferralReportNotification publishes the message that is used to send a referral report notification to the patient via SMS and
// the receiving facility via email
func (ps ServicePubSubMessaging) SendReferralReportNotification(ctx context.Context, data dto.ReferralReportNotification) error {
	return ps.newPublish(ctx, data, common.ReferralReportNotificationTopic, common.ClinicalServiceName)
}

func (ps ServicePubSubMessaging) NotifyCreatePatientCarePlan(ctx context.Context, data dto.CarePlanPayload) error {
	return ps.newPublish(ctx, data, common.CreateCarePlanTopic, common.ClinicalServiceName)
}

func (ps ServicePubSubMessaging) NotifyCreateFollowUpTask(ctx context.Context, data *domain.FHIRTaskInput) error {
	return ps.newPublish(ctx, data, common.FollowUpTaskTopic, common.ClinicalServiceName)
}
