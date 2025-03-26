#!/usr/bin/env python3
"""
Script to plot memory usage from CSV file of following format. 

Usage:
    python plot_memory-usage.py <input.csv> <output.png>

Positional Arguments:
    csv  - CSV file with columns: seconds,virtual_memory_kb,physical_memory_kb
    png  - Output file name for the generated plot

Example:
    python plot_memory.py memory_usage.csv memory_plot.png
"""

import pandas as pd
import matplotlib.pyplot as plt
import sys
import argparse

def plot_memory(csv_file, png_file):
    # Read the CSV file
    try:
        df = pd.read_csv(csv_file)
    except FileNotFoundError:
        print(f"Error: Could not find CSV file '{csv_file}'")
        sys.exit(1)
    except Exception as e:
        print(f"Error reading CSV file: {e}")
        sys.exit(1)

    # Convert KB to MB
    df['virtual_gb'] = df['virtual_memory_kb'] / 1024 / 1024
    df['physical_mb'] = df['physical_memory_kb'] / 1024

    # Create figure with two subplots
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 10))

    # Plot Virtual Memory
    ax1.plot(df['seconds'], df['virtual_gb'], 'b-', linewidth=2)
    ax1.set_title('Virtual Memory Usage Over Time')
    ax1.set_xlabel('Time (seconds)')
    ax1.set_ylabel('Virtual Memory (GiB)')
    ax1.grid(True)

    # Plot Physical Memory
    ax2.plot(df['seconds'], df['physical_mb'], 'g-', linewidth=2)
    ax2.set_title('Physical Memory Usage Over Time')
    ax2.set_xlabel('Time (seconds)')
    ax2.set_ylabel('Physical Memory (MiB)')
    ax2.grid(True)

    # Adjust layout
    plt.tight_layout()

    # Save the plot
    try:
        plt.savefig(png_file)
        print(f"Plot saved as '{png_file}'")
    except Exception as e:
        print(f"Error saving PNG file: {e}")
        sys.exit(1)

    # Display the plot
    plt.show()

def main():
    # Set up argument parser
    parser = argparse.ArgumentParser(description='Plot memory usage from CSV file')
    parser.add_argument('csv', help='Input CSV file containing memory usage data')
    parser.add_argument('png', help='Output PNG file for the plot')

    args = parser.parse_args()

    plot_memory(args.csv, args.png)

if __name__ == "__main__":
    main()

