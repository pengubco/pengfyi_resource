package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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

func bfsWorker(workerID int, g *Graph, nodeChan chan int, visited *sync.Map, wg *sync.WaitGroup, totalScore *int64) {
	for nodeID := range nodeChan {
		node := g.Nodes[nodeID]
		time.Sleep(time.Duration(node.Latency) * time.Second)
		atomic.AddInt64(totalScore, int64(node.Score))
		for _, neighbor := range node.Neighbors {
			if _, found := visited.LoadOrStore(neighbor, true); found {
				continue
			}
			wg.Add(1)
			go func(n int) {
				nodeChan <- n
			}(neighbor)
		}
		wg.Done()
	}
}

func bfsConcurrent(graph *Graph, startNode int, numWorkers int) int {
	var wg sync.WaitGroup
	wg.Add(1)

	var visited sync.Map
	visited.Store(startNode, true)

	nodeChan := make(chan int, 100) // Reasonable buffer size
	nodeChan <- startNode
	var totalScore int64

	for i := range numWorkers {
		go bfsWorker(i, graph, nodeChan, &visited, &wg, &totalScore)
	}

	wg.Wait()
	close(nodeChan)
	return int(totalScore)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <graph_file> <num_workers>")
		fmt.Println("Example: go run main.go graph.txt 4")
		os.Exit(1)
	}

	filename := os.Args[1]
	numWorkers, err := strconv.Atoi(os.Args[2])
	if err != nil || numWorkers < 1 {
		fmt.Printf("Invalid number of workers: %s\n", os.Args[2])
		os.Exit(1)
	}

	fmt.Printf("Loading graph from %s...\n", filename)
	graph, err := loadGraph(filename)
	if err != nil {
		fmt.Printf("Error loading graph: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Graph loaded successfully with %d nodes\n\n", len(graph.Nodes))

	totalScore := bfsConcurrent(graph, 0, numWorkers)

	fmt.Printf("\n=== BFS Complete ===\n")
	fmt.Printf("Total score: %d\n", totalScore)
}
