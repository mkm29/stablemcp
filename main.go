package main

func main() {
	if err := cmd.run(); err != nil {
		// Handle the error appropriately
		panic(err)
	}
}
