package main

import (
	"elevators/controlunit/orderstate"
	"fmt"
	"time"
)

func hasOrder(inputState orderstate.OrderState) bool {
	//TODO: test (Possible riv ruskende here)
	return inputState.LastOrderTime.After(inputState.LastCompleteTime)
}

func testHasOrder() {
	fmt.Println("Testing hasOrder")
	fmt.Println("")

	// var time1 = time.Time{}
	// var time2 = time.Time{}
	orders := []orderstate.OrderState{
		orderstate.OrderState{},
		orderstate.OrderState{LastOrderTime: time.Time{}, LastCompleteTime: time.Time{}},
		orderstate.OrderState{LastOrderTime: time.Now(), LastCompleteTime: time.Time{}},
		orderstate.OrderState{LastOrderTime: time.Now(), LastCompleteTime: time.Now()}}
	for _, order := range orders {
		fmt.Print("Last order: ")
		fmt.Print(order.LastOrderTime)
		fmt.Print(", Last Complete: ")
		fmt.Print(order.LastCompleteTime)
		fmt.Print(" ; hasOrder : ")
		fmt.Println(hasOrder(order))
		fmt.Println("")
	}
}

func main() {
	testHasOrder()
}
