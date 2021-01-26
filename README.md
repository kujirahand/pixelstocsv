# pixelstocsv - 画像一覧をCSVファイルに

画像が入ったフォルダを pixelstocsv.exe にドラッグ＆ドロップすると、out.csv というファイルを生成します。
画像の各ピクセルを、赤,緑,青, 赤,緑,青, 赤,緑,青,... と並べてCSVファイルに変換します。


# how to use

```
$ pixelstocsv (inputdir) (outcsvfile)
```

# how to compile

```
$ go build
```

