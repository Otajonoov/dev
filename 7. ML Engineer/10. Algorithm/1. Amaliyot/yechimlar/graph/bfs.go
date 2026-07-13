package main

import (
	"fmt"
)

func bfs(graph map[int][]int, start int) {
	visited := make(map[int]bool) // Har bir tugun ko‘rilganmi yoki yo‘q
	queue := []int{start}         // FIFO navbat

	for len(queue) > 0 {
		node := queue[0]  // Navbat boshini olamiz
		queue = queue[1:] // Uni navbatdan chiqaramiz

		if visited[node] {
			continue
		}

		fmt.Println("Visited:", node)
		visited[node] = true

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				queue = append(queue, neighbor) // Qo‘shnilarni navbatga qo‘shamiz
			}
		}
	}
}

func main() {
	graph := map[int][]int{
		0: {1, 2},
		1: {0, 3},
		2: {0, 3},
		3: {1, 2, 4},
		4: {},
	}

	fmt.Println("BFS traversal starting from node 0:")
	bfs(graph, 0)
}
