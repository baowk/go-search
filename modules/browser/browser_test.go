package browser

import (
	"fmt"
	"testing"
)

func TestGetProfileAndUa(t *testing.T) {
	prof, ua := GetProfileAndUa("firefox", "135")
	fmt.Println("profile:", prof, "ua", ua)
}
