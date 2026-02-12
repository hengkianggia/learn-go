package main

import (
	"learn/cmd"
	_ "learn/migrations" // Import all migrations
)

func main() {
	cmd.Execute()
}
