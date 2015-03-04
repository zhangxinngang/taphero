# -*- coding:utf-8 -*-
.PHONY: build
.PHONY: test
default:build
get-dep:
	go get -u github.com/go-sql-driver/mysql
	go install github.com/go-sql-driver/mysql
	go get -u github.com/Unknwon/goconfig
	go install  github.com/Unknwon/goconfig
	go get -u github.com/cihub/seelog
	go install github.com/cihub/seelog
	go get -u github.com/stretchr/testify/assert
	go install github.com/stretchr/testify/assert
	go get -u github.com/realint/dbgutil
	go install github.com/realint/dbgutil
	go get -u github.com/deckarep/golang-set
	go install github.com/deckarep/golang-set
	go get -u github.com/garyburd/redigo/redis
	go install github.com/garyburd/redigo/redis
	go get -u github.com/dropbox/godropbox
	go install github.com/dropbox/godropbox
	go get -u github.com/qiniu/api
	go install github.com/qiniu/api
	go get -u github.com/fanngyuan/link
	go install github.com/fanngyuan/link
	go get -u github.com/0studio/storage_key
	go install github.com/0studio/storage_key
	go get -u github.com/0studio/goauth
	go install github.com/0studio/goauth
	go get -u github.com/0studio/bit
	go install github.com/0studio/bit
	go get -u github.com/0studio/cachemap
	go install github.com/0studio/cachemap
	go get -u github.com/0studio/mcstorage
	go install github.com/0studio/mcstorage
	go get -u github.com/0studio/databuffer
	go install github.com/0studio/databuffer
	go get -u github.com/0studio/scheduler
	go install github.com/0studio/scheduler
	go get -u github.com/0studio/idgen
	go install github.com/0studio/idgen
	go get -u github.com/golang/groupcache
	go install github.com/golang/groupcache
	go get -u github.com/tealeg/xlsx
	go install github.com/tealeg/xlsx
	go get -u github.com/hraban/lrucache
	go install github.com/hraban/lrucache

	go get -u github.com/golang/protobuf/proto
	go install github.com/golang/protobuf/proto

	go get github.com/gogo/protobuf/proto
	go install github.com/gogo/protobuf/proto
	go get github.com/gogo/protobuf/protoc-gen-gogo
	go install github.com/gogo/protobuf/protoc-gen-gogo
	go get -u github.com/gogo/protobuf/gogoproto
	go install  github.com/gogo/protobuf/gogoproto

	go get -u github.com/0studio/redisapi
	go install github.com/0studio/redisapi

build:linkdata
	go install	zerogame.info/taphero/defs
	go install	zerogame.info/taphero/utils
	go install	zerogame.info/taphero/conf
	go install	zerogame.info/taphero/net
	go install	zerogame.info/taphero/log
	go install	zerogame.info/taphero/pub
	go install	zerogame.info/taphero/pf
	go install	zerogame.info/taphero/timer
	go install	zerogame.info/taphero/entity
	go install	zerogame.info/taphero/design
	go install	zerogame.info/taphero/resource
	go install	zerogame.info/taphero/service/user_add_attr
	go install	zerogame.info/taphero/service
	go install	zerogame.info/taphero/redis
	go install	zerogame.info/taphero/redis_msg
	go install	zerogame.info/taphero/logic
	go install	zerogame.info/taphero/app
	go install 	zerogame.info/taphero/taphero 
build-dev:
	@if [ !  -d proto ]; then  \
	    echo "please run :svn co https://118.192.76.99/svn/taphero/dev/proto; if failed to get proto/*" ; \
	    svn co https://118.192.76.99/svn/taphero/dev/proto; \
	fi 
	cd proto &&sh ./build.sh && cd ..
	make build
linkdata:
	@echo "ln -s /data/taphero/config"
	@rm -rf /data/taphero/config
	@mkdir -p /data/taphero
	@ln -s -f $(GOPATH)/src/zerogame.info/taphero/config /data/taphero/config

run:
	go run taphero/main.go -locale eng
test:
	go test zerogame.info/taphero/service/user_add_attr
data-test:
	./excel_export.sh test strict
data-eng:
	./excel_export.sh eng strict
package:build
	rm -rf dist
	mkdir -p dist
	cp $(GOPATH)/bin/taphero dist/
	cp  -rf config dist/
	tar -cjf taphero.tar.bz2 dist 
	mv taphero.tar.bz2 /data/ftp/th