# onpage
Takes URLs on stdin and a regex as an argument, returns URLs to pages where a match was found

# usage
   -c int
        Concurrency (default 20)
  -p string
        Regex (The syntax of the regular expressions accepted is the same general syntax used by Perl, Python, and other languages.)
  -proxy string
        Proxy URL (e.g., http://123.321.123.321:8080)
  -r    Match pages that do not contain the supplied pattern

# installation
> go get -u github.com/raverrr/onpage

# use cases
-Find known vulnerable code:
>  cat urls.txt | onpage -p '(VulnerableFunctionName|VulnerableLib.*\\.js)'
  
-Find parameter reflections:
>  cat urls.txt |qsreplace zzzz1 | onpage -p 'zzzz1'
  
-Anything else you might want to check a bunch of pages for:
>  secrets, protection mechanisms or lack of, really anthing you might want to probe a bunch of pages for
  
  
  

