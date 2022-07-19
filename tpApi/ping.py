import sys
import os
from flask import Flask, request, make_response
import gevent.pywsgi
import subprocess
import requests
import imgkit
import base64

app = Flask(__name__)
app.config['PORT'] = 9090
app.config['TIMEOUT'] = 5  # In seconds
app.config['VERSION'] = 1


print("pinging")
@app.route('/tcping', methods=['GET'])
def tcping():
    request_parameters = request.args
    ip = request_parameters.get('ip')
    port = request_parameters.get('port')
    if ip is None:
        return "no ip"
    if port is None:
        return "no port"
    tcping=' '.join(['tcping -n 2 -p',str(port),str(ip)])
    a=os.popen(tcping)
    a=a.read().split("\n")
    a=a[-3].split("=")[-1]
    if a!=" 0ms":
        return a
    else:
        return "error"

@app.route('/curl', methods=['GET'])
def curl():
    try:
        request_parameters = request.args
        url = request_parameters.get('url')
        if url is None:
            return "no url"
        if ("http://" not in url) and ("https://" not in url):
            url="http://"+url
        a=str(requests.get(url).status_code)
    except:
        a="error"
    return a

@app.route('/pic', methods=['GET', 'POST'])
def pic():
    request_parameters = request.args
    url = request_parameters.get('url')
    try:
        img = imgkit.from_url(url, False)
        return str(base64.b64encode(img),'utf-8')
    except:
        return ""
@app.errorhandler(404)
def page_not_found(e):
    return "<p>The resource could not be found.</p>", 404

if __name__ == '__main__':
    app_server = gevent.pywsgi.WSGIServer(('0.0.0.0', app.config['PORT']), app,log=None)
    app_server.serve_forever()
