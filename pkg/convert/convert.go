package convert

import "strconv"

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

//类型转换
func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) UInt32() (uint32, error) {
	v, error := strconv.Atoi(s.String())
	return uint32(v), error
}

func (s StrTo) MustUint32() uint32 {
	v, _ := s.UInt32()
	return v
}
