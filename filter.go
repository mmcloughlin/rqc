package rqc

import (
	"fmt"
	"math"
)

type Range struct {
	Min float64
	Max float64
}

func (r Range) String() string {
	return fmt.Sprintf("(%f,%f)", r.Min, r.Max)
}

func Gt(bound float64) Range {
	return Range{
		Min: bound,
		Max: math.Inf(1),
	}
}

func Lt(bound float64) Range {
	return Range{
		Min: math.Inf(-1),
		Max: bound,
	}
}
