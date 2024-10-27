package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	google "server3/grpc"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

// server es usado para implementar movimientos.Movimientos.
type server struct {
	google.UnimplementedMovimientosServer
}

type Billeteras struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	NroCliente string             `bson:"nro_cliente"`
	Saldo      string             `bson:"saldo"`
	Divisa     string             `bson:"divisa"`
	Nombre     string             `bson:"nombre"`
	Activo     bool               `bson:"activo"`
}

type Mov struct {
	NroClienteOrigen  string `json:"nro_cliente_origen"`
	NroClienteDestino string `json:"nro_cliente_destino"`
	Monto             string `json:"monto"`
	Divisa            string `json:"divisa"`
	Tipo              string `json:"tipo"`
	FechaHora         string `json:"fecha_hora"`
	Id_billetera      string `json:"id_billetera"`
}

func (s *server) RealizarOperacion(ctx context.Context, in *google.MovimientoRequest) (*google.MovimientoResponse, error) {
	// conexion a mongo
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_CONNECTION_STRING")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	billeteraCollection := client.Database("banco").Collection("Billeteras")

	var billetera Billeteras
	err = billeteraCollection.FindOne(ctx, bson.M{"nro_cliente": in.NroClienteOrigen}).Decode(&billetera)
	if err != nil {
		log.Fatalf("No se pudo encontrar la billetera: %v", err)
	}

	id_origen := billetera.ID
	now := time.Now()
	fecha := now.Format("02/01/2006 15:04")

	movimientosCollection := client.Database("banco").Collection("Movimientos")
	Mov := Mov{NroClienteOrigen: in.NroClienteOrigen, NroClienteDestino: in.NroClienteDestino, Monto: in.Monto, Divisa: in.Divisa, Tipo: in.TipoOperacion, FechaHora: fecha, Id_billetera: id_origen.Hex()}
	_, err = movimientosCollection.InsertOne(ctx, Mov)
	if err != nil {
		return &google.MovimientoResponse{Mensaje: "Error al guardar el Mov"}, err
	}
	return &google.MovimientoResponse{Mensaje: "Mov guardado correctamente"}, nil
}

func main() {
	godotenv.Load("var_entorno.env")
	lis, err := net.Listen("tcp", os.Getenv("GRPC_HOST")+":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	google.RegisterMovimientosServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
