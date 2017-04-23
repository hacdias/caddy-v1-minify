# minify

[![Build](https://img.shields.io/travis/hacdias/caddy-minify.svg?style=flat-square)](https://travis-ci.org/hacdias/caddy-minify)
[![community](https://img.shields.io/badge/community-forum-ff69b4.svg?style=flat-square)](https://caddy.community)
[![Go Report Card](https://goreportcard.com/badge/github.com/hacdias/caddy-minify?style=flat-square)](https://goreportcard.com/report/hacdias/caddy-minify)

Caddy plugin that implements minification on-the-fly for CSS, HTML, JSON, SVG and XML. It uses [tdewolff's library](https://github.com/tdewolff/minify) so, let's thank him! You can download this plugin with Caddy on its [official download page](https://caddyserver.com/download).

## Syntax

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
+ **disable** is used to indicate which minifiers to disable; by default, they're all activated.
+ **minifier** sets **value** for **option** on that minifier. When the option is true or false, its omission is trated as `true`. The possible options are described bellow.

### Minifiers options

| Minifier(s)   | Option                    | Value         | Description |
| ------------- |-------------              | ----------    | ----------- |
| css, svg      | decimals                  | number        | Preserves default attribute values. |
| xml, html     | keep_whitespace           | true\|false   | Preserve `html`, `head` and `body` tags. |
| html          | keep_end_tags             | true\|false   | Preserves all end tags. |
| html          | keep_document_tags        | true\|false   | Preserves whitespace between inline tags but still collapse multiple whitespace characters into one. |
| html          | keep_default_attr_vals    | true\|false   | Number of decimals to preserve for numbers, `-1` means no trimming. |

For more information about what does each option and how each minifier work, read the [documentation of tdewolff/minify](https://github.com/tdewolff/minify/blob/master/README.md).

## Examples

Minify all of the supported files of the website:

```
minify
```

Only minify the contents of `/assets` folder:

```
minify /assets
```

Only minify css files:

```
minify {
    disable html svg json xml js
}
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
