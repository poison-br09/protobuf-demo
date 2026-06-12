package main

import (
	"fmt"
	"log"

	"protobuf-demo/generated"

	"google.golang.org/protobuf/proto"
)

func main() {
	user := &generated.User{
		Id:    1,
		Name:  "Abhay",
		Email: "abhay@zenwork.com",
	}

	data, err := proto.Marshal(user)
	if err != nil {
		log.Fatal("Marshaling Error:", err)
	}

	fmt.Println("Binary bytes:", data)
	fmt.Println("Size in bytes:", len(data))

	decoded := &generated.User{}
	err = proto.Unmarshal(data, decoded)
	if err != nil {
		log.Fatal("Unmarshaling error:", err)
	}

	fmt.Println("Decoded User:", decoded)
}
