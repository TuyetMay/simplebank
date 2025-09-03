package api 

import(
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)
//Server serves all HTTP requests for our bank service : chịu trách nhiệm nhận, xử lý và trả lời các request
type Server struct{
	store *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server ans set up routing
func NewServer(store *db.Store) *Server{
	server := &Server{store: store}
	router := gin.Default() // tạo router mặc định của gin framework

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.ListAccount)


	server.router = router
	return server

}

// take an address as input and return error. Start runs the http server on a specific address
func (server *Server) Start ( address string) error{
	return server.router.Run(address)
}
// function will take an error as input and it will  return a gin.H object
func errResponse (err error) gin.H{
	return gin.H{"error": err.Error()}

}

