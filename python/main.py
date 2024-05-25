import argparse
import time
import cv2
import numpy as np

import grpc
import sudoku_pb2
import sudoku_pb2_grpc
from SudokuParser import parse_sudoku

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--port', type=int, default=8080, help='Port number')
    parser.add_argument('--path', type=str, required=True, help='Path to the image')
    parser.add_argument('--solver', type=str, default='backtrack', choices=['backtrack', 'dlx'], help='Solver type')
    args = parser.parse_args()

    t = time.time()
    board = parse_sudoku(cv2.imread(args.path))
    dt = time.time()-t

    channel = grpc.insecure_channel(f'localhost:{args.port}')
    stub = sudoku_pb2_grpc.SudokuStub(channel)
    response = stub.GetSolution(sudoku_pb2.Question(board=board, solverType=args.solver))
    print("result: ", response.board, response.solved)
    print(f"parsing time: {dt}s")
    print(f"solving time: {response.duration}s")
