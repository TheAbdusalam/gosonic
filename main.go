package main

import (
	"bufio"
	"fmt"
	"net"
)

type Server struct {
	host string
	port string
}

type Client struct {
	Conn        net.Conn
	existStatus chan int
	msg         chan string
}

type Conf struct {
	Host string
	Port string
}

func NewTCP(config Conf) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func (c *Client) handleAccept() {
	reader := bufio.NewReader(c.Conn)

	exit := func() {
		fmt.Printf("Client %s, exited", c.Conn.RemoteAddr())
		c.Conn.Close()
		c.existStatus <- 1
	}
	for {
		bytes, _ := reader.ReadBytes(10)
		fmt.Println("bytes, ", bytes)

		content := string(bytes)

		// exit if the client sends exit or the bytes are [13 10] which is the newline character
		if content == "exit\r\n" || bytes[0] == 13 || bytes[0] == 10 {
			c.Conn.Write([]byte("Exiting..."))
			fmt.Println("Client disconnected")
			c.existStatus <- 1
			break
		}

		fmt.Println("Connected to remote client: ", c.Conn.RemoteAddr())

		// broadcast the message to all clients
		fmt.Fprintf(c.Conn, "Message from %s: %s", c.Conn.RemoteAddr(), content)

	}

	<-c.existStatus
	exit()

}

func (s *Server) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		panic(err)
	}

	defer listen.Close()

	for {
		accept, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		client := &Client{
			Conn:        accept,
			existStatus: make(chan int),
		}

		go client.handleAccept()
	}

}

func main() {
	serv := Conf{
		Host: "localhost",
		Port: "3000",
	}

	server := NewTCP(serv)
	server.Start()
}
