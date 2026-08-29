package pubsubmessaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// NATSPubSub publishes clinical integration events over NATS.
//
// It exists so the service can emit its events without Google Cloud. The topic
// names are already dot-delimited and map directly onto NATS subjects, prefixed
// by the publishing service: patient.referral.task.create becomes
// clinical.patient.referral.task.create.
//
// Note that NATS core delivers at-most-once and drops a message when nothing is
// subscribed. That is the right behaviour here — these events are consumed by
// services outside this release, so an operator who has no subscriber sees
// successful publishes and no queue growth. An operator who needs delivery
// guarantees should subscribe, or move this to JetStream.
type NATSPubSub struct {
	conn   *nats.Conn
	prefix string
}

// NewNATSPubSub connects to NATS and returns a publisher.
func NewNATSPubSub(url, prefix string) (*NATSPubSub, error) {
	if url == "" {
		url = nats.DefaultURL
	}

	if prefix == "" {
		prefix = common.ClinicalServiceName
	}

	conn, err := nats.Connect(url,
		nats.Name("empower-clinical"),
		nats.Timeout(5*time.Second),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to NATS at %q: %w", url, err)
	}

	return &NATSPubSub{conn: conn, prefix: prefix}, nil
}

// publish marshals the payload and sends it to the subject for this topic.
func (n *NATSPubSub) publish(_ context.Context, data interface{}, topic string) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data for %q: %w", topic, err)
	}

	subject := fmt.Sprintf("%s.%s", n.prefix, topic)

	if err := n.conn.Publish(subject, payload); err != nil {
		return fmt.Errorf("unable to publish to %q: %w", subject, err)
	}

	// Publishing is asynchronous; flush so a caller that gets no error can rely
	// on the message having reached the server.
	return n.conn.FlushTimeout(5 * time.Second)
}

// Close releases the connection.
func (n *NATSPubSub) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}

// NotifyPatientFHIRIDUpdate publishes a patient's newly assigned FHIR ID.
func (n *NATSPubSub) NotifyPatientFHIRIDUpdate(ctx context.Context, data dto.UpdatePatientFHIRID) error {
	return n.publish(ctx, data, common.AddFHIRIDToPatientProfile)
}

// NotifyFacilityFHIRIDUpdate publishes a facility's newly assigned FHIR ID.
func (n *NATSPubSub) NotifyFacilityFHIRIDUpdate(ctx context.Context, data dto.UpdateFacilityFHIRID) error {
	return n.publish(ctx, data, common.AddFHIRIDToFacility)
}

// NotifyProgramFHIRIDUpdate publishes a program's newly assigned FHIR ID.
func (n *NATSPubSub) NotifyProgramFHIRIDUpdate(ctx context.Context, data dto.UpdateProgramFHIRID) error {
	return n.publish(ctx, data, common.AddFHIRIDToProgram)
}

// NotifySegmentation publishes patient segmentation information.
func (n *NATSPubSub) NotifySegmentation(ctx context.Context, data dto.SegmentationPayload) error {
	return n.publish(ctx, data, common.SegmentationTopicName)
}

// NotifyCreatePatientReferralTask publishes the data used to create a referral task.
func (n *NATSPubSub) NotifyCreatePatientReferralTask(ctx context.Context, data dto.PatientReferralTaskPayload) error {
	return n.publish(ctx, data, common.ReferralTopicName)
}

// NotifyCreatePatientCarePlan publishes a patient's care plan.
func (n *NATSPubSub) NotifyCreatePatientCarePlan(ctx context.Context, data dto.CarePlanPayload) error {
	return n.publish(ctx, data, common.CreateCarePlanTopic)
}

// NotifyCreateFollowUpTask publishes a follow-up task.
func (n *NATSPubSub) NotifyCreateFollowUpTask(ctx context.Context, data *domain.FHIRTaskInput) error {
	return n.publish(ctx, data, common.FollowUpTaskTopic)
}
