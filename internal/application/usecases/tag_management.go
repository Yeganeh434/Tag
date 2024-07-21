package usecases

import (
	"strconv"
	"tag_project/internal/adapters/databases/mysql"
)

func GenerateKey(title string) (string, error) {
	num, err := mysql.TagDB.GetCounter()
	if err != nil {
		return "", err
	}
	stringifiedNum := strconv.Itoa(num)
	key := title + "_" + stringifiedNum
	return key, nil
}
