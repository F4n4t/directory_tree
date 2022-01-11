// Package directory_tree provides a way to generate a directory tree.
//
// Example usage:
//
//	tree, err := directory_tree.NewTree("/home/me")
//
// I did my best to keep it OS-independent but truth be told I only tested it
// on OS X and Debian Linux so YMMV. You've been warned.
package directory_tree

import (
	"os"
	"path/filepath"
	"time"
)

// FileInfo is a struct created from os.FileInfo interface for serialization.
type FileInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"mod_time"`
	IsDir   bool        `json:"is_dir"`
	Ext     string      `json:"extension"`
}

// fileInfoFromInterface Helper function to create a local FileInfo struct from os.FileInfo interface.
func fileInfoFromInterface(v os.FileInfo) *FileInfo {
	return &FileInfo{v.Name(), v.Size(), v.Mode(), v.ModTime(), v.IsDir(), filepath.Ext(v.Name())}
}

// Node represents a node in a directory tree.
type Node struct {
	FullPath string    `json:"path"`
	Info     *FileInfo `json:"info"`
	Children []*Node   `json:"children"`
	Parent   *Node     `json:"-"`
}

// NewTree Create directory hierarchy.
func NewTree(root string) (result *Node, err error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		parents[path] = &Node{
			FullPath: path,
			Info:     fileInfoFromInterface(info),
			Children: make([]*Node, 0),
		}
		return nil
	}
	if err = filepath.Walk(absRoot, walkFunc); err != nil {
		return
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			result = node
		} else {
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		}
	}
	return
}

// GetFiles return a slice with all files.
// Leave ext empty to return all the files.
func (n *Node) GetFiles(ext string) (files []*Node) {
	for _, child := range n.Children {
		if child.Info.IsDir {
			files = append(files, child.GetFiles(ext)...)
		} else if ext == "" || ext == child.Info.Ext {
			files = append(files, child)
		}
	}
	return
}
