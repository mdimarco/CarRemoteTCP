package main

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"time"
)

const THRESHOLD = .01

var PORT string

func main() {
	if len(os.Args) != 3 {
		printUsage()
	} else {
		PORT = os.Args[1]
		speed, _ := strconv.Atoi(os.Args[2])
		upToSpeed(float64(speed))
	}
}

func printUsage() {
	fmt.Println("Error, Usage: ./client portNumber speed")
}

func done(targetSpeed float64) bool {
	return math.Abs(float64(targetSpeed)-getCurrentSpeed()) < THRESHOLD
}

/**
 * Gets the car up to target speed, adjusting pedal in a binary-search-like fashion
 * using approximations to deal with floating points and nature of requests.
 * @param  targetSpeed float64      speed to hopefully end up at
 */
func upToSpeed(targetSpeed float64) {
	start, middle, end := 0.0, 50.0, 100.0
	pedal := middle
	setNewPedal(pedal)
	waitUntilStagnate(targetSpeed)
	for !done(targetSpeed) {
		if getCurrentSpeed() < float64(targetSpeed) {
			pedal = (middle + end) / 2
			start = middle
		} else {
			pedal = (start + middle) / 2
			end = middle
		}
		middle = (start + end) / 2
		setNewPedal(pedal)
		waitUntilStagnate(targetSpeed)
	}
	fmt.Println("Done! current speed:", getCurrentSpeed(), "Pedal at:", pedal)
}

//Give the server a small break
func pause() {
	time.Sleep(500 * time.Millisecond)
}

/**
 * Waits until speed of car evens out or goes over/under the targetted limit
 * this function exiting indicates it is time to change the pedal in order to
 * attain a higher/lower speed
 * @param   targetSpeed int           speed of car we want
 * @param   breakCond   string        "OVER" | "UNDER" type of break condition
 */
func waitUntilStagnate(targetSpeed float64) {
	speed1 := getCurrentSpeed()
	pause()
	speed2 := getCurrentSpeed()
	for math.Abs(speed2-speed1) > .01 {
		speed1 = speed2
		pause()
		speed2 = getCurrentSpeed()

		//If target passed, break out so that don't accelerate the wrong way
		if (speed2-speed1 > 0) && speed2 > float64(targetSpeed) ||
			(speed2-speed1 < 0) && speed2 < float64(targetSpeed) {
			break
		}
	}
}

//Sends GET_SPEED request to server
func getCurrentSpeed() float64 {
	return sendMessage(-1.0)
}

//Sends SET_PEDAL request to server
func setNewPedal(pedal float64) float64 {
	return sendMessage(pedal)
}

//Actual TCP message sending
func sendMessage(pedal float64) float64 {
	conn, err := net.Dial("tcp", "127.0.0.1:"+PORT)
	if err != nil {
		print("Error establishing connection", err)
	}
	if pedal >= 0.0 {
		s64 := strconv.FormatFloat(pedal, 'E', 10, 64)
		fmt.Fprintf(conn, "SET_PEDAL "+s64+"\n")
	} else {
		fmt.Fprintf(conn, "GET_SPEED\n")
	}
	speed, _ := bufio.NewReader(conn).ReadString('\n')
	speedFloat, _ := strconv.ParseFloat(speed[0:len(speed)-1], 64)
	return speedFloat
}
