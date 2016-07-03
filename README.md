# Md2Vim #

Writing technical documentation for [Vim](http://www.vim.org/) extensions is not very exciting; having to manually
convert what you've just written in [Markdown](https://daringfireball.net/projects/markdown/) to the
[vimdoc](http://vimdoc.sourceforge.net/htmldoc/usr_toc.html) help format is even less fun. As I grew tired of doing this
over and over for my [ArgWrap](https://foosoft.net/projects/argwrap/) extension, I decided to formally solve this problem for myself and
others, creating the Md2Vim converter tool.

Md2Vim is implemented in Golang as a custom renderer for the [blackfriday](https://github.com/russross/blackfriday)
Markdown processor. Maintainability, ease of use, and beautiful vimdoc output were primary considerations in its design.

## Installation ##

If you already have the Go environment and toolchain set up, you can get the latest version by running:

```
$ go get github.com/FooSoft/md2vim
```

Otherwise, you can use the pre-built binaries for the platforms below:

*   [md2vim_darwin_386.tar.gz](https://foosoft.net/projects/md2vim/dl/md2vim_darwin_386.tar.gz)
*   [md2vim_darwin_amd64.tar.gz](https://foosoft.net/projects/md2vim/dl/md2vim_darwin_amd64.tar.gz)
*   [md2vim_linux_386.tar.gz](https://foosoft.net/projects/md2vim/dl/md2vim_linux_386.tar.gz)
*   [md2vim_linux_amd64.tar.gz](https://foosoft.net/projects/md2vim/dl/md2vim_linux_amd64.tar.gz)
*   [md2vim_linux_arm.tar.gz](https://foosoft.net/projects/md2vim/dl/md2vim_linux_arm.tar.gz)
*   [md2vim_windows_386.tar.gz](https://foosoft.net/projects/md2vim/dl/md2vim_windows_386.tar.gz)
*   [md2vim_windows_amd64.tar.gz](https://foosoft.net/projects/md2vim/dl/md2vim_windows_amd64.tar.gz)

## Usage ##

Executing Md2Vim with the `-help` command line argument will trigger online help to be displayed. The list below
provides a more detailed description of what the parameters do.

*   **cols**

    The number of columns used for laying out vimdoc files to make them look as good as possible with your content.
    Notice that file contents will not be wrapped to this value; this is purely for such things as horizontal rule
    widths and help tag positioning. This defaults to 80, but that's a bit too narrow for some people.

*   **desc**

    Vim help files are supposed to start with the following two fields on the first line:

    ```
    filename.txt Description of this help file's contents
    ```

    The first field is the filename of the generated vimdoc help file; the second is the description can you provide
    with this parameter.

*   **norules**

    By default, we generate horizontal rules above level 1-2 headings, as shown below:

    ```
    ================================================================================
    Level 1 Heading
    --------------------------------------------------------------------------------
    Level 2 Heading
    ```
    If you don't like the way it looks you can turn it off.

*   **notoc**

    If you don't wish to generate a table of contents you should set this flag to opt out. The table of contents lists
    all of the headings in the document and is always inserted before the beginning of document body.

*   **pascal**

    By default, all help tags get converted to lower case and space delimited words are joined with underscores.

    ```
    rigellians-how_to_cook_for_fourty_humans
    ```

    If you prefer the *PascalCase* way of doing things, set this flag and your output will look like this:

    ```
    Rigellians-HowToCookForFourtyHumans
    ```

*   **tabs**

    If you don't like four space tabs for some reason you can change it to something else with this parameter.

## Screenshots ##

[![Markdown source file](https://foosoft.net/projects/md2vim/img/markdown-thumb.png)](https://foosoft.net/projects/md2vim/img/markdown.png)
[![VimDoc destination file](https://foosoft.net/projects/md2vim/img/vimdoc-thumb.png)](https://foosoft.net/projects/md2vim/img/vimdoc.png)

## License ##

MIT
