import cv2
import time
from face_detect import MTCNN
from align import AlignDlib
import requests
import numpy as np
import threading

# import face_alignment
import queue

stopped = False
frames = queue.Queue(maxsize=10)
result = queue.Queue(maxsize=0)

embedder = cv2.dnn.readNetFromTorch("models/openface_nn4.small2.v1.t7")
alignment = AlignDlib('models\\landmarks.dat')

def align_image(img):
    return alignment.align(96, img, alignment.getLargestFaceBoundingBox(img), 
                           landmarkIndices=AlignDlib.OUTER_EYES_AND_NOSE)

def background():
	while datas.not_empty:
		tracked.put(send.track(datas.get()))
		tracked.task_done()
		time.sleep(3)

def get_frame():
    try:
        vid = cv2.VideoCapture(0)
        while not stopped:
            ret,frame = vid.read()
            if not ret:
                break
            if frames.full():
                continue
            # frame = cv2.resize(frame,(800,600))
            frames.put(frame)
        vid.release()
    except Exception as e:
        pass

def recognize_face(vec):
    try:
        r = requests.post(url="http://localhost:8088/api/v2/user/recognize",json= {"embeddings":np.array(vec[0]).astype("float64").tolist()})
        res = r.json()
        name = res["data"]["name"]
        proba = res["data"]["accuracy"]
        elapsed = res["data"]["elapsed"]
    except Exception as e:
        print(e)
    return proba, name, elapsed
    
def get_vector(aligned):
    try:
        faceBlob = cv2.dnn.blobFromImage(aligned, 1.0 / 255,
            (96, 96), (0, 0, 0), swapRB=True, crop=False)
    except Exception as e:
        pass
    embedder.setInput(faceBlob)
    vec = embedder.forward()
    return vec

def draw_image(img,h0,h1,w0,w1,y,proba,name,elapsed):
    if int(100 - int(proba * 100)) > 50:
        cv2.rectangle(img, (w0, h0), (w1, h1), (0, 255, 0), 2)
        landmark = landmarks[idx]
        for i in range(5):
            pt_h = landmark[i]
            pt_w = landmark[i + 5]
            if 0 <= pt_h and pt_h < h and 0 <= pt_w and pt_w < w:
                cv2.circle(img, (pt_w, pt_h), 1, (0, 255, 0),
                        thickness=3)
        text = name + " " + str(int(100 - int(proba * 100))) + "%"
        cv2.putText(img, text, (w0, y-20), cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 255, 0), 1)
        cv2.putText(img, elapsed, (w0, y), cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 255, 0), 1)
    else:
        cv2.rectangle(img, (w0, h0), (w1, h1), (0, 0, 255), 4)
        text = "Unknown "+str(int(100-int(proba*100)))+"%"
        cv2.putText(img, text, (w0, y-20),cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 0, 255), 1)
        cv2.putText(img, elapsed, (w0, y), cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 0, 255), 1)
        # cv2.putText(img, "RPS: "+elapsed, (10, 50), cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 255, 0), 1)
    return img


def detect(w0,w1,h0,h1,idx, bbox,img):
    score = bbox[4]
    cropped = img[h0 - 10:h1 + 10, w0 -10:w1 + 10]
    aligned = align_image(cropped)
    # If out of bound
    try:
        vec = get_vector(aligned)
        proba, name, elapsed = recognize_face(vec)
        # proba = 0.6
        # name = "test"
        # elapsed = "10ms"
        y = h0 - 10 if h0 - 10 > 10 else h0 + 10
        img = draw_image(img,h0,h1,w0,w1,y,proba,name,elapsed)
    except Exception as e:
        pass
    result.put(cv2.imshow("window", img))

if __name__ == '__main__':

    t2 = threading.Thread(target=get_frame)
    t2.daemon = True
    t2.start()

    pnet_model_path = './models/pnet'
    rnet_model_path = './models/rnet'
    onet_model_path = './models/onet'

    mtcnn = MTCNN(pnet_model_path,
                  rnet_model_path,
                  onet_model_path)

    # vid = cv2.VideoCapture(0)
    while not stopped:
        # _, img = vid.read()
        if not frames.not_empty:
            continue
        img = frames.get()
        bounding_boxes, landmarks = mtcnn.detect(
            img=img, min_size=40, factor=0.709,
            score_threshold=[0.8, 0.8, 0.8]
        )
        h, w, c = img.shape
        if len(bounding_boxes) == 0:
            result.put(cv2.imshow("window", img))
        else:
            threads =[]
            for idx, bbox in enumerate(bounding_boxes):
                h0 = max(int(round(bbox[0])), 0)
                w0 = max(int(round(bbox[1])), 0)
                h1 = min(int(round(bbox[2])), h - 1)
                w1 = min(int(round(bbox[3])), w - 1)
                det = threading.Thread(target=detect(w0,w1,h0,h1,idx,bbox,img))
                det.daemon = True
                det.start()
        if result.not_empty:
            result.get()
        k = cv2.waitKey(5) & 0xFF
        if k == 27:
            stopped = True
            # vid.release()
            cv2.destroyAllWindows()
            break