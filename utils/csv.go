package utils

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/exp/slog"
	"os"
)

type Person struct {
	Date  string
	Name  string
	Phone string
}

// GetDutyFromCSV 从csv值班表文件中获取当天值班人员
func GetDutyFromCSV(f string) ([]Person, error) {
	// 打开 CSV 文件（假设文件名为 data.csv）
	file, err := os.Open(f)
	if err != nil {
		slog.Error(fmt.Sprintf("无法打开文件: %s", err))
		return nil, err
	}
	defer file.Close()

	// 创建一个 CSV Reader
	reader := csv.NewReader(file)

	// 读取 CSV 文件的所有记录
	records, err := reader.ReadAll()
	if err != nil {
		slog.Error(fmt.Sprintf("读取 CSV 文件出错: %s", err))
		return nil, err
	}

	var people []Person

	// 遍历记录并处理
	for _, record := range records {
		// 检查记录是否包含至少 3 个字段
		if len(record) < 3 {
			fmt.Println("无效的记录:", record)
			slog.Error(fmt.Sprintln("无效的记录:", record))
			continue
		}

		person := Person{
			Date:  record[0],
			Name:  record[1],
			Phone: record[2],
		}
		people = append(people, person)
	}
	return people, nil
}
