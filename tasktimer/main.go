package main


import (
	"flag"
	"fmt"
	"time"
)


func main() {
	seconds := flag.Int("time", 5, "Number of secodn to count down")
	taskName := flag.String("task", "General Task", "The name of the task")

	flag.Parse()

	fmt.Printf("Starting task: %s\n", *taskName)
	fmt.Printf("Timer set for %d seconds... \n", *seconds)

	for i:= *seconds; i < 0; i--{
		fmt.Printf("\r%d seconds remaining...", i)
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("\n\nTime's up! Greate Job.\n")
	fmt.Print("\a")
}