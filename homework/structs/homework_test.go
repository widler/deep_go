package main

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func ReplaceBits(pos int, count int, value bool) uint32 {
	var clearMask uint32 = ((1 << count) - 1) << pos
	result := math.MaxUint32 & ^clearMask
	if value {
		var setMask uint32 = ((1 << count) - 1) << pos
		result |= setMask
	}
	return result
}

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		name = name[:42]
		dataNamePtr := (*[42]byte)(unsafe.Pointer(&person.name))
		for i := 0; i < len(name); i++ {
			dataNamePtr[i] = name[i]
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.goldhome &= 1
		person.goldhome |= uint32(gold << 1)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := ReplaceBits(22, 10, false)
		person.flags &= mask
		person.flags |= uint32(mana << 22)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := ReplaceBits(12, 10, false)
		person.flags &= mask
		person.flags |= uint32(health << 12)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := ReplaceBits(8, 4, false)
		person.flags &= mask
		person.flags |= uint32(respect << 8)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := ReplaceBits(4, 4, false)
		person.flags &= mask
		person.flags |= uint32(strength << 4)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := ReplaceBits(4, 0, false)
		person.flags &= mask
		person.flags |= uint32(experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		var mask uint8 = ((1 << 4) - 1) << 4
		clearMask := math.MaxUint8 & ^mask
		person.flags2 &= clearMask
		person.flags2 |= uint8(level << 4)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		mask := math.MaxUint32 - 1
		person.goldhome &= uint32(mask)
		person.goldhome |= uint32(1)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint8(1 << 3)
		person.flags2 |= mask
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint8(1 << 2)
		person.flags2 |= mask
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		var mask uint8 = 3
		clearMask := math.MaxUint8 & ^mask
		person.flags2 &= clearMask
		person.flags2 |= uint8(personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	name     [42]byte
	flags2   uint8
	x, y, z  int32
	goldhome uint32
	flags    uint32
	// need to implement
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, option := range options {
		option(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	if p == nil {
		return ""
	}
	return unsafe.String(&p.name[0], 42)
}

func (p *GamePerson) X() int {
	// need to implement
	return int(p.x)
}

func (p *GamePerson) Y() int {
	// need to implement
	return int(p.y)
}

func (p *GamePerson) Z() int {
	// need to implement
	return int(p.z)
}

func (p *GamePerson) Gold() int {

	return int(p.goldhome >> 1)
}

func (p *GamePerson) Mana() int {
	// need to implement
	return int(p.flags >> 22)
}

func (p *GamePerson) Health() int {
	mask := uint32((1<<10 - 1) << 12)
	return int((p.flags & mask) >> 12)
}

func (p *GamePerson) Respect() int {
	mask := uint32((1<<4 - 1) << 8)
	return int((p.flags & mask) >> 8)
}

func (p *GamePerson) Strength() int {
	mask := uint32((1<<4 - 1) << 4)
	return int((p.flags & mask) >> 4)
}

func (p *GamePerson) Experience() int {
	mask := uint32((1<<4 - 1))
	return int((p.flags & mask))
}

func (p *GamePerson) Level() int {
	return int(p.flags2 >> 4)
}

func (p *GamePerson) HasHouse() bool {
	mask := uint32(1)
	return (p.goldhome & mask) == 1
}

func (p *GamePerson) HasGun() bool {
	mask := uint8(1 << 3)
	return (p.flags2&mask)>>3 == 1
}

func (p *GamePerson) HasFamilty() bool {
	mask := uint8(1 << 2)
	return (p.flags2&mask)>>2 == 1
}

func (p *GamePerson) Type() int {
	var mask uint8 = 3
	return int(p.flags2 & mask)
}

func TestGamePerson(t *testing.T) {
	fmt.Println(unsafe.Sizeof(GamePerson{}))
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = WarriorGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
