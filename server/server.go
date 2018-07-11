package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

// Run starts up the websocket tcp proxy.
func Run(port int, cert string, key string, textMode bool, tcpAddress string) error {
	err := validatePort(port)
	if err != nil {
		return err
	}

	portString := fmt.Sprintf(":%d", port)

	log.Printf("[INFO] Listening on %s\n", portString)

	var proxyFunc = func(ws *websocket.Conn) {
		proxyHandler(ws, textMode, tcpAddress)
	}

	http.Handle("/", websocket.Handler(proxyFunc))

	if cert != "" && key != "" {
		err = http.ListenAndServeTLS(portString, cert, key, nil)
		if err != nil {
			return err
		}
	} else {
		err = http.ListenAndServe(portString, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func proxyHandler(ws *websocket.Conn, textMode bool, tcpAddress string) {
	conn, err := net.Dial("tcp", tcpAddress)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		return
	}

	if !textMode {
		ws.PayloadType = websocket.BinaryFrame
	}

	doneChan := make(chan bool)

	go copyData(conn, ws, doneChan)
	go copyData(ws, conn, doneChan)

	<-doneChan
	conn.Close()
	ws.Close()
	<-doneChan
}

func copyData(dst io.Writer, src io.Reader, doneChan chan<- bool) {
	io.Copy(dst, src)
	doneChan <- true
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("Invalid port requested. Valid values are 1-65535.")
	}

	return nil
}
