package main

func main() {
	store := NewCache()
	a := AppServer{
		Store: store,
		Host:  ":3030",
	}
	a.Run()
}
