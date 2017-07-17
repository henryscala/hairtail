A simple document preparation system mimicking [halibut](https://www.chiark.greenend.org.uk/~sgtatham/halibut/). 

The input format is like halibut. 

In the first step, the output format is markdown/html.

Basically it contains the functions equivalent with markdown, but with math support. 

It also support rust like raw string using `\r~~{}~~` alike syntax. The number of `~` is [0-n], where n is to make sure that the text within the `{}` don't need to escape.

So the meta character of hairtail include `\{}~`. If they need to be shown in plain text, then they neeed to be escaped like `\\` and `\{\}\~`. 

# implementation philosophy 
Efferency is not the first important thing. There may be several passes while handling the input. 

# input grammars
## headings 
\h1 \h2 \h3 \h4 \h5 \h6 defines levels of sections. \h is the same as \h1. 

\h{keyword}, keyword is used to be referenced by other text using \k. The text following \h{keyword} before \n will be the title of the section. E.g. 

\h{intro-hairtail} Introduction of Hairtail 

## inline format 
\e \s inline format command emphasis and strong 

## raw text 
\r~~~{raw text inside, which may be multiple line}~~~, the number of ~ is [0-n]. The intention is that inside the brace, no escape is required. 

## TODO
Handle blank char. It should not be so strict. blanks before or after some keyword or meta chars shall be ignored.

INline format should be with Paragraph, need to combine them 

