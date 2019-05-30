#! /usr/bin/python3

import nmap, socket, struct
from socket import inet_aton
from flask import Flask, jsonify
from flask_caching import Cache


app = Flask(__name__)
cache = Cache(app,config={'CACHE_TYPE': 'simple'})
cache.init_app(app)
SUBNET = '172.16.0.0/24'


def scan_subnet(subnet):
    return nmap.PortScanner().scan(hosts=subnet, arguments="-sP")


@app.route('/')
@cache.cached(timeout=1000)
def index():
    d = {'author': 'zTeeed', 'follow_me': 'https://github.com/zteeed', 
         'paths': {0: '/scan', 1: '/ip'}}
    return jsonify(d)


@app.route('/scan')
@cache.cached(timeout=60)
def scan():
    result = scan_subnet(SUBNET)
    return jsonify(result)


@app.route('/ip')
@cache.cached(timeout=60)
def ip(d = {}):
    result = scan_subnet(SUBNET)
    ips = sorted(result['scan'].keys(), key=lambda ip: socket.inet_aton(ip))
    d = {key: value for (key, value) in enumerate(ips)}
    return jsonify(d)


if __name__ == '__main__':
    app.run(host='0.0.0.0')
