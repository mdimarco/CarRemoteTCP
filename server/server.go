package main

import (
	"fmt"
	"time"
	"net"
	"strconv"
	"strings"
	"CarRemoteTCP/server/car"
	"os"
)

func main() {
	port := "5001"
	if len(os.Args) > 1{
		_, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Error, Usage: ./server [port #], setting default port", port)
		} else {
			port = os.Args[1]
		}
	}

	theCar := car.Car{ 0,  0, time.Now() }
	go updateSpeedTask(&theCar)

	//Initialize TCP server
	initServer(&theCar, port)
}


func updateSpeedTask(theCar *car.Car) {
	for {
		theCar.GetSpeed()
	}
}


func initServer(theCar *car.Car, port string) {
	ln, err := net.Listen("tcp", ":"+port)
	fmt.Println("Listening on port:", port)
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn, theCar)
	}
}

func handleConnection(conn net.Conn, theCar *car.Car) {

	//Read in request and trim after new line
	message := make([]byte, 1024)
	_, err := conn.Read(message)
	parsed := strings.Split(string(message), "\n")[0]
	if err != nil {
		fmt.Println(err)
	}

	if len(parsed) < 9 {
		fmt.Println("Error, incorrect command. Usage: SET_SPEED [0-100] or GET_SPEED ")
	}
	//Parse and set car's pedal
	if string(parsed[0:9]) == "SET_PEDAL" {
		newPedal, err := strconv.ParseFloat(parsed[10:], 64)
		if err != nil || newPedal < 0 || newPedal > 100 {
			fmt.Println("Pedal out of bounds or Error parsing", err)
		}
		theCar.Pedal = newPedal

	}

	//Format and send car's speed
	speed := theCar.GetSpeed()
	s64 := strconv.FormatFloat(speed, 'E', 10, 64)
	fmt.Println("Current car speed",  s64, "pedal", theCar.Pedal)
	conn.Write( append([]byte(s64), '\n') )	
	conn.Close()


}