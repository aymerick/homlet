# homlet

## Build

To build:

```
$ make
```

To generate files:

```
$ make gen
```

To generate, lint and build:

```
$ make release
```

To build a raspberry binary:

```
$ make rpi
```

## Terminal UI

Display all devices directly on you prefered terminal:

```
$ homlet term
```

## Server

Install homlet binary:

```
$ sudo mv homlet /usr/local/bin/
```

Install homlet conf file:

```
$ sudo mkdir /etc/homlet
$ sudo cp configs/config.toml /etc/homlet
```

Install homlet service:

```
$ sudo cp configs/homlet.service /etc/systemd/system/
$ sudo systemctl enable homlet.service
```

Start homlet service:

```
$ sudo systemctl start homlet.service
```

Tail logs:

```
$ journalctl -f -u homlet
```
