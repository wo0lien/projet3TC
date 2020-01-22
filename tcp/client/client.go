/*
! NOT USED
Build it in an other module to use it in parallel
*/

package client

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

const bufferSize = 1024

func main() {

	// Handle host:port, image and filter with flags
	host := flag.String("host", "127.0.0.1", "Nom d'hote du serveur")
	port := flag.Int("port", 8080, "Port du serveur")
	_ = flag.Int("filter", 1, "Filtre à utiliser : 1 = grayscale, 2 = edge") //TODO: add filter choice in the next lines
	filePath := flag.String("path", "", "--REQUIRED-- Chemin relatif vers l'image")

	flag.Parse()

	if *filePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Connect to the server
	connection, err := net.Dial("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal("tcp dial error", err)
	}

	fmt.Println("Connected to server")
	defer connection.Close()
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	filter := fillString(strconv.FormatInt(1, 10), 10)
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filter, filename and filesize!")
	connection.Write([]byte(filter))
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, bufferSize)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")

}

func fillString(returnString string, toLength int) string {
	for {
		lengtString := len(returnString)
		if lengtString < toLength {
			returnString = returnString + ":"
			continue
		}
		break
	}
	return returnString
}
