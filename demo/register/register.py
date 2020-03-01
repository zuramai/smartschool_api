"""Demo code shows how to estimate human head pose.
Currently, human face is detected by a detector from an OpenCV DNN module.
Then the face box is modified a little to suits the need of landmark
detection. The facial landmark detection is done by a custom Convolutional
Neural Network trained with TensorFlow. After that, head pose is estimated
by solving a PnP problem.
"""
from argparse import ArgumentParser
# from multiprocessing import Process, Queue
from queue import Queue
from threading import Thread

import cv2
import numpy as np
import PySimpleGUI as sg
from mark_detector import MarkDetector
from os_detector import detect_os
from pose_estimator import PoseEstimator
from stabilizer import Stabilizer
import os
from align import AlignDlib
import requests
import json

alignment = AlignDlib('models\\landmarks.dat')
def align_image(img):
    return alignment.align(96, img, alignment.getLargestFaceBoundingBox(img), 
                           landmarkIndices=AlignDlib.OUTER_EYES_AND_NOSE)

print("OpenCV version: {}".format(cv2.__version__))

# multiprocessing may not work on Windows and macOS, check OS for safety.
# detect_os()


class Register:
    def __init__(self):
        self.imagelist = []
        self.CNN_INPUT_SIZE = 128
        self.central = 0
        self.right = 0
        self.left = 0
        self.down = 0
        self.up = 0
        self.min_sample = 4
        self.capture = False
        self.path = ""
        self.user_id = ""
        self.protoPath = os.path.join("models", "deploy.prototxt")
        self.modelPath = os.path.join("models",
            "res10_300x300_ssd_iter_140000.caffemodel")
        self.embeddings = []

        self.protoPath = os.path.join("models", "deploy.prototxt")
        self.modelPath = os.path.join("models",
            "res10_300x300_ssd_iter_140000.caffemodel")

        self.detector = cv2.dnn.readNetFromCaffe(self.protoPath, self.modelPath)
        self.embedder = cv2.dnn.readNetFromTorch("models/openface_nn4.small2.v1.t7")

        self.embedder.setPreferableBackend(cv2.dnn.DNN_BACKEND_CUDA)
        self.embedder.setPreferableTarget(cv2.dnn.DNN_TARGET_CUDA)

        self.detector.setPreferableBackend(cv2.dnn.DNN_BACKEND_CUDA)
        self.detector.setPreferableTarget(cv2.dnn.DNN_TARGET_CUDA)

        self.alignment = AlignDlib('models\\landmarks.dat')
        self.failed = False

        sg.change_look_and_feel('Material2')
        layout = [
            [sg.Image(filename='D:\\Athanatius\\Pictures\\faceid.png', key='-IMAGE-', tooltip='Enter NIK or NIS and hit register to start')],
            [sg.Text("Please Enter your NIK or NIS below")],
            [sg.Text('NIS or NIK : '), sg.InputText()],
            [sg.Button('Register'), sg.Button('Cancel')]
        ]
        self.window = sg.Window('Register - CyberNet Facial Recognition', layout, location=(250,250),)  # if trying Qt, you will need to remove this right click menu
        self.window.Finalize()

    def align_image(self,img):
        return alignment.align(96, img, alignment.getLargestFaceBoundingBox(img), 
                           landmarkIndices=AlignDlib.OUTER_EYES_AND_NOSE)
        
    def Verify(self,directory):
        images = os.listdir(directory)
        for path in images:
            frame = cv2.imread(os.path.join(directory,path))
            try:
                (h, w) = frame.shape[:2]
            except:
                continue
            imageBlob = cv2.dnn.blobFromImage(frame, 1.0, (300, 300), swapRB=False, crop=False)
            self.detector.setInput(imageBlob)
            detections = self.detector.forward()
            for i in range(0, detections.shape[2]):
                confidence = detections[0, 0, i, 2]
                if confidence > 0.7:
                    box = detections[0, 0, i, 3:7] * np.array([w, h, w, h])
                    (startX, startY, endX, endY) = box.astype("int")
                    y = startY - 10 if startY - 10 > 10 else startY + 10
                    try:
                        # face = frame[startY-10:endY+10, startX-10:endX+10]
                        (h,w) =frame.shape[:2]
                        aligned = self.align_image(frame)
                        faceBlob = cv2.dnn.blobFromImage(aligned, 1.0 / 255,
                        (96, 96), (0, 0, 0), swapRB=True, crop=False)
                        self.embedder.setInput(faceBlob)
                        vec = self.embedder.forward()
                        print(np.array(vec[0]).astype("float64").tolist())
                        self.embeddings.append(np.array(vec[0]).astype("float64").tolist())
                        cv2.rectangle(frame, (startX, startY), (endX, endY), (0, 255, 0), 4)
                        # self.window['-IMAGE-'].update(data=cv2.imencode('.png', frame)[1].tobytes())
                    except Exception as e:
                        print(e)
                        pass
            k = cv2.waitKey(5) & 0xFF
            if k == 27 or len(self.embeddings) == 20:
                cv2.destroyAllWindows()
                break
        try:
            print(self.embeddings)
            r = requests.post(url="http://192.168.1.12:8088/api/v2/user/verify",json= {"user_id":int(self.user_id),"embeddings":np.array(self.embeddings).astype("float64").tolist()})
            print(r.json())
        except:
            sg.Popup("Register Failed! Api not Online!")
            self.failed = True
            try:
                os.remove(self.path)
            except:
                print("Cannot remove folder! please remove folder manually! :{}".format(self.path))
            self.window['-IMAGE-'].Update(data=self.noimage)

    def get_face(self,detector, img_queue, box_queue):
        """Get face from image queue. This function is used for multiprocessing"""
        while True:
            image = img_queue.get()
            box = detector.extract_cnn_facebox(image)
            box_queue.put(box)

    def get_interface(self,x, y, z):
        if x > 0.45:
            horizontal = "Right"
        elif x < -0.45:
            horizontal = "Left"
        else:
            horizontal= "Center"
        if y > 0.35:
            vertical = "Down"
        elif y < -0.1:
            vertical = "Up"
        else:
            vertical = "Center"
        return vertical, horizontal
        
    def auto_sort_save(self,vertical,horizontal, face):
        if vertical == "Center" and horizontal == "Center" and self.central < self.min_sample:
            self.central +=1
            cv2.imwrite(os.path.join(self.path,"sample_central_{}.jpg".format(self.central)),face)
            self.imagelist.append(face)
        elif vertical == "Down" and horizontal == "Center" and self.down < self.min_sample:
            self.down +=1
            cv2.imwrite(os.path.join(self.path,"sample_down_{}.jpg".format(self.down)), face)
            self.imagelist.append(face)
        elif vertical == "Up" and horizontal == "Center" and self.up < self.min_sample:
            self.up +=1
            cv2.imwrite(os.path.join(self.path, "sample_up_{}.jpg".format(self.up)), face)
            self.imagelist.append(face)
        elif vertical == "Center" and horizontal == "Left" and self.left < self.min_sample:
            self.left +=1
            cv2.imwrite(os.path.join(self.path, "sample_left_{}.jpg".format(self.left)), face)
            self.imagelist.append(face)
        elif vertical == "Center" and horizontal == "Right" and self.right < self.min_sample:
            self.right +=1
            cv2.imwrite(os.path.join(self.path,"sample_right_{}.jpg".format(self.right)),face)
            self.imagelist.append(face)
        # else:
        #     exit(0)
        

    def main(self):
        """MAIN"""
        # Video source from webcam or video file.
        video_src = 0

        cap = cv2.VideoCapture(video_src)
        if video_src == 0:
            cap.set(cv2.CAP_PROP_FRAME_WIDTH, 640)
        _, sample_frame = cap.read()

        # Introduce mark_detector to detect landmarks.
        mark_detector = MarkDetector()

        # Setup process and queues for multiprocessing.
        img_queue = Queue()
        box_queue = Queue()
        img_queue.put(sample_frame)
        box_process = Thread(target=self.get_face, args=(
            mark_detector, img_queue, box_queue,))
        box_process.start()
        # Introduce pose estimator to solve pose. Get one frame to setup the
        # estimator according to the image size.
        height, width = sample_frame.shape[:2]
        pose_estimator = PoseEstimator(img_size=(height, width))
        # Introduce scalar stabilizers for pose.
        pose_stabilizers = [Stabilizer(
            state_num=2,
            measure_num=1,
            cov_process=0.1,
            cov_measure=0.1) for _ in range(6)]

        tm = cv2.TickMeter()
        self.noimage =cv2.imread("D:\\Athanatius\\Pictures\\faceid.png")
        self.noimage =cv2.imencode('.png', self.noimage)[1].tobytes()

        while True:
            # print(img_queue.qsize())
            # print(self.capture)
            event, values = self.window.read(timeout=20)
            if event in (None, "Cancel"):
                os._exit(0)
            if event in (None, "Register"):
                self.path = os.path.join("sample",values[0])
                self.user_id = values[0]
                try:
                    os.mkdir(self.path)
                    self.capture = True
                except:
                    sg.Popup("FaceID of this person already existed!")
                    self.window['-IMAGE-'].Update(data=self.noimage)
                    # pass
            if self.capture:
                
                # Read frame, crop it, flip it, suits your needs.
                frame_got, frame = cap.read()
                if frame_got is False:
                    break
                original = frame
                if video_src == 0:
                    frame = cv2.flip(frame, 2)
                img_queue.put(frame)
                facebox = box_queue.get()
                if facebox is not None:
                    # Detect landmarks from image of 128x128.
                    face_img = frame[facebox[1]: facebox[3],
                                    facebox[0]: facebox[2]]
                    face_img = cv2.resize(face_img, (self.CNN_INPUT_SIZE, self.CNN_INPUT_SIZE))
                    face_img = cv2.cvtColor(face_img, cv2.COLOR_BGR2RGB)

                    tm.start()
                    marks = mark_detector.detect_marks([face_img])
                    tm.stop()

                    # Convert the marks locations from local CNN to global image.
                    marks *= (facebox[2] - facebox[0])
                    marks[:, 0] += facebox[0]
                    marks[:, 1] += facebox[1]

                    # Uncomment following line to show raw marks.
                    mark_detector.draw_marks(
                        frame, marks, color=(0, 255, 0))

                    # Uncomment following line to show facebox.
                    # mark_detector.draw_box(frame, [facebox])

                    # Try pose estimation with 68 points.
                    pose = pose_estimator.solve_pose_by_68_points(marks)

                    # Stabilize the pose.
                    steady_pose = []
                    pose_np = np.array(pose).flatten()
                    for value, ps_stb in zip(pose_np, pose_stabilizers):
                        ps_stb.update([value])
                        steady_pose.append(ps_stb.state[0])
                    steady_pose = np.reshape(steady_pose, (-1, 3))
                    pose_estimator.draw_annotation_box(
                        frame, steady_pose[0], steady_pose[1], color=(128, 255, 128))

                    x = steady_pose[0][0]
                    y = steady_pose[0][1]
                    z = steady_pose[0][2]

                    vertical,horizontal = self.get_interface(x,y,z)
                    face_img = cv2.cvtColor(face_img,cv2.COLOR_RGB2BGR)
                    self.auto_sort_save(vertical,horizontal,face_img)
                    # print(x)
                    cv2.putText(frame,"horizontal : {} vertical : {} ".format(horizontal,vertical),(100,100),cv2.FONT_HERSHEY_TRIPLEX,0.6,(0,255,0),1)
                    # Uncomment following line to draw head axes on frame.
                    # pose_estimator.draw_axes(frame, stabile_pose[0], stabile_pose[1])

                # Show preview.
                self.window['-IMAGE-'].update(data=cv2.imencode('.png', frame)[1].tobytes())
                # cv2.imshow("Preview", frame)
                if self.central == self.min_sample and self.up == self.min_sample and self.left == self.min_sample and self.right == self.min_sample and self.down == self.min_sample:
                    self.capture = False
                    self.up = 0
                    self.down = 0
                    self.left = 0
                    self.right = 0
                    self.Verify(self.path)
                    if not self.failed:
                        self.window['-IMAGE-'].Update(data=self.noimage)
                        sg.Popup("FaceID successfully Registered!\nYou may add more FaceID!")

        # Clean up the multiprocessing process.
        # box_process.terminate()
        box_process.join()
        self.window.Close()


if __name__ == '__main__':
    register = Register()
    register.main()
