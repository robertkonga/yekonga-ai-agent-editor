package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"yekonga-builder/agent"
	"yekonga-builder/icons"
)

// Define a unique service name for your app
const serviceName = "YekongaEditor"
const accountName = "default_user" // You can make this dynamic if your app has multiple accounts

// App struct
type App struct {
	ctx      context.Context
	agent    *agent.Agent
	configMu sync.RWMutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	apiKey, _ := a.LoadConfigFromFile("APIKey")
	apiKeyGemini, _ := a.LoadAPIKey("gemini")
	apiKeyAnthropic, _ := a.LoadAPIKey("anthropic")
	apiKeyDeepseek, _ := a.LoadAPIKey("deepseek")
	ollamaHost, _ := a.LoadOllamaHost()

	apiKeys := agent.ApiKeys{
		ApiKey:          apiKey,
		GeminiApiKey:    apiKeyGemini,
		AnthropicApiKey: apiKeyAnthropic,
		DeepseekApiKey:  apiKeyDeepseek,
	}

	a.ctx = ctx
	a.agent = agent.NewAgent(apiKeys, ollamaHost, &ctx)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ListIcons returns the list of available icons
func (a *App) ListIcons() ([]icons.Icon, error) {
	return icons.ListOfIcons()
}

// AgentChat calls the agent with session support
func (a *App) AgentChat(sessionID string, userInput string, provider string, model string) error {
	return a.agent.AgentChat(sessionID, userInput, provider, model)
}

func (a *App) ListSessions() ([]*agent.Session, error) {
	return a.agent.ListSessions()
}

func (a *App) ListWorkspaceSessions(workspace string) ([]*agent.Session, error) {
	return a.agent.ListWorkspaceSessions(workspace)
}

func (a *App) GetSession(id string) (*agent.Session, error) {
	return a.agent.GetSession(id)
}

// SaveAPIKey receives the key from the frontend and stores it securely
func (a *App) SaveAPIKey(name, value string) error {
	key := fmt.Sprintf("APIKey_%s", name)

	switch name {
	case "gemini":
		a.agent.GeminiApiKey = value
	case "deepseek":
		a.agent.DeepseekApiKey = value
	case "anthropic":
		a.agent.AnthropicApiKey = value
	}

	return a.SaveConfigToFile(key, value)
}

// LoadAPIKey retrieves the key from the OS keyring for backend use
func (a *App) LoadAPIKey(name string) (string, error) {
	key := fmt.Sprintf("APIKey_%s", name)

	return a.LoadConfigFromFile(key)
}

// SaveAPIKey receives the key from the frontend and stores it securely
func (a *App) SaveOllamaHost(value string) error {
	key := "OllamaHost"
	a.agent.OllamaHost = value

	return a.SaveConfigToFile(key, value)
}

// LoadAPIKey retrieves the key from the OS keyring for backend use
func (a *App) LoadOllamaHost() (string, error) {
	key := "OllamaHost"

	return a.LoadConfigFromFile(key)
}

// encryptionKey must be exactly 16, 24, or 32 bytes long to select
// AES-128, AES-192, or AES-256.
// KEEP THIS SECRET. Do not share this key outside your compiled app.
var encryptionKey = []byte("yekonga-editor-ultra-secure-key!")

// SaveConfigToFile updates a specific key-value pair in the encrypted binary file
func (a *App) SaveConfigToFile(key, value string) error {
	a.configMu.Lock()
	defer a.configMu.Unlock()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(configDir, serviceName)
	os.MkdirAll(appDir, os.ModePerm)
	filePath := filepath.Join(appDir, "config.bin")

	configMap := make(map[string]string)

	if encryptedData, err := os.ReadFile(filePath); err == nil {
		if decryptedData, err := decrypt(encryptedData, encryptionKey); err == nil {
			_ = json.Unmarshal(decryptedData, &configMap)
		}
	}

	configMap[key] = value

	jsonData, err := json.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("failed to format config data: %v", err)
	}

	newEncryptedData, err := encrypt(jsonData, encryptionKey)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}

	// Write to a temp file first, then atomically rename — prevents
	// a crash mid-write from leaving a corrupt config.bin
	tmpPath := filePath + ".tmp"
	if err := os.WriteFile(tmpPath, newEncryptedData, 0600); err != nil {
		return fmt.Errorf("failed to write temp config: %v", err)
	}
	return os.Rename(tmpPath, filePath)
}

// LoadConfigFromFile retrieves a specific key from the encrypted binary file
func (a *App) LoadConfigFromFile(key string) (string, error) {
	a.configMu.RLock()
	defer a.configMu.RUnlock()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(configDir, serviceName, "config.bin")

	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("config file not found: %v", err)
	}

	decryptedData, err := decrypt(encryptedData, encryptionKey)
	if err != nil {
		return "", fmt.Errorf("decryption failed (wrong key or corrupted file): %v", err)
	}

	configMap := make(map[string]string)
	if err := json.Unmarshal(decryptedData, &configMap); err != nil {
		return "", fmt.Errorf("failed to read config format: %v", err)
	}

	value, exists := configMap[key]
	if !exists {
		return "", fmt.Errorf("key '%s' not found", key)
	}

	return value, nil
}

// --- Internal Helper Functions for AES-GCM Encryption ---
func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
