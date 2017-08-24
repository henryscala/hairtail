Hairtail is a simple document preparation system mimicking [halibut](https://www.chiark.greenend.org.uk/~sgtatham/halibut/). I read the document of halibut, but I did not read the code. I'm not good at sorting out codes written by other people, yet. 

The input format is like halibut. The only supported output format is html, though more output formats may be supported in future.

# implementation philosophy(or limitation)
Efferency is not the first important thing. There may be several passes while handling the input. For now, it is difficult for me to write a one-pass parser.

# input grammars
The formal grammars are able to be find in the repository with the name [Hairtail.g4](https://github.com/henryscala/hairtail/blob/master/Hairtail.g4). It is in [antlr](https://www.antlr.org/
) format, but I do not guarantee it is able to be fed directly to antlr. It is only to be referred to, and hairtail does not use antlr. 

You may also refer to the `test` folder in the repository to see examples. 

The keywords defined by hairtail all begin with `\`. 

## meta chars 
There are 4 meta characters used by hairtail. They are `\{}#`. If they need to be shown in plain text, then they need to be escaped in the way like `\\` and `\{\}\#`. One way to avoid escaping is to use raw text. 

## raw text 
It supports [rust](https://www.rust-lang.org) like raw string using `\r##{}##` alike syntax. The number of `#` is [0-n], where n is to make sure that the text within the `{}` don't need to escape. 

## comments 

`\--` is used to add comments to the document. It will not show in the compiled document. It is able to comment out other grammar elements, too. Counterpart of html is `<!-- -->`.  

## headings 
`\h1 \h2 \h3 \h4 \h5 \h6` defines levels of sections. `\h` is the same as `\h1`. 

Sections may contain ID. ID is to be referenced by other text using `\k`. The text following `\h*` before `\n` will be the title of the section. E.g. 

```
\h{intro-hairtail} Introduction of Hairtail 
```

The headings will be shown in index. 

## inline format 
`\e` means emphasis. Counterpart of html is `<em></em>`.

`\s` means strong. Counterpart of html is `<strong></strong>`.

`\w` is for hyper links. Counterpart of html is `<a></a>`. E.g. 

```
\w{https://github.com/henryscala}{Henryscala}.
```

`\c` is for inline code. Counterpart of html is `<code></code>`. 

`\t` is for inline math. It requires [mathjax](https://www.mathjax.org/) to function. 

`\a` defines an anchor/mark inside the document, which is able to be referred to. 

`\k` is to refer to elements defined by anchors, sections, tables, etc, everything with IDs.  

`\image` is to add picture to the document. Counterpart of html is `<img/>`. 

`\caption` is to add caption/title to blocks that are with ID but without caption. Blocks(images, tables, code, etc) with caption will be shown in specific indices(image index, table index, code index, etc). 
	
## blocks 
The blocks have ID fields so that they are able to be referred to by `\k`.

It is also able to add caption to the blocks. Once blocks have captions, they are able to be shown in specific indices. 

### math
While `\t` is for inline math, `\tex` is for block math. It requires [mathjax](https://www.mathjax.org/) to function. 

### code 
While `\c` is for inline code, `\code` is for block code. 

### order list and bullet list 
`\ol` counterpart of html is `<ol></ol>`. 

`\ul` Counterpart of html is `<ul></ul>`. 

List items inside `\ol` and `\ul` are defined using `\-`. E.g. 

```
\caption{ul-ID}{example of bullet list}
\ul{ul-ID}{ 
	\- asia 
	\- europe 
	\- america 
	\- africa 
}
```

### table 
`\table` counterpart of html is `<table></table>`.  

Rows inside table are naturally separatd by line ending. Collumns inside table are separated using `\d`. E.g. 

```	
\caption {table-id}{optional caption for table}
\table{table-id}{ 
	name \d score
	henry \d 100
	 
	dily \d 99 
	ruby \d 98 
}	
```

## meta data of the document 
Meta data of the document is able to be specified. And the meta data will be shown in place. 
Each meta data occupies one line. The keywords are separated by `,` herein. 
The meta data `\title \sub-title \author \create-date \modify-date \keywords` are supported. 
E.g. 

```
\title meta data 
\sub-title does meta data works? 

\author somebody 

\create-date 2012-07-07 

\modify-date 2012-07-08

\keywords meta, keyword, author
```
	

## indices 	 
The below indices are supported. Except `\toc`, caption is mandatory for blocks to be shown in respective indices.  

- `\toc` index for sections 
- `\image-index` index for images/figures
- `\table-index` index for tables 
- `\order-list-index` index for order list 
- `\bullet-list-index` index for bullet/unordered list
- `\code-index` index for code block 
- `\math-index` index for math block 

E.g. 

```
main index
\toc

image index 
\image-index

table index 
\table-index

order list index 
\order-list-index

bullet list index 
\bullet-list-index

code block index  
\code-index

\h{h1-1} h1 example 
para graph level 1

para graph level 1

para graph level 1

below is a ordered list. 
\caption{first-ol}{optional caption for first-ol}
\ol{first-ol}{
	\- cat
	\- elephant
	\- dog 
}

below is a bullet list. 
\caption{first-ul}{optional caption for first-ul}
\ul{first-ul}{ 
	\- asia 
	\- europe 
	\- america 
	\- africa 
}

\h2{h2-1} h2 example 
para graph level 2 

para graph level 2 

\caption{image-example1}{example1 of a image from w3c}
\image{image-example1}{http://www.w3school.com.cn/i/eg_tulip.jpg}

para graph level 2 

\caption{code-list-1}{code list example 1}	
this is block code \code{code-list-1}{
		println 1
		println 2 
    }
	
\caption {table-id1}{optional caption for table 1}
\table{table-id1}{ 
	name \d score
	henry \d 100
	 
	dily \d 99 
	ruby \d 98 
}	
	
\h2{h2-2} h2 example 
para graph level 2 

\caption{image-example2}{example2  of a image from w3c}
\image{image-example2}{http://www.w3school.com.cn/i/eg_tulip.jpg}

para graph level 2 

para graph level 2 

\h{h1-2} h1 example 2 

below is a ordered list. 
\caption{first-ol-2}{optional caption for first-ol2}
\ol{first-ol-2}{
	\- cat
	\- elephant
	\- dog 
}

below is a bullet list. 
\caption{first-ul-2}{optional caption for first-ul2}
\ul{first-ul-2}{ 
	\- asia 
	\- europe 
	\- america 
	\- africa 
}

para graph level 1

para graph level 1

para graph level 1

\caption {table-id2}{optional caption for table 2}
\table{table-id2}{ 
	name \d score
	henry \d 100
	 
	dily \d 99 
	ruby \d 98 
}	

\h2{h2-2-2} h2 example 2-2 
para graph level 2 

para graph level 2 

\caption{code-list-2}{code list example 1}	
this is block code 
\code{code-list-2}{
		println 1
		println 2 
}
```

## include other document 
`\include` is used to import other documents to the document, like the `#include` of C language.

The path of the included document is relative to the current file to be compiled or absolute path.		

The support of 	`\include` is quite limited. 
	
Note: If the included file itself contains include keyword, it is still relative to the current file to be compiled(not the included file).

Note: the implementation of include keyword has limitations. It is better the included content does not rely on chunks in other files. Otherwise, surprise may happens.

# TODO

[x] Generate Table

[x] command line argument to specify a html template to put renderred content in

[x] Generate List that may be nested

[x] generate section index 

[x] generate figure index, order-list-index, bullet-list-index, table index, block-code index, math index  

[x] generate Title, SubTitle, CreateDate, ModifyDate, Keywords, author 

[x] numbering table, block-code, list, figure 

[x] inline tex 

[x] block tex 

[] provide a default template, so that the output of all syntax elements(tables, code, etc) looks good 

[] python

[] make some hard coded value configurable, e.g. prefix of image, table, etc. 

[] make the output templates configurable. In this way, we may support multi output formats. 

[x] add `comment` keywords. content in it will not be shown, but in `<!-- -->` block

[x] add `include` keywords. It is to import contents from other documnts to this one like `#include` of C language. (Implementated with some limitation)

[] For ill-formed document, output user friendly error messages  

[] Handle blank char. It should not be so strict. Blanks before or after some keyword or meta chars shall be ignored.

[] Get rid of gDoc. Each time a dedicated doc should be generated per input file 