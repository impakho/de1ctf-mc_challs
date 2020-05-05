#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# @Author: impakho
# @Date:   2020/04/12
# @Github: https://github.com/impakho

from flask import Flask, request, Response, session, render_template_string
import posix, os, sys, signal, random, time, datetime, string, hashlib, json, threading
from uuid import UUID

def rand_str(length=16):
    return ''.join(random.sample(string.ascii_letters + string.digits, length))

app = Flask(__name__)
app.config['SECRET_KEY'] = os.urandom(32)

# Some bad words.
blacklist = ['+', ',', ':', '\'\'', '""', '%', 'lower', 'upper', 'builtin', 'fork', 'exec', 'walk', 'open', 'spawn', 'reload', 'exit', 'bin', 'sh', 'cat', 'config', 'secret', 'key', 'flag']

# Posix is a bad module, filter it all.
for i in dir(posix):
    blacklist.append(i.lower())

random.seed(time.time())

def printableFilter(s):
    return ''.join(filter(lambda x: x in string.printable[:-2], s))


@app.route('/')
def index():
    html = '<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/pure-min.css">\n<div class="pure-g"><div class="pure-u-1-5"></div><div class="pure-u-3-5"><p><h1>Minecraft Log Web Client</h1></p>\n<span style="color:red;">[IMPORTANT WARNING]<br>\n<span style="margin-left: 30px;">For ATTACKER: your log file is public! Please try some tricks to keep your payload in secret, or may leak and stealed by others.</span><br>\n<span style="margin-left: 30px;">For STEALER: you may facing XSS attack by ATTACKER.</span><br>\nHave fun! lol. @impakho</span>\n<p><a class="pure-button" href="/source">Source Code</a></p>\n<table class="pure-table pure-table-horizontal" style="width:100%;">\n<thead><tr><th>Filename</th></tr></thead>\n<tbody>\n'

    filelist = ''

    for root, dirs, files in os.walk('./logs/'):
        for name in files:
            filelist += '<tr><td><a href="read?filename=' + name + '">' + name + '</a></td></tr>\n'

    if len(filelist) <= 0:
        html += '<tr><td><i>empty</i></td></tr>\n'
    else:
        html += filelist

    html += '</tbody>\n</table>\n</div><div class="pure-u-1-5"></div></div>'

    return Response(html, mimetype='text/html')


@app.route('/pow')
def pow():
    text = rand_str()
    result = hashlib.sha256(text.encode()).hexdigest()

    session['text'] = text[:12]
    session['result'] = result

    return Response(json.dumps({'text': text[:12], 'hash': result}), mimetype='application/json')


@app.route('/read')
def read():
    if not checkPoW(session, request.args.get('work')):
        return Response('PoW check fail.', mimetype='text/html')

    filename = request.args.get('filename')

    try:
        val = UUID(filename, version=4)
    except ValueError:
        return Response('Not a valid UUID filename.', mimetype='text/html')

    try:
        fp = open('./logs/' + filename, 'r')
    except:
        return Response('File not exist.', mimetype='text/html')

    binary = printableFilter(fp.read())
    fp.close()

    # Check blacklist
    for i in blacklist:
        if i in binary.lower():
            return Response('Bad log file.', mimetype='text/html')

    # Do some replacement
    binary = binary.replace(' ', '').replace('<', '&lt;').replace('>', '&gt; ').replace('\n', '<br />\n')
    html = '<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/pure-min.css">\n<div class="pure-g"><div class="pure-u-1-5"></div><div class="pure-u-3-5"><p><h1>Minecraft Log Web Client</h1></p>File <' + filename + '><p>\n'

    session['filename'] = filename

    html += binary

    html += '\n</p></div><div class="pure-u-1-5"></div></div>'

    return Response(renderHandler(session, html), mimetype='text/html')


@app.route('/write')
def write():
    if not checkPoW(session, request.args.get('work')):
        return Response('PoW check fail.', mimetype='text/html')

    if 'w' not in session or 'filename' not in session:
        return Response('Select a log file first.', mimetype='text/html')

    text = request.args.get('text')

    if text == None or len(text) <= 0:
        return Response('Text is empty.', mimetype='text/html')

    if len(text) > 512 or not all(c in string.printable for c in text):
        return Response('Invalid text format.', mimetype='text/html')

    try:
        # Write to stdout
        print(session['filename'] + ' ' + text + '\n')

        # Write to log/log.txt
        open('log/log.txt', 'a').write(session['filename'] + ' ' + text + '\n')

        # Write to child
        w = os.fdopen(session['w'], 'w')
        w.write(text)
        w.close()
        del session['w']
    except:
        return Response('Write fail.', mimetype='text/html')

    return Response('Write succ.', mimetype='text/html')


@app.route('/source')
def source():
    html = open(__file__).read()

    return Response(html, mimetype='text/plain')


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


def renderHandler(session, html):
    renderCleanUp(session)

    # For security, fork a child process to render.

    r, w = os.pipe()
    session['w'] = w

    pid = os.fork()

    if pid:
        os.close(r)

        # Check child process status, and wait it to finish
        thread = waitThread(pid)
        thread.start()

    else:
        signal.signal(signal.SIGALRM, kill)
        signal.alarm(30)

        os.close(w)
        r = os.fdopen(r, 'r')
        sys.stdin = r

        try:
            render_template_string(html)
        except:
            pass
        kill(None, None)

    return html


def renderCleanUp(session):
    try:
        os.close(session['w'])
    except:
        pass


def kill(signum, frame):
    os.kill(os.getpid(), signal.SIGKILL)


class waitThread(threading.Thread):
    def __init__(self, pid):
        threading.Thread.__init__(self)
        self.pid = pid


    def run(self):
        count = 0
        while True:
            if count >= 30:
                try:
                    os.kill(self.pid, signal.SIGKILL)
                except:
                    break
            try:
                os.waitpid(self.pid, os.WNOHANG)
            except:
                break
            count += 1
            time.sleep(1)


if __name__ == '__main__':
    app.permanent_session_lifetime = datetime.timedelta(minutes=5)
    app.run(host='0.0.0.0', port=80, debug=False)
