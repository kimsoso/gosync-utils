package dirfile

import (
	"btsync-utils/libs/action"
	"btsync-utils/libs/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func IsValidfile(path string) (fi os.FileInfo, valid bool) {
	if IsValidfilename(path) {
		fi, err := os.Lstat(path)
		return fi, os.IsNotExist(err) || fi.Mode().IsRegular() || fi.Mode().IsDir()
	}
	return nil, false
}

func IsValidfilename(path string) bool {
	return !strings.HasSuffix(path, ".swp") && !strings.HasSuffix(path, "~") && !strings.HasSuffix(path, ".swx") && path != ".DS_Store"
}

// 获取目录中的所有文件结构
func GetDirStruct(basePath, dirPath string) (out []*action.Dirfile, err error) {
	if rootRelPath, err := filepath.Rel(basePath, dirPath); err != nil {
		return nil, err
	} else {
		out = []*action.Dirfile{}

		if des, err := os.ReadDir(dirPath); err != nil {
			return nil, err
		} else {
			for _, de := range des {
				newRealFile := filepath.Join(dirPath, de.Name())
				newRelFile := filepath.Join(rootRelPath, de.Name())

				if IsValidfilename(newRealFile) {
					md5, size := "", int64(0)

					if de.Type().IsRegular() {
						md5, _ = utils.FileMD5(newRealFile)
						_fi, _ := de.Info()
						size = _fi.Size()
					}

					dirfile := action.Dirfile{
						Path:  newRelFile,
						IsDir: de.Type().IsDir(),
						Md5:   md5,
						Size:  size,
						Files: []*action.Dirfile{},
					}

					out = append(out, &dirfile)

					if de.Type().IsDir() {
						if subDirfile, err := GetDirStruct(basePath, newRealFile); err == nil {
							dirfile.Files = subDirfile
						}
					}
				}
			}
		}

		return out, nil
	}
}

// 打印目录文件结构信息
func PrintDirStruct(dirfiles []*action.Dirfile) {
	for _, dirfile := range dirfiles {
		fmt.Println("path:", dirfile.Path, "md5:", dirfile.Md5, "isdir:", dirfile.IsDir, "size:", dirfile.Size)
		if dirfile.IsDir {
			PrintDirStruct(dirfile.Files)
		}
	}
}

// 罗列目录下的所有文件以及目录
func GetAllFiles(pathname string) (files []string) {
	pathname = filepath.Clean(pathname)
	rd, err := os.ReadDir(pathname)
	if err != nil {
		return nil
	}
	newPath := ""
	if err == nil {
		for _, fi := range rd {
			if fi.Name() != ".DS_Store" {
				newPath = filepath.Join(pathname, fi.Name())
				files = append(files, newPath)
				if fi.IsDir() {
					files = append(files, GetAllFiles(newPath)...)
				}
			}
		}
	}
	return
}

// 以提供路径的方式得到结构信息
func DirfileByPath(dirs []*action.Dirfile, path string) *action.Dirfile {
	path = filepath.Clean(path)

	for _, dirfile := range dirs {
		if dirfile.Path == path {
			return dirfile
		} else if strings.HasPrefix(path, dirfile.Path) {
			return DirfileByPath(dirfile.Files, path)
		}
	}
	return nil
}

// 得到两个目录结构的不同文件信息
// map[path]ACT_STAT
func DiffDirs(baseDir, targetDir action.Dirs) (newopts []*action.Operation) {
	adds, ups, dels := baseDir.ChangesLoop(targetDir)

	newopts = make([]*action.Operation, 0, len(adds)+len(ups)+len(dels))

	for _, dfile := range adds {
		newopts = append(newopts, &action.Operation{
			Path:  dfile.Path,
			Act:   uint8(action.ACT_ADD),
			IsDir: dfile.IsDir,
			Md5:   dfile.Md5,
			Size:  dfile.Size,
		})
	}
	for _, dfile := range ups {
		newopts = append(newopts, &action.Operation{
			Path:  dfile.Path,
			Act:   uint8(action.ACT_ADD),
			IsDir: dfile.IsDir,
			Md5:   dfile.Md5,
			Size:  dfile.Size,
		})
	}
	for _, dfile := range dels {
		newopts = append(newopts, &action.Operation{
			Path:  dfile.Path,
			Act:   uint8(action.ACT_ADD),
			IsDir: dfile.IsDir,
			Md5:   dfile.Md5,
			Size:  dfile.Size,
		})
	}

	return
}
