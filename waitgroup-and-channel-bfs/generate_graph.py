#!/usr/bin/env python3
"""
Generate a random undirected connected graph for BFS demonstration.
Each node has a random score between 1-100 and visit latency between 1-10 seconds.
Output format: adjacency list with scores and latency in a text file.
"""

import argparse
import random
import sys


def generate_connected_graph(num_nodes, edge_probability=0.3):
    """
    Generate a random undirected connected graph using a modified approach.
    
    Args:
        num_nodes: Number of nodes in the graph
        edge_probability: Probability of creating an edge between nodes
    
    Returns:
        Dictionary representing adjacency list
    """
    if num_nodes < 1:
        raise ValueError("Number of nodes must be at least 1")
    
    # Initialize adjacency list
    graph = {i: set() for i in range(num_nodes)}
    
    # First, create a spanning tree to ensure connectivity
    # Start with node 0 and add nodes one by one
    unconnected = list(range(1, num_nodes))
    connected = [0]
    
    while unconnected:
        # Pick a random node from connected set
        from_node = random.choice(connected)
        # Pick a random node from unconnected set
        to_node = unconnected.pop(random.randint(0, len(unconnected) - 1))
        
        # Add edge
        graph[from_node].add(to_node)
        graph[to_node].add(from_node)
        connected.append(to_node)
    
    # Add additional random edges based on probability
    for i in range(num_nodes):
        for j in range(i + 1, num_nodes):
            # Skip if edge already exists
            if j in graph[i]:
                continue
            
            # Add edge with given probability
            if random.random() < edge_probability:
                graph[i].add(j)
                graph[j].add(i)
    
    return graph


def generate_visit_times(num_nodes):
    """Generate random scores (1-100) for each node."""
    return {i: random.randint(1, 100) for i in range(num_nodes)}


def generate_latencies(num_nodes):
    """Generate random visit latencies (1-10 seconds) for each node."""
    return {i: random.randint(1, 10) for i in range(num_nodes)}


def write_graph_to_file(graph, visit_times, latencies, filename):
    """
    Write graph to file in the following format:
    
    # Graph with N nodes
    # Format: node_id score latency_seconds neighbor1 neighbor2 ...
    0 42 5 1 2 3
    1 87 3 0 4
    ...
    """
    with open(filename, 'w') as f:
        f.write(f"# Graph with {len(graph)} nodes\n")
        f.write("# Format: node_id score latency_seconds neighbor1 neighbor2 ...\n")
        
        for node in sorted(graph.keys()):
            neighbors = sorted(graph[node])
            score = visit_times[node]
            latency = latencies[node]
            line = f"{node} {score} {latency}"
            if neighbors:
                line += " " + " ".join(map(str, neighbors))
            f.write(line + "\n")


def main():
    parser = argparse.ArgumentParser(
        description='Generate a random undirected connected graph for BFS demonstration.',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog='''
Examples:
  %(prog)s                          # Generate graph with default settings
  %(prog)s -n 15                    # Generate graph with 15 nodes
  %(prog)s -n 20 -p 0.5             # 20 nodes with 0.5 edge probability
  %(prog)s -n 20 -p 0.5 -o my.txt   # Custom output file
  %(prog)s -s 42                    # Use seed 42 for reproducibility
        '''
    )
    
    parser.add_argument('-n', '--nodes', type=int, default=10,
                        help='Number of nodes in the graph (default: 10)')
    parser.add_argument('-p', '--probability', type=float, default=0.3,
                        help='Probability of additional edges (default: 0.3)')
    parser.add_argument('-o', '--output', type=str, default='graph.txt',
                        help='Output file name (default: graph.txt)')
    parser.add_argument('-s', '--seed', type=int, default=None,
                        help='Random seed for reproducibility (optional)')
    
    args = parser.parse_args()
    
    num_nodes = args.nodes
    edge_probability = args.probability
    output_file = args.output
    
    # Set random seed if provided
    if args.seed is not None:
        random.seed(args.seed)
        print(f"Using random seed: {args.seed}")
    
    print(f"Generating graph with {num_nodes} nodes...")
    print(f"Edge probability: {edge_probability}")
    
    # Generate graph
    graph = generate_connected_graph(num_nodes, edge_probability)
    visit_times = generate_visit_times(num_nodes)
    latencies = generate_latencies(num_nodes)
    
    # Calculate statistics
    total_edges = sum(len(neighbors) for neighbors in graph.values()) // 2
    avg_degree = total_edges * 2 / num_nodes
    
    print(f"Generated graph with {total_edges} edges")
    print(f"Average degree: {avg_degree:.2f}")
    
    # Write to file
    write_graph_to_file(graph, visit_times, latencies, output_file)
    print(f"Graph written to {output_file}")
    
    # Print sample
    print("\nSample (first 5 nodes):")
    for node in range(min(5, num_nodes)):
        neighbors = sorted(graph[node])
        print(f"  Node {node} (score: {visit_times[node]}, latency: {latencies[node]}s): neighbors {neighbors}")


if __name__ == "__main__":
    main()
