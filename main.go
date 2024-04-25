package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"sync"
	"time"
)

type Circuit struct {
	ID   int
	Name string
}

var (
	circuits      = make(map[int]*Circuit)
	circuitsMutex sync.RWMutex
)

func main() {
	deepCopyMode := flag.Bool("deep", false, "set to true to use a deep copy of the map")
	updateNameOnly := flag.Bool("name", false, "set to true to only update the name of the circuit")
	flag.Parse()

	// add a bunch of test circuits
	_ = gofakeit.Seed(0)
	for i := 0; i <= 9999; i++ {
		circuits[i] = &Circuit{ID: i, Name: fakeCircuitName()}
	}
	switch *updateNameOnly {
	case false:
		fmt.Println("Replacing entire struct")
		go updaterReplaceStruct()
	case true:
		fmt.Println("Replacing name within struct")
		go updaterNameReplaceOnly()
	}

	count := 1
	switch *deepCopyMode {
	case false:
		fmt.Println("Using shallow map copy")
		fmt.Println("  Count     Size\n----------  -----")
		for {
			// get all circuits
			tempUsers := getAllCircuits()
			if _, err := json.Marshal(tempUsers); err != nil {
				panic(err)
			}
			fmt.Printf("%10d %5d\r", count, len(tempUsers))
			count++
		}

	case true:
		fmt.Println("Using deep map copy")
		fmt.Println("  Count     Size\n----------  -----")
		for {
			// get all circuits
			tempUsers := getAllCircuitsCopy()
			if _, err := json.Marshal(tempUsers); err != nil {
				panic(err)
			}
			fmt.Printf("%10d %5d\r", count, len(tempUsers))
			count++
		}
	}
}

func fakeCircuitName() string {
	return fmt.Sprintf("%s-%08d", gofakeit.Color(), gofakeit.Number(1, 99999999))
}

func getAllCircuits() map[int]*Circuit {
	circuitsMutex.RLock()
	defer circuitsMutex.RUnlock()
	return circuits
}

func getAllCircuitsCopy() map[int]*Circuit {
	circuitsMutex.RLock()
	defer circuitsMutex.RUnlock()
	tempCircuits := make(map[int]*Circuit, len(circuits))
	for k, v := range circuits {
		tempCircuits[k] = v
	}
	return tempCircuits
}

// updaterNameReplaceOnly randomly updates the circuits map/cache by replacing the name of the circuit with a new random name
func updaterNameReplaceOnly() {
	for {
		time.Sleep(time.Millisecond)
		circuitsMutex.Lock()
		// pick random circuit
		i := gofakeit.Number(0, 9999)
		// Get pointer to struct from map/cache
		x := circuits[i]
		// unlock map/cache as we got our pointer out
		circuitsMutex.Unlock()
		// update Name in struct
		x.Name = fakeCircuitName()
	}
}

// updaterReplaceStruct randomly updates the circuits map/cache by replacing the entire Circuit struct with a new one
func updaterReplaceStruct() {
	for {
		time.Sleep(time.Millisecond)
		circuitsMutex.Lock()
		// pick random circuit
		i := gofakeit.Number(0, 9999)
		// update map/cache by replacing the pointer in the map to a new struct
		circuits[i] = &Circuit{ID: i, Name: fakeCircuitName()}
		circuitsMutex.Unlock()
	}
}
