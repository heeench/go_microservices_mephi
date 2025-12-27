package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/time/rate"
)

var (
	// Используем json-iterator для более быстрой сериализации
	jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary
	
	// Пул буферов для JSON encoding
	bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	
	// Пул JSON энкодеров
	encoderPool = sync.Pool{
		New: func() interface{} {
			return jsoniter.NewEncoder(nil)
		},
	}
)

// NewRateLimiter builds a limiter with given rps and burst.
func NewRateLimiter(rps float64, burst int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(rps), burst)
}

// RateLimitMiddleware applies request limiting to handlers.
func RateLimitMiddleware(limiter *rate.Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// WriteJSON оптимизированная функция для записи JSON ответов
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	
	// Используем пул буферов для снижения аллокаций
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)
	
	// Используем json-iterator для быстрой сериализации (в 2-3 раза быстрее стандартного)
	encoder := jsonAPI.NewEncoder(buf)
	if err := encoder.Encode(v); err != nil {
		// Fallback на стандартный encoder если что-то пошло не так
		_ = json.NewEncoder(w).Encode(v)
		return
	}
	
	// Устанавливаем статус и записываем данные
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

// DecodeJSON оптимизированная функция для декодирования JSON
func DecodeJSON(r *http.Request, v interface{}) error {
	// Ограничиваем размер тела запроса для безопасности
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB
	if err != nil {
		return err
	}
	defer r.Body.Close()
	
	// Используем json-iterator для быстрого декодирования
	return jsonAPI.Unmarshal(body, v)
}
