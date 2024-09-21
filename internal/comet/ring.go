package comet

import (
	"github.com/Terry-Mao/goim/api/protocol"
	"github.com/Terry-Mao/goim/internal/comet/conf"
	"github.com/Terry-Mao/goim/internal/comet/errors"
	log "github.com/golang/glog"
)

// Ring ring proto buffer.
type Ring struct {
	// read
	rp   uint64 // 读取位置
	num  uint64
	mask uint64
	// TODO split cacheline, many cpu cache line size is 64
	// pad [40]byte
	// write
	wp   uint64 // 写位置
	data []protocol.Proto
}

// NewRing new a ring buffer.
func NewRing(num int) *Ring {
	r := new(Ring)
	r.init(uint64(num))
	return r
}

// Init init ring.
func (r *Ring) Init(num int) {
	r.init(uint64(num))
}

func (r *Ring) init(num uint64) {
	// =========== 下面代码相当于：2^(N+1) =============
	// 2^N：确保 num 为2的幂，如果不是，则调整为下一个2的幂。如：2、4、8、16
	if num&(num-1) != 0 {
		for num&(num-1) != 0 {
			num &= num - 1
		}
		num <<= 1
	}
	// ===============================================
	r.data = make([]protocol.Proto, num)
	r.num = num
	r.mask = r.num - 1
}

// Get get a proto from ring.
func (r *Ring) Get() (proto *protocol.Proto, err error) {
	if r.rp == r.wp { // 没有数据的情况
		return nil, errors.ErrRingEmpty
	}
	proto = &r.data[r.rp&r.mask]
	return
}

// GetAdv incr read index.
func (r *Ring) GetAdv() {
	r.rp++
	if conf.Conf.Debug {
		log.Infof("ring rp: %d, idx: %d", r.rp, r.rp&r.mask)
	}
}

// Set get a proto to write.
func (r *Ring) Set() (proto *protocol.Proto, err error) {
	if r.wp-r.rp >= r.num {
		return nil, errors.ErrRingFull
	}
	proto = &r.data[r.wp&r.mask]
	return
}

// SetAdv incr write index.
func (r *Ring) SetAdv() {
	r.wp++
	if conf.Conf.Debug {
		log.Infof("ring wp: %d, idx: %d", r.wp, r.wp&r.mask)
	}
}

// Reset reset ring.
func (r *Ring) Reset() {
	r.rp = 0
	r.wp = 0
	// prevent pad compiler optimization
	// r.pad = [40]byte{}
}
