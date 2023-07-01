#!/bin/bash

# iptables v1.8.9 (legacy)
iptables-save > iptables.bak
iptables-restore < iptables.bak
apk del iptables
apk add --repository=https://dl-cdn.alpinelinux.org/alpine/v3.18/main --allow-untrusted iptables=1.8.9-r2


