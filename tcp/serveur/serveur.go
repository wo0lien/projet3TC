package serveur

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const bufferSize = 1024
const host = "127.0.0.1"
const port = "8080"

/*
StartServer function makes the server listen on the desired port
*/
func StartServer(port int) {

	listener, err := net.Listen("tcp", "127.0.0.1"+string(port))
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) string {

	defer connection.Close()
	fmt.Println("Connected to client, start receiving the file name, file size and filter")
	bufferFilter := make([]byte, 10)
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFilter)
	filter := strings.Trim(string(bufferFilter), ":")

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < bufferSize {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+bufferSize)-fileSize))
			break
		}
		io.CopyN(newFile, connection, bufferSize)
		receivedBytes += bufferSize
	}
	fmt.Println("Received file " + fileName + " completely!" + " and filter : " + filter)
	return fileName
}
