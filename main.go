package main
import (
	"log"
	"os"
	"chat-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)
func main() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "1234")
	os.Setenv("DB_NAME", "testdb")


	store.ConnectDatabase()
	store.DB.AutoMigrate(&websockets.OutboundMessage{})
	hub := websockets.NewHub()
	go hub.Run()
	e:= echo.New()
	
	r := gin.Default()
	handlers.RegisterRoutes(e, hub)
	log.Println("Server running on port 8080")
 e.Logger.Fatal(e.Start(":8080"))
	}