package main

import (
	"fmt"
	"log"
	"time"

	e "github.com/nombiezinja/chstub/entities"
	u "github.com/nombiezinja/chstub/utils"
	"github.com/streadway/amqp"
)

var outgoing chan string = make(chan string)

// var c chan string = make(chan string)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	u.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	u.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"stuffs", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	u.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	u.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go enqueuer(outgoing)

	go func() {
		i := 0
		for d := range msgs {
			fmt.Println("Beginning", time.Now().UnixNano())
			d.Ack(true)
			data := e.GojayUnmarshal(d.Body)

			returnResult(data)
			fmt.Printf("End of # %v: %v \n", i, time.Now().UnixNano())
			i++
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func returnResult(data e.Payload) {
	//mock change made to payload
	data.Colour = "rainbow"
	fmt.Println(data)
	json := e.GojayMarshal(&data)
	go pinger(outgoing, string(json))
}

func pinger(c chan string, p string) {
	outgoing <- p
}

func enqueuer(c chan string) {
	fmt.Println("this running")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	u.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	u.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"outgoingStuffs", // namel
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	u.FailOnError(err, "Failed to declare a queue")

	for {
		payload := <-outgoing

		// fulfillmentJson := e.StandardPkgMarshal(fulfillment)

		body := string(payload)

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		u.FailOnError(err, "Failed to publish a message")
	}
}
