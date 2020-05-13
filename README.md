# Goone

Goone is a packing go code utility tool.
Mainly, This help you to create snippet for competitive programming.
This packs the target source code with depended other code share the same package.

![豪腕(Go-One)](https://1.bp.blogspot.com/-E9fHMc86NSQ/V5jKfFeCY-I/AAAAAAAA81g/_rX1b5zSbkkSkR94dR2-cNK07lbJLxGcACLcB/s800/sports_doping_medal.png)


## Install

```sh
go get github.com/masu-mi/goone
```

## Usage

```sh
Usage:
  goone [command]

Available Commands:
  gen         Generate Packed source files
  help        Help about any command
  pack        Pack target source code file with depended files share the package name into single file

Flags:
  -h, --help   help for goone

Use "goone [command] --help" for more information about a command.
```

## Examples

```sh
goone pack ./main.go --package main -o ./example_out.go
goone gen ./         --package main -o ./generated -p snip-compe-
```


## License
See LICENSE.
