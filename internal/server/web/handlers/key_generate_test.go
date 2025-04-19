package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sender/internal/server/web/handlers"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestKeysGenerateHandler(t *testing.T) {
	// Установка Gin в тестовый режим
	gin.SetMode(gin.TestMode)

	// Создание маршрутизатора и маршрута
	r := gin.Default()
	r.GET("/keys", handlers.KeysGenerateHandler)

	// Создание тестового запроса
	req, err := http.NewRequest(http.MethodGet, "/keys", nil)
	assert.NoError(t, err)

	// Создание записи для ответа
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверка кода ответа
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка структуры ответа
	var keys handlers.Keys
	err = json.Unmarshal(w.Body.Bytes(), &keys)
	assert.NoError(t, err)

	// Проверка содержимого ключей
	assert.NotEmpty(t, keys.PublicKey, "PublicKey должен быть непустым")
	assert.NotEmpty(t, keys.PrivateKey, "PrivateKey должен быть непустым")
}
