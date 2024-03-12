package cache

import (
	"container/list"
)

type LRUCache struct {
	capacity   int
	doublyList *list.List
	elementMap map[interface{}]*list.Element // use hash table to check if list node exists
}

// Pair is the value of a list node. in order to contains any type, I use interface{}
type Pair struct {
	key   interface{}
	value interface{}
}

// Constructor initializes a new LRUCache.
func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity:   capacity,
		doublyList: new(list.List),
		elementMap: make(map[interface{}]*list.Element, capacity),
	}
}

// Get a list node from the hash map.
func (c *LRUCache) Get(key interface{}) interface{} {
	// check if list node exists
	if node, ok := c.elementMap[key]; ok {
		val := node.Value.(*list.Element).Value.(Pair).value
		// move node to front
		c.doublyList.MoveToFront(node)
		return val
	}
	return nil
}

// Put key and value in the LRUCache
func (c *LRUCache) Put(key interface{}, value interface{}) {
	// check if list node exists
	if node, ok := c.elementMap[key]; ok {
		// move the node to front
		c.doublyList.MoveToFront(node)
		// update the value of a list node
		node.Value.(*list.Element).Value = Pair{key: key, value: value}
	} else {
		// delete the last list node if the list is full
		if c.doublyList.Len() == c.capacity {
			// get the key that we want to delete; 用Comma-ok断言
			idx := c.doublyList.Back().Value.(*list.Element).Value.(Pair).key
			// delete the node pointer in the hash map by key
			delete(c.elementMap, idx)
			// remove the last list node
			c.doublyList.Remove(c.doublyList.Back())
		}
		// initialize a list node
		node := &list.Element{
			Value: Pair{
				key:   key,
				value: value,
			},
		}
		// push the new list node into the list
		ptr := c.doublyList.PushFront(node)
		// save the node pointer in the hash map
		c.elementMap[key] = ptr
	}
}

// invalid the entry of the cache, if not exist, do nothing
func (c *LRUCache) Invalid(key interface{}) {
	if node, ok := c.elementMap[key]; ok {
		// delete from map
		delete(c.elementMap, key)
		c.doublyList.Remove(node)
	}
}

/*func main() {

	obj := Constructor(2)   // nil
	obj.Put(1, 10)          // nil, linked list: [1:10]
	obj.Put(2, 20)          // nil, linked list: [2:20, 1:10]
	fmt.Println(obj.Get(1)) // 10, linked list: [1:10, 2:20]
	obj.Put(3, 30)          // nil, linked list: [3:30, 1:10]
	fmt.Println(obj.Get(2)) // -1, linked list: [3:30, 1:10]
	obj.Put(4, 40)          // nil, linked list: [4:40, 3:30]
	fmt.Println(obj.Get(1)) // -1, linked list: [4:40, 3:30]
	fmt.Println(obj.Get(3)) // 30, linked list: [3:30, 4:40]
	obj.Invalid(3)
	fmt.Println(obj.Get(3)) // -1, linked list: [ 4:40]
	fmt.Println(obj.Get(4)) // 40, linked list: [4:40]
}*/
