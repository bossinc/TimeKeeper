package main

import (
	"regexp"
	"strings"
)

var reNewItems = regexp.MustCompile(` - \d+ new items`)

// treeNode is one node in the hierarchical window tree.
type treeNode struct {
	label    string
	timeMs   int64
	children []string // child keys
}

// treeData holds the flattened tree. Keys use tab as separator between levels.
type treeData struct {
	roots []string
	nodes map[string]*treeNode
}

// buildTreeData converts a flat WindowTime list into a 3-level tree.
// The rightmost " - " segment is the root (app name), working left toward the leaf.
func buildTreeData(windows []WindowTime) *treeData {
	td := &treeData{nodes: make(map[string]*treeNode)}
	for _, w := range windows {
		parts := strings.Split(reNewItems.ReplaceAllString(w.Label, ""), " - ")

		rootKey := parts[len(parts)-1]
		if td.nodes[rootKey] == nil {
			td.nodes[rootKey] = &treeNode{label: rootKey}
			td.roots = append(td.roots, rootKey)
		}
		if len(parts) == 1 {
			td.nodes[rootKey].timeMs += w.TimeMs
			continue
		}

		subLabel := parts[len(parts)-2]
		subKey := rootKey + "\t" + subLabel
		if td.nodes[subKey] == nil {
			td.nodes[subKey] = &treeNode{label: subLabel}
			td.nodes[rootKey].children = append(td.nodes[rootKey].children, subKey)
		}
		if len(parts) == 2 {
			td.nodes[subKey].timeMs += w.TimeMs
			continue
		}

		leafLabel := strings.Join(parts[:len(parts)-2], " - ")
		leafKey := subKey + "\t" + leafLabel
		if td.nodes[leafKey] == nil {
			td.nodes[leafKey] = &treeNode{label: leafLabel}
			td.nodes[subKey].children = append(td.nodes[subKey].children, leafKey)
		}
		td.nodes[leafKey].timeMs += w.TimeMs
	}
	// Roll up times to parents.
	for _, rk := range td.roots {
		root := td.nodes[rk]
		for _, sk := range root.children {
			sub := td.nodes[sk]
			for _, lk := range sub.children {
				sub.timeMs += td.nodes[lk].timeMs
			}
			root.timeMs += sub.timeMs
		}
	}
	return td
}
