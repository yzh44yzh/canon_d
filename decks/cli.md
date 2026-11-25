# Unix CLI Tools

## Bourne Shell
Developed at `Bell Labs` for early versions of Unix.
The basis for all other shells.
bash `Bourne Again Shell` is the default on most Linux distributions.

Commands can be shell builtins or external programs.

## EOF
`End of File` -- there is no more data to read from the file.
`Ctrl-D` signals EOF when reading from standard input.

## cat & echo
cat `concatenate` takes a list of files and writes their contents to stdout. If no files are specified, it reads from stdin and echoes the input to stdout.

echo writes its arguments to stdout. Useful for seeing how patterns and variables expand.

## symbolic link
symlink, soft link
`ln -s target linkname`

## less
`space` -- forward
`b` -- backward
`/` -- search forward
`?` -- search backward
`n` -- next search result

## tree
tree options dir-list

`-d` list dirs only
`-L 3` level deps
`-h` print human readable size

## Output redirection
`cmd > output` redirect stdout to output
`cmd > output 2>&1` redirect stdout to output, then send stderr to the same destination
`cmd 2> output` redirect stderr to output
`cmd &> output` redirect both stdout and stderr to output

## Brace expansion in Bash
mv some_long_file_name.txt some_long_file_name.md
`mv some_long_file_name.{txt,md}`

mv path/before/to/file path/after/to/file
`mv path/{before,after}/to/file`

rm a.txt b.txt c.txt d.txt e.txt f.txt
`rm {a..f}.txt`

## ps
`ps x` -- show all processes for the current user
`ps ax` -- show all processes for all users
`ps u` -- show detailed process information
`ps w` -- use wider output to show more of the command line
`ps pid` -- show information about a specific process

## grep 

### what it means
`globally search a regular expression and print`

### usage
Usage: `grep [OPTION] PATTERN [FILES]`
Search for PATTERN in each FILE.

### -E -F flags
-E, `--extended-regexp` PATTERNS are extended regular expressions
-F, `--fixed-strings` PATTERNS are strings

`egrep` -- equivalent to grep -E
`fgrep` -- equivalent to grep -F

### -i -v -n flags
-i, `--ignore-case`
-v, `--invert-match` select non-matching lines
-n, `--line-number` print line number with output lines

### -l -c flags
-l, `--files-with-matches` print only names of FILEs with selected lines
-c, `--count` print only a count of selected lines per FILE

### show context
`grep -C 5 pattern files`
show 5 context lines near matching line

-A `--after-context`
-B `--before-context`
-C `--context` before and after

## sudo
`sudo cmd` run command 
`sudo -u username cmd` run command as some user
`sudo -i` root shell 

## scp
Secure Copy Command
`scp options user1@host1:file user2@host2:file`

## tar
type archive
create: `tar czf name.tar.gz dir`
extract: `tar xzf name.tar.gz`
list: `tar tzf name.tar.gz`

-`c` create
-`x` extract
-`z` use gzip
-`f` file (archive name)
-`v` verbose (show file names)
-`p` permissions (preserve file permissions)

## changing file permissions
?
chmod g+r file
chmod o+r file
chmod go+r file
chmod u-r file
chmod 644 file
chmod 755 dir

4 `r--` `read only`
5 `r-x` `read and execute`
6 `rw-` `read and write`
7 `rwx` `read, write, and execute`

## curl
curl -v 'http://localhost:4000/api/something' \
`-H 'Content-Type: application/json'` \
`-H 'Accept: application/json'` \
`-d '{"answer": 42}'`





