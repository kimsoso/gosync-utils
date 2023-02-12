package pool

import (
	"btsync-utils/libs/action"
	"sync"
)

var OperationPool = sync.Pool{
	New: func() interface{} {
		return new(action.Operation)
	},
}

var OperationsPool = sync.Pool{
	New: func() interface{} {
		return new([]action.Operation)
	},
}

var OperationCliPool = sync.Pool{
	New: func() interface{} {
		return new(action.OperationCli)
	},
}

var OperationClisPool = sync.Pool{
	New: func() interface{} {
		return new([]action.OperationCli)
	},
}

var IPool = sync.Pool{
	New: func() interface{} {
		return new([]uint)
	},
}
