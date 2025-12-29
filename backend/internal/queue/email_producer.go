package queue

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishRegisterEmailTask(ctx context.Context, rabbitConn *amqp.Connection, registerToken string) error {

	data := WelcomeEmailData{
		RegisterToken: registerToken,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	queueEvent := QueueEvent{
		Task: registerEmailTask,
		Data: json.RawMessage(jsonData),
	}
	return publishEmailEvent(queueEvent, rabbitConn)
}

func PublishEmailChangeTask(ctx context.Context, rabbitConn *amqp.Connection, emailChangeToken string) error {

	data := EmailChangeEmailData{
		EmailChangeToken: emailChangeToken,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	queueEvent := QueueEvent{
		Task: emailChangeEmailTask,
		Data: json.RawMessage(jsonData),
	}
	return publishEmailEvent(queueEvent, rabbitConn)
}
func PublishPasswordResetTask(ctx context.Context, rabbitConn *amqp.Connection, passwordResetToken string) error {

	data := PasswordResetEmailData{
		PasswordResetToken: passwordResetToken,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	queueEvent := QueueEvent{
		Task: passwordResetEmailTask,
		Data: json.RawMessage(jsonData),
	}
	return publishEmailEvent(queueEvent, rabbitConn)
}

func publishEmailEvent(event QueueEvent, rabbitConn *amqp.Connection) error {
	body, _ := json.Marshal(event)
	ch, err := rabbitConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	setupEmailQueues(ch)

	return ch.Publish(
		"",             // exchange
		emailMainQueue, // routing key (nazwa kolejki)
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

