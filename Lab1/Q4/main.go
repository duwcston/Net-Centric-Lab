package main

import (
	"fmt"
)

func main() {
	var width, height, mines int
	fmt.Println("Q4: Minesweeper")
	fmt.Print("Enter the width of the board: ")
	fmt.Scan(&width)
	fmt.Print("Enter the height of the board: ")
	fmt.Scan(&height)
	fmt.Print("Enter the number of mines: ")
	fmt.Scan(&mines)
	board := createBoard(width, height, mines)
	displayBoard(board)
	fmt.Println("====================================")
}
