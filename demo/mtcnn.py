import os
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3' 
import cv2
import numpy as np
from align import AlignDlib
import requests
import base64
from flask_opencv_streamer.streamer import Streamer
from face_detect import MTCNN
from threading import Thread
import pafy


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
# vid = cv2.VideoCapture(0)
# vid = cv2.VideoCapture("rtsp://admin:AWPZEO@192.168.137.78/0/h264_stream")
frame = None
last_frame=None
embedder = cv2.dnn.readNetFromTorch("models/openface_nn4.small2.v1.t7")

embedder.setPreferableBackend(cv2.dnn.DNN_BACKEND_OPENCV)
embedder.setPreferableTarget(cv2.dnn.DNN_TARGET_OPENCL_FP16)

api = []
frames = []
stopped = False

def get_vector(aligned):        
    try:
        faceBlob = cv2.dnn.blobFromImage(aligned, 1.0 / 255,
            (96, 96), (0, 0, 0), swapRB=True, crop=False)
    except Exception as e:
        pass
    embedder.setInput(faceBlob)
    vec = embedder.forward()
    return vec

def send_api():
    while not stopped:
        # print("get : {}",len(api))
        if len(api) != 0:
            try:
                crop = api.pop(0)
                aligned = align_image(crop)
                cv2.imshow("aligned", aligned)
                k = cv2.waitKey(1) & 0xFF
                if k == 27:
                    pass
                try:
                    faceBlob = cv2.dnn.blobFromImage(aligned, 1.0 / 255, (96, 96), (0, 0, 0), swapRB=True, crop=False)
                except Exception as e:
                    print(e)
                    continue
                embedder.setInput(faceBlob)
                vec = embedder.forward()
                retval,buffer = cv2.imencode(".png",crop )
                string_bytes = base64.b64encode(buffer)
                data = {"image": string_bytes.decode('utf-8'), "camera_id": "5d522f6f16171e56f400246e", "embeddings": np.array(vec[0]).astype("float64").tolist()}
                failed = False
                while failed:
                    print("Api is on hold since connection could not be established")
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
                        print(e)
                        failed = True
            except Exception as e:
                print(e)
                continue

def get_frame():
    try:
        # url = 'https://www.youtube.com/watch?v=b5AurdhTcUc'
        # url = "https://www.youtube.com/watch?v=6qGiXY1SB68"

        # videoPafy = pafy.new(url)
        # best = videoPafy.getbest()
        # vid = cv2.VideoCapture(best.url)
        # vid = cv2.VideoCapture(0)
        # vid = cv2.VideoCapture("C:\\Users\\Athanatius.C\\Downloads\\Video\\LONDON WALK - Oxford Street to Carnaby Street - England.mp4")
        while not stopped:
            print("Video Recaptured")
            vid = cv2.VideoCapture(0)
            # vid = cv2.VideoCapture("rtsp://admin:AWPZEO@192.168.137.212/0/h264_stream")
            while not stopped:
                if stopped:
                    break
                ret,frm = vid.read()
                if not ret:
                    print("Frame get failed")
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

def draw_image(img,h0,h1,w0,w1,y):
    cv2.rectangle(img, (w0, h0), (w1, h1), (0, 255, 0), 2)
    landmark = landmarks[idx]
    for i in range(5):
        pt_h = landmark[i]
        pt_w = landmark[i + 5]
        if 0 <= pt_h and pt_h < h and 0 <= pt_w and pt_w < w:
            cv2.circle(img, (pt_w, pt_h), 1, (0, 255, 0),
                    thickness=3)# cv2.putText(img, "RPS: "+elapsed, (10, 50), cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 255, 0), 1)
    return img

def detect_face(frame,h0,h1,w0,w1):
    y = h0 - 10 if h0 - 10 > 10 else h0 + 10
    frame = draw_image(frame,h0,h1,w0,w1,y)


if __name__ == "__main__":
    port = 4040
    require_login = False
    streamer = Streamer(port, require_login)

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

    while not stopped:
        # _,frame = vid.read()
        try:
            # print(len(frames))
            if len(frames) == 0:
                frame = last_frame
            else:
                frame = frames.pop()
                last_frame = frame.copy()
                # del frames[0]
                (h, w) = frame.shape[:2]
                bounding_boxes, landmarks = mtcnn.detect(
                    img=frame, min_size=80, factor=0.709,
                    score_threshold=[0.8, 0.8, 0.8]
                )
                h, w, c = frame.shape
                for idx, bbox in enumerate(bounding_boxes):
                    h0 = max(int(round(bbox[0])), 0)
                    w0 = max(int(round(bbox[1])), 0)
                    h1 = min(int(round(bbox[2])), h - 1)
                    w1 = min(int(round(bbox[3])), w - 1)
                    threadlist = []
                    try:
                        score = bbox[4]
                        cropped = last_frame[h0 - 30:h1 + 30, w0 - 30:w1 + 30]
                        if len(api) < 10:
                            api.append(cropped)
                            print(len(api))
                            # if api.not_full:
                            #     api.put(cropped)
                            # else:
                            #     pass
                        y = h0 - 10 if h0 - 10 > 10 else h0 + 10
                        frame = draw_image(frame,h0,h1,w0,w1,y)
                    except Exception as e:
                        print(e)
                        pass
        except Exception as e:
            print(e)
            pass
        try:
            # cv2.imshow("window",frame)
            streamer.update_frame(frame)
            if not streamer.is_streaming:
                streamer.start_streaming()
        except Exception as e:
            print(e)
            pass
        k = cv2.waitKey(1) & 0xFF
        if k == 27:
            vid.release()
            # out.release()
            cv2.destroyAllWindows()
            stopped = True
            break
    capture.join()
    background_api.join()