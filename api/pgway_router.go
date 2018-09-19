package api

import (
	"pgway/model"
)

type PgwayRoute struct {
}

type PgwayRouteTree struct {
	Nodes       map[string]PgwayRouteNode
	initialized bool
}

type PgwayRouteNode struct {
	Path                string                    // pathの一部
	IsPathVariable      bool                      //
	PathVariableKeyName string                    //
	Nodes               map[string]PgwayRouteNode //
	IsEndNode           bool                      // is end (it means the routing exists.)
	Handler             interface{}               //
	Exists              bool
}

const PathVariableNodePath = "/" // because / is not use for url node.

func createTree(server *PgwayServer) PgwayRouteTree {
	tree := PgwayRouteTree{}
	for _, api := range server.Apis {
		tree.addRoute(&api)
	}
	return tree
}

func (tree *PgwayRouteTree) tracePath(request *PgwayRequest) (interface{}, map[string]string) {

	pathRunes := []rune(request.Path)
	node := tree.Nodes[request.HTTPMethod]
	pathVariables := map[string]string{}
	targetNode := node.findNext(pathRunes, 0, len(request.Path), pathVariables)
	if targetNode.IsEndNode {
		return targetNode, pathVariables
	} else {
		//
		// request path is added but is not end node.
		return model.ApiNotFound, pathVariables
	}
}

func (node *PgwayRouteNode) findNext(pathRunes []rune, i int, pathLen int, pathVariables map[string]string) PgwayRouteNode {
	if !node.Exists || i == pathLen-1 || (i == pathLen-2 && pathRunes[i] == '/') {
		// finish
		return *node
	}

	if pathRunes[i] != '/' {
		//
		// request path format is invalid.
		return PgwayRouteNode{} //model.ApiNotFound,pathVariables
	}

	var pathPart []rune
	for ; i < pathLen && pathRunes[i] != '/'; i++ {
		pathPart = append(pathPart, pathRunes[i])
	}
	pathPartStr := string(pathPart)
	currentNode := node.Nodes[pathPartStr]
	currentNode.findNext(pathRunes, i, pathLen, pathVariables)
	if currentNode.Exists {
		return currentNode
	}
	//
	// check path(variable type).
	currentNode = node.Nodes[PathVariableNodePath]
	if currentNode.Exists {
		pathVariables[currentNode.PathVariableKeyName] = pathPartStr
		return currentNode.findNext(pathRunes, i, pathLen, pathVariables)
	} else {
		return PgwayRouteNode{}
	}
}

func (tree *PgwayRouteTree) addRoute(api *PgwayApi) {
	pathRunes := []rune(api.Path)
	pathLen := len(api.Path)

	currentNode := tree.Nodes[api.HTTPMethod]
	//
	// config baseNode
	currentNode.Exists = true

	for i := 0; i < pathLen; i++ {
		if pathRunes[i] != '/' {
			panic("invalid path given.")
		}
		i++
		var pathPart []rune
		for ; i < pathLen && pathRunes[i] != '/'; i++ {
			pathPart = append(pathPart, pathRunes[i])
		}

		if pathPart[0] == ':' {
			// path variable
			currentNode = currentNode.Nodes[PathVariableNodePath]
			currentNode.Path = PathVariableNodePath
			pathVariableKeyName := string(pathPart[1:])
			if currentNode.IsPathVariable {
				if currentNode.PathVariableKeyName != pathVariableKeyName {
					panic("[PGWay][route]can't add route: pathvariable name has changed. " + currentNode.PathVariableKeyName + "->" + pathVariableKeyName)
				}
			} else {
				currentNode.PathVariableKeyName = pathVariableKeyName
				currentNode.IsPathVariable = true
			}
		} else {
			key := string(pathPart)
			currentNode = currentNode.Nodes[key]
			currentNode.Path = key
		}
		currentNode.Exists = true
	}
	currentNode.Handler = api.Handler
	currentNode.IsEndNode = true

}
