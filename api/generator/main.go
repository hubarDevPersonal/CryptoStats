package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go generate(&wg)

	wg.Wait()
}

func generate(wg *sync.WaitGroup) {
	defer wg.Done()

	runGoa()
}

func runGoa() {
	start := time.Now()
	log.Printf("Running goa generate for api")

	cmd := exec.Command("goa", "gen", fmt.Sprintf("CryptoStats/api/design"))
	cmd.Dir = fmt.Sprintf("./")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("error running goa generate for api (but that is probably ok): %v", err)
	}
	log.Printf("Finished goa generate for api. Time taken: %s", time.Since(start))
}
