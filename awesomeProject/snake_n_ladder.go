package main

import "awesomeProject/snake_n_ladder/engine"

func snake_n_ladder() {

	//gameEngine := engine.InitEngine(10, 10, 200)
	gameEngine := engine.InitEngine(5, 5, 50)
	gameEngine.AddPlayer("Naruto")
	gameEngine.AddPlayer("Sasuke")
	//gameEngine.AddPlayer("Sakura")
	//gameEngine.AddPlayer("Negi")
	gameEngine.Play()

}
