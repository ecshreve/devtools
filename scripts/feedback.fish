#!/usr/bin/env fish

function feedback
  set -g prompt "Please look at the changes to these files and provide general feedback, such as:
    - Are there any typos?
    - Are there any bugs?
    - Are there any improvements that can be made?
    - Are there any other comments you would like to make?"
  
  echo (git diff --cached | mods -f --fanciness 4 --temp .4 $prompt) > feedback.md
end

feedback