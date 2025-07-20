package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type piano struct {
	obstacles  *[]obstacle
	automatons map[string]*automaton
}

type automaton struct {
	pos  Point
	path Stack //stores the path of the automaton
}

type obstacle struct {
	low  Point
	high Point
}

type Point struct {
	x int
	y int
}

func newPiano() piano {
	var plan piano
	obs := make([]obstacle, 0)
	plan.automatons = make(map[string]*automaton)
	plan.obstacles = &obs
	return plan
}

func crea(p *piano) {
	*p = newPiano()
}

func stato(p *piano, x, y int) {
	if isAvailable(p, Point{x, y}) {
		for _, v := range p.automatons {
			if v.pos.x == x && v.pos.y == y {
				fmt.Println("A")
				return
			}
		}
		fmt.Println("E")
	} else {
		fmt.Println("O")
	}
}

func stampa(p *piano) {
	posizioni(p, "")
	fmt.Println("[")
	for _, v := range *(p.obstacles) {
		fmt.Println("(" + strconv.Itoa(v.low.x) + "," + strconv.Itoa(v.low.y) + ")(" + strconv.Itoa(v.high.x) + "," + strconv.Itoa(v.high.y) + ")")
	}
	fmt.Println("]")
}

func automa(p *piano, a int, b int, w string) {
	var robot automaton

	robot.pos.x = a
	robot.pos.y = b

	if isAvailable(p, robot) {
		p.automatons[w] = &robot
	}
}

func ostacolo(p *piano, a, b, c, d int) {

	var obs obstacle

	obs.low.x = a
	obs.low.y = b

	obs.high.x = c
	obs.high.y = d

	if isAvailable(p, obs) {
		*(p.obstacles) = append(*(p.obstacles), obs)
	}
}

func richiamo(p *piano, a int, b int, w string) {
	selected_aut := selectedAutomata(p, w)
	for i := range selected_aut {
		findPath(p, a, b, i)
	}
}

func posizioni(p *piano, w string) {
	fmt.Println("(")
	selected_aut := selectedAutomata(p, w)
	for n, v := range selected_aut {
		fmt.Println(n + ": " + strconv.Itoa(v.pos.x) + "," + strconv.Itoa(v.pos.y))
	}
	fmt.Println(")")
}

func selectedAutomata(p *piano, w string) map[string]*automaton {
	selected := make(map[string]*automaton)
	for n, v := range p.automatons {
		if strings.HasPrefix(n, w) {
			selected[n] = v
		}
	}
	return selected
}

func esistePercorso(p *piano, a int, b int, w string) {
	aut, exist_aut := p.automatons[w]

	if exist_aut {
		actual_pos := aut.pos
		exist := findPath(p, a, b, w)
		aut.pos = actual_pos

		if exist {
			fmt.Println("SI")
			//stampa_percorso_automa(p, w)

		} else {
			fmt.Println("NO")
		}
	} else {
		fmt.Println("invalid name")
		return
	}
}

func findPath(p *piano, a int, b int, w string) bool {
	goal := Point{a, b}
	aut := p.automatons[w]
	obs := obsInPath(p, aut.pos, goal)
	var exist bool
	exist, aut.path = path(p, obs, aut, goal)

	return exist
}

// finds the shortest path and returns the stack with the path of the automaton taken into consideration
/*
   - p indicates our floor
   - obs indicates the obstacles in our path
   - aut indicates the pointer of our automaton
   - goal indicates the arrival point
*/
func path(p *piano, obs []obstacle, aut *automaton, goal Point) (bool, Stack) {
	path := Stack{nil}

	if couldExist(p, *(p.obstacles), aut.pos, goal) {
		call := 0
		stack_overflow := false
		was_there := make(map[Point]int)
		correctDirection(p, &path, obs, aut, goal, &call, &stack_overflow, &was_there)
		for stack_overflow {
			call = 0
			stack_overflow = false
			correctDirection(p, &path, obs, aut, goal, &call, &stack_overflow, &was_there)
		}
		if path.head.next == nil {
			if aut.pos.x == goal.x && aut.pos.y == goal.y {
				return true, path
			}
			return false, path
		}
		return true, path
	}
	return false, path
}

// preliminary action that verifies the cases in which there is certainly no minimum path
/*
  - p indicates our floor
  - tot_obs indicates the obstacles on the map
  - start the starting position
  - goal the end position
*/
func couldExist(p *piano, tot_obs []obstacle, start, goal Point) bool {
	if !isAvailable(p, goal) { //verify that the recall point is reachable
		return false
	}

	var low_l Point
	var high_r Point
	if start.x > goal.x {
		low_l.x = goal.x
		high_r.x = start.x
	} else {
		low_l.x = start.x
		high_r.x = goal.x
	}
	if start.y > goal.y { //necessary for obstacles with vertices outside the path area
		low_l.y = goal.y
		high_r.y = start.y
	} else {
		low_l.y = start.y
		high_r.y = goal.y
	}
	for _, v := range tot_obs {
		if (v.low.y <= low_l.y && v.high.y >= high_r.y) && (v.low.x >= low_l.x && v.low.x <= high_r.x) || (v.low.x <= low_l.x && v.high.x >= high_r.x) && (v.low.y >= low_l.y && v.low.y <= high_r.y) { //case in which a single obstacle prevents the entire passage
			return false
		}
	}
	return true
}

// updates the current pos in the stack of the path taken, recalled every time a step is taken
/*
  - actions indicates the possible actions for shortest path
   - obs indicates the obstacles in our path
   - aut indicates our automaton
   - goal indicates the end position
   - call indicates the number of recursive calls
   - stack_overflow indicates whether or not it overflowed
   - was_there indicates the spots he has already visited but which were not good
*/
func correctDirection(p *piano, path *Stack, obs []obstacle, aut *automaton, goal Point, call *int, stack_overflow *bool, was_there *map[Point]int) {
	*call++
	if *call > 100000 {
		*stack_overflow = true
	} else {
		path.push(aut.pos, chooseAction(obs, aut, goal)) //each node in the stack contains the current position and future actions in order of relevance
		doStep(p, path, obs, aut, goal, call, stack_overflow, was_there)
	}
}

func doStep(p *piano, path *Stack, obs []obstacle, aut *automaton, goal Point, call *int, stack_overflow *bool, was_there *map[Point]int) {
	var new_pos Point
	if path.head.actions != nil { // if nil we have arrived
		if path.head.counter_action < len(path.head.actions) { //if he exceeds len(path.head.actions) it means that he has run out of available moves and must backtrack
			actual_action := path.head.actions[path.head.counter_action]
			path.head.counter_action += 1
			if actual_action.x == 0 && actual_action.y == 0 { //if the algorithm proposes to stay still as a choice, it means that it has no other options so it has to do backtracking
				path.pop()
				aut.pos = path.head.value

				doStep(p, path, obs, aut, goal, call, stack_overflow, was_there)
			} else {
				new_pos = Point{aut.pos.x + actual_action.x, aut.pos.y + actual_action.y}
				_, controller := (*was_there)[new_pos]
				if !controller {
					if isAvailable(p, new_pos) {
						aut.pos = new_pos
						correctDirection(p, path, obs, aut, goal, call, stack_overflow, was_there)
					} else {
						doStep(p, path, obs, aut, goal, call, stack_overflow, was_there)
					}
				} else {
					doStep(p, path, obs, aut, goal, call, stack_overflow, was_there)
				}
			}
		} else { //if it has finished all available actions it backtracks
			(*was_there)[path.head.value] = 0
			path.pop()
			aut.pos = path.head.value
			if path.head.next != nil { //if it were equal to nil, it means that it has exhausted all available paths and terminates
				doStep(p, path, obs, aut, goal, call, stack_overflow, was_there)
			}
		}
	}
}

// returns a slice of actions in order of quality, established through a heuristic
// that takes into account the distance between obstacles and distance from the goal
func chooseAction(obs []obstacle, aut *automaton, goal Point) []Point {

	distance_to_goal_x := int(math.Abs(float64(goal.x - aut.pos.x)))
	distance_to_goal_y := int(math.Abs(float64(goal.y - aut.pos.y)))

	direction_x := Point{norm(goal.x - aut.pos.x), 0}
	direction_y := Point{0, norm(goal.y - aut.pos.y)}

	if direction_x.x == 0 && direction_y.y == 0 { //if the goal is reached it does not return any action
		return nil
	}

	if direction_x.x == 0 { //if it reaches one side of the minimum path area it returns as a priority action, the one that makes it go to the goal
		return []Point{direction_y, direction_x}
	}
	if direction_y.y == 0 {
		return []Point{direction_x, direction_y}
	}

	distX := 99999
	distY := 99999
	for _, v := range obs { //calculates the distance to a possible obstacle on the x-axis and y-axis
		new_distX := distanceToObsXorY(v, aut.pos, 0, direction_x.x)
		if new_distX > 0 && new_distX < distX {
			distX = new_distX
		}
		new_distY := distanceToObsXorY(v, aut.pos, 1, direction_y.y)
		if new_distY > 0 && new_distY < distY {
			distY = new_distY
		}
	}

	distX += distance_to_goal_x //the distance from the goal influences the choice of action by bringing the automaton towards the goal
	distY += distance_to_goal_y

	if distX > distY { //chooses the priority action based on where it comes closest to an obstacle and closest to the goal
		return []Point{direction_x, direction_y}
	} else if distX < distY {
		return []Point{direction_y, direction_x}
	} else { //if indifferent it chooses the axis that is furthest from the goal
		if (aut.pos.x-goal.x)*direction_x.x > (aut.pos.y-goal.y)*direction_x.y {
			return []Point{direction_x, direction_y}
		}
		return []Point{direction_y, direction_x}
	}
}

// returns any obstacles that are within the shortest path area
// (obstacles with no point in this area but which still hinder the path are not considered since they would make the shortest path non-existent)
func obsInPath(p *piano, aut_pos, goal Point) []obstacle {
	obs := make([]obstacle, 0)
	for _, v := range *(p.obstacles) {
		if isIn(v.low, aut_pos, goal) || isIn(v.high, aut_pos, goal) || isIn(Point{v.low.x, v.high.y}, aut_pos, goal) || isIn(Point{v.high.x, v.low.y}, aut_pos, goal) { // check all the points of the obstacle to see if at least one of them is within the area of ​​our path
			obs = append(obs, v)
		}
	}
	return obs
}

// returns the normalized displacement vector
func norm(n int) int {
	if n > 0 {
		return 1
	} else if n == 0 {
		return 0
	}
	return -1
}

// returns the possible distance x or y of a point from an obstacle
func distanceToObsXorY(obs obstacle, aut Point, xORy, direction int) int { // if the last number i zero it calculates the distance x otherwise y
	x := true
	if xORy == 1 {
		x = false
	}
	if x {
		if aut.y >= obs.low.y && aut.y <= obs.high.y {
			if (obs.high.x-aut.x)*direction >= 0 {
				dist1 := int(math.Abs(float64(aut.x - obs.low.x)))
				dist2 := int(math.Abs(float64(aut.x - obs.high.x)))
				if dist1 > dist2 {
					return dist2
				}
				return dist1
			}
		}
		return -1
	} else {
		if aut.x >= obs.low.x && aut.x <= obs.high.x {
			if (obs.high.y-aut.y)*direction >= 0 {
				dist1 := int(math.Abs(float64(aut.y - obs.low.y)))
				dist2 := int(math.Abs(float64(aut.y - obs.high.y)))
				if dist1 > dist2 {
					return dist2
				}
				return dist1
			}
		}
		return -1
	}

}

func main() {
	p := newPiano() //initializes the initial plan
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		esegui(p, input)
	}
}

// checks if the command is a number and converts it
func parseInt(s string) int {
	n, err := strconv.Atoi(s)

	if err != nil {
		fmt.Println("invalid argument: " + s)
		os.Exit(1)
	}
	return n
}

func esegui(p piano, input string) {

	commands := strings.Split(input, " ")

	action := commands[0]

	switch {
	case action == "f":
		os.Exit(0)
	case action == "c":
		crea(&p)
	case action == "s":
		x := parseInt(commands[1])
		y := parseInt(commands[2])
		stato(&p, x, y)
	case action == "S":
		stampa(&p)
	case action == "a":
		a := parseInt(commands[1])
		b := parseInt(commands[2])
		w := commands[3]
		automa(&p, a, b, w)
	case action == "o":
		a := parseInt(commands[1])
		b := parseInt(commands[2])
		c := parseInt(commands[3])
		d := parseInt(commands[4])
		ostacolo(&p, a, b, c, d)
	case action == "r":
		a := parseInt(commands[1])
		b := parseInt(commands[2])
		w := commands[3]
		richiamo(&p, a, b, w)
	case action == "p":
		w := commands[1]
		posizioni(&p, w)
	case action == "e":
		a := parseInt(commands[1])
		b := parseInt(commands[2])
		w := commands[3]
		esistePercorso(&p, a, b, w)
	case action == "stampapercorso": ////////////// TEST
		w := commands[1]
		stampa_percorso_automa(&p, w)
	}

}

// checks if a given position is inside an area
func isIn(target, point_1, point_2 Point) bool {
	var low_l Point
	var high_r Point

	if point_1.x < point_2.x { //adjust the points to make the area format: bottom left point, top right point
		low_l.x = point_1.x
		high_r.x = point_2.x
	} else {
		low_l.x = point_2.x
		high_r.x = point_1.x
	}
	if point_1.y < point_2.y {
		low_l.y = point_1.y
		high_r.y = point_2.y
	} else {
		low_l.y = point_2.y
		high_r.y = point_1.y
	}

	if target.x >= low_l.x && target.x <= high_r.x && target.y >= low_l.y && target.y <= high_r.y {
		return true
	}
	return false
}

// checks whether the position of the automaton, obstacle or point is free or not
func isAvailable(p *piano, obj any) bool {
	switch v := obj.(type) {
	case obstacle:
		for _, robot := range p.automatons {
			if isIn(robot.pos, v.low, v.high) {
				return false
			}
		}
	case automaton:
		for _, obs := range *(p.obstacles) {
			if isIn(v.pos, obs.low, obs.high) {
				return false
			}
		}
	case Point:
		for _, obs := range *(p.obstacles) {
			if isIn(v, obs.low, obs.high) {
				return false
			}
		}
	}
	return true
}

func stampa_percorso_automa(p *piano, w string) { ///////////// TEST
	for i, v := range p.automatons {
		if i == w && v.path.head != nil {
			curr := v.path.head
			for curr != nil {
				fmt.Println("(" + strconv.Itoa(curr.value.x) + "," + strconv.Itoa(curr.value.y) + ")")
				curr = curr.next
			}
		}
	}
	fmt.Println()
	fmt.Println()
}
