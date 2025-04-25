package flowmanager

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yrn-go/yrn/pkg/yctx"
)

// RedisPluginStatusRepository implementa PluginStatusRepository usando Redis
type RedisPluginStatusRepository struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisPluginStatusRepository cria uma nova instância do repositório Redis
func NewRedisPluginStatusRepository(client *redis.Client, ttl time.Duration) *RedisPluginStatusRepository {
	return &RedisPluginStatusRepository{
		client: client,
		ttl:    ttl,
	}
}

// Save salva o status do plugin no Redis
func (r *RedisPluginStatusRepository) Save(ctx *yctx.Context, status PluginStatus) error {
	// Converte o status para JSON
	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal plugin status: %w", err)
	}

	// Cria a chave para o plugin
	key := fmt.Sprintf("plugin:status:%s", status.PluginID)

	// Salva no Redis com TTL
	err = r.client.Set(ctx.Context(), key, data, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save plugin status to Redis: %w", err)
	}

	return nil
}

// GetByPluginID recupera o status de um plugin específico
func (r *RedisPluginStatusRepository) GetByPluginID(ctx *yctx.Context, pluginID string) (PluginStatus, error) {
	key := fmt.Sprintf("plugin:status:%s", pluginID)

	data, err := r.client.Get(ctx.Context(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return PluginStatus{}, fmt.Errorf("plugin status not found for ID: %s", pluginID)
		}
		return PluginStatus{}, fmt.Errorf("failed to get plugin status from Redis: %w", err)
	}

	var status PluginStatus
	err = json.Unmarshal(data, &status)
	if err != nil {
		return PluginStatus{}, fmt.Errorf("failed to unmarshal plugin status: %w", err)
	}

	return status, nil
}

// GetAll recupera todos os status de plugins
func (r *RedisPluginStatusRepository) GetAll(ctx *yctx.Context) ([]PluginStatus, error) {
	// Busca todas as chaves que correspondem ao padrão
	keys, err := r.client.Keys(ctx.Context(), "plugin:status:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin status keys from Redis: %w", err)
	}

	var statuses []PluginStatus
	for _, key := range keys {
		data, err := r.client.Get(ctx.Context(), key).Bytes()
		if err != nil {
			continue // Ignora erros individuais
		}

		var status PluginStatus
		err = json.Unmarshal(data, &status)
		if err != nil {
			continue // Ignora erros de unmarshal
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// InMemoryPluginStatusRepository implementa PluginStatusRepository usando memória
type InMemoryPluginStatusRepository struct {
	statuses map[string]PluginStatus
	mu       sync.RWMutex
}

// NewInMemoryPluginStatusRepository cria uma nova instância do repositório em memória
func NewInMemoryPluginStatusRepository() *InMemoryPluginStatusRepository {
	return &InMemoryPluginStatusRepository{
		statuses: make(map[string]PluginStatus),
	}
}

// Save salva o status do plugin em memória
func (r *InMemoryPluginStatusRepository) Save(ctx *yctx.Context, status PluginStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.statuses[status.PluginID] = status
	return nil
}

// GetByPluginID recupera o status de um plugin específico
func (r *InMemoryPluginStatusRepository) GetByPluginID(ctx *yctx.Context, pluginID string) (PluginStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	status, exists := r.statuses[pluginID]
	if !exists {
		return PluginStatus{}, fmt.Errorf("plugin status not found for ID: %s", pluginID)
	}

	return status, nil
}

// GetAll recupera todos os status de plugins
func (r *InMemoryPluginStatusRepository) GetAll(ctx *yctx.Context) ([]PluginStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	statuses := make([]PluginStatus, 0, len(r.statuses))
	for _, status := range r.statuses {
		statuses = append(statuses, status)
	}

	return statuses, nil
}
