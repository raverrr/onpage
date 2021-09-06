# onpage
Takes URLs on stdin and a regex as an argument, returns URLs to pages where a match was found

# usage
  -c Concurrency (default 50)
  
  -p Regex (The syntax of the regular expressions accepted is the same general syntax used by Perl, Python, and other languages.)
  
  -r Match anything other than the supplied pattern
  
# use cases
  Find known vulnerable code:
  cat urls.txt | onpage -p '(VulnerableFunctionName|VulnerableLib.*\/.js)'
  
  Find parameter reflections:
  cat urls.txt |qsreplace zzzz1 | onpage -p 'zzzz1'
  
  
  
  

