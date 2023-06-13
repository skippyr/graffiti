<h1>Graffiti</h1>
	<h2>Starting Point</h2>
		<p>Graffiti is a Go library to ease pretty print to standard streams.</p>
		<p>It exists to solve a very common and annoying problem: coloring and formatting the output of programs. This is a hard task, as it envolves dealing with ANSI escape sequences that are hard to debug, as they are invisible when printed.</p>
		<p>Not just that, but if your program is piped to other program or file, all the ANSI sequences will remain unless you deal with them, which will make the output be hard to parse and understand.</p>
		<p>To solve it, Graffiti brings up new <code>fmt.Printf</code>-like functions that can understand all format specifiers <code>fmt.Printf</code> can as well new format specifiers to apply styles and print strings into standard streams.</p>
		<p>Those functions will automatically parse, add or remove ANSI sequences as well their own format specifiers based if the stream is being piped or not, so you can apply styles easily.</p>
		<p>This library is not cross-platform: styles will only be applied in UNIX-like operating systems, while on Windows, they will be removed.</p>
	<h2>Installation</h2>
		<p>Not yet available for production.</p>
	<h2>Usage</h2>
		<h3>Formatters</h3>
			<ul>
				<li><code>@F{&lt;color&gt;}</code>: changes the foreground color.</li>
				<li><code>@K{&lt;color&gt;}</code>: changes the background color.</li>
				<li><code>@B</code>: use bold. Only visible if the font used in the terminal has bold chracters.</li>
				<li><code>@I</code>: use italic. Only visible if the font used in the terminal has italic chracters.</li>
				<li><code>@U</code>: use underline.</li>
				<li><code>@r</code>: remove all styles applied.</li>
			</ul>
			<p>The <code>&lt;color&gt;</code> placeholder must be replaced by the value of a color of the 8 bits palette (values from 0 to 255 - full palette can be found online) or the name of a color of the 3 bits palette:</p>
			<ul>
				<li><code>red</code>: same as value <code>1</code></li>
				<li><code>green</code>: same as value <code>2</code></li>
				<li><code>yellow</code>: same as value <code>3</code></li>
				<li><code>blue</code>: same as value <code>4</code></li>
				<li><code>magenta</code>: same as value <code>5</code></li>
				<li><code>cyan</code>: same as value <code>6</code></li>
				<li><code>white</code>: same as value <code>7</code></li>
			</ul>
			<p>If the value used is invalid, no style will be applied for that format specifier.</p>
			<p>Old terminal emulators, have limited capabilities when rendering fonts and colors. If you want your program to support them, avoid using bold and italic, and prefer to use only colors of the 3 bits palette.</p>
			<p>Some terminal emulators, for instance, <code>st</code>, <code>konsole</code> and <code>linux</code> (the default virtual console of Linux), might render bold or italic with a brighter color.</p>
	<h2>Issues</h2>
		<p>Report issues through the <a href="https://github.com/skippyr/issues">issues tab</a>.</p>
	<h2>Contributions</h2>
		<p>If you want to contribute to this project, check out its <a href="https://skippyr.github.io/materials/pages/contributions_guidelines.html">contributions guidelines</a>.</p>
	<h2>License</h2>
		<p>This project is released under the terms of the MIT license.</p>
		<p>Copyright (c) 2023, Sherman Rofeman. MIT License.</p>
