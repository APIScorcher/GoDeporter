package main

import (
    "fmt"
    "net"
    "sort"
    "os"
    "flag"
)

func usage() {
    fmt.Println(os.Stderr, "usage: %s [URL]\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(2)
}

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf(os.Args[1] + ":%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	flag.Usage = usage
	ports := make(chan int, 10000)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("Port: %d is Open!\n", port)
	}
}