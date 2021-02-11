# Traitor

Local Unix privilege escalation made simple.

Packages up a bunch of methods to exploit local vulnerabilities and misconfigurations (including all of GTFOBins) in order to gain a root shell.

## Usage

Run with no arguments to find potential vulnerabilities/misconfigurations which could allow privilege escalation.

```bash
traitor
```

Run with the `-a`/`--any` flag to find potential vulnerabilities, attempting to exploit each, stopping if a root shell is gained.

```bash
traitor -a
```

Run with the `-e`/`--exploit` flag to attempt to exploit a specific vulnerability and gain a root shell.

```bash
traitor -e docker:writable-socket
```

## Getting Traitor

Grab a binary from the [releases page](https://github.com/liamg/traitor/releases), or use go:

```
go get -u github.com/liamg/traitor/cmd/traitor
```

If the machine you're attempting privesc on cannot reach GitHub to download the binary, and you have no way to upload the binary to the machine over SCP/FTP etc., then you can try base64 encoding the binary on your machine, and echoing the base64 encoded string to `| base64 -d > /tmp/traitor` on the target machine, remembering to `chmod +x` it once it arrives.

## Included Methods

- [x] Writable `docker.sock` (no internet connection or local images required!)
- [ ] sudo:CVE-2021-3156
- [ ] Basic sudo
- [ ] GTFOBins via weak sudo rules
- [ ] Kernel exploits

## TODO

- [ ] Add a whole bunch of methods
- [x] Switch out `/bin/bash` for `traitor shell` as a setuid shell wrapper
