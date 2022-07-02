package aof

import (
	"go-redis/config"
	idatabase "go-redis/core/database/interface"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"go-redis/utils/cmdconv"
	"go.uber.org/zap"
	"io"
	"os"
	"strconv"
)

const aofBufferSize = 1 << 16

type dataStream struct {
	cmdLine [][]byte
	dbIdx   int
}

type Handler struct {
	database    idatabase.Face
	aofChan     chan *dataStream
	aofFile     *os.File
	aofFileName string
	currentDB   int
}

func NewHandler(database idatabase.Face) (*Handler, error) {
	handler := &Handler{
		database:    database,
		aofFileName: config.Config.AppendFilename,
	}
	handler.LoadAof()
	aofFile, err := os.OpenFile(handler.aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile
	handler.aofChan = make(chan *dataStream, aofBufferSize)
	go func() {
		handler.HandleAof()
	}()
	return handler, nil
}

// AddAof Add dataStream(set k v) -> aofChan ---- send data to channel
func (h *Handler) AddAof(dbIdx int, cmd [][]byte) {
	if !config.Config.AppendOnly && h.aofChan == nil {
		return
	}
	h.aofChan <- &dataStream{
		cmdLine: cmd,
		dbIdx:   dbIdx,
	}
}

// HandleAof dataStream(set k v) <- aofChan ---- get data from channel and write to file(disk)
func (h *Handler) HandleAof() {
	h.currentDB = 0
	for stream := range h.aofChan {
		// dbIdx change
		if stream.dbIdx != h.currentDB {
			bytes := reply.NewMultiBulkReply(cmdconv.ToCmdLineStrings("select", strconv.Itoa(stream.dbIdx))).ToBytes()
			if _, err := h.aofFile.Write(bytes); err != nil {
				zap.S().Error(err)
				continue
			}
			h.currentDB = stream.dbIdx
		}
		// dbIdx unchanged
		bytes := reply.NewMultiBulkReply(stream.cmdLine).ToBytes()
		if _, err := h.aofFile.Write(bytes); err != nil {
			zap.S().Error(err)
		}
	}
}

// LoadAof load aof when start server - load data from disk to memory
func (h *Handler) LoadAof() {
	file, err := os.Open(h.aofFileName)
	if err != nil {
		zap.S().Errorf("cannot load aof: %v", zap.Error(err))
		return
	}
	defer func() {
		_ = file.Close()
	}()

	fakeClient := connection.NewEmpty()
	ch := parser.ParseStream(file)
	for stream := range ch {
		if stream.Err != nil {
			if stream.Err == io.EOF {
				break
			}
			zap.S().Error(zap.Error(err))
			continue
		}
		if stream.Data == nil {
			zap.S().Error("empty aof data stream")
			continue
		}

		multiBulkReply, ok := stream.Data.(*reply.MultiBulkReply)
		if !ok {
			zap.S().Error("need multi bulk")
			continue
		}

		execReply := h.database.Exec(fakeClient, multiBulkReply.Args)
		if reply.IsErrorReply(execReply) {
			zap.S().Errorf("cannot execute aof command: %v", execReply)
		}
	}
}
