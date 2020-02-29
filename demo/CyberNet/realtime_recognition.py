import warnings
import pickle
import os
# os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3' 
import cv2
from sklearn.preprocessing import LabelEncoder
import numpy as np
from align import AlignDlib
import requests
import base64
import pafy
# from jcopml.utils import save_model, load_model
from queue import Queue
from threading import Thread


def distance(emb1, emb2):
    return np.sum(np.square(emb1 - emb2))

alignment = AlignDlib('models\\landmarks.dat')
def align_image(img):
    return alignment.align(96, img, alignment.getLargestFaceBoundingBox(img), 
                           landmarkIndices=AlignDlib.OUTER_EYES_AND_NOSE)

vid = cv2.VideoCapture(0)
# vid = cv2.VideoCapture("rtsp://admin:AWPZEO@192.168.137.78/0/h264_stream")

# model = load_model ("C:\\Users\\Athanatius.C\\Tensorflow\\Face Recognition\\Model\\SRCV-core.vcore")
protoPath = os.path.join("models", "deploy.prototxt")
modelPath = os.path.join("models",
    "res10_300x300_ssd_iter_140000.caffemodel")
detector = cv2.dnn.readNetFromCaffe(protoPath, modelPath)
embedder = cv2.dnn.readNetFromTorch("models/openface_nn4.small2.v1.t7")

detector.setPreferableBackend(cv2.dnn.DNN_BACKEND_OPENCV)
detector.setPreferableTarget(cv2.dnn.DNN_TARGET_OPENCL_FP16)

embedder.setPreferableBackend(cv2.dnn.DNN_BACKEND_OPENCV)
embedder.setPreferableTarget(cv2.dnn.DNN_TARGET_OPENCL_FP16)
 
# frames = Queue(maxsize=1000) 
frames = []
api = Queue(maxsize=1000)


stopped = False

def get_frame():
    try:
        # gst = "gst-launch-1.0 rtspsrc location=rtsp://admin:AWPZEO@192.168.137.78/0/h264_stream ! rtph264depay ! h264parse ! omxh264dec ! nveglglessink window-x=100 window-y=100 window-width=640 window-height=360"
        # url = 'https://www.youtube.com/watch?v=6qGiXY1SB68'

        # videoPafy = pafy.new(url)
        # best = videoPafy.getbest()
        # vid = cv2.VideoCapture(best.url)
        # vid = cv2.VideoCapture("rtsp://admin:AWPZEO@192.168.137.212/0/h264_stream")
        vid = cv2.VideoCapture(0)
        # vid = cv2.VideoCapture("C:\\Users\\Athanatius.C\\Downloads\\Video\\LONDON WALK - Oxford Street to Carnaby Street - England.mp4")
        while not stopped:
            ret,frame = vid.read()
            if not ret:
                break
            # if frames.full():
            #     continue
            # frame = cv2.resize(frame,(800,600))
            # frames.put(frame)
            if len(frames) >=1:
                continue
            # frame = cv2.resize(frame, (1280,720), interpolation = cv2.INTER_LANCZOS4)

            frames.append(frame)
        vid.release()
    except Exception as e:
        print(e)
        pass

def send_api():
    while not stopped:
        if api.not_empty:
            try:
                crop = api.get()
                aligned = align_image(crop)
                # cv2.imshow("aligned", aligned)
                # k = cv2.waitKey(1) & 0xFF
                # if k == 27:
                #     continue
                try:
                    faceBlob = cv2.dnn.blobFromImage(aligned, 1.0 / 255, (96, 96), (0, 0, 0), swapRB=True, crop=False)
                except Exception as e:
                    # print(e)
                    continue
                embedder.setInput(faceBlob)
                vec = embedder.forward()
                retval,buffer = cv2.imencode(".png", crop)
                string_bytes = base64.b64encode(buffer)
                data = {"image": string_bytes.decode('utf-8'), "camera_id": "5d522f6f16171e56f400246e", "embeddings": np.array(vec[0]).astype("float64").tolist()}
                failed = True
                while failed:
                    if stopped:
                        break
                    try:
                        if stopped:
                            break
                        # r = requests.post(url="http://172.10.0.57:8088/api/v2/user/recognize", json=data)                   
                        r = requests.post(url="http://localhost:8088/api/v2/user/recognize", json=data)
                        res = r.json()
                        name = res["data"]["name"]
                        accuracy = res["data"]["accuracy"]
                        if accuracy < 0.2:
                            print("Detected : {} Accuracy : {}".format(name,accuracy))      
                        failed = False
                    except Exception as e:
                        # print(e)
                        failed = True
            except Exception as e:
                # print(e)
                continue

if __name__ == "__main__":
    background_api = Thread(target=send_api)
    background_api.start()
    get_frame = Thread(target=get_frame)
    get_frame.start()
    last_frame = None
    frame = None
    while not stopped:
        # print(len(frames))
        if len(frames) == 0:
            frame = last_frame
        else:
            frame = frames[0]
            last_frame = frame.copy()
            del frames[0]
            try:
                (h, w) = frame.shape[:2]
            except:
                continue
            imageBlob = cv2.dnn.blobFromImage(frame, 1.0, (400, 400), swapRB=False, crop=False)
            detector.setInput(imageBlob)
            detections = detector.forward()
            for i in range(0, detections.shape[2]):
                confidence = detections[0, 0, i, 2]
                if confidence > 0.7:
                    box = detections[0, 0, i, 3:7] * np.array([w, h, w, h])
                    (startX, startY, endX, endY) = box.astype("int")
                    y = startY - 10 if startY - 10 > 10 else startY + 10
                    face = frame[startY-10:endY+10, startX-10:endX+10]
                    try:
                        api.put(face)
                        cv2.rectangle(frame, (startX, startY), (endX, endY), (0, 255, 0), 4)
                    except Exception as e:
                        pass
        try:
            cv2.imshow("window", frame)
        except:
            continue
        k = cv2.waitKey(5) & 0xFF
        if k == 27:
            vid.release()
            cv2.destroyAllWindows()
            break
    stopped = True
    get_frame.join()
    background_api.join()