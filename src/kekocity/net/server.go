package net

// <imports>
import (
  "log"
  "fmt"
  "net/http"

  "golang.org/x/net/websocket"

  pnet "kekocity/net/packet"
  cmap "kekocity/misc/concurrentmap"
)

var server *Server

type Server struct {
  port int

  players *cmap.ConcurrentMap
}

func clientConnection(clientsock *websocket.Conn) {
  packet := pnet.NewPacket()
  buffer := make([]uint8, pnet.PACKET_MAXSIZE)

  recv, err := clientsock.Read(buffer)

  if err == nil {
    copy(packet.Buffer[0:recv], buffer[0:recv])

    parseFirstMessage(clientsock, packet)
  } else {
    if err.Error() != "EOF" {
      log.Println("net.server", "Client connection error:", err.Error())
    }
  }
}

func parseFirstMessage(_conn *websocket.Conn, _packet *pnet.Packet) {
  message := _packet.ToString()

  // If the first packet length is < 1 close the socket
  if len(message) < 1 {
    _conn.Close()
    return
  }

  // Create the connection
  connection := NewConnection(_conn)
}

func Listen(_port int) {
  server.port = _port

  log.Println("Listening for new connections...")

  http.Handle("/ws", websocket.Handler(clientConnection))

	err := http.ListenAndServe(fmt.Sprintf(":%d", _port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}