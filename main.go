package main

import (
	"os"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"github.com/labstack/echo/middleware"
	"log"
	"os/signal"
	"syscall"
	"image"
	"image/jpeg"
	"io"
	"image/png"
	"image/gif"
	"errors"
)

var c WeebConfig
var r Registrator

func main() {
	c := LoadConfig()
	fmt.Printf("Starting %v (%v:%v) \n", c.Name, c.Host, c.Port)
	if c.Registration != nil && c.Registration.Enabled {
		println("registration activated")
		r = NewRegistrator(c.Name, c.Env, c.Port, c.Registration.Host, c.Registration.Token)
		err := r.Register()
		if err != nil {
			log.Fatal(err)
		}
	}
	shutdownHandler()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.POST("/", handleImageSubmit)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", c.Host, c.Port)))
}
func shutdownHandler() {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		cleanup()
	}()
}
func cleanup() {
	if r.active == true {
		err := r.Unregister()
		if err != nil {
			println(err)
		}
	}
	println("Shutting down")
	os.Exit(0)
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func handleImageSubmit(c echo.Context) error {
	file, err := c.FormFile("file")
	fileType := c.FormValue("type")
	if err != nil {
		return err
	}
	print(file.Header.Get(""))
	src, err := file.Open()
	if err != nil {
		return err
	}
	src.Close()
	i, err := decodeImage(fileType, src)
	if err != nil {
		return err
	}
	h := HashImage(i)
	return c.JSON(http.StatusOK, h)
}

func decodeImage(contentType string, file io.Reader) (image.Image, error) {
	var imageType = ""
	var err error
	var i image.Image
	switch contentType {
	case "image/jpg":
	case "image/jpeg":
		imageType = "jpg"
		break
	case "image/png":
		imageType = "png"
		break
	case "image/gif":
		imageType = "gif"
		break
	default:
		err = errors.New(fmt.Sprintf("invalid mimetype %v", contentType))
		break
	}
	if err != nil {
		return nil, err
	}
	switch imageType {
	case "jpg":
		i, err = jpeg.Decode(file)
		break
	case "png":
		i, err = png.Decode(file)
		break
	case "gif":
		i, err = gif.Decode(file)
		break
	default:
		err = errors.New(fmt.Sprintf("invalid mimetype %v", contentType))
		break
	}
	if err != nil {
		return nil, err
	}
	return i, nil
}
