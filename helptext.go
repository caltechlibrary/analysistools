package analysistools

const (
	HelpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] ACTION [PARAMS]

# DESCRIPTION

{app_name} performs actions based on phrase checking.

# OPTIONS

-help
: Display help

-version
: Display version

-license
: Display license

# ACTION

check [OPTION] PATTERN_FILE FILE_TO_CHECK [FILE_TO_CHECK ...]
: Report the matches found based on PATTERN_FILE in the FILE_TO_CHECK.
Use '{app_name} check help' to list available options for check.

check-directory [OPTION] PATTERN_FILE PATH [EXCLUDE_LIST_FILENAME]
: Walk the directory indicated by PATH. Check any text files against the
PATTERN_FILE contents. Report matches.
Use '{app_name} check-directory help' to list available options for check.

mimetypes PATH [EXCLUDE_LIST_FILENAME]
: Walk the PATH directory and report file and it's likely mime type

filetypes PATH [EXCLUDE_LIST_FILENAME]
: Walk the PATH directory and aggregate counts by file extension and mime type

tokens FILENAME [FILENAME ...]
: tokenize a file and display the tokens in CSV format (name, token, word number, line number)

PATTERN_FILE
: This holds a list of patterns to match against, one pattern statement per line.

EXCLUDE_LIST_FILENAME
: This is a file contains a list (one entry per line) of path elements to be excluded from the walk.

# EXAMPLE

~~~shell
{app_name} check patterns.txt email.txt
~~~

`

MimeTypeHelp = `%{app_name}-mimetypes(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} mimetypes

# SYNOPSIS

{app_name} mimetypes [OPTIONS] PATH [EXCLUDE_LIST_FILENAME]

# DESCRIPTION

**{app_name} mimetypes** talks a directory and returns a list of files along with
with the common mime types in CSV format. 

# OPTIONS

-h, -help, help
: display this help page

`

TokensHelp = `%{app_name}-tokens(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} tokens

# SYNOPSIS

{app_name} tokens [OPTIONS] FILENAME [FILENAME ...]

# DESCRIPTION

**{app_name} tokens** parses a text file and turns it into a CSV list of
tokens.

# OPTIONS

-h, -help, help
: display this help page

`

FileTypeCountsHelp = `%{app_name}-filetypes(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} filetypes

# SYNOPSIS

{app_name} filetypes [OPTIONS] PATH [EXCLUDE_LIST_FILENAME]

# DESCRIPTION

**{app_name} filetypes** walk the PATH and return a cound of files by
file extension. Files which start with a "." will be considered an extension.

Extentions are start with the last period in the path's basename and continue to the
end of the file's basename.

# OPTIONS

-h, -help, help
: display this help page

`

CheckDirectoryHelp = `%{app_name}-check-directory(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} check-directory

# SYNOPSIS

{app_name} check-directory [OPTIONS] PATTERN_FILE PATH [EXCLUDE_LIST_FILENAME]

# DESCRIPTION

Walk the directory indicated by PATH. Check any text files against the
PATTERN_FILE contents. Report matches. The report is returned after
the directory walk is complete. The report is output to standard output
in CSV format.

PATTERN_FILE
: This holds a list of patterns to match against, one pattern statement per line.

EXCLUDE_LIST_FILENAME
: This is a file contains a list (one entry per line) of path elements to be excluded from the walk.

# OPTIONS

-h, -help, help
: display this help page

-match-one, -1
: stop at first match


`

CheckFileHelp = `%{app_name}-check(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} check

# SYNOPSIS

{app_name} check [OPTIONS] PATTERN_FILE FILE_TO_CHECK [FILE_TO_CHECK ...]

# DESCRIPTION

Report the matches found based on PATTERN_FILE in the FILE_TO_CHECK. If more than
one file is included they will be checked conseculatively and included in the CSV
output.

PATTERN_FILE
: This holds a list of patterns to match against, one pattern statement per line.

EXCLUDE_LIST_FILENAME
: This is a file contains a list (one entry per line) of path elements to be excluded from the walk.


# OPTIONS

-h, -help, help
: display this help page

-match-one, -1
: stop at first match


`

)