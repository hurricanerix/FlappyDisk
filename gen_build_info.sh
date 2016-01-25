#!/bin/bash
# Copyright 2015-2016 Richard Hawkins
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

DATE=`date`
BUILD_HASH=`git log -n 1 | grep commit | cut -d " " -f 2`
VERSION="0.0" # `git describe --abbrev=0`

git --no-pager diff --exit-code "$VERSION" master > /dev/null 2>&1
PRE_RELEASE=$?

# TODO: Hash should not be added when building from a tagged version
#       but it looks like it currently does.  This can be fixed later.
if [ $PRE_RELEASE ]; then
  VERSION="$VERSION.$BUILD_HASH"
fi

CODE="// gen is a generated package, DO NOT EDIT!\n
\n
package gen\n
\n
var Version string = \"$VERSION\"\n
var BuildDate string = \"$DATE\"\n
var BuildHash string = \"$BUILD_HASH\"\n
"

TMP=`go-bindata -version`
if [ $? -ne 0 ]
then
  # go-bindata is required to embed assets into binary
  echo "Downloading go-bindata"
  go get -u github.com/jteeuwen/go-bindata/...
fi

# generate all the files we need
mkdir -p gen
go-bindata -pkg gen -ignore="/*.pyxel" -o gen/assets.go assets/
echo -e $CODE | gofmt > gen/build_info.go
