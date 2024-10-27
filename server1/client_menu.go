// ---- COMANDOS A EJECUTAR ----
// go run path_1/client_menu.go			 Inicializa el menú inicial.
// go run path_2/server_api.go			 Inicializa la API REST.
// go run path_3/rabbitmq_consumer.go	 Inicializa el consumidor de mensajes de RabbitMQ.
// go run path_4/server_grpc.go			 Inicializa el servidor de gRPC.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sesion struct {
	Cliente    string `json:"numero_identificacion,omitempty"`
	Contrasena string `json:"contrasena,omitempty"`
}

type Estado struct {
	Estado string `json:"estado"`
}

type Mov struct {
	Nro_cliente         string `json:"nro_cliente,omitempty"`
	Nro_cliente_origen  string `json:"nro_cliente_origen,omitempty"`
	Nro_cliente_destino string `json:"nro_cliente_destino,omitempty"`
	Monto               string `json:"monto,omitempty"`
	Divisa              string `json:"divisa,omitempty"`
	Tipo                string `json:"tipo,omitempty"`
}

type Cliente struct {
	ID                    primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Nombre                string             `bson:"nombre,omitempty" json:"nombre,omitempty"`
	Contrasena            string             `bson:"contrasena,omitempty" json:"contrasena,omitempty"`
	Fecha_nacimiento      string             `bson:"fecha_nacimiento,omitempty" json:"fecha_nacimiento,omitempty"`
	Direccion             string             `bson:"direccion,omitempty" json:"direccion,omitempty"`
	Numero_identificacion string             `bson:"numero_identificacion,omitempty" json:"numero_identificacion,omitempty"`
	Email                 string             `bson:"email,omitempty" json:"email,omitempty"`
	Telefono              string             `bson:"telefono,omitempty" json:"telefono,omitempty"`
	Genero                string             `bson:"genero,omitempty" json:"genero,omitempty"`
	Nacionalidad          string             `bson:"nacionalidad,omitempty" json:"nacionalidad,omitempty"`
	Ocupacion             string             `bson:"ocupacion,omitempty" json:"ocupacion,omitempty"`
}

func verificarSesion(sesion Sesion) int {
	// Convertir la estructura sesión en un objeto JSON

	jsonData, err := json.Marshal(sesion)
	if err != nil {
		// Manejar el error de conversión
	}

	// Crear una solicitud POST con el objeto JSON como cuerpo de la solicitud
	url := "http://127.0.0.1:3000/api/inicio_sesion"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// Manejar el error de creación de la solicitud
	}

	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Manejar el error de envío de la solicitud
	}

	defer resp.Body.Close()

	// Leer la respuesta de la solicitud
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Manejar el error de lectura de la respuesta
	}

	// Decodificar el JSON en una estructura de datos
	strBody := string(body)
	var estado Estado
	err = json.Unmarshal([]byte(strBody), &estado)

	if err != nil {
		// Manejar el error de decodificación JSON
	}
	if estado.Estado == "exitoso" {
		// Hacer algo si la comparación es exitosa
		return 1

	}
	return 0
}

func depositar(sesion Sesion) {
	var monto string
	fmt.Printf("Ingrese un monto: ")
	fmt.Scan(&monto)

	var deposito Mov
	deposito.Monto = monto
	deposito.Nro_cliente = sesion.Cliente
	deposito.Divisa = "USD"

	jsonData, err := json.Marshal(deposito)
	if err != nil {
		// Manejar el error de conversión
	}

	// Crear una solicitud POST con el objeto JSON como cuerpo de la solicitud
	url := "http://127.0.0.1:3000/api/deposito"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// Manejar el error de creación de la solicitud
	}

	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Manejar el error de envío de la solicitud
	}

	defer resp.Body.Close()

	// Leer la respuesta de la solicitud
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Manejar el error de lectura de la respuesta
	}

	print(string(body))

	fmt.Println("El deposito ha sido enviado correctamente")

}

func girar(sesion Sesion) {
	var monto string
	fmt.Printf("Ingrese un monto: ")
	fmt.Scan(&monto)

	var giro Mov
	giro.Monto = monto
	giro.Nro_cliente = sesion.Cliente
	giro.Divisa = "USD"

	jsonData, err := json.Marshal(giro)
	if err != nil {
		// Manejar el error de conversión
	}

	// Crear una solicitud POST con el objeto JSON como cuerpo de la solicitud
	url := "http://127.0.0.1:3000/api/giro"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// Manejar el error de creación de la solicitud
	}

	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Manejar el error de envío de la solicitud
	}

	defer resp.Body.Close()

	// Leer la respuesta de la solicitud
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Manejar el error de lectura de la respuesta
	}

	// Decodificar el JSON en una estructura de datos
	strBody := string(body)
	var estado Estado
	err = json.Unmarshal([]byte(strBody), &estado)

	if err != nil {
		// Manejar el error de decodificación JSON
	}
	if estado.Estado == "billetera_sin_fondos_suficientes" {
		// Hacer algo si la comparación es exitosa
		fmt.Println("El saldo es insuficiente")
		return
	}

	fmt.Println("El giro ha sido solicitado correctamente")

}

func transferir(sesion Sesion) {
	var destino string
	fmt.Printf("Ingrese la cuenta de destino: ")
	fmt.Scan(&destino)

	var monto string
	fmt.Printf("Ingrese un monto: ")
	fmt.Scan(&monto)

	var transferencia Mov
	transferencia.Monto = monto
	transferencia.Nro_cliente_origen = sesion.Cliente
	transferencia.Nro_cliente_destino = destino

	transferencia.Divisa = "USD"

	jsonData, err := json.Marshal(transferencia)
	if err != nil {
		// Manejar el error de conversión
	}

	// Crear una solicitud POST con el objeto JSON como cuerpo de la solicitud
	url := "http://127.0.0.1:3000/api/transferencia"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// Manejar el error de creación de la solicitud
	}

	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Manejar el error de envío de la solicitud
	}

	defer resp.Body.Close()

	// Leer la respuesta de la solicitud
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Manejar el error de lectura de la respuesta
	}
	strBody := string(body)
	var estado Estado
	err = json.Unmarshal([]byte(strBody), &estado)

	if err != nil {
		// Manejar el error de decodificación JSON
	}
	if estado.Estado == "cliente_no_encontrado" {
		// Hacer algo si la comparación es exitosa
		fmt.Println("La cuenta ingresada no existe")
		return
	}
	if estado.Estado == "billetera_sin_fondos_suficientes" {
		// Hacer algo si la comparación es exitosa
		fmt.Println("Su saldo es insuficiente")
		return
	}

	fmt.Println("La transferencia fue enviada correctamente")

}

func main() {
	fmt.Println("Bienvenido a TrustBank!")
	// Iniciamos un ciclo infinito
	for {
		var opcion int
		fmt.Println("1. Iniciar sesión")
		fmt.Println("2. Salir")
		fmt.Printf("Ingrese una opción: ")
		fmt.Scan(&opcion)

		switch opcion {
		case 1: // INICIAR SESION

			var sesion Sesion
			var numero_identificacion string
			fmt.Printf("Ingrese su número de identificación: ")
			fmt.Scan(&numero_identificacion)

			var password string
			fmt.Printf("Ingrese su contraseña: ")
			fmt.Scan(&password)

			sesion.Cliente = numero_identificacion
			sesion.Contrasena = password

			// VALIDAR DATOS DE INICIO DE SESIÓN
			if verificarSesion(sesion) == 1 {
				fmt.Println("Login exitoso!")
				for {
					var subOpcion int
					fmt.Println("1. Realizar deposito")
					fmt.Println("2. Realizar transferencia")
					fmt.Println("3. Realizar giro")
					fmt.Println("4. Salir")
					fmt.Printf("Ingrese una opción: ")
					fmt.Scan(&subOpcion)

					switch subOpcion {
					case 1: // DEPOSITAR
						depositar(sesion)
					case 2: // TRANSFERIR
						transferir(sesion)
					case 3: // GIRAR
						girar(sesion)
					case 4: // SALIR
						// Rompemos el ciclo del submenú al ingresar la opción 4
						break

					default:
						fmt.Println("Ingrese una opcion valida")
					}

					// Si la subOpción es 4, salimos del ciclo del submenú
					if subOpcion == 4 {
						break
					}
				}
			} else { // SI INGRESA DATOS ERRÓNEOS
				fmt.Println("Número de identificación o contraseña incorrecta")
			}
		case 2: // SALIR
			fmt.Println("Gracias por usar TrustBank!")
			// Rompemos el ciclo al ingresar la opción 2
			break
		default:
			fmt.Println("Ingrese una opcion valida")
		}
		// Si la opción es 2, salimos del ciclo principal
		if opcion == 2 {
			break
		}
	}
}
