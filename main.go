package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-echo/handler"
	"gopkg.in/mgo.v2"
	"log"
	"os"
)



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	e := echo.New()
	//Middleware and Init
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handler.Key),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			list := []string{"/login","/signup","/"}
			for _,b :=range list{
				if c.Path()==b{
					return true
				}
			}
			return false
		},
	}))
	prometheus.NewPrometheus("echo", nil).Use(e)

	db, err := mgo.Dial(os.Getenv("MONGO_URL"))
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Create indices
	if err = db.Copy().DB("twitter").C("users").EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}

	// Initialize handler
	h := &handler.Handler{DB: db}



	// Routes
	e.GET("/",h.TestHome)
	e.GET("/stream",h.TestStream)
	e.POST("/signup", h.Signup)
	e.POST("/login", h.Login)
	e.POST("/follow/:id", h.Follow)
	e.POST("/posts", h.CreatePost)
	e.GET("/feed", h.FetchPost)

	e.Logger.Fatal(e.Start(":"+os.Getenv("PORT")))

}
