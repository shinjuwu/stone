package random

const (
	// Bytes to Number
	_N1 = 0x100
	_N2 = _N1 << 8
	_N3 = _N1 << 16
	_N4 = _N1 << 24

	// Hexadecimals to Decimal
	_D16p0 = 0x1
	_D16p1 = _D16p0 << 4
	_D16p2 = _D16p0 << 8
	_D16p3 = _D16p0 << 12
	_D16p4 = _D16p0 << 16
	_D16p5 = _D16p0 << 20
	_D16p6 = _D16p0 << 24
	_D16p7 = _D16p0 << 28
)

var (
	_Numbers  = []float64{_N1, _N2, _N3, _N4}
	_Decimals = []int64{_D16p7, _D16p6, _D16p5, _D16p4, _D16p3, _D16p2, _D16p1, _D16p0}
)

type Fairness struct {
	ServerSeed string `json:"server_seed"`
	ClientSeed string `json:"client_seed"`
	Nonce      int64  `json:"nonce"`
	Cursor     int64  `json:"cursor"`
}
