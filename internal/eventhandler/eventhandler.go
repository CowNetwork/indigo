package eventhandler

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	pb "github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var (
	client    cloudevents.Client
	sourceUri string
)

func Initialize(brokers []string, topic string, source string) {
	sourceUri = source

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0

	sender, err := kafka_sarama.NewSender(brokers, saramaConfig, topic)
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}
	defer sender.Close(context.Background())

	c, err := cloudevents.NewClient(sender, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	client = c
}

func SendEvent(etype string, message proto.Message) error {
	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetSource(sourceUri)
	event.SetType(etype)

	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	err = event.SetData(cloudevents.ApplicationJSON, data)
	if err != nil {
		return err
	}

	if result := client.Send(
		// Set the producer message key
		kafka_sarama.WithMessageKey(context.Background(), sarama.StringEncoder(event.ID())),
		event,
	); cloudevents.IsUndelivered(result) {
		return fmt.Errorf("failed to send: %v", result)
	}
	return nil
}

func SendRoleUpdateEvent(role *pb.Role, action pb.RoleUpdateEvent_Action) {
	err := SendEvent("cow.indigo.v1.RoleUpdateEvent", &pb.RoleUpdateEvent{
		Role:   role,
		Action: action,
	})
	if err != nil {
		log.Printf("Could not send cloudevent: %v", err)
	}
}
