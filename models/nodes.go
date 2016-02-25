package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

//Node struct
type Node struct {
	Firstseen  time.Time   `json:"firstseen"`
	Lastseen   time.Time   `json:"lastseen"`
	Statistics interface{} `json:"statistics"`
	Nodeinfo   interface{} `json:"nodeinfo"`
	Neighbours interface{} `json:"neighbours"`
}

//Nodes struct: cache DB of Node's structs
type Nodes struct {
	Version   int              `json:"version"`
	Timestamp time.Time        `json:"timestamp"`
	List      map[string]*Node `json:"nodes"` // the current nodemap
	sync.Mutex
}

//NewNodes create Nodes structs (cache DB)
func NewNodes() *Nodes {
	nodes := &Nodes{
		Version: 1,
		List:    make(map[string]*Node),
	}

	return nodes
}

//Get a Node by nodeid
func (nodes *Nodes) Get(nodeID string) *Node {
	now := time.Now()

	nodes.Lock()
	node, _ := nodes.List[nodeID]

	if node == nil {
		node = &Node{
			Firstseen: now,
		}
		nodes.List[nodeID] = node
	}
	nodes.Unlock()

	node.Lastseen = now

	return node
}

//Saver to save the cached DB to json file
func (nodes *Nodes) Saver(outputFile string, saveInterval time.Duration) {
	c := time.Tick(saveInterval)

	for range c {
		nodes.save(outputFile)
	}
}

func (nodes *Nodes) save(outputFile string) {
	nodes.Timestamp = time.Now()

	nodes.Lock()
	data, err := json.Marshal(nodes)
	nodes.Unlock()

	if err != nil {
		log.Panic(err)
	}
	log.Println("saving", len(nodes.List), "nodes")

	tmpFile := outputFile + ".tmp"

	err = ioutil.WriteFile(tmpFile, data, 0644)
	if err != nil {
		log.Panic(err)
	}
	err = os.Rename(tmpFile, outputFile)
	if err != nil {
		log.Panic(err)
	}
}