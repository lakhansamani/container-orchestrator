package providers

// Provider defines current memory store provider
type MemoryStoreProvider interface {
	// SetData sets the data in the memory store
	SetData(key, value string) error
	// GetData gets the data from the memory store
	GetData(key string) (string, error)
	// DeleteData deletes the data from the memory store
	DeleteData(key string) error
}
