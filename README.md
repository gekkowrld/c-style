# C Style

This is a custom implementation of the
[betty coding guideline](https://github.com/alx-tools/Betty/wiki)
in [golang](https://go.dev/).

**NOTES:**

- Most of the regex'es are created, generated and copied from [regex101](https://regex101.com).
- This is not a one-to-one replica of the betty lint tool.
If this tool fails please open an issue where you found the source code from.
- Golang converts all line endings to '\n' before processing the file, so cheking for DOS line endings is not possible.

## Design philosophy

- Produce more "human" readable error messages.
- Easily extensible as the code is broken down into smaller functions.
- Easy to use

## Building

- [go](https://go.dev) - The tool is written in go.

## License

This project is released under the [copyleft GNU GPL version 3](./LICENSE).
You are free to do anything with the source code as permitted by the LICENSE
