A simple document preparation system mimicking [halibut](https://www.chiark.greenend.org.uk/~sgtatham/halibut/). 

The input format is like halibut. 

The output format only support html now. More output formats may be supported in future.

It also support rust like raw string using `\r~~{}~~` alike syntax. The number of `~` is [0-n], where n is to make sure that the text within the `{}` don't need to escape. Thus, the meta character of hairtail are `\{}~`. If they need to be shown in plain text, then they neeed to be escaped like `\\` and `\{\}\~`.

# implementation philosophy(or limitation)
Efferency is not the first important thing. There may be several passes while handling the input. For now, it is difficult for me to write a one-pass parser.

# input grammars
## headings 
\h1 \h2 \h3 \h4 \h5 \h6 defines levels of sections. \h is the same as \h1. 

\h{keyword}, keyword is used to be referenced by other text using \k. The text following \h{keyword} before \n will be the title of the section. E.g. 

\h{intro-hairtail} Introduction of Hairtail 

## inline format 
\e means emphasis. Counterpart of html is `<em></em>`.

\s means strong. Counterpart of html is `<strong></strong>`

\w is for hyper links. \w{https://github.com/henryscala}{Henryscala}.

## raw text 
\r~~~{raw text inside, which may be multiple line}~~~, the number of ~ is [0-n]. The intention is that inside the brace, no escape is required. 

## TODO
Handle blank char. It should not be so strict. Blanks before or after some keyword or meta chars shall be ignored.

Generate Table

Generate List that may be nested
