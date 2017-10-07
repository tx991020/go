package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/shicky/gocelery"
)

// Run Celery Worker First!
// celery -A worker worker --loglevel=debug --without-heartbeat --without-mingle

func main() {

	// create broker and backend
	celeryBroker := gocelery.NewRedisCeleryBroker("localhost:6379", "")
	celeryBackend := gocelery.NewRedisCeleryBackend("localhost:6379", "")

	// AMQP example
	//celeryBroker := gocelery.NewAMQPCeleryBroker("amqp://")
	//celeryBackend := gocelery.NewAMQPCeleryBackend("amqp://")

	// create client
	celeryClient, _ := gocelery.NewCeleryClient(celeryBroker, celeryBackend, 0)

	arg1 := rand.Intn(10)
	arg2 := rand.Intn(10)

	asyncResult, err := celeryClient.Delay("worker.add", arg1, arg2)
	if err != nil {
		panic(err)
	}

	res, err := asyncResult.Get(10 * time.Second)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Result: %v of type: %v\n", res, reflect.TypeOf(res))
	}
}
