package locator

import (
	"fmt"
	"os"
	"testing"
)

// Your API Keys
var (
	googleAPIKey = os.Getenv("GOOGLE_API")
	yandexAPIKey = os.Getenv("YANDEX_API")
)

func TestLocator(t *testing.T) {
	google, err := New(Google, googleAPIKey)
	if err != nil {
		t.Fatal(err)
	}
	mozilla, err := New(Mozilla, "test")
	if err != nil {
		t.Fatal(err)
	}
	yandex, err := New(Yandex, yandexAPIKey)
	if err != nil {
		t.Fatal(err)
	}
	request := Request{
		CellTowers: []*CellTower{
			{250, 2, 7743, 22517, -78, 0, 0},
			{250, 2, 7743, 39696, -81, 0, 0},
			{250, 2, 7743, 22518, -91, 0, 0},
			{250, 2, 7743, 27306, -101, 0, 0},
			{250, 2, 7743, 29909, -103, 0, 0},
			{250, 2, 7743, 22516, -104, 0, 0},
			{250, 2, 7743, 20736, -105, 0, 0},
		},
		WifiAccessPoints: []*WifiAccessPoint{
			{"2:18:E4:C8:38:30", -22, 0, 0, 0},
		},
	}
	resp, err := google.Get(request)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Google:", resp)

	resp, err = mozilla.Get(request)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Mozilla:", resp)

	resp, err = yandex.Get(request)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Yandex:", resp)
}
