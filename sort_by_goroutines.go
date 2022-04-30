package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

/*
Write a program to sort an array of integers. The program should partition the array into 4 parts, each of which is sorted by a different goroutine. Each partition should be of approximately equal size. Then the main goroutine should merge the 4 sorted subarrays into one large sorted array.

The program should prompt the user to input a series of integers. Each goroutine which sorts ¼ of the array should print the subarray that it will sort. When sorting is complete, the main goroutine should print the entire sorted list.
*/

// Binary tree used for the 4 way merge
type Node struct {
	val   int
	left  *Node
	right *Node
}

type BT struct {
	root *Node
}

func (t *BT) insert(s []int) *BT {
	if t.root == nil {
		t.root = &Node{val: s[0], left: nil, right: nil}
	} else {
		t.root.insert(s[0])
	}
	for i := 1; i < len(s); i++ {
		t.root.insert(s[i])
	}
	return t
}

func (n *Node) insert(v int) {
	if n == nil {
		return
	}
	if v <= n.val {
		if n.left == nil {
			n.left = &Node{val: v, left: nil, right: nil}
		} else {
			n.left.insert(v)
		}
	} else {
		if n.right == nil {
			n.right = &Node{val: v, left: nil, right: nil}
		} else {
			n.right.insert(v)
		}
	}
}

func output(n *Node) []int {
	var o = []int{}
	if n == nil {
		return o
	}
	if n.left != nil {
		o = append(output(n.left), o...)
	}
	o = append(o, n.val)
	if n.right != nil {
		o = append(o, output(n.right)...)
	}
	return o
}

/*
// for debugging only
func printTree(node *Node, ns int, ch rune) {
	if node == nil {
		return
	}
	for i := 0; i < ns; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("%c:%v\n", ch, node.val)
	printTree(node.left, ns+2, 'L')
	printTree(node.right, ns+2, 'R')
}
*/

func readInput() []int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	str := strings.TrimSpace(scanner.Text())
	//	fmt.Println(str)
	sls := strings.Split(str, " ")
	s := []int{}
	for _, c := range sls {
		i, err := strconv.Atoi(c)
		if err != nil {
			fmt.Println("Bad character")
			panic(err)
		}
		s = append(s, i)
	}
	return s
}

func splitInto4(s []int) [4][]int {
	var parts [4][]int
	n := len(s)
	partLen := n / 4
	rm := n % 4
	var e int
	var pls int
	i := 0
	for b := 0; b < n; b += partLen {
		if rm > 0 {
			pls = 1
			rm--
		} else {
			pls = 0
		}
		e = b + partLen + pls
		if e > n {
			e = n
		}
		parts[i] = s[b:e]
		b += pls
		i++
	}
	return parts
}

var wg sync.WaitGroup

func partialSort(s *[]int) {
	//	fmt.Println(*s)
	sort.Ints(*s)
	wg.Done()
}

func main() {
	fmt.Println("Enter series of integers separated by space:")
	s := readInput()
	// s := "23 5 58 6 16 57 39 78 63 22 53 17 83 4 8"
	fmt.Print("Integers entered: ")
	fmt.Println(s)
	if len(s) < 4 {
		fmt.Println("Number of elements need to be more than 4")
		return
	}
	parts := splitInto4(s)

	wg.Add(4)
	for i := 0; i < 4; i++ {
		go partialSort(&parts[i])
	}
	wg.Wait()
	fmt.Println()
	// using binary tree to compbine them in approximately O(n•log(4)) time
	tree := &BT{}
	for i := 0; i < 4; i++ {
		fmt.Printf("Sorted by goroutine %d ", i+1)
		fmt.Println(parts[i])
		tree.insert(parts[i])
	}
	fmt.Println()
	sorted := output(tree.root)
	fmt.Print("Sorted Output: ")
	fmt.Println(sorted)
}
