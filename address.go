package aprs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Address struct {
	Call     string
	SSID     int
	Repeated bool
}

func (a Address) EqualTo(b *Address) bool {
	return b != nil && a.Call == b.Call && a.SSID == b.SSID
}

func (a Address) String() string {
	var r = ""

	if a.Repeated {
		r = "*"
	}
	if a.SSID == 0 {
		return a.Call + r
	}
	return fmt.Sprintf("%s-%d%s", a.Call, a.SSID, r)
}

func (a Address) Secret() int16 {
	var h = int16(0x73e2)
	var c = a.Call

	if len(c)%2 > 0 {
		c += "\x00"
	}
	for i := 0; i < len(c); i += 2 {
		h ^= int16(c[i]) << 8
		h ^= int16(c[i+1])
	}
	return h & 0x7fff
}

func ParseAddress(s string) (*Address, error) {
	r := strings.HasSuffix(s, "*")
	if r {
		s = s[:len(s)-1]
	}
	p := strings.Split(strings.ToUpper(s), "-")
	if len(p) == 0 || len(p) > 2 {
		return nil, errors.New("aprs: invalid address")
	}

	a := &Address{Call: p[0], Repeated: r}
	if len(p) == 2 {
		var i int64
		i, err := strconv.ParseInt(p[1], 10, 32)
		if err != nil || i > 16 {
			return nil, errors.New("aprs: invalid address")
		}
		a.SSID = int(i)
	}

	return a, nil
}

func MustParseAddress(s string) *Address {
	a, err := ParseAddress(s)
	if err != nil {
		panic(err)
	}
	return a
}

type Path []*Address

func (p Path) String() string {
	var s = make([]string, len(p))
	for i, a := range p {
		s[i] = a.String()
	}
	return strings.Join(s, ",")
}

func ParsePath(p string) (Path, error) {
	ss := strings.Split(p, ",")

	if len(ss) == 0 {
		return nil, nil
	}

	var err error
	as := make(Path, len(ss))
	for i, s := range ss {
		as[i], err = ParseAddress(s)
		if err != nil {
			return nil, err
		}
	}

	return as, nil
}
