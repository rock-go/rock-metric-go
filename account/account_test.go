package account

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	d := GetAll()

	data, _ := json.Marshal(d)
	fmt.Printf("%s", data)
}
