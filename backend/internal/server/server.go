package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"random_numbers/internal/dto"
	"random_numbers/internal/generator"
	"time"
)

type server struct {
	g            *gin.Engine
	origin       string
	number       string
	clientSecret string
}

func New() *server {
	return &server{
		gin.Default(),
		"http://localhost:5173",
		generator.Generate(),
		"76a69653500ee99eb3606d505d2efe381f24bab6",
	}
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
		AllowOrigins:     []string{s.origin},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin, Content-Type, Accept"},
		AllowWebSockets:  true,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	fmt.Println("Server running on port 8000")
	s.Generate()

	s.g.GET("/ws", func(context *gin.Context) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed to handshake": err.Error()})
			return
		}
		defer func(conn *websocket.Conn) {
			err = conn.Close()
			if err != nil {
				context.JSON(http.StatusOK, gin.H{"Failed to properly close connection": err.Error()})
				return
			}
		}(conn)
		curr := s.number
		err = conn.WriteMessage(websocket.TextMessage, []byte(curr))
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"Connection lost due to message sending failure": err.Error()})
			return
		}
		for {
			if curr != s.number {
				curr = s.number
				err = conn.WriteMessage(websocket.TextMessage, []byte(curr))
				if err != nil {
					context.JSON(http.StatusOK, gin.H{"Connection lost due to message sending failure": err.Error()})
					return
				}
			}
		}
	})

	s.g.POST("/getAccessToken", func(context *gin.Context) {
		var request = new(dto.GetAccessTokenRequest)
		var response = new(dto.GetAccessTokenResponse)
		err := context.ShouldBindJSON(request)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed to parse client request": err.Error()})
			return
		}
		link := request.Url + "?client_id=" + request.Id + "&client_secret=" + s.clientSecret + "&code=" + request.Code
		resp, err := http.Get(link)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed to connect to GitHub": err.Error()})
			return
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed to read response from GitHub": err.Error()})
			return
		}
		bodyString := string(bodyBytes)[13:53]
		response.Token = bodyString
		context.JSON(http.StatusOK, response)
	})

	s.g.POST("/getUserData", func(context *gin.Context) {
		var request = new(dto.GetUserDataRequest)
		var response = new(dto.GetUserDataResponse)
		var gitResponse = new(dto.GetUserDataGithubResponse)
		err := context.ShouldBindJSON(request)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed to parse client request": err.Error()})
			return
		}
		req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed to create request for GitHub": err.Error()})
			return
		}
		req.Header.Add("Authorization", "Bearer "+request.Token)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed to connect to GitHub": err.Error()})
			return
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed to read response from GitHub": err.Error()})
			return
		}
		err = json.Unmarshal(bodyBytes, &gitResponse)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed to unmarshal JSON": err.Error()})
			return
		}
		response.Login = gitResponse.Login
		response.Name = gitResponse.Name
		context.JSON(http.StatusOK, response)
	})

	return s.g.Run(addr)
}
