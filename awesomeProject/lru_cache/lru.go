package lru_cache

type dll struct {
	key  int
	val  int
	next *dll
	prev *dll
}

type LRUCache struct {
	head     *dll
	tail     *dll
	cache    map[int]*dll
	capacity int
}

// ------------------- LRUCache ------------------
// Tail <-> ... <-> Node <-> Head
// Head is most recently used

func Constructor(capacity int) LRUCache {
	cache := make(map[int]*dll)
	lru := LRUCache{
		head:     &dll{0, 0, nil, nil},
		tail:     &dll{0, 0, nil, nil},
		cache:    cache,
		capacity: capacity,
	}
	lru.head.prev = lru.tail
	lru.tail.next = lru.head
	return lru
}

func (this *LRUCache) addNode(node *dll) {
	node.next = this.head
	node.prev = this.head.prev
	this.head.prev.next = node
	this.head.prev = node
}

func (this *LRUCache) removeNode(node *dll) {
	pr, nx := node.prev, node.next
	nx.prev = pr
	pr.next = nx
}

func (this *LRUCache) moveToHead(node *dll) {
	this.removeNode(node)
	this.addNode(node)
}

func (this *LRUCache) moveTail() {
	rm := this.tail.next
	delete(this.cache, rm.key)
	this.removeNode(rm)
}

func (this *LRUCache) Get(key int) int {
	if node, ok := this.cache[key]; ok {
		this.moveToHead(node)
		return node.val
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if node, ok := this.cache[key]; ok {
		node.val = value
		this.moveToHead(node)
		return
	}
	node := &dll{key, value, nil, nil}
	this.cache[key] = node
	if this.capacity > 0 {
		this.addNode(node)
		this.capacity--
		return
	}
	this.addNode(node)
	this.moveTail()
}
