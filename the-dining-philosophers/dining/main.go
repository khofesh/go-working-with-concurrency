package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

/**
 * The Dining Philosophers problem is well known in computer science circles.
 * Five philosophers, numbered from 0 through 4, live in a house where the
 * table is laid for them; each philosopher has their own place at the table.
 * Their only difficulty – besides those of philosophy – is that the dish
 * served is a very difficult kind of spaghetti which has to be eaten with
 * two forks. There are two forks next to each plate, so that presents no
 * difficulty. As a consequence, however, this means that no two neighbours
 * may be eating simultaneously, since there are five philosophers and five forks.
 *
 * This is a simple implementation of Dijkstra's solution to the "Dining
 * Philosophers" dilemma.
 *
 * Philosopher is a struct which stores information about a philosopher.
 */

// stores information about a philosopher
type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// list of philosophers
var philosophers = []Philosopher{
	{
		name:      "Socrates",
		leftFork:  4,
		rightFork: 0,
	},
	{
		name:      "Plato",
		leftFork:  0,
		rightFork: 1,
	},
	{
		name:      "Aristotle",
		leftFork:  1,
		rightFork: 2,
	},
	{
		name:      "Zeno",
		leftFork:  2,
		rightFork: 3,
	},
	{
		name:      "Epictetus",
		leftFork:  3,
		rightFork: 4,
	},
}

// how many times does a person eat ?
var hunger = 3
var eatTime = 1 * time.Second
var thinkTime = 3 * time.Second
var sleepTime = 1 * time.Second

// a mutex for the slice orderFinished;
var orderMutex sync.Mutex

// the order in which philosophers finish dining and leave;
var orderFinished []string

func main() {
	fmt.Println("dining philosophers problem")
	fmt.Println("---------------------------")
	fmt.Println("the table is empty.")

	// pause a little bit
	time.Sleep(sleepTime)

	// start the meal
	dine()

	// finished message
	fmt.Println("the table is empty.")

	time.Sleep(sleepTime)
	fmt.Printf("order finished: %s.\n", strings.Join(orderFinished, ", "))
}

func dine() {
	// eatTime = 0 * time.Second
	// sleepTime = 0 * time.Second
	// thinkTime = 0 * time.Second

	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// a map of all 5 forks
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// start the meal
	for i := 0; i < len(philosophers); i++ {
		// goroutine for the current philosopher
		go diningProblem(philosophers[i], wg, forks, seated)
	}

	wg.Wait()
}

/**
 * diningProblem is the function fired off as a goroutine for each of our philosophers. It takes one
 * philosopher, our WaitGroup to determine when everyone is done, a map containing the mutexes for every
 * fork on the table, and a WaitGroup used to pause execution of every instance of this goroutine
 * until everyone is seated at the table.
 */
func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()

	// seat the philosopher at the table
	fmt.Printf("%s is seated at the table.\n", philosopher.name)
	seated.Done()

	seated.Wait()

	// eat three times
	for i := hunger; i > 0; i-- {
		// get a lock on both forks
		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", philosopher.name)
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", philosopher.name)
		}

		fmt.Printf("\t%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(eatTime)

		fmt.Printf("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("\t%s put down the forks.\n", philosopher.name)
	}

	fmt.Println(philosopher.name, "is satisfied.")
	fmt.Println(philosopher.name, "left the table.")

	orderMutex.Lock()
	orderFinished = append(orderFinished, philosopher.name)
	orderMutex.Unlock()
}
