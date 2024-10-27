package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	//"reflect"
	//"database/sql"
	//"bytes"

	pb "server2/grpc"

	"fmt"

	//"net/http"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var (
	conexion *amqp.Connection
	canal    *amqp.Channel
)

type Billetera struct {
	Nro_cliente string `json:"nro_cliente"`
	Saldo       string `json:"saldo"`
	Divisa      string `json:"divisa"`
	Nombre      string `json:"nombre"`
	Activo      bool   `json:"activo"`
}

type Mov struct {
	Nro_cliente         string `json:"nro_cliente,omitempty"`
	Nro_cliente_origen  string `json:"nro_cliente_origen,omitempty"`
	Nro_cliente_destino string `json:"nro_cliente_destino,omitempty"`
	Monto               string `json:"monto,omitempty"`
	Divisa              string `json:"divisa,omitempty"`
	Tipo                string `json:"tipo,omitempty"`
}

func grpc_girar_depositar(deserializada Mov) {

	// Configurar un cliente gRPC
	conn, err := grpc.Dial(":"+os.Getenv("GRPC_PORT"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()
	c := pb.NewMovimientosClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Crear una estructura de solicitud
	r := &pb.MovimientoRequest{
		NroClienteOrigen:  deserializada.Nro_cliente,
		NroClienteDestino: deserializada.Nro_cliente,
		Monto:             deserializada.Monto,
		Divisa:            deserializada.Divisa,
		TipoOperacion:     deserializada.Tipo,
	}

	// Envía la solicitud al servidor
	resp, err := c.RealizarOperacion(ctx, r)
	if err != nil {
		log.Fatalf("No se pudo realizar la operación: %v", err)
	}

	log.Printf("Respuesta del servidor: %s", resp.GetMensaje())
}

func grpc_transferir(deserializada Mov) {

	// Configurar un cliente gRPC
	conn, err := grpc.Dial(":"+os.Getenv("GRPC_PORT"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()
	c := pb.NewMovimientosClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Crear una estructura de solicitud
	r := &pb.MovimientoRequest{
		NroClienteOrigen:  deserializada.Nro_cliente_origen,
		NroClienteDestino: deserializada.Nro_cliente_destino,
		Monto:             deserializada.Monto,
		Divisa:            deserializada.Divisa,
		TipoOperacion:     deserializada.Tipo,
	}

	// Envía la solicitud al servidor
	resp, err := c.RealizarOperacion(ctx, r)
	if err != nil {
		log.Fatalf("No se pudo realizar la operación: %v", err)
	}

	log.Printf("Respuesta del servidor: %s", resp.GetMensaje())
}

func enCasoDeError(error error, msg string) {
	if error != nil {
		log.Panicf("%s: %s", msg, error)
	}
}

func identificarCaso(body []byte) {

	connectionString := os.Getenv("DB_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("No se encontró la variable de entorno DB_CONNECTION_STRING")
	}

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conexión a la base de datos de MongoDB exitosa!")

	db := client.Database("banco")
	ColleccionBilleteras := db.Collection("Billeteras")
	//ColleccionClientes := db.Collection("Clientes")

	var deserializada Mov
	err = json.Unmarshal(body, &deserializada)
	if err != nil {
		log.Printf("Error al deserializar el mensaje: %s", err)
		return
	}

	switch deserializada.Tipo {
	case "deposito":
		// Definir el filtro para encontrar el documento en la colección
		filter := bson.M{"nro_cliente": deserializada.Nro_cliente}
		// Crear una estructura para almacenar los datos de la billetera
		var billetera Billetera
		err := ColleccionBilleteras.FindOne(context.Background(), filter).Decode(&billetera)
		if err != nil {
			log.Printf("Error al buscar la billetera: %s", err)
			return
		}

		// Actualizar el saldo de la billetera
		saldo_int, _ := strconv.ParseFloat(billetera.Saldo, 64)
		monto_int, _ := strconv.ParseFloat(deserializada.Monto, 64)

		saldo_int += monto_int

		billetera.Saldo = strconv.FormatFloat(saldo_int, byte('f'), -1, 64)

		// Actualizar el documento en la base de datos
		_, err = ColleccionBilleteras.ReplaceOne(context.Background(), filter, billetera)
		if err != nil {
			log.Printf("Error al actualizar la billetera: %s", err)
			return
		}

		fmt.Println("Depósito procesado correctamente")

		grpc_girar_depositar(deserializada)

	case "transferencia":
		// --------------- ORIGEN ---------------
		// Definir el filtro para encontrar el documento en la colección
		filter_origen := bson.M{"nro_cliente": deserializada.Nro_cliente_origen}
		// Crear una estructura para almacenar los datos de la billetera
		var billetera_origen Billetera
		err_origen := ColleccionBilleteras.FindOne(context.Background(), filter_origen).Decode(&billetera_origen)
		if err != nil {
			log.Printf("Error al buscar la billetera: %s", err_origen)
			return
		}

		// Actualizar el saldo de la billetera
		saldo_origen_int, _ := strconv.ParseFloat(billetera_origen.Saldo, 64)
		monto_origen_int, _ := strconv.ParseFloat(deserializada.Monto, 64)

		saldo_origen_int -= monto_origen_int

		billetera_origen.Saldo = strconv.FormatFloat(saldo_origen_int, byte('f'), -1, 64)

		// Actualizar el documento en la base de datos
		_, err = ColleccionBilleteras.ReplaceOne(context.Background(), filter_origen, billetera_origen)
		if err != nil {
			log.Printf("Error al actualizar la billetera: %s", err)
			return
		}
		// --------------- DESTINO ---------------
		// Definir el filtro para encontrar el documento en la colección
		filter_destino := bson.M{"nro_cliente": deserializada.Nro_cliente_destino}
		fmt.Println()
		// Crear una estructura para almacenar los datos de la billetera
		var billetera_destino Billetera
		err_destino := ColleccionBilleteras.FindOne(context.Background(), filter_destino).Decode(&billetera_destino)
		if err != nil {
			log.Printf("Error al buscar la billetera: %s", err_destino)
			return
		}

		// Actualizar el saldo de la billetera
		saldo_destino_int, _ := strconv.ParseFloat(billetera_destino.Saldo, 64)
		monto_destino_int, _ := strconv.ParseFloat(deserializada.Monto, 64)

		saldo_destino_int += monto_destino_int

		billetera_destino.Saldo = strconv.FormatFloat(saldo_destino_int, byte('f'), -1, 64)

		// Actualizar el documento en la base de datos
		_, err = ColleccionBilleteras.ReplaceOne(context.Background(), filter_destino, billetera_destino)
		if err != nil {
			log.Printf("Error al actualizar la billetera: %s", err)
			return
		}

		grpc_transferir(deserializada)

	case "giro":
		// Definir el filtro para encontrar el documento en la colección
		filter := bson.M{"nro_cliente": deserializada.Nro_cliente}
		// Crear una estructura para almacenar los datos de la billetera
		var billetera Billetera
		err := ColleccionBilleteras.FindOne(context.Background(), filter).Decode(&billetera)
		if err != nil {
			log.Printf("Error al buscar la billetera: %s", err)
			return
		}

		// Actualizar el saldo de la billetera
		saldo_int, _ := strconv.ParseFloat(billetera.Saldo, 64)
		monto_int, _ := strconv.ParseFloat(deserializada.Monto, 64)

		saldo_int -= monto_int

		billetera.Saldo = strconv.FormatFloat(saldo_int, byte('f'), -1, 64)

		// Actualizar el documento en la base de datos
		_, err = ColleccionBilleteras.ReplaceOne(context.Background(), filter, billetera)
		if err != nil {
			log.Printf("Error al actualizar la billetera: %s", err)
			return
		}

		fmt.Println("Depósito procesado correctamente")

		grpc_girar_depositar(deserializada)

	default:
		log.Println("Tipo de movimiento no reconocido")
	}
}

func server_suscriptor() {
	error := godotenv.Load("var_entorno.env")
	if error != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	usuario_rabbit := os.Getenv("RABBITMQ_USERNAME")
	contrasena_rabbit := os.Getenv("RABBITMQ_PASSWORD")
	host_rabbit := os.Getenv("RABBITMQ_HOST")
	puerto_rabbit := os.Getenv("RABBITMQ_PORT")
	nombre_cola_rabbit := os.Getenv("RABBITMQ_QUEUE_NAME")

	conexion, error := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", usuario_rabbit, contrasena_rabbit, host_rabbit, puerto_rabbit))
	enCasoDeError(error, "Falla para conectar con RabbitMQ")
	defer conexion.Close()

	canal, error := conexion.Channel()
	enCasoDeError(error, "Falla para abrir un canal")
	defer canal.Close()

	routingKey := "tarea2"

	err := canal.ExchangeDeclare(
		"direct_exchange", // name
		"direct",          // type
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	enCasoDeError(err, "Falla para declarar el intercambio")

	cola, error := canal.QueueDeclare(
		nombre_cola_rabbit, // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	enCasoDeError(error, "Falla para declarar una queue")

	error = canal.QueueBind(
		cola.Name,         // queue name
		routingKey,        // routing key
		"direct_exchange", // exchange
		false,
		nil)
	enCasoDeError(error, "Fallo para bindear la cola")

	mensajes, error := canal.Consume(
		cola.Name, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	enCasoDeError(error, "Falla para registrar al consumidor")

	forever := make(chan os.Signal)
	signal.Notify(forever, os.Interrupt, os.Kill)

	go func() {
		for mensaje := range mensajes {
			identificarCaso(mensaje.Body)
			log.Printf("Mensaje recibido: %s", mensaje.Body)
		}
	}()

	log.Printf(" [*] Esperando un mensaje...")
	<-forever

	log.Printf("Parando al suscriptor.")
}

func main() {
	err := godotenv.Load("var_entorno.env")
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	server_suscriptor()

}
