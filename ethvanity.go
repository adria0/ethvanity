package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func main() {

	if len(os.Args) != 3 {
		log.Fatal("Usage ethvanity <concurrency> <addressprefix>")
	}

	concurrency, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Bad concurrency parameter")
	}
	prefix := os.Args[2]

	var combinations = math.Pow(16, float64(len(prefix)))

	log.Println("Generating key with address starting with", prefix, "using", concurrency, "threads")
	log.Println("Total combinations are", uint64(combinations))

	var counter uint64 = 0
	var finished = make(chan bool)

	for i := 0; i < concurrency; i++ {

		go func() {

			for {

				atomic.AddUint64(&counter, 1)
				if counter%25000 == 0 {
					current := (float64(counter) * 100) / combinations
					log.Printf("%v (%.3f %%)", counter, current)
				}

				k, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
				if err != nil {
					log.Fatal(err)
				}

				address := crypto.PubkeyToAddress(k.PublicKey)
				addressStr := hex.EncodeToString(address[:])

				if strings.HasPrefix(addressStr, os.Args[2]) {

					pvkStr := hex.EncodeToString(crypto.FromECDSA(k))
					log.Println("\n\nFOUND address=[", addressStr, "] pvk=[", pvkStr, "]")

					finished <- true

				}
			}
		}()
	}

	<-finished
}
