package main

import (
	"C"
	"bytes"
	"fmt"
	"sync"
	"unsafe"

	"github.com/yomorun/yomo"
	"github.com/yomorun/yomo/serverless"

	chatbot "github.com/wujunzhuo/yomo-chatbot"
)

var (
	mu      sync.Mutex
	reqChan chan []byte
	resChan chan []byte
)

func Handler(ctx serverless.Context) {
	fmt.Println("[sfn-go] Handler")

	request := ctx.Data()

	mu.Lock()
	defer mu.Unlock()
	reqChan <- request
	result := <-resChan

	ctx.Write(chatbot.TAG_RES, result)
}

//export Init
func Init(addrPtr *C.char, credentialPtr *C.char) int {
	mu = sync.Mutex{}
	reqChan = make(chan []byte)
	resChan = make(chan []byte)

	addr := C.GoString(addrPtr)
	credential := C.GoString(credentialPtr)

	sfn := yomo.NewStreamFunction(
		chatbot.SFN_CHATBOT_NAME, addr,
		yomo.WithSfnCredential(credential),
	)
	sfn.SetObserveDataTags(chatbot.TAG_REQ)
	sfn.SetHandler(Handler)
	sfn.SetErrorHandler(func(err error) {
		fmt.Println("[sfn-go] error:", err)
	})

	err := sfn.Connect()
	if err != nil {
		fmt.Println("[sfn-go] connect:", err)
		return 1
	}

	return 0
}

//export LoadRequest
func LoadRequest(reqPtr *byte, length int) int {
	reqBytes := <-reqChan
	if len(reqBytes) > length {
		return 0
	}
	buf := bytes.NewBuffer(unsafe.Slice(reqPtr, length)[:0])
	buf.Write(reqBytes)
	return len(reqBytes)
}

//export DumpResponse
func DumpResponse(resPtr *C.char) {
	resChan <- []byte(C.GoString(resPtr))
}

func main() {}
