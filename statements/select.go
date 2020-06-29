package statements

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func stdin(input chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input <- scanner.Text()
	}
}

func Select1() {

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			c1 <- "one"
		}
	}()
	go func() {
		for {
			time.Sleep(10 * time.Second)
			c2 <- "two"
		}
	}()

	input := make(chan string)
	go stdin(input)
	fmt.Println("Any string to quit..")
mainLoop:
	for {
		fmt.Println("start select------------------")
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		case <-input:
			break mainLoop
			// default statement makes the select to be non-block operation.
			// default:
			// 	fmt.Println("no message")
		}
		fmt.Println("end select-------------------\n\n")
	}
}
