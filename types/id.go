package types

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const Epoch int64 = 1685577600000 // 2023-06-01 00:00:00 +00:00
const (
	TimeShift = 12
	MaxStep   = 1 << TimeShift
	IDBit     = 54
)

var GlobalIDMake *SummonID

func init() {
	GlobalIDMake = NewSummonID()
}

type SummonID struct {
	mu    sync.Mutex
	time  int64
	step  int64
	Epoch time.Time
}

func NewSummonID() *SummonID {
	return &SummonID{
		Epoch: time.Unix(Epoch/1e3, (Epoch%1e3)*1e6),
		mu:    sync.Mutex{},
		step:  0,
	}
}

func (s *SummonID) Generate() ID {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Since(s.Epoch).Milliseconds()
	if now == s.time {
		s.step = (s.step + 1) &^ MaxStep
		if s.step == 0 {
			for now <= s.time {
				now = time.Since(s.Epoch).Milliseconds()
			}
		}
	} else {
		s.step = 0
	}
	s.time = now

	return ID((now << TimeShift) | s.step)
}

// ID is a 54-bit ID.
type ID int64

func (f ID) Int64() int64   { return int64(f) }
func (f ID) String() string { return strconv.FormatInt(int64(f), 10) }
func (f ID) Base2() string  { return fmt.Sprintf("%0*v", IDBit, strconv.FormatInt(int64(f), 2)) }
func (f ID) Bytes() []byte  { return []byte(f.String()) }
func (f ID) Base64() string { return base64.StdEncoding.EncodeToString(f.Bytes()) }

func (f ID) Time() time.Time {
	t := (int64(f) >> TimeShift) + Epoch
	return time.Unix(t/1e3, (t%1e3)*1e6)
}

// string to ID, if not return -1
func StringToID(s string) ID {
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return ID(i)
}

// for sql
func (f *ID) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		return nil
	case string:
		if src == "" {
			return nil
		}
		*f = StringToID(src)
		return nil
	case int:
		*f = ID(src)
		return nil
	case int64:
		*f = ID(src)
		return nil
	}

	return errors.New("failed to scan ID")
}

// for sql
func (f *ID) Value() (driver.Value, error) {
	return f.Int64(), nil
}

func (f *ID) UnmarshalJSON(data []byte) error {
	now, err := strconv.ParseInt(strings.Trim(string(data), `"`), 10, 64)
	*f = ID(now)
	return err
}

func (f ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}
