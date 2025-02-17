package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	POSTGRES_URI = "postgresql://postgres.imfqyzgimtercyyqeqof:1997Guallaba@aws-0-us-west-1.pooler.supabase.com:6543/postgres"
	CATALOG_URL  = "http://52.91.86.137:5003/models/id/" // URL base del servicio catalog-service
)

type Order struct {
	ID           int     `json:"id"`
	ModelID      string  `json:"model_id"`
	CustomParams string  `json:"custom_params"`
	UserID       string  `json:"user_id"`
	CreatedAt    string  `json:"created_at"`
	CostInitial  float64 `json:"cost_initial"`
	CostFinal    float64 `json:"cost_final"`
	ModelDetails Model   `json:"model_details,omitempty"`
}

type Model struct {
	ModelID     string `json:"model_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Format      string `json:"format"`
	Price       string `json:"price"`
	CreatedBy   string `json:"created_by"`
}

// Función para obtener detalles del modelo desde catalog-service
func fetchModelDetails(modelID string) (Model, error) {
	var modelResponse struct {
		Model Model `json:"model"`
	}

	// Hacer la petición HTTP a catalog-service
	resp, err := http.Get(CATALOG_URL + modelID)
	if err != nil {
		return Model{}, fmt.Errorf("error requesting model: %v", err)
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Model{}, fmt.Errorf("error reading response: %v", err)
	}

	// Si el modelo no se encuentra, devolver vacío
	if resp.StatusCode != http.StatusOK {
		return Model{}, fmt.Errorf("model not found")
	}

	// Parsear JSON
	if err := json.Unmarshal(body, &modelResponse); err != nil {
		return Model{}, fmt.Errorf("error parsing JSON: %v", err)
	}

	return modelResponse.Model, nil
}

// Obtener órdenes filtradas por user_id y agregar detalles del modelo
func GetOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	// Conectar a PostgreSQL
	db, err := sql.Open("postgres", POSTGRES_URI)
	if err != nil {
		log.Println("Error opening database connection:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	// Ejecutar consulta SQL
	rows, err := db.Query("SELECT id, model_id, custom_params, user_id, created_at, cost_initial, cost_final FROM customs WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error querying database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}
	defer rows.Close()

	var orders []Order

	// Procesar resultados
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.ModelID, &order.CustomParams, &order.UserID, &order.CreatedAt, &order.CostInitial, &order.CostFinal)
		if err != nil {
			log.Println("Error scanning order:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning order"})
			return
		}

		// Obtener detalles del modelo desde catalog-service
		modelDetails, err := fetchModelDetails(order.ModelID)
		if err != nil {
			log.Println("Model details not found for model_id:", order.ModelID)
		} else {
			order.ModelDetails = modelDetails
		}

		orders = append(orders, order)
	}

	// Enviar respuesta JSON
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func main() {
	// Configurar Gin en modo release
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://52.91.86.137:8080"}, // Permitir solo la interfaz
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.SetTrustedProxies(nil)

	router.GET("/orders/user/:user_id", GetOrdersByUserID)

	log.Println("Server running on port 5008...")

	router.Run(":5008")
}
