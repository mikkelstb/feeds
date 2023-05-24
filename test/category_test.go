package feeds_test

import (
	"fmt"
	"testing"

	"github.com/mikkelstb/feeds/category/iptctop"
)

func TestCategory(t *testing.T) {

	cat := iptctop.Disaster
	fmt.Println(cat.GetName())
	fmt.Println(cat.GetDescription())
}
