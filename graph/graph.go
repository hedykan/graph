package graph

import (
	"sync"
)

type GraphNode struct {
	Id         int
	EndIdList  []int
	HeadIdList []int
}

type Graph struct {
	Map map[int]*GraphNode
	Mut sync.Mutex
}

func NewGraph() *Graph {
	return &Graph{
		Map: make(map[int]*GraphNode),
	}
}

func NewGraphNode(id int) *GraphNode {
	return &GraphNode{
		Id: id,
	}
}

func (graph *Graph) Add(node *GraphNode) {
	graph.Mut.Lock()
	defer graph.Mut.Unlock()

	graph.Map[node.Id] = node
}

func (graph *Graph) Get(id int) (*GraphNode, bool) {
	data, ok := graph.Map[id]
	return data, ok
}

// 检查是否都为空
func (graph *Graph) checkOwnAndObj(ownId, objId int) bool {
	_, ownOk := graph.Get(ownId)
	_, objOk := graph.Get(ownId)
	return ownOk && objOk
}

// 连在一个物品后，当前的节点连在目标节点后面 obj<-own
func (graph *Graph) Link(ownId int, objId int) bool {
	graph.Mut.Lock()
	defer graph.Mut.Unlock()

	if graph.checkOwnAndObj(ownId, objId) != true {
		return false
	}
	ownNode, _ := graph.Get(ownId)
	objNode, _ := graph.Get(objId)
	// 查看是否已经存在与对方列表中
	if FindInt(ownNode.HeadIdList, objId) && FindInt(objNode.EndIdList, ownId) {
		return false
	}
	ownNode.HeadIdList = append(ownNode.HeadIdList, objId)
	objNode.EndIdList = append(objNode.EndIdList, ownId)
	return true
}

// 取消连接着own的节点，obj<-own
func (graph *Graph) CancelLink(ownId int, objId int) bool {
	headIndex := -1
	endIndex := -1
	ownNode, ok := graph.Get(ownId)
	// 找到当前节点连着的子节点
	if ok {
		headIndex = FindIntIndex(ownNode.HeadIdList, objId)
		ownNode.HeadIdList = DeleteInt(ownNode.HeadIdList, headIndex)
	}
	objNode, ok := graph.Get(objId)
	if ok {
		endIndex = FindIntIndex(objNode.EndIdList, ownId)
		objNode.EndIdList = DeleteInt(objNode.EndIdList, endIndex)
	}

	if headIndex == -1 && endIndex == -1 {
		return false
	}
	return true
}

func (graph *Graph) Delete(id int) (*GraphNode, bool) {
	graph.Mut.Lock()
	defer graph.Mut.Unlock()

	data, ok := graph.Get(id)
	if !ok {
		goto end
	}
	// 取消连接着node的
	// 如果不为空则继续指向
	// 正常循环会因为list变化而导致循环不全
	for {
		if len(data.EndIdList) <= 0 {
			break
		}
		graph.CancelLink(data.EndIdList[0], id)
	}
	// 取消node连接着的
	for {
		if len(data.HeadIdList) <= 0 {
			break
		}
		graph.CancelLink(id, data.HeadIdList[0])
	}
	data, _ = graph.Get(id)
	delete(graph.Map, id)
end:
	return data, ok
}

func (graph *Graph) GetArr(idArr []int) []GraphNode {
	var res []GraphNode
	for _, id := range idArr {
		data, ok := graph.Get(id)
		if ok {
			res = append(res, *data)
		}
	}
	return res
}

// 查找连通图
func (graph *Graph) FindList(nodeId int) []int {
	// for _, v := range swapMap {
	// 	fmt.Println(v)
	// }
	// 初始化栈堆，查找列表
	stack := NewStack[int]()
	stack.Push(nodeId)
	findNode := make([]int, 0)
	// 开始查找
	findNode = graph.listProc(stack, findNode)
	return findNode
}

func (graph *Graph) listProc(stack *Stack[int], findNode []int) []int {
	// 获取处于栈堆顶的节点
	nodeId, _ := stack.Pop()
	// 如果已经找到过了，就弹出
	if !FindInt(findNode, nodeId) {
		findNode = append(findNode, nodeId)
	} else {
		return findNode
	}
	node, _ := graph.Get(nodeId)
	// 尾节点入栈
	for i := 0; i < len(node.EndIdList); i++ {
		stack.Push(node.EndIdList[i])
	}
	if stack.Size > 0 {
		findNode = graph.listProc(stack, findNode)
	}
	return findNode
}

// 查找回路
// 深度优先搜索的同时，如果有和目标点相同的则为有回路
func (graph *Graph) FindCycle(nodeId int) [][]int {
	rec := make([]int, 0)
	res := make([][]int, 0)
	res = graph.cycleProc(res, rec, nodeId, nodeId)
	return res
}

func (graph *Graph) cycleProc(res [][]int, rec []int, nowId int, objId int) [][]int {
	if FindInt(rec, nowId) {
		if nowId == objId { // 检测是否与目标id相同，有则添加列表
			res = append(res, rec)
		}
		return res
	}
	rec = append(rec, nowId) // 记录经过的节点
	node, ok := graph.Get(nowId)
	if !ok {
		return res
	}
	for i := 0; i < len(node.EndIdList); i++ {
		res = graph.cycleProc(res, rec, node.EndIdList[i], objId)
	}
	return res
}

// 获取每层End节点数组
func (graph *Graph) GetEndArr(nodeId, layer int) []int {
	if layer == 0 {
		return nil
	}

	var res []int
	data, ok := graph.Get(nodeId)
	if !ok {
		return res
	}
	for _, v := range data.EndIdList {
		if !FindInt(res, v) {
			res = append(res, v)
			res = append(res, graph.GetEndArr(v, layer-1)...)
		}
	}
	return res
}

// 获取每层Head节点数组
func (graph *Graph) GetHeadArr(nodeId, layer int) []int {
	if layer == 0 {
		return nil
	}

	var res []int
	data, ok := graph.Get(nodeId)
	if !ok {
		return res
	}
	for _, v := range data.HeadIdList {
		if !FindInt(res, v) {
			res = append(res, v)
			res = append(res, graph.GetHeadArr(v, layer-1)...)
		}
	}
	return res
}
