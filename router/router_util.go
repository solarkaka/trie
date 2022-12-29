package router

import (
	"strings"
)

type RouteUtil struct {
	trie *PathTrie
}

func NewRouteUtil() *RouteUtil {
	return &RouteUtil{
		trie:NewPathTrie()}
}

func (util *RouteUtil) MatchService(key string) interface{} {
	if strings.HasPrefix(key, "kpi") {
		svc := util.trie.Get("/" + key)
		if svc != nil {
			return svc
		}
	}
	return util.trie.Get("/" + strings.TrimPrefix(key, "kpi"))
}

// customize func
func (util *RouteUtil) parseHostPath(hosts []string, paths []string, svc interface{}) bool {
	if hosts == nil || paths == nil {
		return false
	}
	var uris []string
	for _, path := range paths {
		if strings.HasPrefix(path, "//") {
			trimPath := strings.TrimPrefix(path, "/")
			util.trie.Put(trimPath, svc)
		} else {
			//切片效率？
			uris = append(uris, path)
		}
	}
	if len(uris) > 0 {
		for _, host := range hosts {
			for _, uri := range uris {
				util.trie.Put("/"+host+uri, svc)
			}
		}
	}
	return true
}
