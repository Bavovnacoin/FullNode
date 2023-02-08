package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Struct struct {
	name string
	age  uint
}

type P struct {
	X, Y, Z int
	Name    string
}

// type Q struct {
// 	X, Y *int32
// 	Name string
// }

func main() {
	//testing.Test1()
	// var db dbController.DataBase
	// println(db.OpenDb())
	// println(db.CloseDb())

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.
	// Encode (send) the value.
	err := enc.Encode(P{3, 2, 5, "Pythagoras"})
	if err != nil {
		log.Fatal("encode error:", err)
	}

	// HERE ARE YOUR BYTES!!!!
	fmt.Println(network.Bytes())

	// Decode (receive) the value.
	var p P
	err = dec.Decode(&p)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	fmt.Printf("%q: {%d,%d}\n", p.Name, p.X, p.Y)
	//db.Db.Put()
}
