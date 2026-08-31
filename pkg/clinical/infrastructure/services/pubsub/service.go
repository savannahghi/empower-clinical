package pubsubmessaging

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub/v2"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

const (
	// HostNameEnvVarName defines the host name
	HostNameEnvVarName = "SERVICE_HOST"

	// TestTopicName is a topic that is used for testing purposes
	TestTopicName = "pubsub.testing.topic"

	TopicVersion = "v2"

	// MyCareHubServiceName defines the service where some of the topics have been created
	MyCareHubServiceName = "mycarehub"
)

type ServicePubsub interface {
	NotifyPatientFHIRIDUpdate(ctx context.Context, data dto.UpdatePatientFHIRID) error
	NotifyFacilityFHIRIDUpdate(ctx context.Context, data dto.UpdateFacilityFHIRID) error
	NotifyProgramFHIRIDUpdate(ctx context.Context, data dto.UpdateProgramFHIRID) error
	NotifySegmentation(ctx context.Context, data dto.SegmentationPayload) error
	NotifyCreatePatientReferralTask(ctx context.Context, data dto.PatientReferralTaskPayload) error
	NotifyCreatePatientCarePlan(ctx context.Context, data dto.CarePlanPayload) error

	// NotifyCreateFollowUpTask is used to record any kind of task. `FHIRTaskInput` domain model has bee used here to give room for creating any kind of task
	// There will be no need to define different method signatures to create any task
	NotifyCreateFollowUpTask(ctx context.Context, data *domain.FHIRTaskInput) error
}

// ServicePubSubMessaging is used to send and receive pubsub notifications
type ServicePubSubMessaging struct {
	client *pubsub.Client
}

// NewServicePubSubMessaging returns a new instance of pubsub
func NewServicePubSubMessaging(
	ctx context.Context,
	client *pubsub.Client,
) (*ServicePubSubMessaging, error) {
	s := &ServicePubSubMessaging{
		client: client,
	}

	if err := s.EnsureTopicsExist(
		ctx,
		s.TopicIDs(),
	); err != nil {
		return nil, err
	}

	if err := s.EnsureSubscriptionsExist(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

// AddPubSubNamespace creates unique topics. The topics will be in the form
// <service name>-<topicName>-<environment>-v1
func (ps ServicePubSubMessaging) AddPubSubNamespace(topicName, serviceName string) string {
	environment := serverutils.GetRunningEnvironment()

	return pubsubtools.NamespacePubsubIdentifier(
		serviceName,
		topicName,
		environment,
		TopicVersion,
	)
}

// TopicIDs returns the known (registered) topic IDs
func (ps ServicePubSubMessaging) TopicIDs() []string {
	return []string{
		ps.AddPubSubNamespace(TestTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.CreatePatientTopic, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.VitalsTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.MedicationTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.AllergyTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.TestOrderTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.OrganizationTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.TenantTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.SegmentationTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.ReferralTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.CreateCarePlanTopic, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.FollowUpTaskTopic, common.ClinicalServiceName),
	}
}

// PublishToPubsub publishes a message to a specified topic
func (ps ServicePubSubMessaging) PublishToPubsub(
	ctx context.Context,
	topicID, serviceName string,
	payload []byte,
) error {
	return pubsubtools.PublishToPubsub(
		ctx,
		ps.client,
		topicID,
		payload,
	)
}

// EnsureTopicsExist creates the topic(s) in the supplied list if they do not
// exist
func (ps ServicePubSubMessaging) EnsureTopicsExist(
	ctx context.Context,
	topicIDs []string,
) error {
	return pubsubtools.EnsureTopicsExist(
		ctx,
		ps.client,
		topicIDs,
	)
}

// EnsureSubscriptionsExist ensures that the subscriptions named in the supplied
// topic:subscription map exist. If any does not exist, it is created.
func (ps ServicePubSubMessaging) EnsureSubscriptionsExist(
	ctx context.Context,
) error {
	hostName, err := serverutils.GetEnvVar(HostNameEnvVarName)
	if err != nil {
		return err
	}

	callbackURL := fmt.Sprintf(
		"%s%s",
		hostName,
		pubsubtools.PubSubHandlerPath,
	)

	return pubsubtools.EnsureSubscriptionsExist(
		ctx,
		ps.client,
		ps.SubscriptionIDs(),
		callbackURL,
	)
}

// SubscriptionIDs returns a map of topic IDs to subscription IDs
func (ps ServicePubSubMessaging) SubscriptionIDs() map[string]string {
	return pubsubtools.SubscriptionIDs(ps.TopicIDs())
}
