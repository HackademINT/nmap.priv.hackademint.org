#! /usr/bin/python3

import nmap, socket, struct
from socket import inet_aton
from flask import Flask, jsonify
from flask_cache import Cache


app = Flask(__name__)
cache = Cache(app,config={'CACHE_TYPE': 'simple'})
cache.init_app(app)
SUBNET = '172.16.0.0/24'


def scan_subnet(subnet):
    return nmap.PortScanner().scan(hosts=subnet, arguments="-sP")


@cache.cached(timeout=1000)
@app.route('/')
def index():
    d = {'author': 'zTeeed', 'follow_me': 'https://github.com/zteeed', 
         'paths': {0: '/scan', 1: '/ip'}}
    return jsonify(d)


@cache.cached(timeout=60)
@app.route('/scan')
def scan():
    result = scan_subnet(SUBNET)
    return jsonify(result)


@cache.cached(timeout=60)
@app.route('/ip')
def ip(d = {}):
    result = scan_subnet(SUBNET)
    ips = sorted(result['scan'].keys(), key=lambda ip: socket.inet_aton(ip))
    d = {key: value for (key, value) in enumerate(ips)}
    return jsonify(d)


if __name__ == '__main__':
    app.run(host='0.0.0.0')
