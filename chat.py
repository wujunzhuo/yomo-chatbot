import argparse
import ctypes
import json
import sys
import torch
from transformers import AutoModelForCausalLM, AutoTokenizer
from transformers.generation.utils import GenerationConfig

BUF_SIZE = 2 * 1024 * 1024


def init_model(model_path: str):
    tokenizer = AutoTokenizer.from_pretrained(
        model_path, use_fast=False, trust_remote_code=True)
    model = AutoModelForCausalLM.from_pretrained(
        model_path, device_map='auto', torch_dtype=torch.bfloat16,
        trust_remote_code=True)
    model.generation_config = GenerationConfig.from_pretrained(model_path)
    return tokenizer, model


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--model_path', type=str,
                        default='baichuan-inc/Baichuan2-13B-Chat-4bits')
    parser.add_argument('--sfn_lib', type=str, default='sfn-lib.so')
    parser.add_argument('--zipper', type=str, default='localhost:9000')
    parser.add_argument('--credential', type=str, default='')
    args = parser.parse_args()

    library = ctypes.cdll.LoadLibrary(args.sfn_lib)
    f_init = library.Init
    f_init.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
    f_init.restype = ctypes.c_int

    f_load_request = library.LoadRequest
    f_load_request.argtypes = [ctypes.c_char_p, ctypes.c_int]
    f_load_request.restype = ctypes.c_int

    f_dump_response = library.DumpResponse
    f_dump_response.argtypes = [ctypes.c_char_p]

    res = f_init(args.zipper.encode(), args.credential.encode())
    if res > 0:
        sys.exit(res)

    toknizer, model = init_model(args.model_path)
    print('[ai-python] init finished')

    buf = ctypes.create_string_buffer(BUF_SIZE)

    while True:
        try:
            n = f_load_request(buf, BUF_SIZE)
            if n <= 0:
                raise ValueError('request buffer oom')

            request = buf[:n]
            message = json.loads(request)
            print('[ai-python] request:', message)

            result = model.chat(toknizer, [message])
            print('[ai-python] result:', result)

            f_dump_response(json.dumps({'answer': result}).encode())
        except Exception as e:
            print(e)
            f_dump_response(json.dumps({'error': str(e)}).encode())
        finally:
            print('[ai-python] done', flush=True)


if __name__ == '__main__':
    main()
