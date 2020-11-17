package protocol

// OperationType this is the type of the request
type OperationType uint8

const (
	// GET is the get request
	GET OperationType = 0
	// DELETE is the delete request
	DELETE OperationType = 1
	// PUT is the put request
	PUT OperationType = 2
	// EMPTY is a reserved request
	EMPTY OperationType = 3
)

// Request ...
type Request struct {
	ReqType OperationType
	Key     string
	Value   string
}

// Response ...
type Response struct {
	Value string
}
