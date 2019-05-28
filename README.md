# nmap.priv.hackademint.org

python api: https://nmap.hackademint.org
go api: https://nmap2.hackademint.org

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

go get -u github.com/Ullaakut/nmap
go build main.go
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
