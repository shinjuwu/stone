package network

type Processor interface {
	// must goroutine safe
	Route(msg interface{}, userData interface{}) error
	// must goroutine safe
	Unmarshal(data []byte) (interface{}, error)
	// must goroutine safe
	Marshal(msg interface{}) ([][]byte, error)
	Package(data interface{}) []byte
	Unpackage(data []byte) ([]byte, error)
	BytePackage(data []byte) []byte
}
