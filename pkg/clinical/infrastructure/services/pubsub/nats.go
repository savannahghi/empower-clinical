package pubsubmessaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

const (
	// StreamName holds every event published by this service.
	StreamName = "CLINICAL"

	// streamMaxAge bounds retention of unconsumed events.
	streamMaxAge = 72 * time.Hour

	// HeaderTopicID carries the namespaced topic identifier, mirroring the
	// topicID attribute on a Cloud Pub/Sub message.
	HeaderTopicID = "Topic-ID"
)

// NATSPubSub publishes clinical integration events over NATS JetStream, so the
// service can run without Google Cloud.
//
// Topic names map onto subjects prefixed by the publishing service:
// patient.referral.task.create becomes clinical.patient.referral.task.create.
//
// JetStream rather than core NATS because five of these topics are consumed by
// this same service and carry clinical work. Core NATS drops a message when no
// subscriber is connected.
type NATSPubSub struct {
	conn   *nats.Conn
	js     jetstream.JetStream
	prefix string
}

// NewNATSPubSub connects to NATS, creates the stream if absent and returns a
// publisher.
func NewNATSPubSub(ctx context.Context, url, prefix string) (*NATSPubSub, error) {
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

	js, err := jetstream.New(conn)
	if err != nil {
		conn.Close()

		return nil, fmt.Errorf("unable to initialise JetStream: %w", err)
	}

	// Created here so a fresh checkout runs with nothing but docker compose up.
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        StreamName,
		Description: "Integration events published by the clinical service",
		Subjects:    []string{prefix + ".>"},
		Storage:     jetstream.FileStorage,
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      streamMaxAge,
	})
	if err != nil {
		conn.Close()

		return nil, fmt.Errorf("unable to create the %q stream: %w", StreamName, err)
	}

	return &NATSPubSub{conn: conn, js: js, prefix: prefix}, nil
}

// publish marshals the payload and sends it to the subject for this topic.
func (n *NATSPubSub) publish(ctx context.Context, data interface{}, topic string) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data for %q: %w", topic, err)
	}

	msg := &nats.Msg{
		Subject: fmt.Sprintf("%s.%s", n.prefix, topic),
		Data:    payload,
		Header:  nats.Header{},
	}

	msg.Header.Set(HeaderTopicID, utils.AddPubSubNamespace(topic, common.ClinicalServiceName))

	// PublishMsg returns once the server has persisted the message.
	if _, err := n.js.PublishMsg(ctx, msg); err != nil {
		return fmt.Errorf("unable to publish to %q: %w", msg.Subject, err)
	}

	return nil
}

// Close drains the connection so in-flight messages are delivered.
func (n *NATSPubSub) Close() {
	if n.conn != nil {
		if err := n.conn.Drain(); err != nil {
			n.conn.Close()
		}
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
