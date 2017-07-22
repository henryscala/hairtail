grammar Hairtail; 

LINE_END : '\r'? '\n' ; 

META_CHAR: '\\\\' | '\\{' | '\\}' | '\\#' ; 

ID : [a-zA-Z_.-][0-9a-zA-Z_.-]* ; 

WS : [ \t] ; 

EMPHASIS : '\\e' ; 

STRONG : '\\s' ; 

HYPER_LINK : '\\w' ; 

IMAGE : '\\image' ; 

RAW : '\\r' ; 

LBRACE : '{' ;

RBRACE : '}' ;

FILLER : '#' ; 

SECTION_LEVEL : [1-6] ; 

SECTION_MARK :  '\\h' ; 

BULLET_LIST :  '\\ul' ; 

ORDER_LIST :  '\\ol' ;

TABLE :  '\\table' ; 

INLINE_TEX :  '\\t' ; 

BLOCK_TEX :  '\\tex' ; 

INLINE_CODE :  '\\c' ; 

BLOCK_CODE :  '\\code' ; 

BLOCK_PYTHON : '\\python' ; 

LIST_ITEM :  '\\-' ; 

PARAGRAPH_DELIM :  WS* LINE_END ;  //blank line 

CELL_DELIM : '\\d' ; 

REFER_TO : '\\k' ; //refer to other keyword 

STRING : .+? ; 

COMMENT : '\\#' ; 

doc : (paragraphs | blocks) sections ;  

line : (inline_block | string)+ (LINE_END | EOF) ; 

paragraph : line+ ; 

paragraphs : paragraph (PARAGRAPH_DELIM paragraph)* ; 

block : image_block | list_block | raw_block  | block_python | block_code |block_tex |table_block  ;  

embraced_id : LBRACE WS* ID WS* RBRACE ;

embraced_block : LBRACE block RBRACE ; 

section_header :  SECTION_MARK SECTION_LEVEL? embraced_block string LINE_END ; 

blocks : block* ; 

section : section_header (paragraphs | blocks); 

sections : section* ; 

emphasis_block :  EMPHASIS embraced_block ; 

strong_block :  STRONG embraced_block ; 

refer_to_block : REFER_TO (LBRACE string RBRACE) ; 

inline_comment_block : COMMENT (LBRACE string RBRACE) ; 

hyper_link_block :  HYPER_LINK (LBRACE string RBRACE) (LBRACE string RBRACE) ; 

image_block :  IMAGE embraced_id (LBRACE string RBRACE) (LBRACE string RBRACE) ; //id, url, caption 

embraced_raw_content : (LBRACE string RBRACE) | (FILLER embraced_raw_content FILLER) ;
 
raw_block :  RAW embraced_raw_content  ;  

inline_block : emphasis_block 
             | strong_block 
             | hyper_link_block 
             | image_block 
             | inline_code 
             | inline_tex
             | refer_to_block 
             | inline_comment 
             ; 

list_block : bullet_list_block | order_list_block ; 

list_item :  LIST_ITEM paragraphs ;

bullet_list_block :  BULLET_LIST embraced_id LBRACE list_item+ RBRACE ;

order_list_block :  ORDER_LIST embraced_id LBRACE list_item+ RBRACE ;

table_row : string ( CELL_DELIM string)* ; //ESC separated values 

table_block :  TABLE embraced_id LBRACE table_row (LINE_END table_row)* RBRACE ; 

inline_tex :  INLINE_TEX raw_block ; 

block_tex :  BLOCK_TEX embraced_id raw_block ; 

inline_code :  INLINE_CODE ((LBRACE string RBRACE) | raw_block) ; 

block_code :  BLOCK_CODE embraced_id ((LBRACE string RBRACE) | raw_block) ; 

block_python :  BLOCK_PYTHON embraced_id ((LBRACE string RBRACE) | raw_block) ; 

string : .+? ; //here . means token not char, ANTLR first tokenize and then parsing 