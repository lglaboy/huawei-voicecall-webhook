package datastore

import (
	"fmt"
	"golang.org/x/exp/slog"
	"sync"
	"time"
)

// Data represents the data to be stored.
type Data struct {
	Count    int
	Value    []string
	ExpireAt time.Time
}

var DataStore = make(map[string]Data)
var Mutex = &sync.RWMutex{}

// RecordedData 存储禁止发送的消息
var RecordedData [][]string
var recordedMutex = &sync.RWMutex{}

func SetData(key string, count int, value []string, ttl time.Duration) {
	Mutex.Lock()
	defer Mutex.Unlock()
	DataStore[key] = Data{
		Count:    count,
		Value:    value,
		ExpireAt: time.Now().Add(ttl),
	}
}

func GetData(key string) (int, []string, bool) {
	Mutex.RLock()
	defer Mutex.RUnlock()
	data, ok := DataStore[key]
	if !ok || time.Now().After(data.ExpireAt) {
		return 0, nil, false // Data has expired or doesn't exist.
	}
	return data.Count, data.Value, true
}

func DeleteData(key string) {
	Mutex.Lock()
	defer Mutex.Unlock()
	delete(DataStore, key)
}

// StartDataCleanup 定时清理Data数据
func StartDataCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			Mutex.Lock()
			now := time.Now()
			slog.Info(fmt.Sprintf("清理数据"))
			for key, data := range DataStore {
				// 清理过期数据
				if now.After(data.ExpireAt) {
					delete(DataStore, key)
				}
			}
			Mutex.Unlock()
		}
	}
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

// containsRecord 检查 RecordedData 中是否已存在相同的记录
func containsRecord(data []string) bool {
	for _, record := range RecordedData {
		if equalSlices(record, data) {
			return true
		}
	}
	return false
}

// AddRecord 添加记录
func AddRecord(data []string) {
	recordedMutex.Lock()
	defer recordedMutex.Unlock()
	// 检查是否已存在相同值
	if !containsRecord(data) {
		RecordedData = append(RecordedData, data)
	}
}

// findRecordIndex 查找 RecordedData 中指定记录的索引
func findRecordIndex(data []string) int {
	for i, record := range RecordedData {
		if equalSlices(record, data) {
			return i
		}
	}
	return -1
}

// RemoveRecord 删除记录
func RemoveRecord(data []string) {
	recordedMutex.Lock()
	defer recordedMutex.Unlock()

	// 检查是否存在该值
	index := findRecordIndex(data)
	if index != -1 {
		// 从切片中删除该值
		RecordedData = append(RecordedData[:index], RecordedData[index+1:]...)
	}
}

// RemoveRecordsByPrefix 根据前两个值查找删除所有
func RemoveRecordsByPrefix(prefix []string) {
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

	// 逆序删除匹配项
	for i := len(indexesToRemove) - 1; i >= 0; i-- {
		idx := indexesToRemove[i]
		RecordedData = append(RecordedData[:idx], RecordedData[idx+1:]...)
	}
}

// RemoveRecordAll 删除所有记录
func RemoveRecordAll() {
	// 清空 RecordedData
	recordedMutex.Lock()
	RecordedData = nil
	recordedMutex.Unlock()
}
