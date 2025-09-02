package in_memory_db

// Transaction type and operations
type OpType int

const (
	OpSet OpType = iota
	OpDelete
)

type Op struct {
	Type       OpType
	Key        string
	Value      []byte
	TTLSeconds int
}

// Tx is a simple transaction that applies ops atomically
type Tx struct {
	oplist []Op
}

func NewTx() *Tx { return &Tx{} }

func (t *Tx) Set(key string, value []byte, ttlSeconds int) {
	t.oplist = append(t.oplist, Op{Type: OpSet, Key: key, Value: append([]byte(nil), value...), TTLSeconds: ttlSeconds})
}
func (t *Tx) Delete(key string) { t.oplist = append(t.oplist, Op{Type: OpDelete, Key: key}) }
