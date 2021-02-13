# Traitor

Linux privilege escalation made easy.

Packages up a bunch of methods to exploit local misconfigurations/vulns (including most of [GTFOBins](https://gtfobins.github.io/)) in order to gain a root shell.

![Demo](demo.gif)

## Usage

Run with no arguments to find potential vulnerabilities/misconfigurations which could allow privilege escalation. Add the `-p` flag if the current user password is known.

```bash
traitor -p
```

Run with the `-a`/`--any` flag to find potential vulnerabilities, attempting to exploit each, stopping if a root shell is gained. Again, add the `-p` flag if the current user password is known.

```bash
traitor -a -p
```

Run with the `-e`/`--exploit` flag to attempt to exploit a specific vulnerability and gain a root shell.

```bash
traitor -p -e docker:writable-socket
```

## Getting Traitor

Grab a binary from the [releases page](https://github.com/liamg/traitor/releases), or use go:

```
CGO_ENABLED=0 go get -u github.com/liamg/traitor/cmd/traitor
```

If the machine you're attempting privesc on cannot reach GitHub to download the binary, and you have no way to upload the binary to the machine over SCP/FTP etc., then you can try base64 encoding the binary on your machine, and echoing the base64 encoded string to `| base64 -d > /tmp/traitor` on the target machine, remembering to `chmod +x` it once it arrives.
