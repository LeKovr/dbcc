#!/bin/bash
set -e

DIST=$1

for f in dbcc* ; do [ -f $f ] && rm $f ; done

if [[ "$DIST" ]] ; then
  gox -os="linux"
else
  gox -osarch="linux/amd64"
fi

sha256sum dbcc* > SHA256SUMS
