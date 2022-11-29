package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients.", barber)

		for {
			// if there are no clients, the barber goes to sleep
			if len(shop.ClientsChan) == 0 {
				color.Yellow("there is nothing to do, so %s takes a nap.", barber)
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientsChan
			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up.", client, barber)
					isSleeping = false
				}

				// cut hair
				shop.cutHair(barber, client)
			} else {
				// shop is closed, so send the barber home and close this gorouting
				shop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *BarberShop) cutHair(barber, client string) {
	color.Green("%s is cutting %s's hair.", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finished cutting %s's hair.", barber, client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home.", barber)
	shop.BarbersDoneChan <- true
}

func (shop *BarberShop) closeShopForDay() {
	color.Cyan("closing shop for the day.")

	close(shop.ClientsChan)
	shop.Open = false

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarbersDoneChan
	}

	close(shop.BarbersDoneChan)

	color.Green("---------------------------------------------------------------------")
	color.Green("the barbershop is now closed for the day, and everyone has gone home.")
}

func (shop *BarberShop) addClient(client string) {
	// randome message
	color.Green("*** %s arrives", client)

	if shop.Open {
		select {
		case shop.ClientsChan <- client:
			color.Yellow("%s takes a seat in the waiting room.", client)
		default:
			color.Red("the waiting room is full, so %s leaves.", client)
		}
	} else {
		color.Red("the shop is already closed, so %s leaves!", client)
	}
}
