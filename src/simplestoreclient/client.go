package simplestoreclient

import (
	"encoding/gob"
	"fmt"
	"net"
	"protocol"
)

// Clientctx ...
type Clientctx struct {
	connection net.Conn
	enc        *gob.Encoder
	dec        *gob.Decoder
}

// NewClientctx Returns a new Clientctx
func NewClientctx(addressWithPort string) *Clientctx {
	conn, err := net.Dial("tcp", addressWithPort)
	checkerror(err)
	ctx := new(Clientctx)
	ctx.connection = conn
	ctx.enc = gob.NewEncoder(conn)
	ctx.dec = gob.NewDecoder(conn)
	return ctx
}

// EndClientCtx ...
func (c *Clientctx) EndClientCtx() {
	checkerror(c.connection.Close())
}

// Get ...
func (c *Clientctx) Get(k string) string {
	var req protocol.Request
	var res protocol.Response

	req.ReqType = protocol.GET
	req.Key = k
	req.Value = ""

	checkerror(c.enc.Encode(&req))
	checkerror(c.dec.Decode(&res))
	return res.Value
}

// Delete ...
func (c *Clientctx) Delete(k string) string {
	var req protocol.Request
	var res protocol.Response

	req.ReqType = protocol.DELETE
	req.Key = k
	req.Value = ""

	checkerror(c.enc.Encode(&req))
	checkerror(c.dec.Decode(&res))
	return res.Value
}

// Put ...
func (c *Clientctx) Put(k string, v string) string {
	var req protocol.Request
	var res protocol.Response

	req.ReqType = protocol.PUT
	req.Key = k
	req.Value = v

	checkerror(c.enc.Encode(&req))
	checkerror(c.dec.Decode(&res))
	return res.Value
}

func checkerror(err error) {
	if err != nil {
		fmt.Println("Error ", err.Error())
		panic(err)
	}
}
