package base

import (
	"container/list"
)

var keyList *list.List = list.New()

func ContainsKey(clientKey string) (bool, *list.Element) {
	compare := false
	var find *list.Element = nil

	for i := keyList.Front(); i != nil; i.Next() {
		if i.Value == clientKey {
			compare = true
			find = i
			break
		}
	}

	return compare, find
}

func RemoveKeyIfExsist(clientKey string) {
	contains, find := ContainsKey(clientKey)

	if contains {
		keyList.Remove(find)
	}
}

func GetKeyOrNil(clientKey string) string {
	_, find := ContainsKey(clientKey)

	if find == nil {
		return ""
	}

	return find.Value.(string)
}

func AddNewKey(clientKey string) bool {
	success := false
	contains, _ := ContainsKey(clientKey)

	if !contains {
		keyList.PushBack(clientKey)
		success = true
	}

	return success
}
