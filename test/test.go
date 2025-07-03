package main

import (
	"fmt"
	"sync"
)

var RecordedData [][]string
var recordedMutex = &sync.RWMutex{}

func addRecord(data []string) {
	recordedMutex.Lock()
	defer recordedMutex.Unlock()

	// 检查是否已存在相同值
	if !containsRecord(data) {
		RecordedData = append(RecordedData, data)
	}
}

func removeRecordsByPrefix(prefix []string) {
	recordedMutex.Lock()
	defer recordedMutex.Unlock()

	// 记录要删除的索引
	indexesToRemove := []int{}

	// 查找匹配项的索引
	for i, record := range RecordedData {
		if len(record) >= 2 && record[0] == prefix[0] && record[1] == prefix[1] {
			indexesToRemove = append(indexesToRemove, i)
		}
	}

	// 删除匹配项
	for _, idx := range indexesToRemove {
		RecordedData = append(RecordedData[:idx], RecordedData[idx+1:]...)
	}
}

// containsRecord 检查 RecordedData 中是否已存在相同的记录
func containsRecord(data []string) bool {
	for _, record := range RecordedData {
		if equalSlices(record, data) {
			return true
		}
	}
	return false
}

// equalSlices 比较两个切片是否相等
func equalSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func main() {
	// 示例1：添加记录
	addRecord([]string{"apple", "banana", "orange"})
	addRecord([]string{"apple", "banana", "grape"})
	addRecord([]string{"apple", "mango", "kiwi"})
	fmt.Println("RecordedData:", RecordedData)

	// 示例2：根据前两个字符串查找并删除匹配项
	removeRecordsByPrefix([]string{"apple", "banana"})
	fmt.Println("RecordedData after removal:", RecordedData)
}
