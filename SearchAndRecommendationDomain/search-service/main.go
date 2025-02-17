package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var esClient *elasticsearch.Client
var mongoClient *mongo.Client

// Carga las variables de entorno y configura la conexión con MongoDB y Elasticsearch
func init() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env, se usarán variables de entorno del sistema.")
	}

	// Obtener las variables de entorno
	mongoURI := os.Getenv("MONGO_URI")
	elasticURI := os.Getenv("ELASTICSEARCH_URI")

	if mongoURI == "" {
		log.Fatal("Error: MONGO_URI no está definido en las variables de entorno")
	}

	// Conectar a MongoDB Atlas
	var err error
	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	// Verificar conexión con MongoDB Atlas
	if err := mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatalf("Error de ping a MongoDB: %v", err)
	}
	log.Println("Conectado a MongoDB Atlas correctamente")

	// Conectar a Elasticsearch
	esClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{elasticURI},
	})
	if err != nil {
		log.Fatalf("Error al conectar a Elasticsearch: %v", err)
	}
	log.Println("Conectado a Elasticsearch correctamente")

	// Crear el índice 'models' si no existe
	createIndexIfNotExists()
}

// Función para crear el índice 'models' si no existe
func createIndexIfNotExists() {
	// Verificar si el índice 'models' existe
	res, err := esClient.Indices.Exists([]string{"models"})
	if err != nil {
		log.Fatalf("Error al verificar el índice: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		// Si el índice no existe, lo creamos
		log.Println("El índice 'models' no existe. Creando el índice...")
		createIndexBody := `{
			"settings": {
				"number_of_shards": 1,
				"number_of_replicas": 0
			},
			"mappings": {
				"properties": {
					"name": { "type": "text" },
					"description": { "type": "text" },
					"created_by": { "type": "keyword" },
					"price": { "type": "float" }
				}
			}
		}`
		createRes, err := esClient.Indices.Create("models", esClient.Indices.Create.WithBody(strings.NewReader(createIndexBody)))
		if err != nil {
			log.Fatalf("Error al crear el índice 'models': %v", err)
		}
		defer createRes.Body.Close()

		if createRes.IsError() {
			log.Fatalf("Error al crear el índice 'models': %s", createRes.String())
		} else {
			log.Println("Índice 'models' creado correctamente.")
		}
	} else {
		log.Println("El índice 'models' ya existe.")
	}
}

// searchModelHandler maneja la búsqueda de modelos en Elasticsearch
func searchModelHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el parámetro de búsqueda desde la URL
	params := mux.Vars(r)
	modelName := params["model_name"]

	// Construir consulta JSON para Elasticsearch
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"name": "%s"
			}
		}
	}`, modelName)

	// Ejecutar la búsqueda en Elasticsearch
	res, err := esClient.Search(
		esClient.Search.WithIndex("models"),
		esClient.Search.WithBody(strings.NewReader(query)),
		esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al realizar la búsqueda: %v", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		http.Error(w, "Error al obtener resultados de búsqueda en Elasticsearch", http.StatusInternalServerError)
		return
	}

	// Decodificar la respuesta de Elasticsearch
	var searchResults map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResults); err != nil {
		http.Error(w, fmt.Sprintf("Error al leer la respuesta de Elasticsearch: %v", err), http.StatusInternalServerError)
		return
	}

	// Devolver la respuesta como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResults)
}

func main() {
	// Crear el enrutador
	r := mux.NewRouter()

	// Definir los endpoints
	r.HandleFunc("/search/{model_name}", searchModelHandler).Methods("GET")

	// Iniciar el servidor en el puerto 5005
	fmt.Println("Iniciando el microservicio de búsqueda en el puerto 5005...")
	log.Fatal(http.ListenAndServe(":5005", r))
}
