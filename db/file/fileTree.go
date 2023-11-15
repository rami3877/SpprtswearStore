package file

import (
	"encoding/binary"
	"io"
	"os"
)

type Node struct {
	Name string
	Size uint64
	Data []byte
	Next *Node
	Per  *Node
}

type FileNode struct {
	Len      int
	HeadFile *Node
	Tail     *Node
}

func InitFileNode() *FileNode {
	return &FileNode{}
}

func (n *FileNode) Append(node *Node) {
	if n.HeadFile == nil {
		node.Per = node
		n.HeadFile = node
		n.Tail = node
	} else {
		node.Per = &*n.Tail
		n.Tail.Next = &*node
		n.Tail = &*node
	}
	n.Len++
}
func (n *FileNode) InstrByIndexNode(index int, node *Node) {
	n.InsrtByIndex(index, node.Name, node.Data)
}

func (n *FileNode) InsrtByIndex(index int, name string, data []byte) {
	newNode := new(Node)
	newNode.Data = data
	newNode.Size = uint64(len(data))
	newNode.Name = name
	currentNode := n.HeadFile

	for i := 0; i != index; i++ {
		if currentNode == nil || currentNode.Next == nil {
			break
		}
		currentNode = currentNode.Next
	}
	if currentNode == nil {
		n.HeadFile = newNode
		n.Tail = newNode
		n.Len++
	} else if currentNode.Next == nil {
		newNode.Per = currentNode
		currentNode.Next = newNode
		n.Tail = newNode
		n.Len++
	} else if index == 0 {
		n.HeadFile.Per = newNode
		newNode.Next = &*n.HeadFile
		n.HeadFile = &*newNode
		n.Len++
	} else {
		newNode.Next = &*currentNode.Next
		newNode.Per = &*currentNode
		currentNode.Next = &*newNode

	}
}

func (node *FileNode) DeleteByName(name string) *Node {
	root := node.HeadFile
	for i := 0; i < node.Len; i++ {
		if root.Next == nil || root.Name == name {
			break
		}
		root = root.Next
	}
	if root.Name != name {
		return nil
	} else if root == node.HeadFile {
		node.HeadFile = &*root.Next
		root.Next.Per = &*root.Next
	} else if root == node.Tail {
		root.Per.Next = nil
	} else {
		root.Per.Next = *&root.Next
		root.Next.Per = *&root.Per
		return root
	}
	return root
}

func (node *FileNode) WritToFile(writer io.Writer) []byte {
	root := node.HeadFile
	alldata := []byte{}
	for root != nil {
		alldata = append(alldata, root.Name...)
		alldata = append(alldata, 0)
		alldata = append(alldata, binary.LittleEndian.AppendUint64([]byte{}, root.Size)...)
		alldata = append(alldata, root.Data...)
		root = root.Next
	}
	writer.Write(binary.LittleEndian.AppendUint64([]byte{}, uint64(node.Len)))
	writer.Write(alldata)

	return alldata
}

func (node *FileNode) DeleteByIndex(index int) *Node {
	if node.Len == 0 {
		return nil
	}
	if index == 0 {

		root := &*node.HeadFile
		if root.Next == nil {
			node.HeadFile = nil
		}
		node.HeadFile = node.HeadFile.Next
		node.HeadFile.Per = nil
		node.Len--
		return root
	} else if index >= node.Len-1 {
		root := &*node.Tail
		node.Tail.Per.Next = nil
		node.Tail = node.Tail.Per
		node.Len--
		return root
	} else {
		root := node.HeadFile
		for i := 0; i < index; i++ {
			if root.Next == nil {
				break
			}
			root = root.Next
		}
		root.Per.Next = root.Next
		root.Next.Per = root.Per
		node.Len--
		return root
	}

}

func NewFileNodeFormFile(File *os.File) *FileNode {
	info, _ := File.Stat()
	if info.Size() == 0 {
		return nil
	} else {
		return makeNodes(File)
	}

}

var offset = 8

func makeNodes(File *os.File) *FileNode {
	root := InitFileNode()
	File.Seek(0, 0)

	data, _ := io.ReadAll(File)
	Len := int(binary.LittleEndian.Uint64(data[:8]))
	for i := 0; i < Len; i++ {
		root.Append(&*nextNode(data))
		if offset >= len(data){
			 break
		}
	}
	offset = 8
	return root
}

func nextNode(data []byte) *Node {
	fileNode := Node{}
	for _, v := range data[offset:] {
		if v == 0 {
			break
		}
		fileNode.Name += string(v)
		offset++
	}
	offset += 1
	fileNode.Size = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	fileNode.Data = data[offset : offset+int(fileNode.Size)]
	offset += int(fileNode.Size)

	return &fileNode
}
