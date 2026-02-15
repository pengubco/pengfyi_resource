package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	ID        int
	Score     int
	Latency   int // in seconds
	Neighbors []int
}

type Graph struct {
	Nodes map[int]*Node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[int]*Node),
	}
}

func (g *Graph) AddNode(id, score, latency int, neighbors []int) {
	g.Nodes[id] = &Node{
		ID:        id,
		Score:     score,
		Latency:   latency,
		Neighbors: neighbors,
	}
}

func loadGraph(filename string) (*Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	graph := NewGraph()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		nodeID, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid node ID: %s", parts[0])
		}

		score, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid score: %s", parts[1])
		}

		latency, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid latency: %s", parts[2])
		}

		neighbors := make([]int, 0, len(parts)-3)
		for i := 3; i < len(parts); i++ {
			neighbor, err := strconv.Atoi(parts[i])
			if err != nil {
				return nil, fmt.Errorf("invalid neighbor ID: %s", parts[i])
			}
			neighbors = append(neighbors, neighbor)
		}

		graph.AddNode(nodeID, score, latency, neighbors)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return graph, nil
}

func bfs(graph *Graph, startNode int) int {
	if _, exists := graph.Nodes[startNode]; !exists {
		fmt.Printf("Start node %d does not exist in graph\n", startNode)
		return 0
	}

	visited := make(map[int]bool)
	queue := []int{startNode}
	visited[startNode] = true
	totalScore := 0

	fmt.Printf("BFS traversal starting from node %d:\n", startNode)
	startTime := time.Now()

	for len(queue) > 0 {
		currentID := queue[0]
		queue = queue[1:]

		current := graph.Nodes[currentID]

		// Sleep to simulate node visit latency
		fmt.Printf("Visiting node %d (score: %d, latency: %ds)...\n",
			current.ID, current.Score, current.Latency)
		time.Sleep(time.Duration(current.Latency) * time.Second)

		totalScore += current.Score
		fmt.Printf("  Completed node %d (total score so far: %d)\n",
			current.ID, totalScore)

		for _, neighborID := range current.Neighbors {
			if !visited[neighborID] {
				visited[neighborID] = true
				queue = append(queue, neighborID)
			}
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nTotal time elapsed: %v\n", elapsed)

	return totalScore
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <graph_file>")
		fmt.Println("Example: go run main.go graph.txt")
		os.Exit(1)
	}

	filename := os.Args[1]

	fmt.Printf("Loading graph from %s...\n", filename)
	graph, err := loadGraph(filename)
	if err != nil {
		fmt.Printf("Error loading graph: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Graph loaded successfully with %d nodes\n\n", len(graph.Nodes))

	totalScore := bfs(graph, 0)

	fmt.Printf("\n=== BFS Complete ===\n")
	fmt.Printf("Total score: %d\n", totalScore)
}
