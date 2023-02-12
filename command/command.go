package command

type REPLY = uint8
type COMMAND = uint8

// socket 返回命令格式
const (
	// 需要发送文件数据
	NEEDFILE REPLY = 200 + iota

	// 线程完毕
	THREAD_DONE
)

// socket 返回命令格式,提示性信息
const (
	// 正确
	OK REPLY = 0 + iota
	// 参数错误
	BADPARAM
	// 数据库错误
	DBERROR
	// 无效的操作类型
	INVLIDACT
	// 客户端写入磁盘错误
	DISKERROR
)
const (
	// BLOCK数据格式为：{BLOCK_OK}{DATA}
	// 数据块都存在
	BLOCK_FULL REPLY = 10 + iota

	// 数据块存在
	BLOCK_OK

	// 数据块不存在
	BLOCK_NO

	// 存在未被同步的CLIID, {FETCH_OK}[]operationClis
	FETCH_OK

	// 不存在未被同步的CLIID
	FETCH_NO

	// 不需要同步数据
	DETECT_OPT_NONEED

	// 需要同步数据 {DETECT_OPT_NEED}{LASTID:uint32}
	DETECT_OPT_NEED
)

const (
	// 随后是目录操作
	OPERATION COMMAND = 30 + iota

	// 客户单之间同步操作{FETCH}{uint32(STARTID)}
	FETCH

	//  侦测活跃客户端
	DETECT

	// 请求文件块,数据格式为：{REQ_BLOCK}{BLOCKINDEX}{filename}
	BLOCK

	// 询问是否需要与服务器同步数据{DETECT_OPT}{LASTID:uint32}
	DETECT_OPT
)

const (
	FILE_ABDT = 40 + iota // file already been deleted
)

const (
	CBLOCKS_START = 50 + iota
	CBLOCKS_DATA

	CBLOCKS_REVDATAS  = 54 // 传输组合命令数据块
	CBLOCKS_FSTART    = 55 // 文件开始传输
	CBLOCKS_FDONE     = 56 // 文件传输完毕
	CBLOCKS_FDATA     = 57 //文件块
	CBLOCKS_TRANSFILE = 58 //需要传输文件
	CBLOCKS_END       = 59
)
