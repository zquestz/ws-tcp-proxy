package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

var binaryProxy = true
var address = ""

// Run starts up the websocket tcp proxy.
func Run(port int, cert string, key string, binaryMode bool, tcpAddress string) error {
	err := validatePort(port)
	if err != nil {
		return err
	}

	binaryProxy = binaryMode
	address = tcpAddress

	portString := fmt.Sprintf(":%d", port)

	log.Printf("[INFO] Listening on %s\n", portString)

	http.Handle("/", websocket.Handler(proxyHandler))

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

func proxyHandler(ws *websocket.Conn) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		return
	}

	if binaryProxy {
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
