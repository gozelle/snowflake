package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	DefaultTsBits       = 44            // 558 years
	DefaultNodeBits     = 7             // 128 servers
	DefaultSequenceBits = 12            // 4096 ids/millisecond
	DefaultEpoch        = 1136214245000 // 2006-01-02 15:04:05
)

var (
	SnowZero, _ = NewSnowflake(0)
)

type ID struct {
	id int64
}

func (i ID) Int64() int64 {
	return i.id
}

func (i ID) String() string {
	return strconv.FormatInt(i.id, 10)
}

type SnowOptions struct {
	Node         uint16
	TsBits       uint8
	NodeBits     uint8
	SequenceBits uint8
	Epoch        int64
}

type Snowflake struct {
	mu          sync.Mutex
	options     SnowOptions
	tick        int64
	sequence    uint16
	maxTs       int64
	maxNode     uint16
	maxSequence uint16
}

func (g *Snowflake) NewID() (*ID, error) {

	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixNano() / 1e6

	if g.tick == now {
		g.sequence++
		if g.sequence > g.maxSequence {
			for now <= g.tick {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		g.sequence = 0
	}

	ts := now - g.options.Epoch

	if ts > g.maxTs {
		return nil, errors.New("time use up")
	}

	g.tick = now

	return &ID{
		id: ts<<(g.options.NodeBits+g.options.SequenceBits) |
			int64(g.options.Node)<<g.options.SequenceBits |
			int64(g.sequence),
	}, nil
}

func (g *Snowflake) EndAt() string {
	endAt := (g.options.Epoch + int64(1)<<g.options.TsBits) / 1e3
	return time.Unix(endAt, 0).Format("2006-01-02 15:04:05")
}

func NewSnowflakeWithOptions(ops SnowOptions) (*Snowflake, error) {

	maxNode := uint16(1<<ops.NodeBits - 1)

	if ops.Node > maxNode {
		return nil, errors.New("node id greater")
	}

	return &Snowflake{
		mu:          sync.Mutex{},
		options:     ops,
		tick:        0,
		sequence:    0,
		maxTs:       1<<ops.TsBits - 1,
		maxNode:     maxNode,
		maxSequence: 1<<ops.SequenceBits - 1,
	}, nil
}

func NewSnowflake(node uint16) (*Snowflake, error) {
	ops := SnowOptions{
		Node:         node,
		TsBits:       DefaultTsBits,
		NodeBits:     DefaultNodeBits,
		SequenceBits: DefaultSequenceBits,
		Epoch:        DefaultEpoch,
	}

	return NewSnowflakeWithOptions(ops)
}
