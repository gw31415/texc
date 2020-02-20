# これは何
LaTeXのオンラインコンパイラのtex.amas.devを使う際のスクリプト.
# インストール
## システム要件
* Bourne Shell互換のシェル (sh, bash, zshなど. fishやcmdはムリ)
* cURLのコマンド
## インストール方法
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
sudo mv ./texc /usr/local/bin/
```
# 使い方の例
引数が必要. 以下を実行する
``` bash
texc ./path/to/example.tex
```
現在のディレクトリに `example.pdf` が出てきます
# 環境
* texlive2019
* uplatex + dvipdfmx を latexmkで括ったもの
* `.latexmkrc`は以下みたいな感じ
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
