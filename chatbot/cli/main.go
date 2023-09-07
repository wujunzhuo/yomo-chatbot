package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/yomorun/yomo"
	"github.com/yomorun/yomo/serverless"

	chatbot "github.com/wujunzhuo/yomo-chatbot"
)

func Handler(ctx serverless.Context) {
	var res chatbot.Result
	err := json.Unmarshal(ctx.Data(), &res)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Answer:", res.Answer)
}

func main() {
	addr := os.Getenv("YOMO_ZIPPER_ADDR")
	if addr == "" {
		addr = "localhost:9000"
	}

	credential := os.Getenv("YOMO_ZIPPER_CREDENTIAL")
	if credential == "" {
		credential = "localhost:9000"
	}

	sfn := yomo.NewStreamFunction(chatbot.SFN_SINK_NAME, addr, yomo.WithSfnCredential(credential))
	sfn.SetObserveDataTags(chatbot.TAG_RES)
	sfn.SetHandler(Handler)
	err := sfn.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer sfn.Close()

	source := yomo.NewSource("source", addr, yomo.WithCredential(credential))
	err = source.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer source.Close()

	fmt.Println("please input your question:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		content := scanner.Text()
		if len(content) == 0 {
			fmt.Println("Bye")
			return
		}

		message := chatbot.Message{
			Role:    "user",
			Content: content,
		}

		buf, err := json.Marshal(&message)
		if err != nil {
			log.Fatalln(err)
		}

		err = source.Write(chatbot.TAG_REQ, buf)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
