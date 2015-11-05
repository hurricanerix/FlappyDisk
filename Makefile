# Copyright 2015 Richard Hawkins
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
EXE=flappy

default: bin/$(EXE)


run: bin/$(EXE)
	./bin/$(EXE)

bin/$(EXE):	main.go gen/assets.go
	go build main.go
	mkdir -p bin
	mv main bin/$(EXE)

gen/assets.go: $(shell find assets -type f)
	mkdir -p gen
	go-bindata -pkg gen -o gen/assets.go assets/

clean:
	rm -rf bin gen
