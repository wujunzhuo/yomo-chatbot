# YoMo Chatbot

这个示例展示了利用[YoMo](https://github.com/yomorun/yomo)传输数据，运行LLM大语言模型chatbot

## 整体架构

- Linux CPU服务器:
  [Zipper](https://yomo.run/docs/glossary#zipper-service)，负责数据传输
- Linux GPU服务器:
  [SFN](https://yomo.run/docs/glossary#streamfunction)，运行LLM-AI推理
- Local Macbook: 在命令行中进行提问，并收取LLM运行的答案

## 安装 YoMo

https://yomo.run/docs#install-cli

## 运行 YoMo Zipper

复制`config.yaml`到CPU服务器，然后运行

```sh
yomo serve -c config.yaml
```

## 编译 YoMo SFN 动态库

复制`chatbot`目录到GPU服务器，然后运行

```sh
go build -buildmode=c-shared -o sfn-lib.so sfn/lib.go
```

## 准备Python AI环境

确保GPU服务器上已安装CUDA环境：NVIDIA驱动、CUDA Toolkit、cuDNN SDK

复制`chat.py`和`requirements.txt`文件到GPU服务器，然后运行

```sh
pip install -r requirements.txt
```

下载大语言模型（如百川2）

```sh
git lfs install
git clone https://huggingface.co/baichuan-inc/Baichuan2-13B-Chat-4bits
```

## 运行Python AI推理程序

```sh
python chat.py \
    --sfn-lib chatbot/sfn-lib.so \
    --zipper ${YOUR_CPU_SERVER}:29000 \
    --model-path ${YOUR_MODEL_PATH}
```

## 在本机进行提问

```sh
YOMO_ZIPPER_ADDR=${YOUR_CPU_SERVER}:29000 go run cli/main.go
```
