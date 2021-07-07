package main

func main() {
	a := App{}
	a.Initialize("avm", "testing123", "task_organizer")
	a.Run(":8082")
}
