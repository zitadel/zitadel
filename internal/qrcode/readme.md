# QR Code to SVG

This package is a copy of https://github.com/aaronarduino/goqrsvg with the difference of creating the svg with `class` attribute instead of inline `style`:

```go
s.Rect(currX, currY, qs.blockSize, qs.blockSize, "class=\"color\"")
```

and not
```go
s.Rect(currX, currY, qs.blockSize, qs.blockSize, "fill:black;stroke:none")
```

This allows the svg to be styled by css more easily and does not compromise Content Security Policy (CSP).

