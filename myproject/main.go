package main

import (
	"context"
	"log"
	"strings"
	"time"

	"html/template"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Структура для криптопроекта
type CryptoProject struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Status    string    `bson:"status"`
	Urgency   string    `bson:"urgency"`
	Deadline  time.Time `bson:"deadline"`
	CreatedAt time.Time `bson:"created_at"`
}

var client *mongo.Client
var collection *mongo.Collection

// Инициализация подключения к MongoDB
func initDB() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("cryptoDB").Collection("projects")
}

// Функция для создания нового проекта
func createProject(name, status, urgency string, deadline time.Time) error {
	project := CryptoProject{
		Name:      name,
		Status:    status,
		Urgency:   urgency,
		Deadline:  deadline,
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), project)
	return err
}

// Функция для получения всех проектов
func getProjects() ([]CryptoProject, error) {
	var projects []CryptoProject
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var project CryptoProject
		err := cursor.Decode(&project)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

// Функция для получения информации о конкретном проекте
func getProjectByID(id string) (*CryptoProject, error) {
	var project CryptoProject
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// Функция для преобразования строки в нижний регистр
func lower(s string) string {
	return strings.ToLower(s)
}

func main() {
	initDB()
	defer client.Disconnect(context.Background())

	// Создаем новый рутер Gin
	router := gin.Default()

	// Регистрируем функцию в шаблонах
	router.SetFuncMap(template.FuncMap{
		"lower": lower, // Регистрация функции для использования в шаблонах
	})

	// Загружаем шаблоны
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Главная страница с отображением проектов
	router.GET("/", func(c *gin.Context) {
		projects, err := getProjects()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.HTML(200, "index.html", gin.H{
			"projects": projects,
		})
	})

	// Страница для отображения деталей проекта
	router.GET("/project/:id", func(c *gin.Context) {
		id := c.Param("id")
		project, err := getProjectByID(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.HTML(200, "project.html", gin.H{
			"Project": project,
		})
	})

	// Обработка формы для добавления нового проекта
	router.POST("/add", func(c *gin.Context) {
		name := c.DefaultPostForm("name", "")
		status := c.DefaultPostForm("status", "")
		urgency := c.DefaultPostForm("urgency", "")
		deadlineStr := c.DefaultPostForm("deadline", "")

		// Преобразуем строку дедлайна в time.Time
		deadline, err := time.Parse("2006-01-02", deadlineStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format"})
			return
		}

		err = createProject(name, status, urgency, deadline)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Перенаправляем на главную страницу
		c.Redirect(302, "/")
	})

	// Запуск веб-сервера на порту 8080
	router.Run(":8080")
}
