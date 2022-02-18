package cabstate

type ElevatorBehaviour int

const (
	Idle ElevatorBehaviour = iota
	DoorOpen
	Moving
)

type Direction int

const (
	Up Direction = iota
	Down
)

type CabState struct {
	aboveOrAtFloor int
	betweenFloors  bool
	doorOpen       bool
	motorRunning   bool
	motorDirection Direction
	behaviour      ElevatorBehaviour
}

var Cab CabState

func InitCabState() {
	Cab := new(CabState)
	_ = Cab
}

func FSMNewOrder(orderFloor int) {
	switch Cab.behaviour {
	case Idle:

	}
}
