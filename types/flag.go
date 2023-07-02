package types

import (
	"database/sql/driver"
	"encoding/base64"
	"strconv"
)

type Flag uint32

func (p Flag) Int64() int64   { return int64(p) }
func (p Flag) String() string { return strconv.FormatInt(int64(p), 10) }
func (p Flag) Binary() string { return strconv.FormatInt(int64(p), 2) }
func (p Flag) Bytes() []byte  { return []byte(p.String()) }
func (p Flag) Base64() string { return base64.StdEncoding.EncodeToString(p.Bytes()) }

// if self has the same or fewer permissions as other.
func (p Flag) Subset(other Flag) bool { return (p & other) == p }

// if self has the same or more permissions as other.
func (p Flag) Superset(other Flag) bool { return (p | other) == p }

// if the permissions on other are a strict subset of those on self.
func (p Flag) StrictSubset(other Flag) bool { return p.Subset(other) && p != other }

// if the permissions on other are a strict superset of those on self.
func (p Flag) StrictSuperset(other Flag) bool { return p.Superset(other) && p != other }

// string to Permission, if not return 0
func StringToFlag(p string) Flag {
	i, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return Flag(i)
}

// for sql
func (p *Flag) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		return nil
	case string:
		if src == "" {
			return nil
		}

		*p = StringToFlag(src)
	case int:
		*p = Flag(src)
	}

	return nil
}

// for sql
func (f *Flag) Value() (driver.Value, error) {
	return f.Int64(), nil
}
