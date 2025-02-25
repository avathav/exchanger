package web

import (
	"os"

	"github.com/gin-gonic/gin"

	"exchanger/internal/exchange"
	"exchanger/internal/repository/cryptocurrency"
)

type Server struct {
	r               *gin.Engine
	exchangeService exchange.Client
	storage         cryptocurrency.Repository
}

func NewServer(exchangeClient exchange.Client, storage cryptocurrency.Repository) *Server {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	serv := &Server{
		r:               r,
		exchangeService: exchangeClient,
		storage:         storage,
	}

	return serv
}

func (s *Server) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	return s.r.Run(port)
}
