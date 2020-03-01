import requests,time,pafy,cv2,os
import numpy as np
from align import AlignDlib
from base64 import b64encode
from flask_opencv_streamer.streamer import Streamer
from face_detect import MTCNN
from cv2 import rectangle
from threading import Thread
from queue import Queue
from multiprocessing import Process, Lock, Pool
from aiohttp import ClientSession
import asyncio
import argparse
from itertools import repeat

os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3' 

class CyberNet:
    def __init__(self):
        parser = argparse.ArgumentParser(description='Input the Application Config.')
        parser.add_argument('--camera_id',required=True,help='sum the integers (default: find the max)')
        parser.add_argument('--rtsp_link',required=True,help='an integer for the accumulator')
        parser.add_argument('--stream_port',required=True, type=int,help='an integer for the accumulator')
        parser.add_argument('--end_point',required=True,help='an integer for the accumulator')
        args = vars(parser.parse_args())
        
        self.camera_id =args["camera_id"]
        self.rtsp_link = args["rtsp_link"]
        self.port = args["stream_port"]
        self.endpoint = args["end_point"]
        self.api = Queue(maxsize=10)
        self.frames = Queue(maxsize=10)
        self.frame = None
        self.last_frame = None
        self.stopped = False
        self.require_login = False
        self.use_gpu = True
        # self.helper = []
        self.embedder = cv2.dnn.readNetFromTorch("models/openface_nn4.small2.v1.t7")
        self.alignment = AlignDlib('models\\landmarks.dat')
        self.streamer = Streamer(self.port, self.require_login)

        # if self.use_gpu:
        self.embedder.setPreferableBackend(cv2.dnn.DNN_BACKEND_CUDA)
        self.embedder.setPreferableTarget(cv2.dnn.DNN_TARGET_CUDA)
        # self.main()


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

    def postreq(self, data):
        # requests.post(url="http://172.10.0.31:8088/api/v2/user/recognize", json=data)
        requests.post(url="http://{}/api/v2/user/recognize".format(self.endpoint), json=data)
        
    def send(self,crop):
        try:
            # lock.acquire
            aligned = self.align_image(crop)
            try:
                faceBlob = cv2.dnn.blobFromImage(aligned, 1.0 / 255, (96, 96), (0, 0, 0), swapRB=True, crop=False)
            except:
                return
            self.embedder.setInput(faceBlob)
            vec = self.embedder.forward()
            retval,buffer = cv2.imencode(".png",crop )
            string_bytes = b64encode(buffer)
            data = {"image": string_bytes.decode('utf-8'), "camera_id": self.camera_id, "embeddings": np.array(vec[0]).astype("float64").tolist()}
            # self.postreq(data)
            t = Thread(target=self.postreq(data))
            # t.setDaemon(True)
            t.daemon = True
            t.start()
        except Exception as e:
            print(e)
            return
        # lock.release

    def get_frame(self):
        try:
            '''url = "https://www.youtube.com/watch?v=6qGiXY1SB68"
            videoPafy = pafy.new(url)
            best = videoPafy.getbest()
            vid = cv2.VideoCapture(best.url)'''
            while not self.stopped:
                print("Video Recaptured")
                try:
                    vid = cv2.VideoCapture(int(self.rtsp_link))
                except:
                    vid = cv2.VideoCapture(self.rtsp_link)
                # vid = cv2.VideoCapture("rtsp://admin:AWPZEO@192.168.137.166/0/h264_stream")
                while not self.stopped:
                    if self.stopped:
                        break
                    ret,frm = vid.read()
                    if not ret:
                        print("Frame get failed")
                        break
                    if self.frames.not_full:
                        self.frames.put(frm)
                vid.release()
        except Exception as e:
            pass
        
    def helper(self,idx, bbox,w,h,landmarks):
        h0 = max(int(round(bbox[0])), 0)
        w0 = max(int(round(bbox[1])), 0)
        h1 = min(int(round(bbox[2])), h - 1)
        w1 = min(int(round(bbox[3])), w - 1)
        try:
            score = bbox[4]
            y = h0 - 10 if h0 - 10 > 10 else h0 + 10
            rectangle(self.frame, (w0, h0), (w1, h1), (255, 255, 255), 1)
            landmark = landmarks[idx]
            for i in range(5):
                pt_h = landmark[i]
                pt_w = landmark[i + 5]
                if 0 <= pt_h and pt_h < h and 0 <= pt_w and pt_w < w:
                    cv2.circle(self.frame, (pt_w, pt_h), 1, (255, 255, 255),
                            thickness=1)
        except Exception as e:
            print(e)
            pass

    def main(self):
        capture = Thread(target=self.get_frame)
        capture.start()

        pnet_model_path = './models/pnet'
        rnet_model_path = './models/rnet'
        onet_model_path = './models/onet'
        mtcnn = MTCNN(pnet_model_path,
                        rnet_model_path,
                        onet_model_path)

        while not self.stopped:
            if self.frames.not_empty:
                self.frame = self.frames.get()
                (h, w) = self.frame.shape[:2]
                bounding_boxes, landmarks = mtcnn.detect(
                    img=self.frame, min_size=80, factor=0.709,
                    score_threshold=[0.8, 0.8, 0.8]
                )
                # send_queue = []
                # lock = Lock()
                for idx, bbox in enumerate(bounding_boxes):
                    self.helper(idx, bbox,w,h,landmarks)
                    h0 = max(int(round(bbox[0])), 0)
                    w0 = max(int(round(bbox[1])), 0)
                    h1 = min(int(round(bbox[2])), h - 1)
                    w1 = min(int(round(bbox[3])), w - 1)
                    cropped = self.frame[h0 - 30:h1 + 30, w0 - 30:w1 + 30]
                    self.send(cropped)
                #     p = Process(target=self.send(cropped,lock,))
                #     send_queue.append(p)
                # for process in send_queue:
                #     process.start()
                # for process in send_queue:
                #     process.join()
                # for process in send_queue:
                #     process.terminate()
                # send_queue.clear()
            try:
                # cv2.imshow("window",self.frame)
                self.streamer.update_frame(self.frame)
                if not self.streamer.is_streaming:
                    self.streamer.start_streaming()
            except Exception as e:
                # print(e)
                pass
            k = cv2.waitKey(1) & 0xFF
            if k == 27:
                cv2.destroyAllWindows()
                self.stopped = True
                break
        capture.join()
        # background_api.join()

if __name__ == "__main__":
    cybernet = CyberNet
    cybernet().main()