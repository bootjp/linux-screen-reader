# linux-screen-reader

## これはなに？
Linuxデスクトップ環境でgoogle text to speech を用いてスクリーンリーダーを提供します．


## コンポーネント

golang実装のAPIサーバーとJavaScript実装のフロントエンドから構成されます．

APIはhttpサーバーとクリップボード監視で読み上げる文章を待機します．

## 依存関係

`go.mod` 及び `bin/clipnotify` 

[`clipnotify`](https://github.com/cdown/clipnotify) は public domain のバイナリを同梱

## 実行

