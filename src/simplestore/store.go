package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"protocol"
	"sync"
	"sync/atomic"
)

type requestAndConnection struct {
	con *gob.Encoder
	req protocol.Request
}

func handleRequest(queue chan requestAndConnection, ht *sync.Map, done *int32) {
	for atomic.LoadInt32(done) == 0 {
		r := <-queue
		var resp protocol.Response
		switch r.req.ReqType {
		case protocol.GET:
			val, b := ht.Load(r.req.Key)
			if b {
				v, ok := val.(string)
				if !ok {
					panic("Not able to make assertion about the value")
				}
				resp.Value = v
			} else {
				resp.Value = ""
			}
		case protocol.DELETE:
			ht.Delete(r.req.Key)
			resp.Value = ""
		case protocol.PUT:
			val, b := ht.LoadOrStore(r.req.Key, r.req.Value)
			if b {
				v, ok := val.(string)
				if !ok {
					panic("Not able to make assertion about the value")
				}
				resp.Value = v
			} else {
				v, ok := val.(string)
				if !ok {
					panic("Not able to make assertion about the value")
				}
				resp.Value = v
			}
		}
		if r.req.ReqType != protocol.EMPTY {
			checkerror(r.con.Encode(&resp))
		}
	}
}

func serveConnection(enc *gob.Encoder, dec *gob.Decoder, queue chan requestAndConnection, done *int32, threadpoolSize int) {
	//fmt.Println("Starting serving the connection")
	var err error = nil
	for err == nil {
		var req protocol.Request
		err = dec.Decode(&req)
		if err == nil {
			queue <- requestAndConnection{con: enc, req: req}
		}
	}
	atomic.StoreInt32(done, 1)
	for i := 0; i < threadpoolSize; i++ {
		queue <- requestAndConnection{req: protocol.Request{ReqType: protocol.EMPTY}}
	}

	//fmt.Println("Ended serving the connection")
}

func main() {
	address := "localhost:8080"

	l, err := net.Listen("tcp", address)
	checkerror(err)

	queue := make(chan requestAndConnection)

	ht := new(sync.Map)

	done := new(int32)

	atomic.StoreInt32(done, 0)

	threadpoolSize := 4
	for i := 0; i < threadpoolSize; i++ {
		go handleRequest(queue, ht, done)
	}

	fmt.Println("Accepting connections")
	con, err := l.Accept()
	checkerror(err)
	serveConnection(gob.NewEncoder(con), gob.NewDecoder(con), queue, done, threadpoolSize)
	checkerror(l.Close())

}

func checkerror(err error) {
	if err != nil {
		fmt.Println("Error ", err.Error())
		panic(err)
	}
}
