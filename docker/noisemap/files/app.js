/**
* @Author: impakho
* @Date: 2020/04/12
* @Github: https://github.com/impakho
*/

const http = require('http');
const path = require('path');
const fs = require('fs');
const gen = require('random-seed');
const { Cluster } = require('puppeteer-cluster');
const sharp = require('sharp');
const TextToSVG = require('text-to-svg');
const { secret } = require('./secret');

const DEBUG = false;
const PORT = 80;

if (typeof secret !== 'string' || secret.length % 2 != 0) {
  console.log('invalid secret');
  return;
}

var textToSVG = TextToSVG.loadSync('./font.ttf');
const fontOptions = {x: 0, y: 0, fontSize: 6, anchor: 'top', attributes: { fill: '#fff', stroke: '#fff' }};

var reqWhiteList = [
  '/index.html',
  '/assets/jquery.min.js',
  '/assets/noisemap.js',
  '/assets/p5.dom.js',
  '/assets/p5.js'
];

var lock = false;

function randStr(length) {
   var result           = '';
   var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
   var charactersLength = characters.length;
   for ( var i = 0; i < length; i++ ) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
   }
   return result;
}

function genMaps() {
  if (lock) {
    setTimeout(genMaps, 15000);
    return;
  }

  var directory = 'www/maps/';
  fs.readdir(directory, (err, files) => {
    if (err) return;

    for (const file of files) {
      fs.unlink(path.join(directory, file), err => {
        if (err) return;
      });
    }
  });

  var time = Math.floor(Date.now() / 1000);

  var rand = gen.create(time);

  var bytes = Buffer.from(secret);
  for (var i = 0; i < bytes.length; i++) {
    bytes[i] ^= rand.range(0x100);
  }

  var groups = [];
  for (var i = 0; i < bytes.length / 2; i++) {
    var high = (bytes[2 * i + 1] + bytes[2 * i]) & 0xff;
    var low = (bytes[2 * i + 1] - bytes[2 * i]) & 0xff;
    groups.push((high << 8) + low);
  }

  if (DEBUG) {
    console.log(time, bytes, groups);
  }

  lock = true;

  (async () => {
    const cluster = await Cluster.launch({
      concurrency: Cluster.CONCURRENCY_PAGE,
      maxConcurrency: 4,
      puppeteerOptions: {
        args: ['--no-sandbox', '--disable-setuid-sandbox'],
        defaultViewport: {
          width: 520,
          height: 520
        }
      }
    });

    await cluster.task(async ({ page, data }) => {
      await page.goto('http://127.0.0.1:' + PORT + '/map?seed=' + data.seed);
      var font = await Buffer.from(textToSVG.getSVG(randStr(32), fontOptions));
      await sharp(await page.screenshot())
        .composite([{ input: font, tile: true }])
        .removeAlpha()
        .toFile('www/maps/' + data.id + '.webp', (err, info) => {
          if (err) return console.log(err);
        });
    });

    for (var i = 0; i < groups.length; i++) {
      cluster.queue({id: i, seed: groups[i]});
    }

    await cluster.idle();
    await cluster.close();

    lock = false;
    setTimeout(genMaps, 20000 + rand.range(20000));
  })();

}

function reqFile(res, pathname, html) {
  fs.readFile(pathname,
    function (err, data) {
      if (err) {
        res.writeHead(404);
        return res.end('404 not found');
      }

      var ext = path.extname(pathname);

      var typeExt = {
        '.html': 'text/html',
        '.js':   'text/javascript',
        '.css':  'text/css',
        '.webp':  'image/webp'
      };

      var contentType = typeExt[ext] || 'text/plain';

      res.writeHead(200, { 'Content-Type': contentType });
      res.write(data);
      res.end(html);
    }
  );
}

function handleRequest(req, res) {
  var pathname = req.url;
  console.log(pathname);
  
  if (pathname === '/') {
    pathname = '/index.html';
  }

  if (reqWhiteList.includes(pathname)) {
    var html = '';

    if (pathname === '/index.html') {
      html += '<script>document.getElementById("maps").innerHTML=`\n';

      if (lock) {
        html += '<h3>Waiting for map generator...</h3>';
      }else{
        for (var i = 0; i < secret.length / 2; i++) {
          html += '<img src="maps/' + i + '.webp" width="256" height="256" style="float:left;">\n'
        }
      }

      html += '`;\nsetTimeout(function() {location.reload();}, 10000);\n</script>\n';
    }

    reqFile(res, 'www' + pathname, html);
    return;
  }

  if (pathname.substring(0, 10) === "/map?seed=") {
    var seed = pathname.substring(10);

    seed = parseInt(seed);

    var html = '';

    html = '<script>var seed = ' + seed + ';</script>';

    reqFile(res, 'www/map.html', html);
    return;
  }

  if (pathname.substring(0, 6) === '/maps/' && pathname.substring(pathname.length - 5) === '.webp') {
    var id = pathname.substring(6, pathname.length - 5);

    id = parseInt(id);

    var html = '';

    reqFile(res, 'www/maps/' + id + '.webp', html);
    return;
  }

  if (pathname === "/source") {
    fs.readFile(__filename,
      function (err, data) {
        if (err) {
          res.writeHead(404);
          return res.end('404 not found');
        }

        res.writeHead(200, { 'Content-Type': 'text/plain' });
        res.write(data);
        res.end(html);
      }
    );
    return;
  }

  res.writeHead(404);
  res.end('404 not found');
}

var server = http.createServer(handleRequest);
server.listen(PORT);

console.log('Server started on port ' + PORT);

genMaps();
