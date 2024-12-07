package main

import (
	"fmt"
	"sort"

	utils "example.com/aoc2024"
)

type Page struct {
	data []int
}

type Rules struct {
	rules map[int]map[int]bool
}

func (rules Rules) AddRule(lhs int, rhs int) {
	_, ok := rules.rules[lhs]
	if !ok {
		rules.rules[lhs] = make(map[int]bool)
	}
	rules.rules[lhs][rhs] = true
}

func (rules Rules) Dfs(v int, visited *map[int]bool, answer *[]int, subset *[]int) {
	(*visited)[v] = true

	intersection := []int{}

	adj, hasAdj := rules.rules[v]
	if hasAdj {
		for _, u := range *subset {
			_, hasRule := adj[u]
			if hasRule {
				intersection = append(intersection, u)
			}
		}
	}

	for _, u := range intersection {
		_, inVisited := (*visited)[u]
		if !inVisited {
			rules.Dfs(u, visited, answer, subset)
		}
	}

	(*answer) = append((*answer), v)
}

func (rules Rules) TopologicalSort(subset *[]int) map[int]int {
	visited := map[int]bool{}
	answer := []int{}

	for _, v := range *subset {
		_, inVisited := visited[v]
		if !inVisited {
			rules.Dfs(v, &visited, &answer, subset)
		}
	}

	order := map[int]int{}
	for reversedOrder, value := range answer {
		order[value] = len(answer) - reversedOrder
	}

	return order
}

func ReadRulesAndPages() (Rules, []Page) {
	rules := Rules{rules: make(map[int]map[int]bool)}

	reader := utils.Read(5, false)
	arrayReader := utils.ArrayReader(reader, "|")
	for {
		line, err := arrayReader()
		if err != nil {
			break
		}
		rules.AddRule(line[0], line[1])
	}

	pages := []Page{}

	arrayReader = utils.ArrayReader(reader, ",")
	for {
		data, err := arrayReader()
		if err != nil {
			break
		}
		page := Page{data}
		pages = append(pages, page)
	}

	return rules, pages
}

func Part1() {
	rules, pages := ReadRulesAndPages()
	ans := 0
	for _, page := range pages {
		order := rules.TopologicalSort(&page.data)

		pageOrdered := []int{}
		for _, p := range page.data {
			pageOrdered = append(pageOrdered, order[p])
		}
		isSorted := sort.IntsAreSorted(pageOrdered)
		if isSorted {
			ans += page.data[len(page.data)/2]
		}
	}

	fmt.Println(ans)
}

func Part2() {
	rules, pages := ReadRulesAndPages()
	ans := 0
	for _, page := range pages {
		order := rules.TopologicalSort(&page.data)

		pageOrdered := []int{}
		for _, p := range page.data {
			pageOrdered = append(pageOrdered, order[p])
		}
		isSorted := sort.IntsAreSorted(pageOrdered)
		if !isSorted {
			sort.Ints(pageOrdered)
			for k, v := range order {
				if v == pageOrdered[len(pageOrdered)/2] {
					ans += k
					break
				}
			}
		}
	}

	fmt.Println(ans)
}

func main() {
	Part1()
	Part2()
}
