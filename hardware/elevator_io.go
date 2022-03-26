package hardware

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const _pollRate = 20 * time.Millisecond

var _initialized bool = false
var _numFloors int = 4
var _mtx sync.Mutex
var _conn net.Conn

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down MotorDirection = -1
	MD_Stop MotorDirection = 0
)

type DoorAction int

const (
	DS_Close DoorAction = iota
	DS_Open_Up
	DS_Open_Down
	DS_Open_Cab
	DS_Do_Nothing
)

type ButtonType int

const (
	BT_HallUp ButtonType = iota
	BT_HallDown
	BT_Cab
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

func Init(
	addr string,
	numFloors int) {

	if _initialized {
		fmt.Println("Driver already initialized!")
		return
	}
	_numFloors = numFloors
	_mtx = sync.Mutex{}
	var err error
	_conn, err = net.Dial(
		"tcp",
		addr)
	if err != nil {
		panic(err.Error())
	}
	TurnOffAllLamps()
	_initialized = true
}

func TurnOffAllLamps() {
	for f := 0; f < FloorCount; f++ {
		for b := ButtonType(0); b < 3; b++ {
			SetButtonLamp(
				b,
				f,
				false)
		}
	}
	SetDoorOpenLamp(false)
	SetStopLamp(false)

}

func SetMotorDirection(dir MotorDirection) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write(
		[]byte{1,
			byte(dir),
			0,
			0})
}

func SetButtonLamp(
	button ButtonType,
	floor int,
	value bool) {

	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write(
		[]byte{2,
			byte(button),
			byte(floor),
			toByte(value)})
}

func SetFloorIndicator(floor int) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write(
		[]byte{3,
			byte(floor),
			0,
			0})
}

func SetDoorOpenLamp(value bool) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write(
		[]byte{4,
			toByte(value),
			0,
			0})
}

func SetStopLamp(value bool) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write(
		[]byte{5,
			toByte(value),
			0,
			0})
}

func PollButtons(receiver chan<- ButtonEvent) {
	prev := make(
		[][3]bool,
		_numFloors)
	for {
		time.Sleep(_pollRate)
		for f := 0; f < _numFloors; f++ {
			for b := ButtonType(0); b < 3; b++ {
				v := getButton(
					b,
					f)
				if v != prev[f][b] &&
					v {
					receiver <- ButtonEvent{f,
						ButtonType(b)}
				}
				prev[f][b] = v
			}
		}
	}
}

func PollFloorSensor(
	arrival_receiver chan<- int,
	leave_reciever chan<- bool) {

	prev := -1
	for {
		time.Sleep(_pollRate)
		v := getFloor()
		if v != prev {
			if v == -1 {
				leave_reciever <- true
			} else {
				arrival_receiver <- v
			}
		}
		prev = v
	}
}

func PollStopButton(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := getStop()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

func PollObstructionSwitch(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := getObstruction()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

func getButton(
	button ButtonType,
	floor int) bool {

	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write(
		[]byte{6,
			byte(button),
			byte(floor),
			0})
	var buf [4]byte
	_conn.Read(buf[:])
	return toBool(buf[1])
}

func getFloor() int {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write(
		[]byte{7,
			0,
			0,
			0})
	var buf [4]byte
	_conn.Read(buf[:])
	if buf[1] != 0 {
		return int(buf[2])
	} else {
		return -1
	}
}

func getStop() bool {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write(
		[]byte{8,
			0,
			0,
			0})
	var buf [4]byte
	_conn.Read(buf[:])
	return toBool(buf[1])
}

func getObstruction() bool {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write(
		[]byte{9,
			0,
			0,
			0})
	var buf [4]byte
	_conn.Read(buf[:])
	return toBool(buf[1])
}

func toByte(a bool) byte {
	var b byte = 0
	if a {
		b = 1
	}
	return b
}

func toBool(a byte) bool {
	var b bool = false
	if a != 0 {
		b = true
	}
	return b
}
