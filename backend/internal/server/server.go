package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"random_numbers/internal/dto"
	"random_numbers/internal/generator"
	"time"
)

type server struct {
	g            *gin.Engine
	number       string
	clientSecret string
}

func New() *server {
	return &server{gin.Default(), generator.Generate(), "76a69653500ee99eb3606d505d2efe381f24bab6"}
}

func (s *server) Generate() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				s.number = generator.Generate()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *server) Start(addr string) error {
	s.g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin, Content-Type, Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	fmt.Println("Server running on port 8000")
	s.Generate()
	s.g.GET("/getNumber", func(context *gin.Context) {
		var response = new(dto.GetNumberResponse)
		response.Generated = s.number
		context.JSON(http.StatusOK, response)
	})
	s.g.POST("/getAccessToken", func(context *gin.Context) {
		var request = new(dto.GetAccessTokenRequest)
		var response = new(dto.GetAccessTokenResponse)
		err := context.ShouldBindJSON(request)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		link := request.Url + "?client_id=" + request.Id + "&client_secret=" + s.clientSecret + "&code=" + request.Code
		fmt.Println(link)
		//resp, err := http.Get(link)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//fmt.Println(resp)
		response.Token = "fdsfdfdf"
		context.JSON(http.StatusOK, response)
	})
	return s.g.Run(addr)
}
