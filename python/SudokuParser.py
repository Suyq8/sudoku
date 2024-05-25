import cv2
import tesserocr
from imutils import perspective
import numpy as np
from PIL import Image

def split_cells(image: np.ndarray):
    length = image.shape[0]
    box_length = length//9
    rows = np.vsplit(image, 9)
    boxes = []
    percent = []
    size = (box_length-10)**2

    for row in rows:
        boxs = np.hsplit(row, 9)
        for box in boxs:
            box = box[6:box_length-4, 6:box_length-4]
            boxes.append(Image.fromarray(box))
            percent.append(np.sum(box==255)/size)

    return boxes, percent

def parse_sudoku(img: np.ndarray):
    img = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    binary_image = cv2.adaptiveThreshold(img, 255, cv2.ADAPTIVE_THRESH_GAUSSIAN_C, cv2.THRESH_BINARY_INV, 11, 2)

    contours, _ = cv2.findContours(binary_image, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    length = cv2.arcLength(contours[0], True)
    approx = cv2.approxPolyDP(contours[0], 0.002 * length, True)

    wrapped = perspective.four_point_transform(binary_image, approx.squeeze())
    length = 450
    wrapped = cv2.resize(wrapped, (length, length))
    wrapped = cv2.adaptiveThreshold(wrapped, 255, cv2.ADAPTIVE_THRESH_GAUSSIAN_C, cv2.THRESH_BINARY, 91, 3)

    cells, percent = split_cells(wrapped)

    sudoku = []

    for sudoku_cell, p in zip(cells, percent):
        if p>0.03:
            number_text = tesserocr.image_to_text(sudoku_cell, psm=10)
            try:
                sudoku.append(int(number_text))
            except:
                sudoku.append(0)
        else:
            sudoku.append(0)

    return sudoku