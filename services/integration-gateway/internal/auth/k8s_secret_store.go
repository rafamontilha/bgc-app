package auth

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// KubernetesSecretStore implementa SecretStore usando Kubernetes Secrets API
type KubernetesSecretStore struct {
	clientset *kubernetes.Clientset
	namespace string
	cache     map[string]secretCacheEntry
	cacheMu   sync.RWMutex
	cacheTTL  time.Duration
}

type secretCacheEntry struct {
	value     string
	expiresAt time.Time
}

// NewKubernetesSecretStore cria um novo store conectado ao cluster Kubernetes
func NewKubernetesSecretStore(namespace string) (*KubernetesSecretStore, error) {
	// Tenta configuração in-cluster primeiro
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	// Valida acesso ao namespace
	_, err = clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to access namespace %s: %w", namespace, err)
	}

	store := &KubernetesSecretStore{
		clientset: clientset,
		namespace: namespace,
		cache:     make(map[string]secretCacheEntry),
		cacheTTL:  5 * time.Minute, // Cache secrets por 5 minutos
	}

	// Inicia goroutine para limpar cache expirado
	go store.cleanExpiredCache()

	return store, nil
}

// GetSecret busca um secret do Kubernetes
// ref formato: "secret-name/key-name" ou apenas "key-name" (usa env var como fallback)
func (s *KubernetesSecretStore) GetSecret(ref string) (string, error) {
	// Verifica cache primeiro
	if value, found := s.getFromCache(ref); found {
		return value, nil
	}

	// Parse ref
	parts := strings.Split(ref, "/")

	var secretName, keyName string
	if len(parts) == 2 {
		// Formato: "secret-name/key-name"
		secretName = parts[0]
		keyName = parts[1]
	} else if len(parts) == 1 {
		// Formato antigo (env var): "my-api-key" → SECRET_MY_API_KEY
		// Tenta primeiro como env var para backward compatibility
		envVar := "SECRET_" + strings.ToUpper(strings.ReplaceAll(ref, "-", "_"))
		if value := os.Getenv(envVar); value != "" {
			s.putInCache(ref, value)
			return value, nil
		}
		return "", fmt.Errorf("invalid secret ref format: %s (expected: secret-name/key-name)", ref)
	} else {
		return "", fmt.Errorf("invalid secret ref format: %s (expected: secret-name/key-name)", ref)
	}

	// Busca secret no Kubernetes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	secret, err := s.clientset.CoreV1().Secrets(s.namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s from namespace %s: %w", secretName, s.namespace, err)
	}

	// Extrai key
	valueBytes, exists := secret.Data[keyName]
	if !exists {
		return "", fmt.Errorf("key %s not found in secret %s", keyName, secretName)
	}

	value := string(valueBytes)

	// Armazena em cache
	s.putInCache(ref, value)

	return value, nil
}

// getFromCache busca valor do cache se ainda válido
func (s *KubernetesSecretStore) getFromCache(ref string) (string, bool) {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()

	entry, exists := s.cache[ref]
	if !exists {
		return "", false
	}

	// Verifica se expirou
	if time.Now().After(entry.expiresAt) {
		return "", false
	}

	return entry.value, true
}

// putInCache armazena valor no cache com TTL
func (s *KubernetesSecretStore) putInCache(ref, value string) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	s.cache[ref] = secretCacheEntry{
		value:     value,
		expiresAt: time.Now().Add(s.cacheTTL),
	}
}

// cleanExpiredCache limpa entradas expiradas do cache periodicamente
func (s *KubernetesSecretStore) cleanExpiredCache() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.cacheMu.Lock()
		now := time.Now()
		for ref, entry := range s.cache {
			if now.After(entry.expiresAt) {
				delete(s.cache, ref)
			}
		}
		s.cacheMu.Unlock()
	}
}

// InvalidateCache remove uma entrada específica do cache
func (s *KubernetesSecretStore) InvalidateCache(ref string) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	delete(s.cache, ref)
}

// InvalidateAllCache limpa todo o cache
func (s *KubernetesSecretStore) InvalidateAllCache() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.cache = make(map[string]secretCacheEntry)
}
