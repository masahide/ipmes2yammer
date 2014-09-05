ipmes2yammer
=========

IP Messengerのログファイルをtailで追っかけて、更新があるとYammerの指定スレッドに投稿します。



使い方
======

```
ipmes2yammer.exe -to_id=投稿スレッドID -file=IPmesのログファイル

```


例:
```
ipmes2yammer.exe -to_id=332256272 -file=C:\Users\hoge\Documents\ipmsg.log
```

