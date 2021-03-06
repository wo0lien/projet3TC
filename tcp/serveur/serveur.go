package serveur

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/wo0lien/projet3TC/filters/edge"
	"github.com/wo0lien/projet3TC/filters/grayscale"
	"github.com/wo0lien/projet3TC/filters/negative"
	"github.com/wo0lien/projet3TC/filters/noise"
	"github.com/wo0lien/projet3TC/imagetools"
	"image"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const bufferSize = 1024

/*
StartServer function makes the server listen on the desired port
*/
func StartServer(port int, host string, concurrent bool) {

	var _ = edge.FSobel
	var _ = grayscale.GrayFilter
	var _ = noise.Fmean
	var _ = negative.NegativeFilter

	listener, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		// le serveur stop si cette erreur est raised
		log.Fatal("tcp server listener error:", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %+v", err)
			continue
		}

		go handleConnection(conn, concurrent)
	}
}

func handleConnection(connection net.Conn, concurrent bool) {

	fmt.Println("Handling connection")

	defer connection.Close()

	fileName, filter, err := receiveFile(connection)

	if err != nil {
		fmt.Printf("Error: %+v", err)
		return
	}

	//----------------------Traitement de l'image reçue------------------

	img, err := imagetools.Open(fileName)

	if err != nil {
		fmt.Printf("Error: %+v", err)
		return
	}

	var imgFiltered image.Image

	if concurrent == true {
		// apply filter
		switch filter {
		case "1":
			imgFiltered = negative.ConcurrentNegFilter(img)
		case "2":
			imgFiltered = grayscale.ConcurrentGrayFilter(img)
		case "3":
			imgFiltered = edge.ConcurrentEdgeFilter(img)
		case "4":
			imgFiltered = noise.ConcurrentFmediane(img, 3)
		case "5":
			imgFiltered = noise.ConcurrentFmean(img, 3)
		default:
			fmt.Println("Error: Not a good index")
			return
		}
	} else {
		// apply filter
		switch filter {
		case "1":
			imgFiltered = negative.NegativeFilter(img)
		case "2":
			imgFiltered = grayscale.GrayFilter(img)
		case "3":
			imgFiltered = edge.FSobel(img)
		case "4":
			imgFiltered = noise.Fmediane(img, 3)
		case "5":
			imgFiltered = noise.Fmean(img, 3)
		default:
			fmt.Println("Error: Not a good index")
			return
		}
	}

	imgFilteredFileName := "f_" + fileName

	err = imagetools.Export(imgFiltered, imgFilteredFileName)

	if err != nil {
		fmt.Printf("Error: %+v", err)
		return
	}

	//----------------------Renvoi -----------------------------------

	err = sendFileBack(imgFilteredFileName, connection)
	if err != nil {
		fmt.Printf("Error: %+v", err)
		return
	}

	fmt.Println("Closing the connnection")

}

func receiveFile(connection net.Conn) (string, string, error) {

	fmt.Println("Start receiving the file name, file size and filter")
	bufferFilter := make([]byte, 10)
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFilter)
	filter := strings.Trim(string(bufferFilter), ":")

	connection.Read(bufferFileSize)
	fileSize, err := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	if err != nil {
		return "", "", errors.Wrap(err, "File receiver failed with error :")
	}

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)

	if err != nil {
		return "", "", errors.Wrap(err, "File receiver failed with error :")
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
	fmt.Println("Received file " + fileName + ", client want to compute on it with filter : " + filter)
	return fileName, filter, nil
}

func sendFileBack(fileName string, connection net.Conn) error {
	//---------------------Renvoi de l'image-----------------------------

	fmt.Println("Starting to send back file")

	file, err := os.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "File sender failed with error :")
	}
	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		return errors.Wrap(err, "File sender failed with error :")
	}

	fileBackSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileBackName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize")
	connection.Write([]byte(fileBackSize))
	connection.Write([]byte(fileBackName))
	sendBuffer := make([]byte, bufferSize)
	fmt.Println("Start sending file")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent")
	return nil
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
