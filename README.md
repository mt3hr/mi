# TODOアプリケーション「mi」
miはTODOアプリケーションです。  
![mi](https://raw.githubusercontent.com/mt3hr/mi/main/document/img/mi.png)  

## ダウンロード
[miダウンロード](https://github.com/mt3hr/mi/releases/latest)  
Windowsにてウイルスであると検知されます。  
おそらく誤検知ですが、不安な方はインストールしないでください。（[参考](https://zenn.dev/akb428/articles/a7fdcd2614e366)）  
対応予定はありません。  

## 実行
「mi.exe」または「mi_server.exe」をダブルクリック  
（mi.server.exeの場合は起動後「[http://localhost:2734](http://localhost:2734)」にアクセス）  

<details>
<summary>開発者向け</summary>

開発者向けと言いつつ自分向けです。  
ビルドに必要パッケージを公開していないのでビルド不可能だと思います。  

### 開発環境

### セットアップ
1. Golang バージョン1.20の開発環境を用意する  
2. Cコンパイラを用意する（cgo使用のため）  
3. Node.js バージョン18.12.1の開発環境を用意する  
4. 以下のスクリプトを実行する  
```
npm i
```

### ビルド・インストール

アプリケーションインストール  
```
npm run go_mod
npm run install_app
```

サーバインストール  
```
npm run go_mod
npm run install_build
```
</details>