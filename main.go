package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"net"
	"net/http"
	"strconv"
)

// TrackHandler -
func TrackHandler(response http.ResponseWriter, request *http.Request) {
	rgb := image.NewRGBA(image.Rect(0, 0, 1, 1))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(rgb, rgb.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	var img image.Image = rgb

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	response.Header().Set("Content-Type", "image/jpeg")
	response.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := response.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}

	ip, _, _ := net.SplitHostPort(request.RemoteAddr)

	log.Println(ip, "Ip")
	log.Println(request)
}

func main() {
	http.HandleFunc("/googlepix.jpg", TrackHandler)
	http.ListenAndServe(":5000", nil)
}
