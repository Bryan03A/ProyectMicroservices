// models/model.go
package models

type Modelo3D struct {
    ModelID    string `json:"model_id"`
    Nombre     string `json:"nombre"`
    Descripcion string `json:"descripcion"`
    Archivo    string `json:"archivo"`
}
