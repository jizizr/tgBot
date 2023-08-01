#from paddleocr import PaddleOCR

import json
import wordcloud
from flask import Flask, request
import gevent.pywsgi
from multiprocessing import cpu_count, Process
import imageio
from matplotlib import colors
app = Flask(__name__)

app.config['TIMEOUT'] = 5  # In seconds

app.config['VERSION'] = 1
mk = imageio.imread("../source/mask.png")
def isLine(pos1,pos2,fontSize):
    return abs((pos1[0][1]+pos1[1][1])/2 - (pos2[0][1]+pos2[1][1])/2) < fontSize

def touchEdge(line,edge):
    return line[1][0] > edge

def pretty(result):
    fontSize = 0
    edge = 0
    points = []
    for i in range(len(result)):
        pos = result[i][0]
        fontSize += (pos[3][1]-pos[0][1]+pos[2][1]-pos[1][1])/2
        points.append([((pos[0][0]+pos[3][0])/2,(pos[0][1]+pos[3][1])/2),((pos[1][0]+pos[2][0])/2,(pos[1][1]+pos[2][1])/2)])
        if points[i][1][0]>edge:
            edge = points[i][1][0]
    fontSize = fontSize/len(result)
    edge = edge - 2 * fontSize
    line = [points[0]]
    lines = []
    for i in range(1,len(points)):
        if isLine(points[i-1],points[i],fontSize):
            line.append(points[i])
        else:
            lines.append(line)
            line = [points[i]]
    lines.append(line)
    # print(lines)
    i=0
    r=""
    for line in lines:
        if len(line) == 1:
            r = f"{r}{result[i][-1][0]}"
            if not touchEdge(line[0],edge):
                r = f"{r}\n\n"
            i+=1
        else:
            for j in range(len(line)-1):
                sep = "   "
                r = f"{r}{result[i][-1][0]}{sep}"
                i+=1
            r = f"{r}{result[i][-1][0]}\n"
            i+=1
    return r
#ocr = PaddleOCR(use_angle_cls=True,det_model_dir='/home/zr/.paddleocr/whl/det/ch/ch_PP-OCRv3_det_infer')
#def ocr_img(img):
    result = ocr.ocr(img, cls=True)
    result = result[0]
    return pretty(result)
# @app.route('/ocr', methods=['GET'])
# def ocr():
#    request_parameters = request.args
#    url = request_parameters.get('url')
#    return ocr_img(url)
color_list=['#10357B','#A6BEEC','#6C81B0','#092F70','#A7AAD3','#758BC4']
# 调用
colormap = colors.ListedColormap(color_list)
@app.route('/wc', methods=['GET'])
def wordCloud():
    request_parameters = request.args
    words = request_parameters.get('words')
    group = request_parameters.get('group')
    words = json.loads(words)
    wordcloud.WordCloud(width=400,
                            height=400,
                            background_color='white',
                            font_path='../source/font.ttf',
                            mask=mk,
                            colormap=colormap,
                            scale=5).generate_from_frequencies(words).to_file(f"../result/{group}.png") 
    return f"result/{group}.png"


@app.errorhandler(404)
def page_not_found(e):
    return "<p>The resource could not be found.</p>", 404

if __name__ == '__main__':
    app_server = gevent.pywsgi.WSGIServer(('127.0.0.1', 1222), app,log=None)
    app_server.start()
    def server_forever():
        app_server.start_accepting()
        app_server._stop_event.wait()
    for i in range(cpu_count()):
        p = Process(target=server_forever)
        p.start()
    #app_server.serve_forever()
