package hasedbuffer

import (
	"encoding/hex"
	"fmt"
	"github.com/AppImageCrafters/zsync/rollinghash"
	"github.com/glycerine/rbuf"
	"golang.org/x/crypto/md4"
	"os"
)

type HashedRingBuffer struct {
	hash *rollinghash.RollingHash
	rBuf *rbuf.FixedSizeRingBuf
}

func NewHashedBuffer(size int) *HashedRingBuffer {
	/* Calculate bit-shift for blocksize */
	var blockShift uint16
	for i := uint16(0); i < 32; i++ {
		if size <= (1 << i) {
			blockShift = i
			break
		}
	}

	return &HashedRingBuffer{
		hash: rollinghash.NewRollingHash(blockShift),
		rBuf: rbuf.NewFixedSizeRingBuf(size),
	}
}

func (h *HashedRingBuffer) Write(p []byte) (n int, err error) {
	pSize := len(p)
	evictedSize := (h.rBuf.Readable + pSize) - h.rBuf.N
	if evictedSize < 0 {
		evictedSize = 0
	}

	for i := 0; i < pSize; i++ {
		if i < evictedSize {
			evicted := uint16(h.rBuf.A[h.rBuf.Use][h.rBuf.Beg])
			h.hash.Update(uint16(p[i]), evicted)

			h.rBuf.Advance(1)
		} else {
			h.hash.Update(uint16(p[i]), 0)
		}
	}

	return h.rBuf.Write(p)
}

func (h HashedRingBuffer) Bytes() []byte {
	return h.rBuf.Bytes()
}

func (h HashedRingBuffer) RollingSumHex() string {
	sum := h.RollingSum()

	return hex.EncodeToString(sum)
}

func (h HashedRingBuffer) RollingSum() []byte {
	sum := make([]byte, 4)
	h.hash.PutSum(sum)
	return sum
}

func (h HashedRingBuffer) CheckSum() []byte {
	sumBuilder := md4.New()
	sumBuilder.Write(h.rBuf.Bytes())
	sum := sumBuilder.Sum(nil)

	return sum
}

func (h HashedRingBuffer) CheckSumHex() string {
	sum := h.CheckSum()

	return hex.EncodeToString(sum)
}

func (h *HashedRingBuffer) readNFrom(input *os.File, bytes int64) (int64, error) {
	return -1, fmt.Errorf("not implemented yet")
}
