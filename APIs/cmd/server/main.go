package main

import "github.com/valdir-alves3000/postgraduate-challenges-go-expert/APIs/configs"

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
