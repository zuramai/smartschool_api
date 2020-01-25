import warnings
import pickle
import cython
import os
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3' 
import cv2
from sklearn.preprocessing import LabelEncoder
import numpy as np
from align import AlignDlib
import requests
import base64
# from jcopml.utils import save_model, load_model
from face_detect import MTCNN
from queue import Queue
from threading import Thread
import pafy
# from dotenv import load_dotenv

# load_dotenv(".env")

def distance(emb1, emb2):
    return np.sum(np.square(emb1 - emb2))

alignment = AlignDlib('models\\landmarks.dat')
def align_image(img):
    return alignment.align(96, img, alignment.getLargestFaceBoundingBox(img), 
                           landmarkIndices=AlignDlib.OUTER_EYES_AND_NOSE)

url = 'https://www.youtube.com/watch?v=6qGiXY1SB68'

videoPafy = pafy.new(url)
best = videoPafy.getbest()
vid = cv2.VideoCapture(best.url)
out = cv2.VideoWriter('output.avi', -1, 20.0, (640,480))
# vid = cv2.VideoCapture(0)
# vid = cv2.VideoCapture("rtsp://admin:AWPZEO@192.168.137.78/0/h264_stream")

embedder = cv2.dnn.readNetFromTorch("models/openface_nn4.small2.v1.t7")

embedder.setPreferableBackend(cv2.dnn.DNN_BACKEND_OPENCV)
embedder.setPreferableTarget(cv2.dnn.DNN_TARGET_OPENCL_FP16)

api = Queue(maxsize=10)
frames = []
stopped = False

def get_vector(aligned):
    try:
        faceBlob = cv2.dnn.blobFromImage(aligned, 1.0 / 255,
            (96, 96), (0, 0, 0), swapRB=True, crop=False)
    except Exception as e:
        print(e)
        pass
    embedder.setInput(faceBlob)
    vec = embedder.forward()
    return vec

def get_frame():
    try:
        # url = 'https://www.youtube.com/watch?v=b5AurdhTcUc'

        # videoPafy = pafy.new(url)
        # best = videoPafy.getbest()
        # vid = cv2.VideoCapture(best.url)
        # vid = cv2.VideoCapture("C:\\Users\\Athanatius.C\\Downloads\\Video\\LONDON WALK - Oxford Street to Carnaby Street - England.mp4")
        vid = cv2.VideoCapture(0)
        # vid = cv2.VideoCapture("rtsp://admin:AWPZEO@192.168.137.78/0/h264_stream")
        while not stopped:
            ret,frm = vid.read()
            if not ret:
                break
            # frame = cv2.resize(frame,(800,600))
            if len(frames) >=10:
                continue
            # frm = cv2.resize(frm, (1280,720), interpolation = cv2.INTER_LANCZOS4)
            frames.append(frm)
        vid.release()
    except Exception as e:
        print(e)
        pass

def send_api():
    while not stopped:
        if stopped:
            os._exit(0)
        if api.not_empty:
            try:

                aligned = None
                aligned = align_image(api.get())
                vec = None
                if np.shape(aligned) != ():
                    vec = get_vector(aligned)
                # else:
                #     print("reinitialized Vid Cap")
                #     cap = cv.VideoCapture(0)
                #     continue
                buffer = None
                _,buffer = cv2.imencode(".png", aligned)
                string_bytes = None
                string_bytes = base64.b64encode(buffer)
                data = None
                data = {"image": string_bytes.decode('utf-8'), "camera_id": "5d522f6f16171e56f400246e", "embeddings": np.array(vec[0]).astype("float64").tolist()}
                failed = True
                while failed:
                    if stopped:
                        os._exit(0)
                    try:
                        # r = requests.post(url="http://172.10.0.57:8088/api/v2/user/recognize", json=data)                   
                        r = requests.post(url="http://localhost:8088/api/v2/user/recognize", json=data)
                        res = None
                        res = r.json()
                        name,accuracy = None    
                        name = res["data"]["name"]
                        accuracy = res["data"]["accuracy"]
                        if accuracy < 0.2:
                            print("Detected : {} Accuracy : {}".format(name,accuracy))      
                        failed = False
                    except Exception as e:
                        print(e)
                        failed = True
            except Exception as e:
                print(e)
                continue

def draw_image(img,h0,h1,w0,w1,y):
    cv2.rectangle(img, (w0, h0), (w1, h1), (0, 255, 0), 2)
    # landmark = landmarks[idx]
    # for i in range(5):
    #     pt_h = landmark[i]
    #     pt_w = landmark[i + 5]
    #     if 0 <= pt_h and pt_h < h and 0 <= pt_w and pt_w < w:
    #         cv2.circle(img, (pt_w, pt_h), 1, (0, 255, 0),
    #                 thickness=3)# cv2.putText(img, "RPS: "+elapsed, (10, 50), cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 255, 0), 1)
    return img

def detect_face(frame,idx,bbox):
    h0 = max(int(round(bbox[0])), 0)
    w0 = max(int(round(bbox[1])), 0)
    h1 = min(int(round(bbox[2])), h - 1)
    w1 = min(int(round(bbox[3])), w - 1)
    threadlist = []
    try:
        score = bbox[4]
        cropped = frame[h0 - 30:h1 + 30, w0 - 30:w1 + 30]        
        api.put(cropped)
        y = h0 - 10 if h0 - 10 > 10 else h0 + 10
        frame = draw_image(frame,h0,h1,w0,w1,y)
    except Exception as e:
        print(e)
        pass

if __name__ == "__main__":
    background_api = Thread(target=send_api)
    background_api.start()

    capture = Thread(target=get_frame)
    capture.start()

    pnet_model_path = './models/pnet'
    rnet_model_path = './models/rnet'
    onet_model_path = './models/onet'

    mtcnn = MTCNN(pnet_model_path,
                  rnet_model_path,
                  onet_model_path)

    frame = None
    last_frame=None
    while not stopped:
        # _,frame = vid.read()
        if stopped:
            os._exit(0)
        if len(frames) == 0:
            frame = last_frame
        else:
            frame = frames[0]
            last_frame = frame.copy()
            del frames[0]
            try:
                (h, w) = frame.shape[:2]
            except Exception as e:
                print(e)
                continue
            bounding_boxes, landmarks = mtcnn.detect(
                img=frame, min_size=40, factor=0.709
            )
            h, w, c = frame.shape
            thread_list = []
            for idx, bbox in enumerate(bounding_boxes):
                thread = Thread(target=detect_face(frame,idx,bbox))
                thread_list.append(thread)
            for thread in thread_list:
                thread.start()
            for thread in thread_list:
                thread.join()
        try:
            out.write(frame)
            cv2.imshow("window",frame)
        except Exception as e:
            print(e)
            continue
        k = cv2.waitKey(1) & 0xFF
        if k == 27:
            stopped = True
            vid.release()
            out.release()
            cv2.destroyAllWindows()
            break
    capture.join()
    background_api.join()