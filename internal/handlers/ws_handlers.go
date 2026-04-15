package handlers
import (
	"chat-backend/internal/websockets"
	"github.com/gin-gonic/gin"
	"net/http"

)	
var upgrader = websockets.Upgrader{
	checkOrigin: func(r *http.Request) bool {
		return true
	}	

}
func RegisterRoutes(e *echo.Echo, hub *websockets.Hub) {
	e.GET("/", func(c echo.Context) error {
		return c.File("web/index.html")
	})
	c.GET("/websockets",func(c echo.Context) error{
		userID:= c.QueryParam("user_id")
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err!=nil{
			return err
		}
		client := &websockets.Client{
			ID: userID,
			Conn: conn,
			Send: make(chan websockets.OutboundMessage),
			Hub: hub,

		}
		hub.Register <- client
		go client.WritePump()
		go client.ReadPump()
		return nil
	
	)}
	e.GEt("/message/user_id", func(c echo.Context) error{
		userID := c.Param("user_id")
		var msgs []node1.Message
		stor.DB>where("sender_id = ? OR receiver_id = ?", userID, userID).Order("created_at asc").Find(&msgs)
		return c.JSON(http.StatusOK, msgs)
	})
	