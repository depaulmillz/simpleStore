package main

import (
	"fmt"
	"math/rand"
	"simplestoreclient"
	"strconv"
	"time"
)

func main() {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	ctx := simplestoreclient.NewClientctx("localhost:8080")
	initial := 1000
	count := 0

	for count < initial {
		if ctx.Put(strconv.Itoa(seededRand.Int()), "2") != "" {
			count++
		}
	}

	start := time.Now()

	ops := 10000

	for i := 0; i < ops; i++ {
		ctx.Get(strconv.Itoa(seededRand.Int()))
	}
	end := time.Now()

	elapsed := end.Sub(start)

	fmt.Printf("Throughput %f\n", (float64(ops))/elapsed.Seconds())
	fmt.Printf("Avg Latency %f (ms)\n", (elapsed.Seconds()*1e3)/(float64(ops)))

	ctx.EndClientCtx()
}

func checkerror(err error) {
	if err != nil {
		fmt.Println("Error ", err.Error())
		panic(err)
	}
}
