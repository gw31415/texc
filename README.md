<a href="https://asciinema.org/a/2z2ddsoonJZQW3gjGyFFL9VhK" target="_blank"><img src="https://asciinema.org/a/2z2ddsoonJZQW3gjGyFFL9VhK.svg" /></a>
# これは何
LaTeXのオンラインコンパイラのtex.amas.devを使う際のスクリプト.
# インストール
## システム要件
* Bourne Shell互換のシェル (sh, bash, zshなど. fishやcmdはムリ)
* tarとcURLのコマンド
## インストール方法
### これだけ
``` bash
curl -o ~/.texc https://raw.githubusercontent.com/gw31415/texc/master/texc ; chmod +x ~/.texc ; sudo ln -s ~/.texc /usr/local/bin/texc
```

### 詳細
1. ダウンロードする
``` bash
curl -O https://raw.githubusercontent.com/gw31415/texc/master/texc
```
2. (必要なら)実行権限をつける
``` bash
chmod +x ./texc
```
3. (必要なら)パスを通す, 通ったところに移動する.
``` bash
mv ./texc ~/.texc
sudo ln -s ~/.texc /usr/local/bin/texc
```
# 使い方の例
## texをpdfにする
``` bash
texc ./path/to/example.tex
```
現在のディレクトリに `example.pdf` が出てきます
* 環境は uplatex + dvipdfmx を latexmkで括ったもの
* latexmkを使います
## 複数ファイルにまたがる、または独自の設定を用いる場合
``` bash
cd path/to/project
texc -l ./main.tex
```
プロジェクトのディレクトリに `main.pdf` が出てきます
* 必ずプロジェクトのディレクトリまで移動してから行ってください
* デフォルトは uplatex + dvipdfmx を latexmkで括ったもの
* latexmkを使います
* `.latexmk` を同梱することができます
# 環境
* texlive2019
* デフォルトの`.latexmkrc`は以下みたいな感じ
``` perl
#!/usr/bin/env perl
$latex            = 'uplatex -halt-on-error -interaction=nonstopmode';
$latex_silent     = 'uplatex -halt-on-error -interaction=batchmode';
$bibtex           = 'upbibtex';
$biber            = 'biber --bblencoding=utf8 -u -U --output_safechars';
$dvipdf           = 'dvipdfmx %O -o %D %S';
$makeindex        = 'mendex %O -o %D %S';
$max_repeat       = 5;
$pdf_mode         = 3;
```
