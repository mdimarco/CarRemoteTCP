package car

import "time"

const MAX_SPEED = 50 // Max speed car can get at full pedal
const DAMP = 1.2 // Dampener for acceleration (to avoid reaching max too fast)

type Car struct {
	Pedal float64
	Speed float64
	LastModified time.Time
}


// Gets acceleration based on pedal (accel coming from car)
// and wind resistance from speed, dampen to decrease rate
func getAcceleration(pedal float64, speed float64) float64{
	return (MAX_SPEED * pedal/100 - speed) / DAMP
}

// Updates speed of car based on time elapsed, previous speed, and car's pedal
// Returns speed
func (car *Car) GetSpeed() float64 {
	dur := time.Since(car.LastModified)
	seconds := time.Duration.Seconds(dur)

	//Update car's speed
	car.Speed += getAcceleration(car.Pedal, car.Speed) * seconds
	if car.Speed < 0 {
		car.Speed = 0
	}
	car.LastModified = time.Now()
	return car.Speed
}