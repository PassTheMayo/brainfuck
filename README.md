# Brainfuck
An extremely optimized and fast [Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) interpreter written entirely in Go. Supports all commands as well as different methods for handling [out-of-range memory](https://en.wikipedia.org/wiki/Brainfuck#Portability_issues:~:text=When%20the%20pointer%20moves%20outside%20the%20bounds%20of%20the%20array%2C%20some%20implementations%20will%20give%20an%20error%20message%2C%20some%20will%20try%20to%20extend%20the%20array%20dynamically%2C%20some%20will%20not%20notice%20and%20will%20produce%20undefined%20behavior%2C%20and%20a%20few%20will%20move%20the%20pointer%20to%20the%20opposite%20end%20of%20the%20array.) errors.

## Installation

You can either choose to download the source code as a ZIP file from the repository, or you can clone the repository using the command listed below.

```bash
git clone https://github.com/PassTheMayo/brainfuck.git
```

You will then need to move the working directory into the `brainfuck` folder you just cloned.

```bash
cd brainfuck
```

To compile the interpreter, you can either use [GNU Make](https://www.gnu.org/software/make/) which should come pre-installed on Ubuntu as well as many other Unix-based operating systems, or you can use the Go compiler itself.

```bash
make
# or
go build -o bin/main src/*.go # Unix
go build -o .\bin\main.exe src\*.go # Windows
```

You can now use the binary executable `main` or `main.exe` (OS-dependent) that resides in the `bin` folder to run any Brainfuck program. Use `--help` to list all options available.

## Sample Programs

There is an unofficial website for Brainfuck that has a lot of sample programs that you can use to test this interpreter. You can find that [at this link](http://brainfuck.org/).

## License

[MIT License](https://github.com/PassTheMayo/brainfuck/blob/main/LICENSE)