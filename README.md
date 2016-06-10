# minify - a caddy plugin

[![community](https://img.shields.io/badge/community-forum-ff69b4.svg?style=flat-square)](https://forum.caddyserver.com)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/hacdias/caddy-minify)

This package is a plugin for [Caddy](https://caddyserver.com) webserver that implements a minifier that is able to compress CSS, HTML, JS, JSON, SVG and XML.

#Syntax

```
minify  {
	only foo...
	exclude bar...
}
```

+ **foo** (optional) are space separated single file paths or folders to include on minifying. By default the whole website will be minified. If this directive is set, only the files on the specified paths will be minified.
+ **bar** (optional) are space separated single file paths or folders to exclude from minifying.
