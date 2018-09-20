package api

type PgwayRoute struct {
}

type PgwayRouteTree struct {
	Nodes       map[string]*PgwayRouteNode
	initialized bool
}

type PgwayRouteNode struct {
	Path                string                     // pathの一部
	IsPathVariable      bool                       //
	PathVariableKeyName string                     //
	Nodes               map[string]*PgwayRouteNode //
	IsEndNode           bool                       // is end (it means the routing exists.)
	Api                 *PgwayApi                  //
}

const pathVariableNodePath = "/" // because / is not use for url node.

func ShowNodes(nodes map[string]*PgwayRouteNode) {
	if len(nodes) == 0 {
		return
	}

	for key, node := range nodes {
		println("key:" + key + "->" + node.Api.Path)
		ShowNodes(node.Nodes)
	}
}

func (tree *PgwayRouteTree) tracePath(request *PgwayRequest) (*PgwayApi, map[string]string) {

	pathRunes := []rune(request.Path)
	node := tree.Nodes[request.HTTPMethod]
	pathVariables := map[string]string{}
	targetNode := findNext(node, pathRunes, 0, len(request.Path), pathVariables)
	if targetNode != nil && targetNode.IsEndNode {
		return targetNode.Api, pathVariables
	} else {
		//
		// request path is added but is not end node.
		return nil, pathVariables
	}
}

func findNext(node *PgwayRouteNode, pathRunes []rune, i int, pathLen int, pathVariables map[string]string) *PgwayRouteNode {
	if node == nil {
		return nil
	}

	if i == pathLen || (i == pathLen-1 && pathRunes[i] == '/') {
		// finish
		return node
	}

	if pathRunes[i] != '/' {
		//
		// request path format is invalid.
		return nil
	}
	i++

	var pathPart []rune
	for ; i < pathLen && pathRunes[i] != '/'; i++ {
		pathPart = append(pathPart, pathRunes[i])
	}
	pathPartStr := string(pathPart)
	currentNode := node.Nodes[pathPartStr]
	currentNode = findNext(currentNode, pathRunes, i, pathLen, pathVariables)
	if currentNode != nil {
		return currentNode
	}
	//
	// check path(variable type).
	currentNode = node.Nodes[pathVariableNodePath]
	if currentNode != nil {
		pathVariables[currentNode.PathVariableKeyName] = pathPartStr
		return findNext(currentNode, pathRunes, i, pathLen, pathVariables)
	} else {
		return nil
	}
}

func getOrCreateNode(key string, nodes map[string]*PgwayRouteNode) *PgwayRouteNode {
	if node, ok := nodes[key]; ok {
		return node
	} else {
		newNode := &PgwayRouteNode{
			IsPathVariable: false,
			IsEndNode:      false,
			Nodes:          map[string]*PgwayRouteNode{},
		}
		nodes[key] = newNode
		return newNode
	}
}

func (tree *PgwayRouteTree) addRoute(api PgwayApi) {
	pathRunes := []rune(api.Path)
	pathLen := len(api.Path)

	currentNode := getOrCreateNode(api.HTTPMethod, tree.Nodes)

	for i := 0; i < pathLen; {
		if pathRunes[i] != '/' {
			panic("[PGWay]{Router]invalid path given.")
		}
		i++
		var pathPart []rune
		for ; i < pathLen && pathRunes[i] != '/'; i++ {
			pathPart = append(pathPart, pathRunes[i])
		}

		if pathPart[0] == ':' {
			// path variable
			currentNode = getOrCreateNode(pathVariableNodePath, currentNode.Nodes)
			currentNode.Path = pathVariableNodePath
			pathVariableKeyName := string(pathPart[1:])
			if currentNode.IsPathVariable {
				if currentNode.PathVariableKeyName != pathVariableKeyName {
					panic("[PGWay][Router]can't add route: path variable name has changed. " + currentNode.PathVariableKeyName + "->" + pathVariableKeyName)
				}
			} else {
				currentNode.PathVariableKeyName = pathVariableKeyName
				currentNode.IsPathVariable = true
			}
		} else {
			key := string(pathPart)
			currentNode = getOrCreateNode(key, currentNode.Nodes)
			currentNode.Path = key
		}
	}
	if currentNode.IsEndNode {
		panic("[PGWay][Router] duplicate path")
	} else {
		currentNode.Api = &api
		currentNode.IsEndNode = true
	}
}
