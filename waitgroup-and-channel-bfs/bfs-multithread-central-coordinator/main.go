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

func bfsWorker(workerID int, graph *Graph, toVisitChan <-chan int, visitedChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d started\n", workerID)

	for nodeID := range toVisitChan {
		node := graph.Nodes[nodeID]

		// Sleep to simulate node visit latency
		fmt.Printf("Worker %d: Visiting node %d (score: %d, latency: %ds)...\n",
			workerID, node.ID, node.Score, node.Latency)
		time.Sleep(time.Duration(node.Latency) * time.Second)

		fmt.Printf("Worker %d: Completed node %d\n", workerID, node.ID)

		// Send node ID back to coordinator
		visitedChan <- nodeID
	}

	fmt.Printf("Worker %d finished\n", workerID)
}

func bfsConcurrent(graph *Graph, startNode int, numWorkers int) int {
	if _, exists := graph.Nodes[startNode]; !exists {
		fmt.Printf("Start node %d does not exist in graph\n", startNode)
		return 0
	}

	visited := make(map[int]bool)
	visited[startNode] = true

	toVisitChan := make(chan int, 100)
	visitedChan := make(chan int, 100)

	var wg sync.WaitGroup
	var totalScore int64

	fmt.Printf("Starting concurrent BFS from node %d with %d workers\n", startNode, numWorkers)
	startTime := time.Now()

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go bfsWorker(i, graph, toVisitChan, visitedChan, &wg)
	}

	// Goroutine to close resultChan when all workers are done
	go func() {
		wg.Wait()
		close(visitedChan)
	}()

	// Send the start node
	toVisitChan <- startNode
	nodesInFlight := 1

	// Process results and send new nodes (COORDINATOR)
	for nodeID := range visitedChan {
		nodesInFlight--

		// Look up node data from graph
		node := graph.Nodes[nodeID]
		atomic.AddInt64(&totalScore, int64(node.Score))
		fmt.Printf("Coordinator: Processed node %d, total score: %d, nodes in flight: %d\n",
			nodeID, atomic.LoadInt64(&totalScore), nodesInFlight)

		// Add unvisited neighbors to the channel
		for _, neighborID := range node.Neighbors {
			if !visited[neighborID] {
				visited[neighborID] = true
				toVisitChan <- neighborID
				nodesInFlight++
				fmt.Printf("Coordinator: Added node %d to channel (nodes in flight: %d)\n", neighborID, nodesInFlight)
			}
		}

		// If no more nodes in flight, we're done
		if nodesInFlight == 0 {
			close(toVisitChan)
			break
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nTotal time elapsed: %v\n", elapsed)

	return int(atomic.LoadInt64(&totalScore))
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
