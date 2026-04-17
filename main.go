package main

func main() {
	Server := NewServer("127.0.0.1", 6666)
	Server.Start()
}
