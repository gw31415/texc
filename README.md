#これは何
LaTeXのオンラインコンパイラのtex.amas.devを使う際のスクリプト.
# インストール
## システム要件
* Bourne Shell互換のシェル (sh, bash, zshなど, fishやcmdはムリ)
* cURLのコマンド
## インストール方法
1. ダウンロードする
1. (必要なら)実行権限をつける (`chmod +x ./texc` とか)
1. (必要なら)パスを通す
# 使い方の例
以下を実行する
```
texc ./path/to/example.tex
```
現在のディレクトリに `example.pdf` が出てきます
# 環境
* texlive2019
* uplatex + dvipdfmx を latexmkで括ったもの
* `.latexmkrc`は以下みたいな感じ
``` perl
#!/usr/bin/env perl
$latex            = 'uplatex -halt-on-error';
$latex_silent     = 'uplatex -halt-on-error -interaction=batchmode';
$dvipdf           = 'dvipdfmx %O -o %D %S';
$makeindex        = 'mendex %O -o %D %S';
$max_repeat       = 5;
$pdf_mode         = 3;
$clean_ext		  = 'dvi';
```
