package generates

import (
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func GenerateUniqueID(input string) string {
	id := ""
	for _, c := range input {
		id += strconv.Itoa(int(c))
	}
	return id
}
