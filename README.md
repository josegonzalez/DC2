# dc2 [![CircleCI](https://circleci.com/gh/josegonzalez/dc2.svg?style=svg)](https://circleci.com/gh/josegonzalez/dc2)

A replacement api service for the [DC2](https://github.com/hardtware/DC2) built in golang.

## requirements

- golang 1.11+

## usage

> For a prebuilt binary, see the [github releases page](https://github.com/josegonzalez/dc2/releases).

```shell
# build the binary
make build

# copy to your server via scp
scp build/linux/dc2 jose@dc2.local:/tmp/

# for systemd systems, copy the dc2.service
scp init/upstart/dc2.service jose@dc2.local:/tmp/

# for upstart systems, copy the dc2.service
scp init/upstart/dc2.conf jose@dc2.local:/tmp/

# from the dc2 server, change ownership on the files
sudo chown root:root /tmp/dc2*

# copy the files into place (mv the correct init file as well)
sudo mv /tmp/dc2 /usr/local/bin/dc2
sudo mv /tmp/dc2.service /etc/systemd/system/dc2.service
sudo mv /tmp/dc2.conf /etc/init/dc2.conf

# start the service and enable it at boot
sudo systemctl start dc2.service
sudo systemctl enable dc2.service

# curl the software
curl dc2.local:8765
```

## differences

- Doesn't ping an external discovery service.
    - The official version requests `http://dc2.hardtware.com/version/package.json`, but this no longer responds with a valid response.
- Requests latest version from Github Releases.
