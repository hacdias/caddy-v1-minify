# minify - a caddy plugin

[![community](https://img.shields.io/badge/community-forum-ff69b4.svg?style=flat-square)](https://forum.caddyserver.com)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/hacdias/caddy-minify)

This package is a plugin for [Caddy](https://caddyserver.com) webserver that implements a minifier that is able to compress CSS, HTML, JS, JSON, SVG and XML using [github.com/tdewolff/minify](https://github.com/tdewolff/minify).

#Syntax

```
minify paths...  {
	if    	a cond b
   	if_op 	[and|or]
}
```

+ **paths** are space separated file paths to minify. If nothing is specified, the whole website will be minified.
+ **if** specifies a condition. Multiple ifs are AND-ed together by default. **a** and **b** are any string and may use [request placeholders](https://caddyserver.com/docs/placeholders). **cond** is the condition, with possible values explained in [rewrite](https://caddyserver.com/docs/rewrite#if) (which also has an `if` statement).
+ **if_op** specifies how the ifs are evaluated; the default is `and`.
