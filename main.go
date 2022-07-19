package main

import (
	"math/rand"
	"time"

	"github.com/kbrgl/flapioca/internal"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	internal.Execute()
}
