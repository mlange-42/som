package plot

type ConstantValues struct {
	Val    float64
	Length int
}

func (vs *ConstantValues) Len() int {
	return vs.Length
}

func (vs *ConstantValues) Value(i int) float64 {
	return vs.Val
}
