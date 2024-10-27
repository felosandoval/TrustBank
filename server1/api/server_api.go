package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() { // Carga las variables de entorno desde el archivo .env

	// ---------------- CONFIGURACIONES (INICIO)----------------
	err := godotenv.Load("var_entorno.env")
	if err != nil {
		log.Fatal("Error al cargar el archivo .env de server_api")
	}
	// Obtiene el valor de CONNECTION_STRING del archivo .env
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("No se encontró la variable de entorno DB_CONNECTION_STRING")
	}
	// ConnectionString a de MongoDB
	clientOptions := options.Client().ApplyURI(connectionString)
	// Conectarse a la base de datos.
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Comprobar que la conexión es correcta.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err, "1")
	}
	fmt.Println("Conexión a la base de datos de MongoDB exitosa!")

	db := client.Database("banco")
	router := gin.Default()
	ColleccionBilleteras := db.Collection("Billeteras")
	ColleccionClientes := db.Collection("Clientes")
	//ColleccionMovimientos := db.Collection("Movimientos")
	// ---------------- CONFIGURACIONES (FIN) ----------------

	// ---------------- ENDPOINTS (INICIO) ----------------
	router.GET("/api/cliente", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		getCliente(ctx, ColleccionClientes, c)
	})

	router.POST("/api/inicio_sesion", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		iniciarSesion(ctx, ColleccionClientes, c)
	})

	router.POST("/api/deposito", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		depositar(ctx, ColleccionBilleteras, c)
	})

	router.POST("/api/transferencia", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		transferir(ctx, ColleccionBilleteras, c)
	})

	router.POST("/api/giro", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		girar(ctx, ColleccionBilleteras, c)
	})
	// ---------------- ENDPOINTS (FIN) ----------------

	puerto := os.Getenv("HTTP_PORT")
	if puerto == "" {
		log.Fatal("No se encontró el puerto PORT")
	}
	router.Run(":" + puerto)
}
