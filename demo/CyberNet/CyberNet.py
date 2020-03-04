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
        parser.add_argument('--mode',default="delayed",help='an integer for the accumulator')
        args = vars(parser.parse_args())
        
        self.camera_id =args["camera_id"]
        self.rtsp_link = args["rtsp_link"]
        self.port = args["stream_port"]
        self.endpoint = args["end_point"]
        if args["mode"] == "delayed":
            self.delayed = True
        else:
            self.delayed = False 
        if self.delayed:
            self.frames = Queue(maxsize=10)
        else:
            self.frames = []
        self.api = Queue(maxsize=10)
        self.post_queue = Queue(maxsize=10)
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

        self.sender_thread = Thread(target=self.sender)
        self.sender_thread.start()

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

    def sender(self):
        while not self.stopped:
            if self.post_queue.not_empty:
                try:
                    requests.post(url="http://{}/api/v2/user/recognize".format(self.endpoint), json=self.post_queue.get())
                except:
                    print("Time Out")
        # requests.post(url="http://172.10.0.31:8088/api/v2/user/recognize", json=data)
        
    def prepare_data(self,crop):
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
            if self.post_queue.not_full:
                self.post_queue.put(data)
        except Exception as e:
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
                    if self.delayed:
                        if self.frames.not_full:
                            self.frames.put(frm)
                    else:
                        if len(self.frames) < 10:
                            self.frames.append(frm)
                vid.release()
        except Exception as e:
            pass
        
    def drawer(self,idx, bbox,w,h,landmarks):
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
            if self.delayed:
                if not self.frames.not_empty:
                    continue
                self.frame = self.frames.get()
            else:
                if len(self.frames) ==0 :
                    continue
                self.frame = self.frames.pop()
            (h, w) = self.frame.shape[:2]
            bounding_boxes, landmarks = mtcnn.detect(
                img=self.frame, min_size=80, factor=0.709,
                score_threshold=[0.8, 0.8, 0.8]
            )
            send_queue = []
            # lock = Lock()
            for idx, bbox in enumerate(bounding_boxes):
                h0 = max(int(round(bbox[0])), 0)
                w0 = max(int(round(bbox[1])), 0)
                h1 = min(int(round(bbox[2])), h - 1)
                w1 = min(int(round(bbox[3])), w - 1)
                # start = time.time()
                cropped = self.frame[h0 - 30:h1 + 30, w0 - 30:w1 + 30]
                # print("took {}".format(time.time()-start))
                self.prepare_data(cropped)
                # self.drawer(idx, bbox,w,h,landmarks)
                t = Thread(target=self.drawer,args=(idx, bbox,w,h,landmarks,))
                send_queue.append(t)

            for thread in send_queue:
                thread.start()
            for thread in send_queue:
                thread.join()
                # send_queue.append(t)
            #     t = Thread(target=self.prepare_data,args=(cropped,))
            #     send_queue.append(t)
            #     t.start()
            # for thread in send_queue:
            #     thread.join()
            # with Pool(8) as pool:
            #     pool.map_async(self.prepare_data,send_queue)
            try:
                cv2.imshow("window",self.frame)
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
        self.sender_thread.join()
        # background_api.join()

if __name__ == "__main__":
    cybernet = CyberNet
    cybernet().main()