package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Pipeline service started")
	fmt.Println("Service: pipeline")
	fmt.Println("Status: running")
	fmt.Println("Message: Hello from Pipeline (Data Pipeline)")

	// Mantém o serviço rodando
	for {
		time.Sleep(10 * time.Second)
		fmt.Println("Pipeline heartbeat:", time.Now().Format("2006-01-02 15:04:05"))
	}
}


