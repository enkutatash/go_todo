package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

type Todo struct {
    ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    Completed bool               `json:"completed" bson:"completed"`
    Body      string             `json:"body" bson:"body"`
}

func main() {
	print("Hello, World!")

	
	

	err := godotenv.Load(".env")

	if err != nil {

		log.Fatalf("Error loading .env file")
	}

	mongo_url := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongo_url)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil{
		log.Fatalf("Error connecting to MongoDB")
	}
	defer client.Disconnect(context.Background())
	err = client.Ping(context.Background(), nil)
	if err != nil{
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Get("/todos",allTodos)
	app.Post("/todos",addTodo)
	app.Patch("/todos/:id",updateTodo)
	app.Delete("/todos/:id",deleteTodo)
	

	port := os.Getenv("PORT")
	if port == ""{
		port = "3000"
	}
	log.Fatal(app.Listen("0.0.0.0:"+port))


	// os.Getenv("PORT")
	

}

func allTodos(c *fiber.Ctx) error{

	var todos []Todo
	cursor , err := collection.Find(context.Background(),bson.M{})
	if err != nil{
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()){
		var todo Todo
		if err := cursor.Decode(&todo); err != nil{
			return err
		}
		
		todos = append(todos,todo)
	}
	return c.Status(200).JSON(todos)

}

func addTodo(c *fiber.Ctx) error{
	log.Printf("Raw Request Body: %s", string(c.Body()))
	 todo:= new(Todo)
	 if err := c.BodyParser(todo);err != nil{
		log.Printf("Error parsing request body: %v", err)
        return c.Status(422).JSON(fiber.Map{
            "error": "Invalid request body",
        })
	 }
	 if todo.Body == ""{
		 return c.Status(400).SendString("Body is required")
	 }
	 insertResult, err := collection.InsertOne(context.Background(),todo)
	 if err != nil{
		log.Printf("Failed to insert todo: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create todo",
		})
	 }
	 todo.ID = insertResult.InsertedID.(primitive.ObjectID)
	 return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error{
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return c.Status(400).SendString("Invalid ID")
	}
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil{
		return c.Status(422).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	filter := bson.M{"_id": objectId}

	update := bson.M{"$set":bson.M{"completed": todo.Completed}}
	// update := bson.M{"$set": bson.M{"completed": todo.Completed,},}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update todo",
		})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Todo updated successfully"})
}

func deleteTodo(c *fiber.Ctx) error{
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil{
		return c.Status(400).SendString("Invalid ID")
	}
	filter := bson.M{"_id": objectId}

	_,err = collection.DeleteOne(context.Background(),filter)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete todo",
		})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Todo deleted successfully"})
}