<h1>Graffiti</h1>
	<h2>Starting Point</h2>
		<p>Graffiti is a Go library to ease pretty print to standard streams.</p>
		<p>It exists to solve a very common and annoying problem: coloring and formatting the output of programs. This is a hard task, as it envolves dealing with ANSI escape sequences that are hard to debug, as they are invisible when printed.</p>
		<p>Not just that, but if your program is piped to other program or file, all the ANSI sequences will remain unless you deal with them, which will make the output be hard to parse and understand.</p>
		<p>To solve it, Graffiti brings up new <code>fmt.Printf</code>-like functions that can understand all format specifiers <code>fmt.Printf</code> can as well new format specifiers to apply styles and print strings into standard streams.</p>
		<p>Those functions will automatically parse, add or remove ANSI sequences as well their own format specifiers based if the stream is being piped or not, so you can apply styles easily.</p>
	<h2>Installation</h2>
		<p>Not yet available for production.</p>
	<h2>Issues</h2>
		<p>Report issues through the <a href="https://github.com/skippyr/issues">issues tab</a>.</p>
	<h2>Contributions</h2>
		<p>If you want to contribute to this project, check out its <a href="https://skippyr.github.io/materials/pages/contributions_guidelines.html">contributions guidelines</a>.</p>
	<h2>License</h2>
		<p>This project is released under the terms of the MIT license.</p>
		<p>Copyright (c) 2023, Sherman Rofeman. MIT License.</p>
