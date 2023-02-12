package main

import (
	"btsync-utils/libs/action"
	"btsync-utils/libs/utils"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// func main() {
// 	gdb := dao.NewDB()
// 	lastrow := &action.Operation{}
// 	if err := gdb.Last(lastrow).Error; err == nil {
// 		fmt.Println(lastrow.ID)
// 	} else {
// 		fmt.Println(err)
// 	}
// }

// func main() {
// 	writer := bytes.NewBuffer([]byte{})

// 	if err := binary.Write(writer, binary.BigEndian, uint8(200)); err != nil {
// 		fmt.Println(err)
// 	}

// 	if err := binary.Write(writer, binary.BigEndian, int64(65535)); err != nil {
// 		fmt.Println(err)
// 	}

// 	data := writer.Bytes()

// 	fmt.Println(data[0])

// 	fmt.Println(binary.BigEndian.Uint64(data[1:]))
// }

// func main() {
// 	opt := &action.Operation{
// 		ID:    1,
// 		Path:  "a.txt",
// 		Act:   0,
// 		IsDir: false,
// 		Md5:   "adafdsaf",
// 		Size:  0,
// 	}

// 	if data, err := msgpack.Marshal(opt); err != nil {
// 		println(data)
// 	} else {
// 		fmt.Println(data)
// 		opt1 := &action.Operation{}

// 		msgpack.Unmarshal(data, opt1)
// 		fmt.Printf("%v", opt1)
// 	}
// }

// func main() {
// 	ctx1, cc1 := context.WithCancel(context.Background())
// 	go r1(ctx1, 1)

// 	ctx2, cc2 := context.WithCancel(context.Background())
// 	go r1(ctx2, 2)

// 	time.Sleep(time.Second * 5)
// 	cc1()
// 	time.Sleep(time.Second * 5)
// 	cc2()
// 	time.Sleep(time.Second)
// }

// func r1(ctx context.Context, n int) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println(n, " exit")
// 			return
// 		default:
// 			fmt.Println(n, "running")
// 			time.Sleep(time.Second)
// 		}
// 	}
// }

var (
	a action.Dirs
	b action.Dirs
)

func AppendDirNode(nodes *action.Dirs, node *action.Dirfile) {
	dirnames := strings.Split(node.Path, string(filepath.Separator))
	currentDir := ""

	for _, dir := range dirnames {
		currentDir = filepath.Join(currentDir, dir)
		nodes = appendDirNode(nodes, &action.Dirfile{
			Path:  currentDir,
			IsDir: utils.If(currentDir == node.Path, node.IsDir, true),
			Md5:   utils.If(currentDir == node.Path, node.Md5, ""),
			Size:  utils.If(currentDir == node.Path, node.Size, 0),
			Files: node.Files,
		})
	}
}

func appendDirNode(nodes *action.Dirs, inNode *action.Dirfile) (out *action.Dirs) {
	for _, node := range *nodes {
		if node.Path == inNode.Path {
			return &node.Files
		}
	}
	*nodes = append(*nodes, inNode)
	sort.Sort(*nodes)
	return &inNode.Files
}

// func main() {
// 	// dirnames := strings.Split("opts.txt", string(filepath.Separator))
// 	// fmt.Println(dirnames)
// 	d1 := &action.Dirfile{
// 		Path:  "aaa.txt",
// 		IsDir: false,
// 	}
// 	// d2 := &action.Dirfile{
// 	// 	Path:  "bbb.txt",
// 	// 	IsDir: false,
// 	// }
// 	d3 := &action.Dirfile{
// 		Path:  "ccc",
// 		IsDir: true,
// 	}
// 	d4 := &action.Dirfile{
// 		Path:  "ddd.txt",
// 		IsDir: false,
// 	}
// 	d5 := &action.Dirfile{
// 		Path:  "eee",
// 		IsDir: true,
// 	}

// 	d6 := &action.Dirfile{
// 		Path:  "ccc",
// 		IsDir: true,
// 	}
// 	d7 := &action.Dirfile{
// 		Path:  "ddd.txt",
// 		IsDir: false,
// 	}
// 	d8 := &action.Dirfile{
// 		Path:  "eee.txt",
// 		IsDir: false,
// 	}
// 	d9 := &action.Dirfile{
// 		Path:  "eee",
// 		IsDir: true,
// 	}

// 	d6.Files = append(d6.Files, d7, d8, d9)
// 	d5.Files = append(d5.Files, d6)
// 	d3.Files = append(d3.Files, d4, d5)

// 	a = append(a, d1)

// 	fmt.Println("org:")
// 	for _, f := range a.ListLastNodes() {
// 		fmt.Println(f.Path, utils.If(f.IsDir, "[D]", "F"))
// 	}

// 	AppendDirNode(&a, &action.Dirfile{
// 		Path:  "ccc/ddd/opts.txt",
// 		IsDir: false,
// 		Md5:   "asdfasdfafds",
// 		Size:  39654,
// 		Files: action.Dirs{},
// 	})

// 	fmt.Println("current:")
// 	for _, f := range a.ListLastNodes() {
// 		fmt.Println(f.Path, utils.If(f.IsDir, "[D]", "F"))
// 	}

// 	gdb := dao.NewDB()
// 	act := action.New(gdb)

// 	if err := act.SaveDirStructToDB(a); err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("saved")
// 	}

// 	// a = append(a, d1, d2)
// 	// b = append(a, d11, d2)

// 	// aaaa := funk.Join(a, b, funk.InnerJoin)
// 	// for _, aa := range aaaa.([]*action.Dirfile) {
// 	// 	fmt.Println(aa.Path, aa.IsDir)
// 	// }
// }

func walkFunc(path string, info os.DirEntry, err error) error {
	fmt.Println(path, info.IsDir(), info.Name(), info.Type(), err)
	return nil
}

// func main() {
// 	// filepath.WalkDir("/home/soso/go/src/btsync-utils/libs", walkFunc)
// 	bm := bitmap.NewBitMap(5)

// 	fmt.Println(bm.Max())
// 	if bm.IsFull() {
// 		fmt.Println("1")
// 	} else {
// 		fmt.Println("0")
// 	}

// 	fmt.Println(bm.Serialize())

// 	os.WriteFile("a.bitmap", bm.Serialize(), os.ModePerm)
// }

type AS struct {
	ID   int
	Name string
}

// func main() {
// 	aaa := []*AS{}
// 	aaa = append(aaa,
// 		&AS{
// 			ID:   3,
// 			Name: "3",
// 		}, &AS{
// 			ID:   9,
// 			Name: "9",
// 		}, &AS{
// 			ID:   12,
// 			Name: "12",
// 		}, &AS{
// 			ID:   13,
// 			Name: "13",
// 		})

// 	bbb := []*AS{}
// 	bbb = append(bbb,
// 		&AS{
// 			ID:   1,
// 			Name: "1",
// 		}, &AS{
// 			ID:   5,
// 			Name: "5",
// 		}, &AS{
// 			ID:   20,
// 			Name: "20",
// 		}, &AS{
// 			ID:   16,
// 			Name: "16",
// 		})

// 	ccc := append(aaa, bbb...)

// 	sort.Slice(ccc, func(i int, j int) bool {
// 		return ccc[i].ID < ccc[j].ID
// 	})

//		for _, c := range ccc {
//			fmt.Println(c.ID, c.Name)
//		}
//	}
//
// ADDED: exclude symlink for temporary
func isInvalidfiles(path string) bool {
	if strings.HasSuffix(path, ".swp") || strings.HasSuffix(path, "~") || strings.HasSuffix(path, ".swx") || path == ".DS_Store" {
		return true
	} else {
		if fi, err := os.Lstat(path); err == nil {
			mod := fi.Mode()
			return !mod.IsRegular() && !mod.IsDir()
		} else {
			return false
		}
	}
}

var (
	filescount = 0
)

func walkDirFunc(path string, info os.DirEntry, err error) error {
	if err == nil && !isInvalidfiles(path) {
		filescount++
		fmt.Println(path)
	}
	return nil
}

// func main() {
// 	filepath.WalkDir("perl5", walkDirFunc)
// 	fmt.Println(filescount)
// 	// os.Symlink("bbb.txt", "blink.txt")
// 	// if fi, err := os.Lstat("alink.txt"); err != nil {
// 	// 	fmt.Println("can't get link")
// 	// } else {
// 	// 	switch mod := fi.Mode(); {
// 	// 	case mod.IsRegular():
// 	// 		fmt.Println("is regular file")
// 	// 	case mod.IsDir():
// 	// 		fmt.Println("is dir")
// 	// 	case mod&os.ModeSymlink != 0:
// 	// 		fmt.Println("is symlink")
// 	// 	case mod&os.ModeNamedPipe != 0:
// 	// 		fmt.Println("named pipe")
// 	// 	}
// 	// }
// }

// func main() {
// rd, _ := os.ReadDir(".")
// for _, ri := range rd {
// 	if !ri.IsDir() && !ri.Type().IsRegular() {
// 		fmt.Println(ri.Name())
// 	}
// }
// }

func main() {
	a := []uint{}

	done := map[uint]bool{}

	for i := 2; i < 80; i++ {
		if i%2 == 1 {
			done[uint(i)] = true
		}
	}

	for i := 0; i < 100; i++ {
		a = append(a, uint(i))
	}

	fmt.Println(a)

	for i, num := range a {
		if _, ok := done[uint(num)]; ok {
			a = append(a[:i], a[i+1:]...)
		}
	}

	fmt.Println(a)
}
