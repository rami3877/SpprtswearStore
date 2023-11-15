package file 

/*
var offset int = 8

type FileNode struct {
	name string
	size uint64
	data []byte
	next *FileNode
	per  *FileNode
}

func (node *FileNode) InsrtByIndex(i int, name string, data []byte) {
	Node := new(FileNode)
	Node.data = data
	Node.name = name
	Node.size = uint64(len(data))
	root := &*node
	for s := 0; s != i; s++ {
		if root.next == nil {
			break
		}
		root = root.next
	}
	if root == node {
		copyFromRoot := FileNode{node.name, node.size, node.data, node.next, Node}
		root.name = Node.name
		root.data = Node.data
		root.size = Node.size
		root.next = &copyFromRoot
		root.per = Node

	} else if root.next == nil {
		Node.next = root
		Node.per = root.per
		root.per.next = Node
		root.per = Node
	} else {
		Node.per = root.per
		Node.next = root
		fmt.Println("root : ", root)
		fmt.Println("per :", root.per)

		fmt.Println("next :", root.next)
		fmt.Println()

	}

}

func (node *FileNode) DeleteByname(name string) *FileNode {

	root := &*node
	for {
		if root.next == nil || root.name == name {
			break
		}
		root = root.next
	}
	if root.name != name {
		return nil
	}

	if root == node {
		if root.next == nil {
			node.name = ""
			node.size = 0
			node.data = []byte{}
			node.next = nil
		} else {
			node.name = root.next.name
			node.size = root.next.size
			node.data = root.next.data
			node.next = root.next.next
			node.per = root.next
		}
	} else if root.next == nil {
		root.per.next = nil
	} else {
		root.per.next = root.next
		root.next.per = root.per
	}

	return root
}

func (node *FileNode) DeleteByIndex(i int) *FileNode {
	root := &*node
	for n := 0; n != i; n++ {
		if root.next == nil {
			break
		}
		root = root.next
	}

	if root == node {
		if root.next == nil {
			node.name = ""
			node.size = 0
			node.data = []byte{}
			node.next = nil
		} else {
			node.name = root.next.name
			node.size = root.next.size
			node.data = root.next.data
			node.next = root.next.next
			node.per = root.next
		}
	} else if root.next == nil {
		root.per.next = nil
	} else {
		root.per.next = root.next
		root.next.per = root.per
	}

	return root

}

func (node *FileNode) WritToFile(writer io.Writer) []byte {
	root := node
	alldata := []byte{}
	lenNode := 0
	for root != nil {
		alldata = append(alldata, root.name...)
		alldata = append(alldata, 0)
		alldata = append(alldata, binary.LittleEndian.AppendUint64([]byte{}, uint64(len(root.data)))...)
		alldata = append(alldata, root.data...)
		root = root.next
		lenNode++
	}
	writer.Write(binary.LittleEndian.AppendUint64([]byte{}, uint64(lenNode)))
	writer.Write(alldata)

	return alldata
}

func (node *FileNode) BytesALL() []byte {
	root := node
	alldata := []byte{}
	for root != nil {
		alldata = append(alldata, root.name...)
		alldata = append(alldata, 0)
		alldata = append(alldata, binary.LittleEndian.AppendUint64([]byte{}, uint64(len(root.data)))...)
		alldata = append(alldata, root.data...)
		root = root.next
	}

	return alldata
}

func NewFileNodeFormFile(File *os.File) *FileNode {
	info, _ := File.Stat()
	if info.Size() == 0 {
		return nil
	} else {
		return makeNodes(File)
	}

}

func NewFileRoot(name string, data []byte) *FileNode {
	root := new(FileNode)
	root.name = name
	root.size = uint64(len(data))
	root.data = data
	root.per = root
	return root
}

func (fileNode *FileNode) Insrt(name string, data []byte) {
	root := fileNode

	for root.next != nil {
		root = root.next
	}
	root.next = new(FileNode)
	root.next.name = name
	root.next.size = uint64(len(data))
	root.next.data = data
	root.next.per = root
}
func makeNodes(File *os.File) *FileNode {

	data, _ := io.ReadAll(File)
	lenLinkList := (binary.LittleEndian.Uint64(data[:8]))
	root := nextNode(data)
	root.per = root
	rootNext := root
	for i := 0; i < int(lenLinkList); i++ {
		rootNext.next = nextNode(data)
		rootNext.next.per = rootNext
		rootNext = rootNext.next

	}
	offset = 8
	return root
}

func nextNode(data []byte) *FileNode {
	fileNode := FileNode{}
	for _, v := range data[offset:] {
		if v == 0 {
			break
		}
		fileNode.name += string(v)
		offset++
	}
	offset += 1
	fileNode.size = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	fileNode.data = data[offset : offset+int(fileNode.size)]
	offset += int(fileNode.size)

	return &fileNode
}

func (fileNode *FileNode) Data() (string, uint64, []byte) {
	return fileNode.name, fileNode.size, fileNode.data
}

func (node *FileNode) Next() *FileNode {
	return node.next
}
func (node *FileNode) Per() *FileNode {
	return node.per
}
*/






