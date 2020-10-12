# linux-screen-reader

## これはなに？
Linuxデスクトップ環境でgoogle text to speech を用いてスクリーンリーダーを提供します．


## コンポーネント

golang実装のAPIサーバーとJavaScript実装のフロントエンドから構成されます．

APIはhttpサーバーとクリップボード監視で読み上げる文章を待機します．

## 依存関係

### 実行時

* `alsa-lib` or `libasound2-dev`
* [`clipnotify`](https://github.com/cdown/clipnotify) 

### 開発時
* `go mod download` 

## 実行
* go run ./cmd or binary download by release page
