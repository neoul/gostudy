package samples

import (
	"fmt"
	"time"
)

func Switch1() {
	i := "korea"
	switch i {
	case "korea":
		fmt.Println("korea")
	case "usa":
		fmt.Println("usa")
	case "japan":
		fmt.Println("japan")
	}
}

func Switch2() {
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("It's before noon")
	default:
		fmt.Println("It's after noon")
	}
}

func Switch3(c rune) bool {
	switch c {
	case ' ', '\t', '\n', '\f', '\r':
		return true
	}
	return false
}

func Switch4_goto() {
	/* local variable definition */
	var a int = 10

	/* do loop execution */
LOOP:
	for a < 15 {
		if a == 12 {
			/* skip the iteration */
			a = a + 1
			goto LOOP
		}
		fmt.Printf("value of a: %d\n", a)
		a++
	}
}
