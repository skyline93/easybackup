package stor

import (
	"encoding/json"
	"os"
)

type JsonNode struct {
	Id      string `json:"id"`
	GroupId string `json:"group_id"`
	Data    []byte `json:"data"`
}

type JsonGroup struct {
	Id string `json:"id"`
}

type JsonCollection struct {
	Id     string      `json:"id"`
	Nodes  []JsonNode  `json:"nodes"`
	Groups []JsonGroup `json:"groups"`
}

func Serialize(col *Collection, path string) error {
	var nodes []JsonNode
	var groups []JsonGroup

	for _, g := range col.Groups {
		for _, n := range g.Nodes {
			nodes = append(nodes, JsonNode{
				Id:      n.Id,
				GroupId: n.GroupId,
				Data:    n.Data,
			})
		}

		groups = append(groups, JsonGroup{
			Id: g.Id,
		})
	}

	c := &JsonCollection{
		Id:     col.Id,
		Nodes:  nodes,
		Groups: groups,
	}

	d, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, d, 0664)
}

func Deserialize(col *Collection, path string) error {
	d, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var jsonCol JsonCollection
	if err := json.Unmarshal(d, &jsonCol); err != nil {
		return err
	}

	nodes := make(map[string]*Node)

	for _, j := range jsonCol.Groups {
		g := &group{Id: j.Id, Nodes: make([]*Node, 0)}

		for _, n := range jsonCol.Nodes {
			if n.GroupId == g.Id {
				node := &Node{
					Id:      n.Id,
					GroupId: g.Id,
					Data:    n.Data,
				}
				g.addNode(node)
				nodes[n.Id] = node
			}
		}
		col.addGroup(g)
	}

	col.Id = jsonCol.Id
	col.nodes = nodes

	return nil
}
