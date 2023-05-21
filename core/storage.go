package core

type Storage interface{
	Put(*Block) error
}


// MemoryStore is a simple in-memory data store.
type MemoryStore struct {
    // Add any necessary fields here.
}


// NewMemoryStorage returns a new instance of MemoryStore.
func NewMemoryStorage() *MemoryStore {
    return &MemoryStore{}
}



// Put saves the given block to the memory store.
func (s *MemoryStore) Put(b *Block) error {
    // TODO: Implement actual logic here.
    return nil
}
