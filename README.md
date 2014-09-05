ipmes2yammer
=========

IP Messengerのログファイルをtailで追っかけて、更新があるとYammerの指定スレッドに投稿します。


ビルド方法
=====


### クロスコンパイルの準備

Windows用なのでLinuxやMacでbiudする場合は必要です。

```bash:
$ cd $(go env GOROOT)/src # go をインストールしたディレクトリに移動
$ GOOS=linux GOARCH=amd64 ./make.bash
```

### clientIDとClientSecretの設定

http://developer.yammer.com/introduction/ に従いアプリ登録してclientIDとClientSecretを得てから

```bash:
cp yammer/var.go.example yammer/var.go
```

`vi yammer/var.go` で以下の箇所を書き換える

```go:
	clientId     = "xxxxxxxxxxxxfa"
	clientSecret = "falsdfjaslfkjsalfkjaslfkjasldfkjalj"
```
### buid


```bash:
make
```


使い方
======

```
ipmes2yammer.exe -to_id=投稿スレッドID -file=IPmesのログファイル

```


例:
```
ipmes2yammer.exe -to_id=332256272 -file=C:\Users\hoge\Documents\ipmsg.log
```

