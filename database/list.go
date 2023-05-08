package database

import (
	List "go_redis/datastruct/list"
	databaseInterface "go_redis/interface/database"
	"go_redis/interface/resp"
	"go_redis/lib/utils"
	"go_redis/resp/reply"
	"strconv"
)

func init() {
	RegisterCommand("LLEN", execLLen, 2)
	RegisterCommand("LINDEX", execLIndex, 3)
	RegisterCommand("RPUSH", execRPush, -3)
	RegisterCommand("RPUSHX", execRPushX, -3)
	RegisterCommand("LPUSH", execLPush, -3)
	RegisterCommand("LPUSHX", execLPushX, -3)
	RegisterCommand("LPOP", execLPop, 2)
	RegisterCommand("RPOP", execRPop, 2)
	RegisterCommand("LSET", execLSet, 4)
	RegisterCommand("LREM", execLRem, 4)
	RegisterCommand("RPOPLPUSH", execRPopLpush, 3)
	RegisterCommand("LRANGE", execLRange, 4)
}

func (DB *DB) getAsList(key string) (list *List.Link, errReply reply.ErrorReply) {
	entity, ok := DB.GetEntity(key)
	if !ok {
		return nil, nil
	}
	link, ok := entity.Data.(*List.Link)
	if !ok {
		return nil, &reply.WrongTypeErrReply{}
	}
	return link, nil
}

func (DB *DB) getOrInitList(key string) (list *List.Link, isNew bool, errReply reply.ErrorReply) {
	list, err := DB.getAsList(key)
	if err != nil {
		return nil, false, err
	}
	if list == nil {
		list = List.Make()
		DB.PutEntity(key, &databaseInterface.DataEntity{Data: list})
		isNew = true
	}
	return list, isNew, nil
}

func execLLen(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	list, err := db.getAsList(key)
	if err != nil {
		return err
	}
	if list == nil {
		return reply.MakeIntReply(0)
	}
	len := list.Len()
	return reply.MakeIntReply(int64(len))
}

func execLIndex(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	index64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeNullBulkReply()
	}
	index := int(index64)
	size := list.Len()
	if index < -1*size {
		return reply.MakeNullBulkReply()
	} else if index < 0 {
		index += size
	} else if index >= size {
		return reply.MakeNullBulkReply()
	}
	val, _ := list.Get(index).([]byte)
	return reply.MakeBulkReply(val)
}

func execLPop(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeNullBulkReply()
	}
	val, _ := list.Remove(0).([]byte)
	if list.Len() == 0 {
		db.Remove(key)
	}
	return reply.MakeBulkReply(val)
}

// execLpush 将一个或多个值插入到列表头部
// LPUSH KEY_NAME VALUE1.. VALUEN
func execLPush(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	values := args[1:]
	list, _, errReply := db.getOrInitList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeNullBulkReply()
	}
	for _, value := range values {
		list.Insert(0, value)
	}
	return reply.MakeIntReply(int64(list.Len()))
}

// execlPushx 和LPush相似但是只有当key存在时才会执行LPush
// LPUSHX KEY_NAME VALUE1.. VALUEN
func execLPushX(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	values := args[1:]
	list, isNew, errReply := db.getOrInitList(key)
	if errReply != nil {
		return errReply
	}
	if isNew {
		return reply.MakeIntReply(0)
	}
	for _, value := range values {
		list.Insert(0, value)
	}
	return reply.MakeIntReply(int64(list.Len()))
}

func execLRange(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	start64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	start := int(start64)
	end64, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	stop := int(end64)
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeNullBulkReply()
	}
	size := list.Len()
	if start < -1*size {
		start = 0
	} else if start < 0 {
		start += size
	} else if start >= size {
		return reply.MakeNullBulkReply()
	}
	if stop < -1*size {
		stop = 0
	} else if stop < size {
		stop = size + 1
	} else if stop < 0 {
		stop += size + 1
	} else {
		stop = size
	}
	if stop < start {
		stop = start
	}
	slice := list.Range(start, stop)
	result := make([][]byte, len(slice))
	for i, val := range slice {
		result[i] = val.([]byte)
	}
	return reply.MakeMultiBulkReply(result)
}

func execLRem(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	count64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	value := string(args[2])
	count := int(count64)
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeNullBulkReply()
	}
	var rmCount int
	if count == 0 {
		rmCount = list.RemoveAllByVal(func(val interface{}) bool {
			return utils.Equals(value, val)
		})
	} else if count > 0 {
		//从左到右删除
		rmCount = list.RemoveByVal(func(val interface{}) bool {
			return utils.Equals(value, val)
		}, count)
	} else {
		//从右到左删除
		rmCount = list.ReverseRemoveByVal(func(val interface{}) bool {
			return utils.Equals(value, val)
		}, count)
	}
	if list.Len() == 0 {
		db.Remove(key)
	}
	return reply.MakeIntReply(int64(rmCount))
}

func execLSet(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	index64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}
	index := int(index64)
	value := args[2]

	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeErrReply("ERR no such key")
	}
	size := list.Len()
	if index < -1*size {
		return reply.MakeErrReply("ERR index out of range")
	} else if index >= size {
		return reply.MakeErrReply("ERR index out of range")
	} else if index < 0 {
		index += size
	}
	list.Set(index, value)
	return reply.MakeOkReply()
}

func execRPop(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeNullBulkReply()
	}
	val := list.RemoveLast()
	if list.Len() == 0 {
		db.Remove(key)
	}
	return reply.MakeBulkReply(val.([]byte))
}

// execRPopLpush 移除列表的最后一个元素，并将该元素添加到另一个列表并返回
func execRPopLpush(db *DB, args [][]byte) (resp resp.Reply) {
	sourceKey := string(args[0])
	targetKey := string(args[1])
	sourceList, errReply := db.getAsList(sourceKey)
	if errReply != nil {
		return errReply
	}
	if sourceList == nil {
		return reply.MakeNullBulkReply()
	}
	targetList, _, errReply := db.getOrInitList(targetKey)
	if errReply != nil {
		return errReply
	}
	val := sourceList.RemoveLast()
	if sourceList.Len() == 0 {
		db.Remove(sourceKey)
	}
	targetList.Insert(0, val)
	return reply.MakeBulkReply(val.([]byte))
}

func execRPush(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	values := args[1:]
	list, _, errReply := db.getOrInitList(key)
	if errReply != nil {
		return errReply
	}
	for _, value := range values {
		list.Add(value)
	}
	return reply.MakeIntReply(int64(list.Len()))
}

func execRPushX(db *DB, args [][]byte) (resp resp.Reply) {
	key := string(args[0])
	values := args[1:]
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return reply.MakeIntReply(0)
	}
	for _, value := range values {
		list.Add(value)
	}
	return reply.MakeIntReply(int64(list.Len()))
}
