#!/bin/bash
set -e

rm dbcc*
#gox -osarch="linux/amd64"
gox -os="linux"
sha256sum dbcc* > SHA256SUMS
