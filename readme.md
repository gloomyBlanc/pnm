# pnm
PNM(PBM/PGM/PPM)形式の画像に対するエンコーダ・デコーダです。

## 使い方
```
import "github.com/gloomyBlanc/pnm"
```

### func [Decode](pnm/reader.go#28)
<pre>
func Decode(r <a href="https://pkg.go.dev/io">io</a>.<a href="https://pkg.go.dev/io#Reader">Reader</a>) (<a href="https://pkg.go.dev/image">image</a>.<a href="https://pkg.go.dev/image#Image">Image</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)
</pre>
PNM画像をReaderから読み込み、image.Imageとして返します。

### func [DecodeConfig](pnm/reader.go#33)
<pre>
func DecodeConfig(r <a href="https://pkg.go.dev/io">io</a>.<a href="https://pkg.go.dev/io#Reader">Reader</a>) (<a href="https://pkg.go.dev/image">image</a>.<a href="https://pkg.go.dev/image#Config">Config</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)
</pre>
PNM画像をReaderから読み込み、image.Configのみを返します。

### func [Encode](pnm/writer.go#17)
<pre>
func Encode(w <a href="https://pkg.go.dev/io">io</a>.<a href="https://pkg.go.dev/io#Writer">Writer</a>, img <a href="https://pkg.go.dev/image">image</a>.<a href="https://pkg.go.dev/image#Image">Image</a>) <a href="https://pkg.go.dev/builtin#error">error</a>
</pre>
imgをPGM(P5)もしくはPPM(P6)としてwに出力します。

### func [EncodeWithType](pnm/writer.go#23)
<pre>
func EncodeWithType(w <a href="https://pkg.go.dev/io">io</a>.<a href="https://pkg.go.dev/io#Writer">Writer</a>, img <a href="https://pkg.go.dev/image">image</a>.<a href="https://pkg.go.dev/image#Image">Image</a>, magic string) <a href="https://pkg.go.dev/builtin#error">error</a>
</pre>
imgをwにmagicで指定の形式("P1"~"P6")で出力します。
