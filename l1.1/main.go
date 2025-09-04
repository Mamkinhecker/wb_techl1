package main

import (
	"fmt"
	"time"
)

type Human struct {
	Name    string
	Surname string
}

func (h *Human) Sleep() {
	fmt.Println("I am sleeping right now!")
}

func (h *Human) Coding() {
	fmt.Println("starting PC")
	time.Sleep(2 * time.Second)
	fmt.Println("checking logs")
	time.Sleep(2 * time.Second)
	fmt.Println("eat")
}

type Action struct {
	Human
	ActionType string
}

func main() {
	fmt.Println("what are you doing right now?")

	var action string
	fmt.Scanf("%s\n", &action)

	new_action := Action{
		Human: Human{
			Name:    "alexey",
			Surname: "pavlovich",
		},
		ActionType: action,
	}

	if new_action.ActionType == "sleep" {
		new_action.Sleep()
	} else {
		new_action.Coding()
	}
}
