// Implementation of an R-Way Trie data structure.
//
// A Trie has a root Node which is the base of the tree.
// Each subsequent Node has a letter and children, which are
// nodes that have letter values associated with them.
//
//customize by zhou chang
//

package router

import (
	"strings"
)

// StringSegmenter takes a string key with a starting index and returns
// the first segment after the start and the ending index. When the end is
// reached, the returned nextIndex should be -1.
// Implementations should NOT allocate heap memory as Trie Segmenters are
// called upon Gets. See PathSegmenter.
type StringSegmenter func(key string, start int) (segment string, nextIndex int)

// PathSegmenter segments string key paths by slash separators. For example,
// "/a/b/c" -> ("/a", 2), ("/b", 4), ("/c", -1) in successive calls. It does
// not allocate any heap memory.
func PathSegmenter(path string, start int) (segment string, next int) {
	if len(path) == 0 || start < 0 || start > len(path)-1 {
		return "", -1
	}
	end := strings.IndexRune(path[start+1:], '/') // next '/' after 0th rune
	if end == -1 {
		return path[start:], -1
	}
	return path[start : start+end+1], start + end + 1
}

// PathTrie is a trie of paths with string keys and interface{} values.

// PathTrie is a trie of string keys and interface{} values. Internal nodes
// have nil values so stored nil values cannot be distinguished and are
// excluded from walks. By default, PathTrie will segment keys by forward
// slashes with PathSegmenter (e.g. "/a/b/c" -> "/a", "/b", "/c"). A custom
// StringSegmenter may be used to customize how strings are segmented into
// nodes. A classic trie might segment keys by rune (i.e. unicode points).
type PathTrie struct {
	segmenter StringSegmenter // key segmenter, must not cause heap allocs
	value     interface{}
	children  map[string]*PathTrie
}

// PathTrieConfig for building a path trie with different segmenter
type PathTrieConfig struct {
	Segmenter StringSegmenter
}

// NewPathTrie allocates and returns a new *PathTrie.
func NewPathTrie() *PathTrie {
	return &PathTrie{
		segmenter: PathSegmenter,
	}
}

// NewPathTrieWithConfig allocates and returns a new *PathTrie with the given *PathTrieConfig
func NewPathTrieWithConfig(config *PathTrieConfig) *PathTrie {
	segmenter := PathSegmenter
	if config != nil && config.Segmenter != nil {
		segmenter = config.Segmenter
	}

	return &PathTrie{
		segmenter: segmenter,
	}
}

// newPathTrieFromTrie returns new trie while preserving its config
func (trie *PathTrie) newPathTrie() *PathTrie {
	return &PathTrie{
		segmenter: trie.segmenter,
	}
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (trie *PathTrie) Get(key string) interface{} {
	node := trie
	var anyMatchNode PathTrie
	trimQueryParams := strings.Split(key, "?")
	trimKey := trimQueryParams[0]
	for part, i := trie.segmenter(trimKey, 0); part != ""; part, i = trie.segmenter(trimKey, i) {
		if node.children["/*"] != nil {
			anyMatchNode = *node.children["/*"]
		} else if node.children["/**"] != nil {
			anyMatchNode = *node.children["/**"]
		}
		node = node.children[part]
		//i=-1代表路径末尾
		if node == nil || (node.value == nil && i == -1) {
			return anyMatchNode.value
		}
	}
	return node.value
}

// Put inserts the value into the trie at the given key, replacing any
// existing items. It returns true if the put adds a new value, false
// if it replaces an existing value.
// Note that internal nodes have nil values so a stored nil value will not
// be distinguishable and will not be included in Walks.
func (trie *PathTrie) Put(key string, value interface{}) bool {
	node := trie
	for part, i := trie.segmenter(key, 0); part != ""; part, i = trie.segmenter(key, i) {
		child, _ := node.children[part]
		if child == nil {
			if node.children == nil {
				node.children = map[string]*PathTrie{}
			}
			child = trie.newPathTrie()
			node.children[part] = child
		}
		node = child
	}
	// does node have an existing value?
	isNewVal := node.value == nil
	node.value = value
	return isNewVal
}

// Delete removes the value associated with the given key. Returns true if a
// node was found for the given key. If the node or any of its ancestors
// becomes childless as a result, it is removed from the trie.
func (trie *PathTrie) Delete(key string) bool {
	var path []nodeStr // record ancestors to check later
	node := trie
	for part, i := trie.segmenter(key, 0); part != ""; part, i = trie.segmenter(key, i) {
		path = append(path, nodeStr{part: part, node: node})
		node = node.children[part]
		if node == nil {
			// node does not exist
			return false
		}
	}
	// delete the node value
	node.value = nil
	// if leaf, remove it from its parent's children map. Repeat for ancestor path.
	if node.isLeaf() {
		// iterate backwards over path
		for i := len(path) - 1; i >= 0; i-- {
			parent := path[i].node
			part := path[i].part
			delete(parent.children, part)
			if !parent.isLeaf() {
				// parent has other children, stop
				break
			}
			parent.children = nil
			if parent.value != nil {
				// parent has a value, stop
				break
			}
		}
	}
	return true // node (internal or not) existed and its value was nil'd
}

// PathTrie node and the part string key of the child the path descends into.
type nodeStr struct {
	node *PathTrie
	part string
}

func (trie *PathTrie) isLeaf() bool {
	return len(trie.children) == 0
}


