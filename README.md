# Graffiti

## Starting Point

Graffiti is a Go library to ease pretty print to standard streams.

It exists to solve a very common and annoying problem: coloring and formatting the output of programs. This is a hard task, as it envolves dealing with ANSI escape sequences that are hard to debug, as they are invisible when printed.

Not just that, but if your program is piped to other program or file, all the ANSI sequences will remain unless you deal with them, which will make the output be hard to parse and understand.

To solve it, Graffiti brings up new fmt.Printf-like functions that can understand all format specifiers fmt.Printf can as well new format specifiers to apply styles and print strings into standard streams.

Those functions will automatically parse, add or remove ANSI sequences as well their own format specifiers based if the stream is being piped or not, so you can apply styles easily.

This library is not cross-platform: styles will only be applied in UNIX-like operating systems, while on Windows, they will be removed.

## Installation

It is not available for production yet.

## Usage

### Functions

Graffiti offers some functions for you do to your work:

* `graffiti.Print`: prints to `stdout`.
* `graffiti.Println`: prints to `stdout` and appends a new line character in the end.
* `graffiti.Eprintln`: prints to `stderr`.
* `graffiti.Eprintln`: prints to `stderr` and appends a new line character in the end.

Those functions are wrapper of the `fmt.Sprintf` function, which means that you can them to format data just like you normally do with the `fmt.Printf`. The difference is that they can interpret new format specifiers to apply styles.

They will automatically replace those format specifiers with styles sequences, or will remove them automatically of the string if it detects that the stream is not a terminal.

### Formatter Specifiers

* `@F{<color>}`: changes the foreground color.
* `@K{<color>}`: changes the background color.
* `@B`: uses bold. Only visible if font contains bold characters.
* `@I`: uses italic. Only visible if font contains italic characters.
* `@U`: uses underline.
* `@r`: removes all styles applied.

The `<color>` placeholder must be replaced by the value of a color of the 8 bits palette (values from 0 to 255 - full palette can be found online) or the name of a color of the 3 bits palette:

* `red`: same as value `1`.
* `green`: same as value `2`.
* `yellow`: same as value `3`.
* `blue`: same as value `4`.
* `magenta`: same as value `5`.
* `cyan`: same as value `6`.
* `white`: same as value `7`.

If the value used is invalid, no style will be applied for that format specifier.

You do not need to reset your styles in the end of the string, as Graffiti automatically does it for you if detects that you have used a style.

Old terminal emulators, have limited capabilities when rendering fonts and colors. If you want your program to support them, avoid using bold and italic, and prefer to use only colors of the 3 bits palette.

Some terminal emulators, for instance, st, konsole and linux (the default virtual console of Linux), might render bold or italic with a brighter color.

### Examples

Let's create a simple program to test Graffiti's capabilities.

```go
// File: main.go

package main

import (
	"github.com/skippyr/graffiti"
)

func main() {
	// Prints a colorful "Hello World!".
	graffiti.Println("@F{yellow}Hello @F{green}world@F{magenta}!")

	// Formats and prints a colorful error message.
	errorMsg := "No Such File Or Directory"
	errorOsCode := 2
	graffiti.Eprintln(
		"@F{red}@BError:@r @F{yellow}%s @r(os code @F{red}%d@r).",
		errorMsg,
		errorOsCode,
	)
}
```

Now, let's run this program and see its output:

```bash
go run main.go
```

![](images/preview.png)

To see if the sequences will be removed, let's check out what will be put in a file if the output of the program is redirected:

```bash
go run main.go &>output.txt; cat output.txt
```

![](images/preview_pipeline.png)

## Issues

Report issues through the [issues tab](https://github.com/skippyr/graffiti/issues).

## Contributions

If you want to contribute to this project, check out its [contributions guidelines](https://skippyr.github.io/materials/pages/contributions_guidelines.html).

## License

This project is released under the terms of the MIT license.

Copyright (c) 2023, Sherman Rofeman. MIT License.

