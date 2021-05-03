lmt.1
====
## NAME
lmt - Literate Markdown Tangle -- convert .md files to code
## SYNOPSIS
lmt file.md [file.md] ...
## DESCRIPTION
lmt tangles a group of .md files, extracting the named code
sections and producing compilable files from them.

Files need to contain at least one output filename to
write to and alt least one code section.

A code section begins with ```
`<200b>```language filename +=`
and ends with `<200b>``` `

A code section is inserted in the output file by enclosing its 
name in <<< and >>>

That looks like, for example
```cpp hello.cpp +=
<<<copyright>>>
<<<includes>>>

int main() {
    <<<body of main>>>
}
```
where "body of main" contains

```cpp "body of main"
std::cout << "Hello, werld!" << std::endl;
```

## BOOTSTRAPPING
If you wish to bootstrap from and existing program you need at least 
one file containing a lmt section with a filename to write to,
and probably a second containing some lmt blocks

I used README.md, and inserted
```go main.go += 
<<<body of main>>>
```
in order to get a main.go program that contained anything in the <<<mainline>>> block.

I called the other file Mainline.md, and stuck this in before the first line of go
```go "body of main" +=
and then closed the file off with ```.

When I ran lmt README.md Mainline.md, I got a main.go file containing
the contents of Mainline.md


### BUGS

The <<< and >>> must be right beside the name, no white space is allowed. If you 
write <<< body of maine >>>, you will get 
    expected statement, found '<<'
White space can occur inside the name of the section, though

A newline must appaera at the end of the ``` closing lines or you get 
    expected statement, found '<<'

Selecting the output file to be main.go seems to only work in the README.md file
We end up with an empoty file and get
main.go:1:2: expected 'package', found 'EOF'

running them through lmt in the oposite order fails identically

doesn't see that go "fred+= is an error




