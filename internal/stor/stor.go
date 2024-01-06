package stor

import (
	"errors"

	"github.com/google/uuid"
)

type Node struct {
	Id      string
	Prev    *Node
	Next    *Node
	Data    []byte
	GroupId string
}

type group struct {
	Id    string
	Prev  *group
	Next  *group
	Nodes []*Node
}

type Collection struct {
	Id     string
	Groups []*group
	nodes  map[string]*Node
}

func newNode(id string, data []byte) *Node {
	return &Node{Id: id, Data: data}
}

func newGroup() *group {
	return &group{
		Id:    uuid.New().String(),
		Nodes: make([]*Node, 0),
	}
}

func NewCollection() *Collection {
	return &Collection{
		Id:     uuid.New().String(),
		Groups: make([]*group, 0),
		nodes:  make(map[string]*Node),
	}
}

func (g *group) addNode(node *Node) {
	if len(g.Nodes) == 0 {
		g.Id = node.Id
		g.Nodes = append(g.Nodes, node)
	} else {
		tail := g.Nodes[len(g.Nodes)-1]
		node.Prev = tail
		tail.Next = node
		g.Nodes = append(g.Nodes, node)
	}

	node.GroupId = g.Id
}

func (c *Collection) addGroup(g *group) {
	if len(c.Groups) == 0 {
		c.Groups = append(c.Groups, g)
	} else {
		tail := c.Groups[len(c.Groups)-1]
		g.Prev = tail
		tail.Next = g
		c.Groups = append(c.Groups, g)
	}
}

func (c *Collection) NewNode(id string, data []byte, isNewGroup bool) (*Node, error) {
	_, ok := c.nodes[id]
	if ok {
		return nil, errors.New("the id is exists already")
	}

	node := newNode(id, data)

	if isNewGroup {
		g := newGroup()
		g.addNode(node)
		c.addGroup(g)
	} else {
		c.Groups[len(c.Groups)-1].addNode(node)
	}

	c.nodes[id] = node
	return node, nil
}

func (c *Collection) GetNode(nodeId string) *Node {
	return c.nodes[nodeId]
}

func (c *Collection) GetBeforeNodes(nodeId string) []*Node {
	for _, g := range c.Groups {
		var nodes []*Node

		for _, n := range g.Nodes {
			nodes = append(nodes, n)
			if n.Id == nodeId {
				return nodes
			}
		}
	}

	return nil
}

func (c *Collection) GetStartNode(nodeId string) *Node {
	for _, g := range c.Groups {
		for _, n := range g.Nodes {
			if n.Id == nodeId {
				return g.Nodes[0]
			}
		}
	}

	return nil
}

func (c *Collection) GetLastNode() *Node {
	g := c.Groups[len(c.Groups)-1]
	return g.Nodes[len(g.Nodes)-1]
}

func (c *Collection) GetAllNodes() []*Node {
	var nodes []*Node

	for _, g := range c.Groups {
		nodes = append(nodes, g.Nodes...)
	}

	return nodes
}
