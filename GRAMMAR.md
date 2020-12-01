# Grammar

```
expression
    : term
    | expression '+' term
    | expression '-' term
    ;

term
    : factor
    | term '*' factor
    | term '/' factor
    | term '%' factor
    ;

factor
    : primary
    | '-' factor
    | '+' factor
    | factor '^' factor
    ;

primary
    : IDENTIFIER
    | NUMBER
    | '(' expression ')'
    | FUNCTION '(' args ')'
    ;

args
    : expression
    | expression ',' expression
    ;
```