package main

type Movimiento struct {
	Nro_cliente  string  `json:"nro_cliente"`
	Monto        float32 `json:"monto"`
	Divisa       string  `json:"divisa"`
	Tipo         string  `json:"tipo"`
	Fecha_hora   string  `json:"fecha_hora"`
	Id_billetera string  `json:"id_billetera"`
}
