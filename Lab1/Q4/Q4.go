package main

import (
	"fmt"
	"math/rand"
)

// create a board with random mines position
func createBoard(width int, height int, mines int) [][]string {
	board := make([][]string, height+1)
	minesCoordinates := make([][]int, mines)
	for i := range board {
		board[i] = make([]string, width+1)
	}
	for i := range minesCoordinates {
		minesCoordinates[i] = make([]int, 2)
	}
	for i := range board {
		for j := range board[i] {
			board[i][j] = ". "
		}
	}
	for i := 0; i < mines; i++ {
		x := rand.Intn(width) + 1
		y := rand.Intn(height) + 1
		for j := 0; j < i; j++ {
			if minesCoordinates[j][0] == x && minesCoordinates[j][1] == y {
				x = rand.Intn(width) + 1
				y = rand.Intn(height) + 1
				j = 0
			}
		}
		minesCoordinates[i][0] = x
		minesCoordinates[i][1] = y
		board[x][y] = "* "
	}
	return board
}

var neighbourOffset = [8][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func checkNeighbour(x, y int, board [][]string) int {
	count := 0
	for i := range neighbourOffset {
		posX := x + neighbourOffset[i][0]
		posY := y + neighbourOffset[i][1]
		if posX >= 0 && posX < len(board) && posY >= 0 && posY < len(board) {
			if board[posX][posY] == "* " {
				count++
			}
		}
	}
	return count
}

func displayBoard(board [][]string) {
	for i := range board {
		for j := range board[i] {
			if board[i][j] == ". " {
				if checkNeighbour(i, j, board) == 0 {
					board[i][j] = ". "
				} else {
					board[i][j] = fmt.Sprintf("%d ", checkNeighbour(i, j, board))
				}
			}
		}
	}
	for i := range board {
		fmt.Println(board[i])
	}
}
