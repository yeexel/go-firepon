package helpers

import "fmt"

func GetDocPath(col string, docID string) string {
	return fmt.Sprintf("%s/%s", col, docID)
}
