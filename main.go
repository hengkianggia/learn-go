package main

import (
	_ "learn/migrations" // Import all migrations
	"learn/cmd"
)

func main() {
	cmd.Execute()
}
