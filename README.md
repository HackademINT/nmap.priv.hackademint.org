# nmap.priv.hackademint.org

- python api: https://nmap.priv.hackademint.org
- go api: https://nmap2.priv.hackademint.org

## Install 

### python

```bash
sudo apt install python3 python3-pip
pip3 install -r requirements.txt
```

### go

```
wget https://dl.google.com/go/go1.12.2.linux-amd64.tar.gz
tar -xvf go1.12.2.linux-amd64.tar.gz
sudo mv go /usr/local
export GOROOT=/usr/local/go
export GOPATH=/root/nmap.priv.hackademint.org
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

go get -u github.com/HackademINT/nmap.priv.hackademint.org/handler
go get -u github.com/HackademINT/nmap.priv.hackademint.org/gateway
go build main.go
```

## Benchmarking

tool: [https://github.com/wg/wrk](https://github.com/wg/wrk)

### python

```bash
wrk -t12 -c400 -d30s http://localhost:5000/ip
Running 30s test @ http://localhost:5000/ip
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    82.80ms   71.96ms   1.74s    98.30%
    Req/Sec   149.10     93.70   444.00     59.76%
  51391 requests in 30.10s, 260.20MB read
  Socket errors: connect 0, read 55, write 0, timeout 40
Requests/sec:   1707.60
Transfer/sec:      8.65MB

```

### go

```bash
wrk -t12 -c400 -d30s http://localhost:5001/ip
Running 30s test @ http://localhost:5001/ip
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    18.30ms   24.65ms 338.61ms   86.28%
    Req/Sec     3.83k   606.50     7.80k    68.25%
  1376244 requests in 30.08s, 4.95GB read
Requests/sec:  45745.51
Transfer/sec:    168.40MB
```

> No Socket errors

## systemd

### python

`/etc/systemd/system/nmap.service`

```
[Unit]
Description=Nmap

[Service]
WorkingDirectory=/root/nmap.priv.hackademint.org
ExecStart=/usr/bin/python3 main.py
User=root
Group=root

[Install]
WantedBy=default.target
```

```bash
systemctl enable nmap.service
systemctl start nmap.service
```

### go


`cat /etc/systemd/system/nmap2.service`
```bash
[Unit]
Description=Nmap2

[Service]
ExecStart=/root/nmap.priv.hackademint.org/main
User=root
Group=root

[Install]
WantedBy=default.target
```

```bash
systemctl enable nmap2.service
systemctl start nmap2.service
```

## Troubleshooting

### No module named 'flask.ext'

IF

```bash
Traceback (most recent call last):
  File "./main.py", line 9, in <module>
    cache = Cache(app,config={'CACHE_TYPE': 'simple'})
  File "/usr/local/lib/python3.5/dist-packages/flask_cache/__init__.py", line 121, in __init__
    self.init_app(app, config)
  File "/usr/local/lib/python3.5/dist-packages/flask_cache/__init__.py", line 156, in init_app
    from .jinja2ext import CacheExtension, JINJA_CACHE_ATTR_NAME
  File "/usr/local/lib/python3.5/dist-packages/flask_cache/jinja2ext.py", line 33, in <module>
    from flask.ext.cache import make_template_fragment_key
ImportError: No module named 'flask.ext'
```

THEN

```bash
sed -i 's|flask.ext.cache|flask_cache|g' /usr/local/lib/python3.5/dist-packages/flask_cache/jinja2ext.py 
```

FI
