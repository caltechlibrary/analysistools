%phrasecheck(1) user manual | version 0.0.0 6a0f789
% R. S. Doiel
% 2025-11-14

# NAME

phrasecheck

# SYNOPSIS

phrasecheck [OPTIONS] ACTION [PARAMS]

# DESCRIPTION

phrasecheck performs actions based on phrase checking.

# OPTIONS

-help
: Display help

-version
: Display version

-license
: Display license

# ACTION

check [OPTION] PATTERN_FILE FILE_TO_CHECK [FILE_TO_CHECK ...]
: Report the matches found based on PATTERN_FILE in the FILE_TO_CHECK

filetypes PATH
: Walk the PATH directory and report file type counts based on file extension

prune PATTERN_FILE PATH_TO_CHECK [PATH_TO_CHECK ...]
: Walk the paths provided and remove files that match what is in the PATTERN_FILE.

Check a file(s) against the contents of the PATTERN_FILE. Report matches for line numbers and phrases matched

# EXAMPLE

~~~shell
phrasecheck check patterns.txt email.txt
~~~


