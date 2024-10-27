package main

type Billetera struct {
	Nro_cliente string `json:"nro_cliente"`
	Saldo       string `json:"saldo"`
	Divisa      string `json:"divisa"`
	Nombre      string `json:"nombre"`
	Activo      bool   `json:"activo"`
}
