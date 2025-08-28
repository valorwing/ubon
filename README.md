# UBON (Universal Binary Object Notation)

UBON is a highly efficient, binary serialization format designed to be a simpler and more lightweight alternative to Protocol Buffers, while maintaining the ease of use and flexibility of JSON. UBON offers significant data compression, making it an ideal solution for IoT and other data-intensive applications.

## Features

- **Easy Integration**: Minimal codebase and low entry threshold, designed for developers of all levels.
- **Versatile**: Supports all standard data types, including nested objects and arrays.
- **Cross-Platform**: Initially implemented in Go, with plans for support in Python, C++, C, C#, JavaScript, PHP, and Rust.
- **Open Source**: Released under the GNU GPL license, ensuring it is free to use and modify.
- **Auto-Layout**: Data auto include data structure layout like JSON
- **No-Byte-Aligment**: Use bit level data manipulation no byte aligment and 3 bit flags
- **Compact**: Compress rate vs JSON ~ 0.4 (60% compression) vs CBOR ~ 0.75 (25% compression)

## Limitations

- **Arrays**: Only multidimensional arrays can be heterogeneous but in one dimenitonal scope object's or primitives must be single typed

## Build status

Zero build (proof-of-work) main brach all worked without array support now branch first-release with full refactor and stabilize binary protocol definition (see RUBICON)

## RUBICON

[29/12/25] All R&D operations finished. Make first production version with usage plan 9 assembly accelerations. Branch develop status in brogress branch - first-release