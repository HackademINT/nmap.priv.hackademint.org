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
wrk -t12 -c400 -d30s https://nmap.priv.hackademint.org/ip
Running 30s test @ https://nmap.priv.hackademint.org/ip
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   349.17ms  350.91ms   2.00s    83.16%
    Req/Sec    92.69     35.54   232.00     69.16%
  27006 requests in 30.10s, 9.56MB read
  Socket errors: connect 0, read 0, write 0, timeout 410
Requests/sec:    897.25
Transfer/sec:    325.08KB
```

### go

```bash
Running 30s test @ https://nmap2.priv.hackademint.org
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   313.44ms   50.50ms   1.09s    70.78%
    Req/Sec   105.29     45.24   282.00     64.71%
  36346 requests in 30.09s, 8.94MB read
Requests/sec:   1208.01
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
