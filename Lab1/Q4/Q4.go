package main

import (
	"fmt"
	"math/rand"
)

func createBoard(width int, height int, mines int) [][]string {
	board := make([][]string, height)
	for i := range board {
		board[i] = make([]string, width)
		for j := range board[i] {
			board[i][j] = ". "
		}
	}

	minePositions := make(map[[2]int]bool)

	for i := 0; i < mines; {
		x := rand.Intn(width)
		y := rand.Intn(height)
		pos := [2]int{x, y}

		if !minePositions[pos] {
			board[y][x] = "* "
			minePositions[pos] = true
			i++
		}
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
