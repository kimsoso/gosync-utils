// 处理目录结构相关的函数
package action

type Dirs []*Dirfile

func (d Dirs) Len() int {
	return len(d)
}

func (d Dirs) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Dirs) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return d[i].Path < d[j].Path
}

func (d Dirs) ListLastNodes() (out Dirs) {
	out = Dirs{}
	for _, dfile := range d {
		if !dfile.IsDir || len(dfile.Files) == 0 {
			out = append(out, dfile)
		} else {
			out = append(out, dfile.Files.ListLastNodes()...)
		}
	}
	return out
}

// 获取两个目录结构的所有变化
func (d Dirs) ChangesLoop(a Dirs) (adds, ups, dels Dirs) {
	_, _, adds, ups, dels = d.ListLastNodes().Changes(a.ListLastNodes())
	return
}

// 目的目录与当前目录的变化.
// sames: afile, adds: afile, ups: afile, dels: dfile
func (d Dirs) Changes(a Dirs) (dsames, sames, adds, ups, dels Dirs) {
	dsames, sames, adds, ups, dels = Dirs{}, Dirs{}, Dirs{}, Dirs{}, Dirs{}

	dmap := make(map[string]*Dirfile, len(d))
	for _, dd := range d {
		dmap[dd.Path] = dd
	}

	for _, afile := range a {
		if dfile, ok := dmap[afile.Path]; ok {
			if afile.IsDir != dfile.IsDir || afile.Size != dfile.Size || afile.Md5 != dfile.Md5 {
				if dfile.IsDir == afile.IsDir {
					// 这种情况一定是文件，目录不存在更新的情况。
					ups = append(ups, afile)
				} else {
					dels = append(dels, dfile)
					adds = append(adds, afile)
				}
			} else {
				dsames = append(dsames, dfile)
				sames = append(sames, afile)
			}
		} else {
			adds = append(adds, afile)
		}

		delete(dmap, afile.Path)
	}

	for _, dfile := range dmap {
		dels = append(dels, dfile)
	}

	return
}
