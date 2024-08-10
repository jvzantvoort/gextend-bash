[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)

# gextend-bash

Go based extra tools for shell scripting.

## Installation

Either download the release and unpack it some place. Or:

```sh
curl -Ls https://raw.githubusercontent.com/jvzantvoort/tmux-project/main/tmux-project-update | bash
```

## Update

Updating can be done with:

```sh
tmux-project-update
```

## Contents


### cprint

Print a colored status.

Usages:
```
cprint [red|oke|magenta|green|nok|workspace|blue|cyan|white|yellow|debug|black|profile|platform|warn|ok] message
```

### logging

WIP

### mvx

Move a file with avoiding duplicates or overwriting different files:

Example:

Generate a random file:
```
# tr -dc A-Za-z0-9 </dev/urandom | head -c 13 > tempfile
```

Copy the file to a backup:
```
# cp tempfile tempfile.dup
```

Move the original to a ``bck/`` folder:
```
# mvx tempfile bck/
              target found /home/jvzantvoort/tmp/tempfile                                [  SUCCESS   ]
                target found /home/jvzantvoort/tmp/bck/tempfile                          [  FAILURE   ]
          move tempfile                                                                  [  SUCCESS   ]
```

*Note*: the error points to no dups found.

Copy the duplicate to the original and repeat, effectively moving the same file twice:

```
# mv tempfile.dup tempfile
# mvx tempfile bck/
              target found /home/jvzantvoort/tmp/tempfile                                [  SUCCESS   ]
                target found /home/jvzantvoort/tmp/bck/tempfile                          [  SUCCESS   ]
          move tempfile                                                                  [  SUCCESS   ]
```

Create a new tempfile and move it:

```
# tr -dc A-Za-z0-9 </dev/urandom | head -c 13 > tempfile
# mvx tempfile bck/
              target found /home/jvzantvoort/tmp/tempfile                                [  SUCCESS   ]
                target found /home/jvzantvoort/tmp/bck/tempfile                          [  SUCCESS   ]
                target found /home/jvzantvoort/tmp/bck/tempfile.1                        [  FAILURE   ]
          move tempfile                                                                  [  SUCCESS   ]
```


The result is the original moved and the second version moved to a ``.1`` file:
```
# ls -l bck
total 8
-rw-r--r--. 1 jvzantvoort jvzantvoort 13 Aug 10 21:03 tempfile
-rw-r--r--. 1 jvzantvoort jvzantvoort 13 Aug 10 21:04 tempfile.1


### path_clean

Returns a cleaned up version of the path

### print_status

Alias for cprint

### today

Returns "vYYWW.<dow>" as a timestamp
