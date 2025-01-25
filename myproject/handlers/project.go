package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"your_project/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

// Инициализация MongoDB клиента
func InitDB() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	collection = client.Database("cryptoprojects").Collection("projects")
}

// Обработчик для получения всех проектов
func GetProjects(w http.ResponseWriter, r *http.Request) {
	var projects []models.Project
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var project models.Project
		if err := cursor.Decode(&project); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		projects = append(projects, project)
	}
	fmt.Fprintf(w, "Projects: %+v", projects)
}

// Обработчик для добавления нового проекта
func AddProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project.CreatedAt = time.Now().Format(time.RFC3339)
	_, err = collection.InsertOne(context.Background(), project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Project added successfully")
}
