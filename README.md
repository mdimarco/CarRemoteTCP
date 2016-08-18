
**Car Remote Control**

Client and server for a virtual car's speed. 

*Server* hosts the "car" and manages it's speed, taking in TCP requests to alter how far down the pedal is pressed (0-100). The maximum speed the car will settle at (at 100% pedal) is currently 50, as is described in car.go. The server accepts two types of requests:
	1. "GET_SPEED", which will respond with the current speed
	2. "SET_PEDAL" which will change the % the pedal is pressed on the car and then return the current speed

Usage: `./server [Port Number]`

*Client* repeatedly sends network requests to get the server up to speed, using a binary-search-like algorithm to narrow the window of correct pedal % for the given speed

Usage: `./client portNumber targetSpeed`