package main

import (
	"fmt"
	"time"
)

type User interface {
	PrintName()
	PrintDetails()
}

type Person struct {
	FirstName, LastName string
	Dob                 time.Time
	Email, Location     string
}

//A person method
func (p Person) PrintName() {
	fmt.Printf("\n%s %s\n", p.FirstName, p.LastName)
}

//A person method
func (p Person) PrintDetails() {

	fmt.Printf("[Date of Birth: %s, Email: %s, Location: %s ]\n", p.Dob.String(), p.Email, p.Location)
}

func main() {

	alex := Person{"Alex", "John", time.Date(1970, time.January, 10, 0, 0, 0, 0, time.UTC), "alex@email.com", "New York"}
	shiju := Person{"Shiju", "Varghese", time.Date(1979, time.February, 17, 0, 0, 0, 0, time.UTC), "shiju@email.com", "Kochi"}
	users := []User{alex, shiju}
	for _, v := range users {

		v.PrintName()

		v.PrintDetails()
	}

}
