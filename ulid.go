package ulid

import (
	"database/sql/driver"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

const (
	BINARY16 ValueType = iota
	VARCHAR26
)

type (
	ValueType uint8

	UID struct {
		id  ulid.ULID
		bin bool
	}

	UIDGen struct {
		ofs uint64
		erp *entropyReaderPool
		vt  ValueType
	}

	entropyReaderPool struct {
		p sync.Pool
	}
)

func (u UID) Bytes() []byte {
	return u.id[:]
}

func (u UID) String() string {
	return u.id.String()
}

func (u *UID) Scan(src interface{}) error {
	return u.id.Scan(src)
}

func (u *UID) Value() (driver.Value, error) {
	if u.bin {
		return u.id.Value()
	}
	return u.id.String(), nil
}

func (g *UIDGen) UID() UID {
	ms := ulid.Timestamp(time.Now()) - g.ofs
	r := g.erp.get()

	lid, err := ulid.New(ms, r)
	g.erp.put(r)
	if err != nil {
		panic(fmt.Sprintf("uidgen: %s", err))
	}

	return UID{id: lid, bin: g.vt == BINARY16}
}

func (g *UIDGen) Parse(s string) (UID, bool) {
	lid, err := ulid.Parse(s)
	if err != nil {
		return UID{}, false
	}
	return UID{id: lid, bin: g.vt == BINARY16}, true
}

func (p *entropyReaderPool) get() io.Reader {
	return p.p.Get().(io.Reader)
}

func (p *entropyReaderPool) put(r io.Reader) {
	p.p.Put(r)
}

func New(offset uint64, vt ValueType) *UIDGen {
	g := UIDGen{
		ofs: offset,
		erp: newEntropyReaderPool(),
		vt:  vt,
	}

	return &g
}

func newEntropyReaderPool() *entropyReaderPool {
	return &entropyReaderPool{
		p: sync.Pool{
			New: func() interface{} {
				t := time.Now()
				rnr := rand.New(rand.NewSource(t.UnixNano()))
				return ulid.Monotonic(rnr, 0)
			},
		},
	}
}
