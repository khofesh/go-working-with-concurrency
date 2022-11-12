package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++

	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("received order #%d\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("making pizza #%d. It will take %d seconds...\n", pizzaNumber, delay)

		// delay
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** we ran out of ingredients for pizza #%d\n", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** the cook quit while making pizza #%d\n", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("pizza order #%d is ready\n", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

		return &p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	var i = 0

	// run forever or until we receive a quit notification
	// try making pizzas
	for {
		currentPizza := makePizza(i)

		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			// try making a pizza - sending something to data channel (a chan PizzaOrder)
			case pizzaMaker.data <- *currentPizza:

			// quit, send pizzaMaker.quit to the quitChan (a chan error)
			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// print out a message
	color.Cyan("the pizzeria is open for business!")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzeria(pizzaJob)

	// create and run consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("order #%d is out for delivery!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("the customer is really mad!")
			}
		} else {
			color.Cyan("done making pizzas")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("--- error closing channel", err)
			}
		}
	}

	// print out the ending message
	color.Cyan("-----")
	color.Cyan("done.")

	color.Cyan("we made %d pizzas, but failed to make %d, with %d attempts in total.",
		pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("awful day.")
	case pizzasFailed >= 6:
		color.Red("not a very good day.")
	case pizzasFailed >= 4:
		color.Yellow("it was an okay day.")
	case pizzasFailed >= 2:
		color.Yellow("it was a pretty good day")
	default:
		color.Green("it was a great day.")
	}
}
