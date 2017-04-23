# minify - a caddy plugin

[![community](https://img.shields.io/badge/community-forum-ff69b4.svg?style=flat-square)](https://forum.caddyserver.com)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/hacdias/caddy-minify)

This package is a plugin for [Caddy](https://caddyserver.com) webserver that implements a minifier that is able to compress CSS, HTML, JS, JSON, SVG and XML using [github.com/tdewolff/minify](https://github.com/tdewolff/minify).

# Syntax

```
minify paths...  {
    if          a cond b
    if_op       [and|or]
    disable     [js|css|html|json|svg|xml]
    minifier    option value
}
```

+ **paths** are space separated file paths to minify. If nothing is specified, the whole website will be minified.
+ **if** specifies a condition. Multiple ifs are AND-ed together by default. **a** and **b** are any string and may use [request placeholders](https://caddyserver.com/docs/placeholders). **cond** is the condition, with possible values explained in [rewrite](https://caddyserver.com/docs/rewrite#if) (which also has an `if` statement).
+ **if_op** specifies how the ifs are evaluated; the default is `and`.
+ **disable** is used to indicate which minifiers to disable. By default, they're all activated.
+ **minifier** sets **value** for **option** on that minifier. When the option is true or false, its omission is trated as `true`. The possible options are described bellow.

```
html    keep_default_attr_vals  [true|false]
html    keep_document_tags      [true|false]
html    keep_end_tags           [true|false]
html    keep_whitespace         [true|false]

xml     keep_whitespace         [true|false]
css     decimals integer
svg     decimals integer
```

For more information about what does each option and how each minifier work, read the [documentation of tdewolff/minify](https://github.com/tdewolff/minify/blob/master/README.md).

### Examples

Minify all of the supported files of the website:

```
minify
```

Only minify the contents of `/assets` folder:

```
minify /assets
```

Minify the whole website except `/api`:

```
minify  {
    if {path} not_match ^(\/api).*
}
```

Minify the files of `/assets` folder except `/assets/js`:

```
minify /assets {
    if {path} not_match ^(\/assets\/js).*
}
```
