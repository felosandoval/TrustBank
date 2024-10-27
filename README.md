# TrustBank
* Felipe Sandoval
* Lucas Bernard

##### Por reglas de firewall, ejecutaremos la app en "Local" debido a errores en las VMs a la hora de conectarse a la base de datos de la VM1

#### ¿Cómo ejecutar la aplicación?
* Ingresar los siguientes 4 comandos, en el mismo orden:
1. go run server1/./api -> Inicializa la API REST.
2. go run server2/rabbitmq_consumer.go -> Inicializa el consumidor de mensajes de RabbitMQ.
3. go run server3/server_grpc.go -> Inicializa el servidor de gRPC
4. go run server1/client_menu.go -> Inicializa el menú inicial.
