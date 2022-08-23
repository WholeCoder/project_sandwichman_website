package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/iancoleman/orderedmap"
)

// Used

// map[string]string
func compressText(encodingHash *orderedmap.OrderedMap, originalText string) ([]byte, int, int) { // compressed, byteCount, bitCount
	compressed := InitNewByteset(make([]byte, 1))
    countOfBits := 0
    lastByte := 1
	for _, letter := range originalText {
		value, exists := (*encodingHash).Get(string(letter))

		if exists {
            for _,zero_or_one_string := range value.(string) {
		        compressed.SetBit(countOfBits, string(zero_or_one_string) == "1")
                countOfBits++;
                if countOfBits / 8 == lastByte {
                    lastByte++
                    compressed = append(compressed, byte(0))
                }

            }
	    }
    }

	return compressed, lastByte, countOfBits
}

func printOutPathOfNodeToRoot(node *Node) {

	fmt.Println("\n^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ START")
	count := 1
	for node != nil {
		fmt.Println(count, ": ", node.Letter_s)
		node = node.Parent
		count++
	}
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^END\n")
}

func check(err error) {
	if err != nil {
		log.Fatalf("failed to open file:  %s", err)
	}
}

// map[string]string
func printEncodingHash(encodingHash orderedmap.OrderedMap) {
	fmt.Println("------------------")

	for _, key := range encodingHash.Keys() {
		value, _ := encodingHash.Get(key)
		fmt.Println("encodingHash[ ", key, " ] = ", value)
	}

	fmt.Println("------------------")
}

//*map[string]*Node   returns *map[string]string
func buildEncodingHash(hashForEncoding orderedmap.OrderedMap) *orderedmap.OrderedMap {
	encodingHash := orderedmap.New() //map[string]string{}

	for _, key := range hashForEncoding.Keys() {
		if len(key) == 1 {
			value, _ := hashForEncoding.Get(key)
			encodingHash.Set(key, buildEncoding(value.(*Node)))
		}
	}

	return encodingHash
}

func buildEncoding(node *Node) string {
	encoding := ""
	n := node

	count := 1
	for n != nil {
		encoding = n.ChildNodeRorL + encoding
		n = n.Parent
		count++
	}
	return encoding
}

// Used
// *map[string]Node, *map[string]string
func initBinaryTree(hash *orderedmap.OrderedMap, encodingHash *orderedmap.OrderedMap) *Node {

	for len((*hash).Keys()) > 1 {
		// findFreeMinNode will remove the nodes from the hash
		nextNode := findFreeMinNode(hash)
		nextNode.ChildNodeRorL = "0"

		secondNode := findFreeMinNode(hash)
		secondNode.ChildNodeRorL = "1"

		newNode := createNewNodeFrom(nextNode, secondNode)
		//(*hash)[newNode.Letter_s] = *newNode
		(*hash).Set(newNode.Letter_s, *newNode)
	}

	var n Node
	for _, key := range (*hash).Keys() { // Runs Once.
		node_as_interface_type, _ := (*hash).Get(key)
		n = node_as_interface_type.(Node)
	}

	fixBinaryTree(&n) // sorry folks!!
	fixEncodingHash(&n, encodingHash)
	return &n
}

// Used
// map[string]string
func fixEncodingHash(node *Node, encodingHash *orderedmap.OrderedMap) {
	if node == nil {
		return
	}

	if len(node.Letter_s) == 1 {
		(*encodingHash).Set(node.Letter_s, buildEncoding(node))
	}
	fixEncodingHash(node.Left, encodingHash)
	fixEncodingHash(node.Right, encodingHash)
}

// Used
func fixBinaryTree(n *Node) {
	if n.Left != nil {
		n.Left.Parent = n
		fixBinaryTree(n.Left)
	}
	if n.Right != nil {
		n.Right.Parent = n
		fixBinaryTree(n.Right)
	}
}

func printNodeDetails(n *Node) {
	if n == nil {
		return
	}
	printNodeDetails(n.Parent)
	fmt.Println("\n\n")
}

func findAndReturnANode(n *Node, nodeName string) *Node {
	if n == nil {
		return nil
	}
	if n.Letter_s == nodeName {
		return n
	}
	n1 := findAndReturnANode(n.Left, nodeName)
	if n1 != nil {
		return n1
	}
	n2 := findAndReturnANode(n.Right, nodeName)
	if n2 != nil {
		return n2
	}
	return nil
}

func printOutWholeTreeInOrder(n *Node, encodingHash *map[string]string) {
	if n == nil {
		return
	}
	printOutWholeTreeInOrder(n.Left, encodingHash)
	fmt.Println()
	//fmt.Println(n.Letter_s)

	if len(n.Letter_s) == 1 {
		(*encodingHash)[n.Letter_s] = buildEncoding(n)

		//printOutPathOfNodeToRoot(n)
	}

	if n.Left != nil {
		//fmt.Println("\t",n.Left.Letter_s)
	}

	if n.Right != nil {
		//fmt.Println("\t",n.Right.Letter_s)
	}
	if n.Parent != nil {
		//fmt.Println("Parent: ", n.Parent.Left, n.Parent.Data, n.Parent.Letter_s, n.Parent.Right, n.Parent, n.ChildNodeRorL, n.AlreadyUsedToBuildBinaryTree)
	}

	printOutWholeTreeInOrder(n.Right, encodingHash)
}

func debugCountHowManyLeftNodes(node *Node) {
	fmt.Println("----------")

	for node != nil {
		fmt.Println("\t", node.Letter_s)
		if node.Right != nil && !true {
			node = node.Right
		} else {
			node = node.Left
		}
	}

	fmt.Println("----------")
}

// Find and remove node from hash.  Return the node
//   *map[string]Node
func findFreeMinNode(hash *orderedmap.OrderedMap) *Node {
	var minKey string
	var minValue Node

	for _, key := range (*hash).Keys() {
		minKey = key
		minValueInterfacetype, _ := (*hash).Get(minKey)
		minValue = minValueInterfacetype.(Node)
		break
	}

	for _, key := range (*hash).Keys() {
		valueInterface, _ := (*hash).Get(key)
		value := valueInterface.(Node)
		if value.Data < minValue.Data { // Data is the Probability
			minKey = key
			minValue = value
		}
	}

	hashMinValueInterface, _ := (*hash).Get(minKey)
	hashMinValue := hashMinValueInterface.(Node)
	nodeMinValue := Node{Left: hashMinValue.Left,
		Data:     hashMinValue.Data, // Data is Probability
		Letter_s: minKey,
		Right:    hashMinValue.Right,

		Parent: nil,

		ChildNodeRorL: hashMinValue.ChildNodeRorL}
	(*hash).Delete(minKey)

	return &nodeMinValue
}

// Used
func createNewNodeFrom(node1, node2 *Node) *Node {
	newNode := Node{Left: node1, Data: node1.Data + node2.Data, Letter_s: node1.Letter_s + node2.Letter_s, Right: node2, Parent: nil}
	return &newNode
}

func convertNilToZero(valueInterface interface{}) int {

	if valueInterface == nil {
		return 0
	}

	return valueInterface.(int)
}

// Used
// map[string]Node
func initFrequencyHash(fileName string) orderedmap.OrderedMap {

	dat, err := ioutil.ReadFile(fileName)
	check(err)
	asString := string(dat)

	hash := orderedmap.New() //map[string]int{}

	for _, ch := range asString {
		valueInterface, _ := hash.Get(string(ch))
		value2 := convertNilToZero(valueInterface)

		value := value2
		hash.Set(string(ch), value+1)
		//hash[string(ch)] += 1
	}

	totalLetters := 0
	for _, key := range hash.Keys() {
		valueInterface, _ := hash.Get(key)
		value := valueInterface.(int)
		totalLetters += value
	}

	freqNodemap := orderedmap.New() //map[string]Node{}

	for _, key := range hash.Keys() {
		valueInterface, _ := hash.Get(key)
		var value int
		if valueInterface == nil {
			value = 0
		} else {
			value = valueInterface.(int)
		}
		freqNodemap.Set(key, Node{Data: float64(value) / float64(totalLetters), AlreadyUsedToBuildBinaryTree: false})

		//freqNodemap[key] = Node{Data: float64(value) / float64(totalLetters), AlreadyUsedToBuildBinaryTree: false}
	}

	return *freqNodemap
}

// Used
// returns map[string]float64
func initFrequencyHashWithFloat64ForValues(fileName string) orderedmap.OrderedMap {

	dat, err := ioutil.ReadFile(fileName)
	check(err)
	asString := string(dat)

	hash := orderedmap.New() //map[string]int{}

	for _, ch := range asString {
		valueInterface, _ := hash.Get(string(ch))

		var value int
		if valueInterface == nil {
			value = 0
		} else {
			value = valueInterface.(int)
		}
		hash.Set(string(ch), value+1)
		// hash[string(ch)] += 1
	}

	totalLetters := 0
	for _, key := range hash.Keys() {
		valueInterface, _ := hash.Get(key)
		value := valueInterface.(int)
		totalLetters += value
	}

	freqNodemap := orderedmap.New() //map[string]float64{}

	for _, key := range hash.Keys() {
		valueInterface, _ := hash.Get(key)
		var value int
		if valueInterface == nil {
			value = 0.0
		} else {
			value = valueInterface.(int)
		}
		freqNodemap.Set(key, float64(value)/float64(totalLetters))
		// freqNodemap[key] = float64(value) / float64(totalLetters)
	}

	return *freqNodemap
}
