import os
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3' 
import cv2
from cv2.dnn import blobFromImage
import numpy as np
from align import AlignDlib
import requests
from base64 import b64encode
from flask_opencv_streamer.streamer import Streamer
from face_detect import MTCNN
from cv2 import rectangle
from threading import Thread
import time
import pafy

class CyberNet:
    def __init__(self):
        self.api = []
        self.frames = []
        self.port = 4040
        self.frame = None
        self.last_frame = None
        self.stopped = False
        self.require_login = False
        self.embedder = cv2.dnn.readNetFromTorch("models/openface_nn4.small2.v1.t7")
        self.embedder.setPreferableBackend(cv2.dnn.DNN_BACKEND_OPENCV)
        self.embedder.setPreferableTarget(cv2.dnn.DNN_TARGET_OPENCL_FP16)
        self.alignment = AlignDlib('models\\landmarks.dat')
        self.streamer = Streamer(self.port, self.require_login)
        self.main()


    def align_image(self,img):
        return self.alignment.align(96, img, self.alignment.getLargestFaceBoundingBox(img), 
                            landmarkIndices=AlignDlib.OUTER_EYES_AND_NOSE)

    def draw_image(self,img,h0,h1,w0,w1,y,h,landmarks,idx):
        rectangle(img, (w0, h0), (w1, h1), (255, 255, 255), 1)
        landmark = landmarks[idx]
        for i in range(5):
            pt_h = landmark[i]
            pt_w = landmark[i + 5]
            if 0 <= pt_h and pt_h < h and 0 <= pt_w and pt_w < w:
                cv2.circle(img, (pt_w, pt_h), 1, (255, 255, 255),
                        thickness=2)
        return img

    def background_sender(self):
        while not self.stopped:
            if len(self.api) != 0:
                try:
                    crop = self.api.pop(0)
                    aligned = self.align_image(crop)
                    # cv2.imshow("aligned", aligned)
                    # k = cv2.waitKey(1) & 0xFF
                    # if k == 27:
                    #     pass
                    faceBlob = blobFromImage(aligned, 1.0 / 255, (96, 96), (0, 0, 0), swapRB=True, crop=False)
                    self.embedder.setInput(faceBlob)
                    vec = self.embedder.forward()
                    retval,buffer = cv2.imencode(".png",crop )
                    string_bytes = b64encode(buffer)
                    data = {"image": string_bytes.decode('utf-8'), "camera_id": "5d522f6f16171e56f400246e", "embeddings": np.array(vec[0]).astype("float64").tolist()}
                    failed = False
                    while failed:
                        print("Api is on hold since connection could not be established")
                        if self.stopped:
                            break
                        try:
                            if self.stopped:
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

    def get_frame(self):
        try:
            '''url = "https://www.youtube.com/watch?v=6qGiXY1SB68"
            videoPafy = pafy.new(url)
            best = videoPafy.getbest()
            vid = cv2.VideoCapture(best.url)'''
            while not self.stopped:
                print("Video Recaptured")
                vid = cv2.VideoCapture(0)
                # vid = cv2.VideoCapture("rtsp://admin:AWPZEO@192.168.137.212/0/h264_stream")
                while not self.stopped:
                    if self.stopped:
                        break
                    ret,frm = vid.read()
                    if not ret:
                        print("Frame get failed")
                        break
                    if len(self.frames) >=10:
                        continue
                    self.frames.append(frm)
                vid.release()
        except Exception as e:
            # print(e)
            pass

    def main(self):
        background_api = Thread(target=self.background_sender)
        background_api.start()

        capture = Thread(target=self.get_frame)
        capture.start()

        pnet_model_path = './models/pnet'
        rnet_model_path = './models/rnet'
        onet_model_path = './models/onet'
        mtcnn = MTCNN(pnet_model_path,
                        rnet_model_path,
                        onet_model_path)

        while not self.stopped:
            try:
                if len(self.frames) == 0:
                    self.frame = self.last_frame
                else:
                    self.frame = self.frames.pop()
                    self.last_frame = self.frame.copy()
                    # del frames[0]
                    (h, w) = self.frame.shape[:2]
                    start = time.time()
                    bounding_boxes, landmarks = mtcnn.detect(
                        img=self.frame, min_size=80, factor=0.709,
                        score_threshold=[0.8, 0.8, 0.8]
                    )
                    print("Mtcnn took {}".format(time.time()-start))
                    for idx, bbox in enumerate(bounding_boxes):
                        h0 = max(int(round(bbox[0])), 0)
                        w0 = max(int(round(bbox[1])), 0)
                        h1 = min(int(round(bbox[2])), h - 1)
                        w1 = min(int(round(bbox[3])), w - 1)
                        # threadlist = []
                        try:
                            score = bbox[4]
                            cropped = self.last_frame[h0 - 30:h1 + 30, w0 - 30:w1 + 30]
                            # if len(self.api) < 10:
                            #     self.api.append(cropped)
                            y = h0 - 10 if h0 - 10 > 10 else h0 + 10
                            # self.frame = self.draw_image(self.frame,h0,h1,w0,w1,y,h,w,landmarks,idx)
                            rectangle(self.frame, (w0, h0), (w1, h1), (255, 255, 255), 1)
                            # landmark = landmarks[idx]
                            # for i in range(5):
                            #     pt_h = landmark[i]
                            #     pt_w = landmark[i + 5]
                            #     if 0 <= pt_h and pt_h < h and 0 <= pt_w and pt_w < w:
                            #         cv2.circle(self.frame, (pt_w, pt_h), 1, (255, 255, 255),
                            #                 thickness=2)
                        except Exception as e:
                            print(e)
                            pass
            except Exception as e:
                print(e)
                pass
            try:
                cv2.imshow("window",self.frame)
                self.streamer.update_frame(self.frame)
                if not self.streamer.is_streaming:
                    self.streamer.start_streaming()
            except Exception as e:
                print(e)
                pass
            k = cv2.waitKey(1) & 0xFF
            if k == 27:
                cv2.destroyAllWindows()
                self.stopped = True
                break
        capture.join()
        background_api.join()

if __name__ == "__main__":
    cybernet = CyberNet
    cybernet()