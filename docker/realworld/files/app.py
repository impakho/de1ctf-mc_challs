#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# @Author: impakho
# @Date:   2020/04/15
# @Github: https://github.com/impakho

from flask import Flask, request, Response, session, send_file
import os, random, time, datetime, string, hashlib, json
from config import *

def rand_str(length=16):
    return ''.join(random.sample(string.ascii_letters + string.digits, length))

app = Flask(__name__)
app.config['SECRET_KEY'] = os.urandom(32)

random.seed(time.time())

@app.route('/')
def index():
    html = '<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/pure-min.css">\n<div class="pure-g"><div class="pure-u-1-5"></div><div class="pure-u-3-5"><p><h1>Minecraft Real World</h1></p>\n<span style="color:red;">[IMPORTANT WARNING]<br>\n<span style="margin-left: 30px;">The binary is vulnerable. Please DO NOT run in production environment, or may be attacked by others.</span><br>\n<span style="margin-left: 30px;">Try some tricks to get flag!</span><br>\nHave fun! lol. @impakho</span>\n<p>Info: the game service is running on another server.</p>\n<p><a class="pure-button" href="/source">Source Code</a></p>\n</div><div class="pure-u-1-5"></div></div>'

    return Response(html, mimetype='text/html')


@app.route('/deploy')
def deploy():
    if not checkPoW(session, request.args.get('work')):
        return Response('PoW check fail.', mimetype='text/html')

    player_name = 'bot-' + rand_str(8)
    deployClient(player_name)

    html = 'Deploy succ. The bot player\'s name is "' + player_name + '".'

    return Response(html, mimetype='text/html')


@app.route('/pow')
def pow():
    text = rand_str()
    result = hashlib.sha256(text.encode()).hexdigest()

    session['text'] = text[:12]
    session['result'] = result

    return Response(json.dumps({'text': text[:12], 'hash': result}), mimetype='application/json')


@app.route('/source')
def source():
    html = open(__file__).read()

    return Response(html, mimetype='text/plain')


@app.route('/server')
def server():
    return send_file('server.zip', attachment_filename='server.zip')


@app.route('/client')
def client():
    return send_file('client.zip', attachment_filename='client.zip')


def checkPoW(session, work):
    if 'text' not in session or 'result' not in session or work == None:
        return False

    text = session['text']
    result = session['result']

    del session['text']
    del session['result']

    if len(text) != 12 or len(result) != 64 or len(work) != 4:
        return False

    if hashlib.sha256((text + work).encode()).hexdigest() != result:
        return False

    return True


if __name__ == '__main__':
    app.permanent_session_lifetime = datetime.timedelta(minutes=5)
    app.run(host='0.0.0.0', port=80, debug=False)
