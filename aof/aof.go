package aof

import (
	resphandler "go-redis/resp/handler"
	"os"
)

type stream struct {
	cmdLine [][]byte
	dbIdx   int
}

type Handler struct {
	database    resphandler.DatabaseInterface
	aofChan     chan *stream
	aofFile     *os.File
	aofFileName string
	currentDB   int
}

// New
// Add stream(set k v) -> aofChan
// HandleAof stream(set k v) <- aofChan  (to disk)
// LoadAof
