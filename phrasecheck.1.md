%phrasecheck(1) user manual | version 0.0.2 19fa525
% R. S. Doiel
% 2025-11-20

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

check-file [OPTION] PATTERN_FILE FILE_TO_CHECK [FILE_TO_CHECK ...]
: Report the matches found based on PATTERN_FILE in the FILE_TO_CHECK.

check-directory [OPTION] PATTERN_FILE PATH [EXCLUDE_LIST_FILENAME]
: Walk the directory indicated by PATH. Check any text files against the
PATTERN_FILE contents. Report matches.

filetypes PATH [EXCLUDE_LIST_FILENAME]
: Walk the PATH directory and report file and it's likely mime type

filetype-counts PATH [EXCLUDE_LIST_FILENAME]
: Walk the PATH directory and aggregate counts by file extension and mime type

prune PATTERN_FILE PATH [EXCLUDE_LIST_FILENAME]
: Walk the PATH directory and remove files that have at least one match in the PATTERN_FILE.

PATTERN_FILE
: This holds a list of patterns to match against, one pattern statement per line.

EXCLUDE_LIST_FILENAME
: This is a file contains a list (one entry per line) of path elements to be excluded from the walk.

# EXAMPLE

~~~shell
phrasecheck check patterns.txt email.txt
~~~


