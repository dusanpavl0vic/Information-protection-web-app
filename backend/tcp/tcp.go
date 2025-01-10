package main

import (
	"backend-ZI/hash"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct{}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server started on :3000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := new(bytes.Buffer)

	var fileNameLen int32
	if err := binary.Read(conn, binary.LittleEndian, &fileNameLen); err != nil {
		log.Println("Error reading file name length:", err)
		return
	}

	fileName := make([]byte, fileNameLen)
	if _, err := io.ReadFull(conn, fileName); err != nil {
		log.Println("Error reading file name:", err)
		return
	}
	fmt.Printf("Received file name: %s\n", string(fileName))

	var fileSize int64
	if err := binary.Read(conn, binary.LittleEndian, &fileSize); err != nil {
		log.Println("Error reading file size:", err)
		return
	}
	fmt.Printf("Received file size: %d bytes\n", fileSize)

	var hashLen int32
	if err := binary.Read(conn, binary.LittleEndian, &hashLen); err != nil {
		log.Println("Error reading hash length:", err)
		return
	}
	fmt.Printf("Received hash length: %d bytes\n", hashLen)

	hashh := make([]byte, hashLen)
	if _, err := io.ReadFull(conn, hashh); err != nil {
		log.Println("Error reading hash:", err)
		return
	}
	fmt.Printf("Received hash: %x\n", hashh)

	if _, err := io.CopyN(buf, conn, fileSize); err != nil && err != io.EOF {
		log.Println("Error reading file data:", err)
		return
	}

	fileContent := buf.Bytes()

	hashTest := hash.TigerHash(fileContent)

	if bytes.Equal(hashh, hashTest[:]) {
		fmt.Println("Hashes match!")
	} else {
		fmt.Println("Hashes don't match")
		fmt.Printf("Expected hash: %x\n", hashh)
		fmt.Printf("Calculated hash: %x\n", hashTest)
	}

	fmt.Println("Hashes 1: ", hashh)
	fmt.Println("Hashes 2: ", hashTest[:])

	if bytes.Equal(hashh, hashTest[:]) {
		fmt.Println("Hashes match!")
	} else {
		fmt.Println("Hashes don't match")
	}

	// FIXME: Dekodiranje fajla koji je primljen
	fmt.Printf("Received %d bytes of file content\n", len(fileContent))
	fmt.Println("File content:")
	fmt.Println(string(fileContent))
}

func sendFile(fileName string, data []byte) error {

	// FIXME: Kodiranje fajla pre slanja
	// FIXME: Smisliti po kom algoritmu cu da kodiram na koji nacin da prosledim
	hash := hash.TigerHash(data)
	hashLen := int32(len(hash))

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}
	defer conn.Close()

	fileNameBytes := []byte(fileName)
	if err := binary.Write(conn, binary.LittleEndian, int32(len(fileNameBytes))); err != nil {
		return err
	}
	if _, err := conn.Write(fileNameBytes); err != nil {
		return err
	}
	fmt.Printf("Written %d bytes of file name\n", len(fileNameBytes))

	fileSize := int64(len(data))
	if err := binary.Write(conn, binary.LittleEndian, fileSize); err != nil {
		return err
	}
	fmt.Printf("Written %d bytes of file size\n", 8)

	if err := binary.Write(conn, binary.LittleEndian, hashLen); err != nil {
		return err
	}
	fmt.Printf("Written %d bytes of hash length\n", 4)

	if _, err := conn.Write(hash[:]); err != nil {
		return err
	}
	fmt.Printf("Written %d bytes of hash data\n", hashLen)

	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	fmt.Printf("Written %d bytes of file data\n", n)

	return nil
}

func main() {
	size := 4000
	file := make([]byte, size)
	for i := 0; i < size; i++ {
		file[i] = byte((i % 10) + '0') // Generate '0' to '9' repeatedly
	}

	go func() {
		time.Sleep(4 * time.Second)
		if err := sendFile("numbers.txt", file); err != nil {
			log.Fatal("Error sending file:", err)
		}
	}()

	go func() {
		time.Sleep(15 * time.Second)
		if err := sendFile("numbers.txt", file); err != nil {
			log.Fatal("Error sending file:", err)
		}
	}()

	server := &FileServer{}
	server.start()
}
