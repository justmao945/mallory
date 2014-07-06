#!/bin/sh
grep "code.google.com/p/go.crypto/ssh" ssh/ -R|awk '{print $1}'|cut -d ':' -f 1|xargs sed -i '' -e 's/code.google.com\/p\/go.crypto\/ssh/github.com\/justmao945\/mallory\/ssh/g'
