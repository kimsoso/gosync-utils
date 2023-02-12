// 处理目录结构相关的函数
package action

import (
	"btsync-utils/libs/utils"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vmihailenco/msgpack"
	"gorm.io/gorm/clause"
)

// 目录结构
type Dirfile struct {
	Path  string `yaml:"pt"`
	IsDir bool   `yaml:"id"`
	Md5   string `yaml:"md5"`
	Size  int64  `yaml:"sz"`
	Files Dirs   `yaml:"fs"`
}

var (
	DirStruct = []*Dirfile{}
)

// 读取目录结构
func (act *Action) GetDirStruct() Dirs {
	if len(DirStruct) > 0 {
		return DirStruct
	}

	dirfiles := &File{}
	if err := act.db.First(dirfiles, 1).Error; err == nil {
		dirbytes := utils.S2B(dirfiles.Content)
		if err := msgpack.Unmarshal(dirbytes, &DirStruct); err == nil {
			return DirStruct
		}
	}

	return DirStruct
}

func (act *Action) PrintDirStruct(dirs []*Dirfile) {
	for _, dir := range dirs {
		log.Println("|", dir.Path, "---isDir:", dir.IsDir)
		act.PrintDirStruct(dir.Files)
	}
}

func (act *Action) SaveCurrentDir() {
	act.SaveDirStructToDB(act.GetDirStruct())
}

// 保存目录结构
func (act *Action) SaveDirStructToDB(dirfile Dirs) (err error) {
	content, err := msgpack.Marshal(&dirfile)
	if err != nil {
		return err
	}

	dirfiles := &File{ID: 1, Content: utils.B2S(content)}

	return act.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(dirfiles).Error
}

// 目录中加入一个节点，并保存目录结构
func (act *Action) RemoveDirNode(node *Dirfile) {
	dirstruct := act.GetDirStruct()
	nodes := &dirstruct
LOOP:
	for i, nnode := range *nodes {
		if nnode.Path == node.Path {
			*nodes = append((*nodes)[:i], (*nodes)[i+1:]...)
		} else if strings.HasPrefix(node.Path, nnode.Path) {
			nodes = &nnode.Files
			goto LOOP
		}
	}

	act.SaveDirStructToDB(dirstruct)

	DirStruct = []*Dirfile{}
}
func (act *Action) AppendDirNode(node *Dirfile) {
	dirstruct := act.GetDirStruct()
	dirnames := strings.Split(node.Path, string(filepath.Separator))
	currentDir := ""
	nodes := &dirstruct

	for _, dir := range dirnames {
		currentDir = filepath.Join(currentDir, dir)
		nodes = act.appendDirNode(nodes, &Dirfile{
			Path:  currentDir,
			IsDir: utils.If(currentDir == node.Path, node.IsDir, true),
			Md5:   utils.If(currentDir == node.Path, node.Md5, ""),
			Size:  utils.If(currentDir == node.Path, node.Size, 0),
			Files: node.Files,
		})
	}

	act.SaveDirStructToDB(dirstruct)

	DirStruct = []*Dirfile{}
}

func (act *Action) appendDirNode(nodes *Dirs, inNode *Dirfile) (out *Dirs) {
	for _, node := range *nodes {
		if node.Path == inNode.Path {
			return &node.Files
		}
	}
	*nodes = append(*nodes, inNode)
	sort.Sort(*nodes)
	return &inNode.Files
}

// 从暂存的目录结构表中搜索提供的文件
func (act *Action) GetDirFile(relPath string) *Dirfile {
	relPath = filepath.Clean(relPath)
	return act._GetDirFile(act.GetDirStruct(), relPath)
}

func (act *Action) _GetDirFile(dirStruct Dirs, relPath string) *Dirfile {
	for _, dir := range dirStruct {
		if dir.Path == relPath {
			return dir
		} else if strings.HasPrefix(relPath, dir.Path) {
			return act._GetDirFile(dir.Files, relPath)
		}
	}
	return nil
}
