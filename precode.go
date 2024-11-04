package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var cafeList = map[string][]string{
	"moscow": {"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

// Тесты для mainHandle
func TestMainHandlerWhenCountExceedsAvailableCafes(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", responseRecorder.Code)
	}
	cafes := strings.Split(responseRecorder.Body.String(), ",")
	if len(cafes) != 4 {
		t.Errorf("Expected 4 cafes, got %d", len(cafes))
	}
}

func TestMainHandlerWhenCityIsNotSupported(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=2&city=unknown_city", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "wrong city value" {
		t.Errorf("Expected 'wrong city value', got %s", responseRecorder.Body.String())
	}
}

func TestMainHandlerWhenCountIsMissing(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?city=moscow", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "count missing" {
		t.Errorf("Expected 'count missing', got %s", responseRecorder.Body.String())
	}
}
