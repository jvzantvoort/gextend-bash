[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)

# gextend-bash

Go based extra tools for shell scripting.

<!-- toc -->

- [gextend-bash](#gextend-bash)
  - [Installation](#installation)
  - [Update](#update)
  - [Releases](#releases)
  - [Contents](#contents)

<!-- tocstop -->

## Installation

Download a release archive and unpack it somewhere in your `PATH`, or run
the installer, which picks the archive for your platform and verifies its
checksum:

```sh
curl -fsSL https://raw.githubusercontent.com/jvzantvoort/gextend-bash/main/gextend-bash-update | bash
```

It installs into `/usr/local/bin` when run as root, otherwise into
`$(go env GOBIN)` or `~/bin`.

With a Go toolchain you can also build and install from source:

```sh
go install github.com/jvzantvoort/gextend-bash/cmd/...@latest
```

## Update

```sh
gextend-bash-update
```

Every tool reports its version with `--version`.

## Releases

Merging a pull request that changes the code (Go sources, `go.mod`/`go.sum`
or `gextend-bash-update`) releases a new version automatically. The version
is derived from the [conventional commit](https://www.conventionalcommits.org)
messages since the last tag: `feat:` bumps the minor version, other types
the patch version. Release notes are generated with
[git-cliff](https://git-cliff.org). Merges that only touch documentation or
workflows do not create a release; run the "Tag and release" workflow by
hand to release them anyway.

## Contents

| Tool | Purpose |
| ---- | ------- |
| [achelper](#achelper) | bash completion for ssh, sftp, scp and rsync hosts |
| [center_text](#center_text) | center text in a `#` comment block |
| [confirm](#confirm) | yes/no prompt; the answer is the exit code |
| [countdown](#countdown) | full screen countdown timer |
| [gbanner](#gbanner) | print text inside a box |
| [jinja_comment](#jinja_comment) | wrap text in a jinja comment block |
| [mdtoc](#mdtoc) | table of contents for a markdown file |
| [mvx](#mvx) | move files without overwriting different ones |
| [namegen](#namegen) | random `adjective-noun` name |
| [path_clean](#path_clean) | remove duplicate and missing `PATH` entries |
| [pathmunge](#pathmunge) | add or remove `PATH` entries without duplicates |
| [pprint](#pprint) | colored one-line message with a symbol |
| [print_status](#print_status) | status line on stderr |
| [retry](#retry) | rerun a command until it succeeds |
| [run_status](#run_status) | run a command and print its status line |
| [test_result](#test_result) | status line for an exit code |
| [today](#today) | ISO week based timestamp |

### achelper

Generates bash completion for host names, taken from `/etc/hosts` and
`~/.ssh/config` (following `Include` lines). `--known-hosts` also adds the
unhashed entries of `~/.ssh/known_hosts`. `scp` and `rsync` fall back to
file name completion when no host matches.

```sh
eval "$(achelper ssh)"
eval "$(achelper ssh --known-hosts)"
```

### center_text

Centers text (arguments or stdin) in an 80 column `#` comment block,
wrapping long text.

```console
$ center_text "Section title"
#                                Section title                                 #
```

### confirm

Asks a yes/no question on stderr and exits 0 for yes, 1 for no. `-d`
sets the answer for an empty reply, end of input or a timeout (`-t`);
without it those count as no.

```sh
if confirm -d no -t 30s "Delete the old backups?"; then
    rm -rf backups/old
fi
```

### countdown

Full screen countdown; takes seconds or a duration like `90s`, `5m` or
`1h30m`. Press `q` to quit. Rings the terminal bell when done, unless
`--no-bell` is given.

```sh
countdown 5m
```

![countdown](docs/countdown/example.gif)

### gbanner

Print text inside a box, optionally colored.

```sh
gbanner [-c|--color name] [-x|--hpad n] [-y|--vpad n] [-l|--left] [text...]
echo "some text" | gbanner --color=rainbow
gbanner --hpad=4 --vpad=1 "Padded banner"
gbanner --left "Left-aligned text"
```

Available colors: black, red, green, yellow, blue, magenta, cyan, white,
their `hi<color>` bright variants, and `rainbow` which alternates the
color of every character.

`--hpad`/`-x` controls the padding (columns) kept between the terminal
width and the box; `--vpad`/`-y` controls the padding (blank lines) kept
between the box border and the text; `--left`/`-l` left-aligns the text
instead of centering it.

### jinja_comment

Wraps text (arguments or stdin) in an 80 column jinja comment block.

```console
$ jinja_comment "This template is managed by Ansible"
{#- This template is managed by Ansible                                      -#}
```

### mdtoc

Generates a table of contents for a markdown file. Headings in fenced code
blocks are ignored and anchors match the ones GitHub generates.

```sh
mdtoc -file README.md -depth 3     # print the TOC
mdtoc -file README.md -w           # rewrite it in place
```

With `-w` the TOC is written between `<!-- toc -->` and `<!-- tocstop -->`
marker lines, as in this README.

### mvx

Moves files into a directory without overwriting different files: when the
target name is taken by a file with other content, the file is stored as
`name.1.ext`, `name.2.ext`, ... When the target has the same content, the
source is removed. Moves across filesystems are done by copying.

```console
$ echo a > tempfile; mvx tempfile bck/
          move tempfile                                                [  SUCCESS   ]
$ echo a > tempfile; mvx tempfile bck/    # identical: source removed
          move tempfile                                                [  SUCCESS   ]
$ echo b > tempfile; mvx tempfile bck/    # different: stored as tempfile.1
          move tempfile                                                [  SUCCESS   ]
$ ls bck
tempfile  tempfile.1
```

`--debug` shows what it checks. It exits non-zero when any move fails.

### namegen

Prints a random `adjective-noun` name, e.g. `bitter-dream`.

### path_clean

Prints `PATH` (or the path given as argument) without duplicates, missing
directories and relative entries.

```sh
PATH="$(path_clean)"
```

### pathmunge

Prints `PATH` with directories prepended (default), appended (`-a`) or
removed (`-r`). Directories already present stay where they are. `-v`
works on another variable, such as `MANPATH`.

```sh
PATH="$(pathmunge /opt/tool/bin)"
PATH="$(pathmunge -a "$HOME/bin")"
MANPATH="$(pathmunge -v MANPATH -r /opt/old/man)"
```

### pprint

Print a colored message with a leading symbol.

```sh
pprint [black|blue|cyan|debug|err|error|green|magenta|nok|ok|oke|platform|profile|red|warn|white|workspace|yellow] message
```

### print_status

Print a formatted status line to stderr, with a colored, aligned status
label. Aliases are normalized to their canonical name (`ok`/`oke` →
`SUCCESS`, `nok`/`fail`/`failed` → `FAILURE`, `info` → `NOTICE`, `warn` →
`WARNING`, `undefined` → `UNKNOWN`).

```sh
print_status [success|failure|warning|notice|ok|oke|fail|nok|info|warn|undefined] message...
```

```console
$ print_status success "deployment finished"
deployment finished..........................................................[ SUCCESS ]
```

Shell completion for it is available through `print_status completion bash`
(or `zsh`, `fish`, `powershell`).

### retry

Reruns a command until it succeeds, `-n` times at most (default 3), with a
delay of `-d` between attempts (default 1s). `-b` doubles the delay after
every failure, up to `-max-delay`. Commands that cannot be run at all are
not retried. Exits with the command's last exit code.

```sh
retry -n 5 -d 2s -b -- curl -fsS https://example.com/health
```

### run_status

Runs a command and prints its status line, replacing the
`command; test_result $? "message"` pattern. Exits with the command's exit
code. `-q` hides the command's output unless it fails.

```console
$ run_status -m "fetch sources" git fetch
fetch sources................................................................... [ SUCCESS ]
$ run_status -q make test
```

### test_result

Prints the status line for an exit code.

```sh
command
test_result "$?" "command has executed successfully"
```

### today

Prints a timestamp based on the ISO week: `wk<yy><ww>.<d>`, where `<d>` is
the day of the week (Sunday is 0), e.g. `wk2642.4`.
