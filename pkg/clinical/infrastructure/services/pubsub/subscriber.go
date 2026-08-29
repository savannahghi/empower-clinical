package pubsubmessaging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
)

const (
	// consumerName is the durable consumer this service reads with. Durability
	// means events published while it was down are delivered on restart.
	consumerName = "clinical-worker"

	// handlerTimeout bounds processing of a single event. Shorter than ackWait
	// so a stuck handler releases the message for redelivery.
	handlerTimeout = 45 * time.Second

	// ackWait is how long JetStream waits for an acknowledgement before
	// redelivering.
	ackWait = 60 * time.Second

	// maxDeliver caps redelivery attempts.
	maxDeliver = 5
)

// EventHandler routes an event to the usecase that consumes it. It is satisfied
// by HandlePubSubEvent, which the Cloud Pub/Sub push endpoint also uses.
type EventHandler func(ctx context.Context, topicID string, payload []byte) error

// Subscriber consumes events until stopped.
type Subscriber struct {
	consume jetstream.ConsumeContext
}

// Stop ends consumption. In-flight handlers are left to finish.
func (s *Subscriber) Stop() {
	if s.consume != nil {
		s.consume.Stop()
	}
}

// Subscribe consumes every subject in the stream, dispatching each event through
// handle.
//
// One wildcard subscription covers all topics, so a new topic needs only a case
// in the dispatch switch. Events meant for other services arrive here too and
// are acknowledged and ignored.
func (n *NATSPubSub) Subscribe(ctx context.Context, handle EventHandler) (*Subscriber, error) {
	stream, err := n.js.Stream(ctx, StreamName)
	if err != nil {
		return nil, fmt.Errorf("unable to look up the %q stream: %w", StreamName, err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       consumerName,
		Description:   "Clinical's consumers of the events it publishes",
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       ackWait,
		MaxDeliver:    maxDeliver,
		DeliverPolicy: jetstream.DeliverAllPolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create the %q consumer: %w", consumerName, err)
	}

	consume, err := consumer.Consume(func(msg jetstream.Msg) {
		n.dispatch(ctx, msg, handle)
	})
	if err != nil {
		return nil, fmt.Errorf("unable to start consuming: %w", err)
	}

	slog.Info("clinical event subscriber started",
		"stream", StreamName,
		"consumer", consumerName,
		"subjects", n.prefix+".>",
	)

	return &Subscriber{consume: consume}, nil
}

// dispatch handles a single message and decides its acknowledgement.
func (n *NATSPubSub) dispatch(ctx context.Context, msg jetstream.Msg, handle EventHandler) {
	topicID := msg.Headers().Get(HeaderTopicID)
	if topicID == "" {
		// Published by another producer, or before this header existed.
		topicID = n.topicIDFromSubject(msg.Subject())
	}

	msgCtx, cancel := context.WithTimeout(ctx, handlerTimeout)
	defer cancel()

	err := handle(msgCtx, topicID, msg.Data())

	switch {
	case err == nil:
		ack(msg)

	case errors.Is(err, common.ErrUnknownTopic):
		// Traffic for another service on the same wildcard subscription.
		slog.Debug("ignoring event this service does not consume",
			"topic", topicID,
			"subject", msg.Subject(),
		)
		ack(msg)

	default:
		slog.Error("failed to handle event",
			"topic", topicID,
			"subject", msg.Subject(),
			"error", err,
		)

		if nakErr := msg.Nak(); nakErr != nil {
			slog.Error("failed to negatively acknowledge event", "error", nakErr)
		}
	}
}

// ack acknowledges a message, logging if it cannot.
func ack(msg jetstream.Msg) {
	if err := msg.Ack(); err != nil {
		slog.Error("failed to acknowledge event", "subject", msg.Subject(), "error", err)
	}
}

// topicIDFromSubject recovers the namespaced topic identifier from a subject.
func (n *NATSPubSub) topicIDFromSubject(subject string) string {
	topic := strings.TrimPrefix(subject, n.prefix+".")

	return utils.AddPubSubNamespace(topic, common.ClinicalServiceName)
}
