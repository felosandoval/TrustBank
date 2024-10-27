package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mov struct {
	Nro_cliente         string `json:"nro_cliente,omitempty"`
	Nro_cliente_origen  string `json:"nro_cliente_origen,omitempty"`
	Nro_cliente_destino string `json:"nro_cliente_destino,omitempty"`
	Monto               string `json:"monto,omitempty"`
	Divisa              string `json:"divisa,omitempty"`
	Tipo                string `json:"tipo,omitempty"`
}

type Sesion struct {
	Cliente    string `json:"numero_identificacion,omitempty"`
	Contrasena string `json:"contrasena,omitempty"`
}

func getCliente(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	numero_identificacion := c.Query("numero_identificacion")
	// Realizamos la búsqueda en la base de datos
	var cliente Cliente

	if err := collection.FindOne(ctx, bson.M{"numero_identificacion": numero_identificacion}, options.FindOne().SetProjection(bson.M{"contrasena": 0})).Decode(&cliente); err != nil {

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Retornamos el cliente encontrada
	c.JSON(http.StatusOK, gin.H{"cliente": cliente})
}

func iniciarSesion(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	var reqBody Sesion

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(reqBody)
	// Creamos el filtro para la búsqueda
	var cliente Cliente
	if err := collection.FindOne(ctx, bson.M{"numero_identificacion": reqBody.Cliente, "contrasena": reqBody.Contrasena}).Decode(&cliente); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"estado": "no_exitoso"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Retornamos la reserva encontrada
	c.JSON(http.StatusOK, gin.H{"estado": "exitoso"})
}

func depositar(ctx context.Context, collection *mongo.Collection, c *gin.Context) {

	var reqBody Mov

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cliente Billetera

	if err := collection.FindOne(ctx, bson.M{"nro_cliente": reqBody.Nro_cliente}).Decode(&cliente); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"estado": "cliente_no_encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := collection.FindOne(ctx, bson.M{"nro_cliente": reqBody.Nro_cliente, "activo": true}).Decode(&cliente); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"estado": "billetera_no_encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var deposito Mov

	deposito.Nro_cliente = reqBody.Nro_cliente
	deposito.Monto = reqBody.Monto
	deposito.Divisa = reqBody.Divisa
	deposito.Tipo = "deposito"

	server_productor(deposito)
}

func transferir(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	var reqBody Mov

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//<-----------------------------------------ORIGEN-------------------------------------------------->
	var origen Billetera

	if err := collection.FindOne(ctx, bson.M{"nro_cliente": reqBody.Nro_cliente_origen}).Decode(&origen); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"estado": "cliente_no_encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := collection.FindOne(ctx, bson.M{"nro_cliente": reqBody.Nro_cliente_origen, "activo": true}).Decode(&origen); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"estado": "billetera_no_encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//<------------------------------------------DESTINO------------------------------------------------->
	var destino Billetera

	if err := collection.FindOne(ctx, bson.M{"nro_cliente": reqBody.Nro_cliente_destino}).Decode(&destino); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"estado": "cliente_no_encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := collection.FindOne(ctx, bson.M{"nro_cliente": reqBody.Nro_cliente_destino, "activo": true}).Decode(&destino); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"estado": "billetera_no_encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	saldo, _ := strconv.ParseFloat(destino.Saldo, 64)
	monto, _ := strconv.ParseFloat(reqBody.Monto, 64)
	fmt.Println(saldo, monto)

	if monto > saldo {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"estado": "billetera_sin_fondos_suficientes"})
		return
	}

	//<--------------------------------------ENVIO----------------------------------------------->

	var transferencia Mov

	transferencia.Nro_cliente_origen = reqBody.Nro_cliente_origen
	transferencia.Nro_cliente_destino = reqBody.Nro_cliente_destino
	transferencia.Monto = reqBody.Monto
	transferencia.Divisa = reqBody.Divisa
	transferencia.Tipo = "transferencia"

	server_productor(transferencia)
}

func girar(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	var reqBody Mov

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cliente Billetera

	if err := collection.FindOne(ctx, bson.M{"nro_cliente": reqBody.Nro_cliente}).Decode(&cliente); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"estado": "cliente_no_encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := collection.FindOne(ctx, bson.M{"nro_cliente": reqBody.Nro_cliente, "activo": true}).Decode(&cliente); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"estado": "billetera_no_encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	saldo, _ := strconv.ParseFloat(cliente.Saldo, 64)
	monto, _ := strconv.ParseFloat(reqBody.Monto, 64)
	fmt.Println(saldo, monto)

	if monto > saldo {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"estado": "billetera_sin_fondos_suficientes"})
		return
	}

	var giro Mov

	giro.Nro_cliente = reqBody.Nro_cliente
	giro.Monto = reqBody.Monto
	giro.Divisa = reqBody.Divisa
	giro.Tipo = "giro"

	server_productor(giro)
}
