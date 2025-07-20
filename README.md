# Automaton Pathfinding and Environment Management

This Go project implements a system for managing a 2D environment with obstacles and autonomous agents (automatons). It includes functionalities for creating and querying the environment, as well as a pathfinding algorithm for the automatons to navigate around obstacles.

## Table of Contents

- [Automaton Pathfinding and Environment Management](#automaton-pathfinding-and-environment-management)
  - [Table of Contents](#table-of-contents)
  - [Project Overview](#project-overview)
  - [Core Components](#core-components)
    - [`piano` (Environment)](#piano-environment)
    - [`automaton`](#automaton)
    - [`obstacle`](#obstacle)
    - [`Point`](#point)
    - [`Stack`](#stack)
  - [Algorithm: `percorso_automa` (Automaton Path)](#algorithm-percorso_automa-automaton-path)
    - [Logic](#logic)
    - [Heuristic](#heuristic)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Building and Running](#building-and-running)
  - [Usage](#usage)

## Project Overview

The project simulates a 2D grid-based environment where "automatons" (robots or agents) need to move while avoiding "obstacles." It provides commands to initialize the environment, add obstacles and automatons, query the state of a position, and instruct automatons to find paths to a target. The pathfinding is implemented using a depth-first-like search approach, potentially guided by a heuristic.

## Core Components

### `piano` (Environment)
Represents the 2D plane or environment. It holds a slice of `obstacle` objects and a map of `automaton` objects, keyed by their names.
* `obstacles`: A slice of rectangular obstacle areas.
* `automatons`: A map of automatons, allowing access by their unique name.

### `automaton`
Represents an autonomous agent within the `piano`.
* `pos`: The current `Point` (x, y coordinate) of the automaton.
* `path`: A `Stack` that stores the path taken by the automaton.

### `obstacle`
Defines a rectangular area that automatons cannot enter.
* `low`: A `Point` representing the bottom-left corner of the obstacle.
* `high`: A `Point` representing the top-right corner of the obstacle.

### `Point`
A simple struct representing a 2D coordinate.
* `x`: X-coordinate.
* `y`: Y-coordinate.

### `Stack`
A custom stack implementation used to store the path of an automaton.
* `head`: Pointer to the top `StackNode`.
* `push(value Point, actions []Point)`: Adds a new node to the top of the stack.
* `pop() (Point, []Point, int)`: Removes and returns the top node's value, actions, and counter.
* `top() Point`: Returns the value of the top node without removing it.
* `isEmpty() bool`: Checks if the stack is empty.

## Algorithm: `percorso_automa` (Automaton Path)

The `percorso_automa` function, typically called within a `muovi` (move) command, is responsible for guiding an automaton from its current position to a target destination while avoiding obstacles. The algorithm appears to be a form of **Depth-First Search (DFS)** or a similar tree-traversal algorithm, enhanced with a heuristic to prioritize movement towards the target.

### Logic

The algorithm uses a `Stack` to keep track of the current path being explored and to backtrack when a dead-end or obstacle is encountered.

1.  **Initialization:**
    * The `Stack` for the automaton's path is initialized with the starting position.
    * A `visited` map (or similar mechanism) is used to prevent revisiting the same points and falling into infinite loops.
2.  **Exploration Loop:**
    * The automaton continuously tries to move from its current position (`pos`) towards the `target`.
    * It determines possible next moves (up, down, left, right).
3.  **Heuristic Application:**
    * Before choosing the next move, the algorithm applies a heuristic to evaluate potential next positions.
    * The heuristic calculates the **Manhattan distance** (or Chebyshev distance, though Manhattan is more common for grid movement) from each potential next point to the `target`.
    * It prioritizes moves that **decrease the distance to the target**.
    * If multiple moves result in the same minimal distance, a specific order (e.g., trying Y-axis movements before X-axis movements, or a fixed order like up, down, left, right) might be followed.
4.  **Movement and Backtracking:**
    * If a valid, unvisited, and available (not an obstacle or occupied by another automaton) next position is found and it improves the heuristic score, the automaton moves to that position, and the position is pushed onto the stack.
    * If no such move can be made from the current position (e.g., surrounded by obstacles, visited points, or moved away from the target), the automaton **backtracks** by popping elements from the stack until it reaches a point from which a new, valid, and improving path can be explored.
    * The `counter_action` within `StackNode` suggests that the algorithm keeps track of how many moves have been tried from a specific node before backtracking, potentially for exploring alternative paths upon backtracking.

The algorithm effectively searches for a path by attempting to move greedily towards the target, but with the ability to backtrack and explore other directions when the greedy approach fails.

### Heuristic

The heuristic used is likely a **distance-based metric**, such as the Manhattan distance, which calculates the sum of the absolute differences of the x and y coordinates between the current point and the target point:

`H(current, target) = |current.x - target.x| + |current.y - target.y|`

This heuristic provides an estimate of the remaining distance to the target. By prioritizing moves that minimize this distance, the algorithm attempts to find a path efficiently without exhaustively searching all possible paths, making it a form of **greedy best-first search** when combined with a depth-first traversal strategy and backtracking.

## Getting Started

### Prerequisites

* Go programming language installed (version 1.16 or higher recommended). You can download it from [https://golang.org/dl/](https://golang.org/dl/).

### Building and Running

1.  **Save the files:** Ensure both `stack.go` and `29669A_veneroni_andrea.go` are in the same directory (e.g., `automaton_project`).
2.  **Navigate to the directory:** Open your terminal or command prompt and change your current directory to `automaton_project`.
    ```bash
    cd path/to/automaton_project
    ```
3.  **Run the program:**
    ```bash
    go run 29669A_veneroni_andrea.go stack.go
    ```
    (Note: The `29669A_veneroni_andrea.go` file likely contains the `main` function and should be specified first.)

## Usage

The program appears to be driven by user input, reading commands from standard input. Based on the functions like `crea`, `stato`, `stampa`, `posizioni`, `ostacolo`, `automa`, and `muovi`, the following commands are likely supported:

* `crea`: Initializes a new empty environment.
* `stato <x> <y>`: Checks the state of the position (`x`, `y`). Returns `A` for automaton, `E` for empty, `O` for obstacle.
* `stampa`: Prints information about the environment, possibly including obstacles.
* `posizioni [<automaton_name>]`: Lists the positions of automatons. If `automaton_name` is provided, it might show details for that specific automaton.
* `ostacolo <x1> <y1> <x2> <y2>`: Adds a new obstacle rectangle defined by two points (`x1, y1`) and (`x2, y2`).
* `automa <name> <x> <y>`: Adds a new automaton with `name` at position (`x`, `y`).
* `muovi <automaton_name> <target_x> <target_y>`: Commands the specified automaton to move to the target position (`target_x`, `target_y`), attempting to find a path around obstacles.
* `path <automaton_name>`: Prints the path taken by the specified automaton.
* `fine`: Terminates the program.

You would interact with the program by typing these commands into the console.
