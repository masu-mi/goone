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
  -c, --config string   config file path (default: ~/.config/goone/config.toml)
  -h, --help            help for goone

Use "goone [command] --help" for more information about a command.
```

## Examples

```sh
goone pack ./main.go --package main -o ./example_out.go
goone gen ./         --package main -o ./generated -p snip-compe-
```

### Switch template

You can set template file path in config.toml like this.

```toml
templatefile = "/Users/masumi/.config/goone/template.go"
```

Template is required to match go's text/template format like this.

```
package {{ .Package }}
// packed from {{ .SrcFiles }} with goone.
// {{"{{_cursor_}}"}}

{{ .Imports }}

{{ .Decls }}
```

Usable attributes are [here](https://github.com/masu-mi/goone/blob/master/model/packed_code.go).


## License
See LICENSE.
