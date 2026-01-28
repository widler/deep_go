package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
type OrderedMap struct { ... }

func NewOrderedMap() OrderedMap                      // создать упорядоченный словарь
func (m *OrderedMap) Insert(key, value int)          // добавить элемент в словарь
func (m *OrderedMap) Erase(key int)                  // удалить элемент из словари
func (m *OrderedMap) Contains(key int) bool          // проверить существование элемента в словаре
func (m *OrderedMap) Size() int                      // получить количество элементов в словаре
func (m *OrderedMap) ForEach(action func(int, int))  // применить функцию к каждому элементу словаря от меньшего к большему
*/
// go test -v homework_test.go

type node struct {
	key   int
	value int
	left  *node
	right *node
}

type OrderedMap struct {
	root *node
	size int
}

// NewOrderedMap создает новый упорядоченный словарь
func NewOrderedMap() OrderedMap {
	return OrderedMap{
		root: nil,
		size: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.root = &node{key: key, value: value}
		m.size++
		return
	}

	current := m.root
	var parent *node

	// Поиск места для вставки или обновления
	for current != nil {
		parent = current
		if key < current.key {
			current = current.left
		} else if key > current.key {
			current = current.right
		} else {
			// Ключ уже существует - обновляем значение
			current.value = value
			return
		}
	}

	// Вставка нового узла
	newNode := &node{key: key, value: value}
	if key < parent.key {
		parent.left = newNode
	} else {
		parent.right = newNode
	}
	m.size++
}
func (m *OrderedMap) findNode(key int) (*node, *node) {
	var parent *node
	current := m.root

	for current != nil {
		if key == current.key {
			return current, parent
		}

		parent = current
		if key < current.key {
			current = current.left
		} else {
			current = current.right
		}
	}

	return nil, nil
}

func minNode(n *node) (*node, *node) {
	if n == nil {
		return nil, nil
	}

	var parent *node
	current := n

	for current.left != nil {
		parent = current
		current = current.left
	}

	return current, parent
}

func (m *OrderedMap) Erase(key int) {
	nodeToDelete, parent := m.findNode(key)

	if nodeToDelete == nil {
		// Элемент не найден
		return
	}

	// Случай 1: У узла нет детей или только один ребенок
	if nodeToDelete.left == nil || nodeToDelete.right == nil {
		var child *node
		if nodeToDelete.left != nil {
			child = nodeToDelete.left
		} else {
			child = nodeToDelete.right
		}

		if parent == nil {
			// Удаляем корень
			m.root = child
		} else {
			if parent.left == nodeToDelete {
				parent.left = child
			} else {
				parent.right = child
			}
		}
		m.size--
		return
	}

	// Случай 2: У узла два ребенка
	// Находим преемника (минимальный элемент в правом поддереве)
	successor, successorParent := minNode(nodeToDelete.right)

	// Если преемник не является прямым правым ребенком
	if successorParent == nil {
		successorParent = nodeToDelete
	}

	// Копируем ключ и значение преемника
	nodeToDelete.key = successor.key
	nodeToDelete.value = successor.value

	// Удаляем преемника
	if successorParent.left == successor {
		successorParent.left = successor.right
	} else {
		successorParent.right = successor.right
	}

	m.size--
}

func (m *OrderedMap) Contains(key int) bool {
	node, _ := m.findNode(key)
	return node != nil
}

func (m *OrderedMap) Size() int {
	return m.size // need to implement
}

func inOrderTraversal(n *node, action func(int, int)) {
	if n == nil {
		return
	}

	inOrderTraversal(n.left, action)
	action(n.key, n.value)
	inOrderTraversal(n.right, action)
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	inOrderTraversal(m.root, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
