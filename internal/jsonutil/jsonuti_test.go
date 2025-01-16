package jsonutil_test

import (
	"bytes"
	"sender/internal/jsonutil"
	"testing"
)

func TestObjectToJson(t *testing.T) {
	object := struct {
		ID     int
		Title  string
		Artist string
	}{
		ID:     101,
		Title:  "Hot Rats",
		Artist: "Frank Zappa",
	}

	resultObjectJson := `{"ID":101,"Title":"Hot Rats","Artist":"Frank Zappa"}`

	testText, err := jsonutil.ToJSON(object)
	if err != nil {
		t.Fatalf("Failed create object to JSON: %v", err)
	}

	if !bytes.Equal(testText, []byte(resultObjectJson)) {
		t.Fatal("Json incorrect")
	}
}

func TestNilToJson(t *testing.T) {
	text, err := jsonutil.ToJSON(nil)
	if err == nil {
		t.Fatalf("Failed create object to JSON: %v Text: %v", err, text)
	}
}

func TestObjectFromJson(t *testing.T) {
	resultObject := struct {
		ID     int
		Title  string
		Artist string
	}{
		ID:     101,
		Title:  "Hot Rats",
		Artist: "Frank Zappa",
	}

	objectJson := `{"ID":101,"Title":"Hot Rats","Artist":"Frank Zappa"}`
	var object struct {
		ID     int
		Title  string
		Artist string
	}

	err := jsonutil.FromJSON([]byte(objectJson), &object)
	if err != nil {
		t.Fatalf("Failed create object from JSON: %v", err)
	}

	if resultObject.ID != object.ID || resultObject.Title != object.Title || resultObject.Artist != object.Artist {
		t.Fatal("Object not equals")
	}
}

func TestEmptyJson(t *testing.T) {
	var object struct {
		ID     int
		Title  string
		Artist string
	}

	err := jsonutil.FromJSON([]byte(""), &object)
	if err == nil {
		t.Fatalf("Failed create object from JSON: %v, Object: %v", err, object)
	}
}

func TestNilObject(t *testing.T) {
	objectJson := `{"ID":101,"Title":"Hot Rats","Artist":"Frank Zappa"}`
	err := jsonutil.FromJSON([]byte(objectJson), nil)
	if err == nil {
		t.Fatalf("Failed create nil object from JSON: %v", err)
	}
}
