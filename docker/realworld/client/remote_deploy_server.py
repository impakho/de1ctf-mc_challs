from flask import Flask, Response
import os

app = Flask(__name__)
app.config['SECRET_KEY'] = os.urandom(32)

pipe_file = '/var/lib/docker/volumes/de1ctf-realworld-pipe/_data/pipe.txt'

@app.route('/v2L4qFXhGU4AFiGv/deploy')
def deploy():
    pipe = open(pipe_file, 'r').read()
    open(pipe_file, 'w').write('')
    return Response(pipe, mimetype='text/plain')


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=443, debug=False)
