package sane

// A Range is a set of discrete integer or fixed-point values. Value x is in
// the range if there is an integer k >= 0 such that Min <= k*Quant <= Max.
// The type of Min, Max and Quant is either int or float64 for all three.
type Range struct {
	Min   interface{} // minimum value
	Max   interface{} // maximum value
	Quant interface{} // quantization step
}
