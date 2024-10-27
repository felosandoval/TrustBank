package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conexion *amqp.Connection
	canal    *amqp.Channel
)

func enCasoDeError(error error, msg string) {
	if error != nil {
		log.Panicf("%s: %s", msg, error)
	}
}

func server_productor(movimiento Mov) {
	error := godotenv.Load("var_entorno.env")
	if error != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	usuario_rabbit := os.Getenv("RABBITMQ_USERNAME")
	contrasena_rabbit := os.Getenv("RABBITMQ_PASSWORD")
	host_rabbit := os.Getenv("RABBITMQ_HOST")
	puerto_rabbit := os.Getenv("RABBITMQ_PORT")

	conexion, error := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", usuario_rabbit, contrasena_rabbit, host_rabbit, puerto_rabbit))
	enCasoDeError(error, "Falla para conectar con RabbitMQ")
	defer conexion.Close()

	canal, error := conexion.Channel()
	enCasoDeError(error, "Falla para abrir un canal")
	defer canal.Close()

	routingKey := "tarea2"

	error = canal.ExchangeDeclare(
		"direct_exchange", // name
		"direct",          // type
		false,             // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	enCasoDeError(error, "Error al declarar un exchange")
	jsonData, error := json.Marshal(movimiento)
	error = canal.Publish(
		"direct_exchange", // exchange
		routingKey,        // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonData),
		})
	enCasoDeError(error, "Error al publicar un mensaje")
	log.Printf(" [x] Enviado %s", jsonData)
}
