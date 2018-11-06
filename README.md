# micromsg

## 实现
考虑到时间问题，为了快速实现，主要选择了自己相对熟悉了方案，前端采用了vuejs&muse-ui,后端使用了websocket，自己的通信类库pillx


## install
### client
* npm install
* npm run build
### server
* go build

## run
### client
* npm run dev
### server
1. start etcd ./etcd
2. start gateway ./gateway ../bin/etc/config.toml
3. start worker ./worker ../bin/etc/config.toml

## enjoy fun
