# Calcutron-33 The Decimal RISC CPU, Now in Go!

Calcutron-33 is the name of a made-up computer with a RISC-like microprocessor which operates with decimal numbers rather than binary numbers. While this may sound crazy it would not be entirely impossible to build such as computer. For instance the legendary Analytical Engine designed by Charles Babbage in the 1800s was a mechanical computer operating on decimal numbers rather than binary numbers.

This project contains a simulator which lets you run machine code programs on this imaginary computer. It also comes bundled with an assembler to turn assembly code into machine code which will run on this virtual machine.

The obvious question is: What it is the point? Primary purpose is as an educational tool. When teaching of microprocessors and computers work we hit the problem that most people unfamiliar with binary numbers. All modern digital computers operate with binary numbers rather than decimal numbers.

My motivation for making this virtual machine, was to offer a stepping stone towards RISC-V. RISC-V offers a great modern instruction-set for anyone who wants to get into assembly coding. Yet, I think that despite the simplicity of RISC-V it can be challenging for beginners.

A simple computer architecture and assembly language for beginners already exists. It is called the _Little Man Computer_. One of the problems I see it that _Little Man_ is quite far away from modern RISC microprocessor architectures.

Calcutron-33 tries to follow many design elements which are common for modern RISC processors:

- Operations such as add and subtract only happen on registers
- Special load and store instructions are used to get data to operate on and store results in memory
- A standard instruction takes 3 registers as operands. One for storing the results and two containing the values combined.
- Register x0 is always zero
- Small instruction-set but with many pseudo instruction to making coding easier

# Usage
Starting with Calcutron-33 version 2.0 we have bundled all commands into a single executable called `cutron` with subcommands for assembly, disassembly, simulation and debugging. Run it with the help subcommand to learn what options you have:

    ❯ cutron help
    NAME:
    cutron - Tool to assemble, disassemble and run Calcutron-33 assembly code

    USAGE:
    cutron [global options] command [command options] [arguments...]

    COMMANDS:
    assemble, asm        assemble a calcutron-33 assembly code file
    disassemble, disasm  disassemble a file containing calcutron-33 machine code
    run, simulate, sim   run a calcutron-33 machine code file
    debug, dbg           debug calcutron program
    help, h              Shows a list of commands or help for one command

    GLOBAL OPTIONS:
    --help, -h  show help (default: false)

You can also use `help` to get info about individual subcommands.

In the `examples` subdirectory you can find examples of assembly source code with the extension `.ct33` and assembled machine code files with extension `.machine`.

You can use the `assemble` subcommand to turn source code into machinecode. Here is an example of assembling the `simplemult.ct33` file:

    ❯ /cutron asm --sourcecode examples/simplemult.ct33
    loop:
    5109 LOAD x1, x0, -1
    5209 LOAD x2, x0, -1
    1300 CLR  x3
    multiply:
    1331 ADD  x3, x1
    2299 DEC  x2, -1
    9208 BGT  x2, x0, -2
    7309 STOR x3, x0, -1
    8000 JMP  x0, loop
    0000 HLT

You will notice we use the `--sourcecode` switch to show the original source code next to the generated 4-digit machine code.

The `sim` subcommand is used to run the simulator and actually execute the machine code. When you run the simulator it will read inputs on STDIN. In this example I am writing some inputs and hiting Ctrl-D when I am done.

    ❯ cutron sim examples/maximizer.machine
    2
    8
    6
    1
    00 5109 LOAD x1, x0, -1
    01 5209 LOAD x2, x0, -1
    02 9123 BGT  x1, x2, 3
    03 7209 STOR x2, x0, -1
    04 8000 JMP  x0, 0
    00 5109 LOAD x1, x0, -1
    01 5209 LOAD x2, x0, -1
    02 9123 BGT  x1, x2, 3
    05 7109 STOR x1, x0, -1
    06 8000 JMP  x0, 0
    00 5109 LOAD x1, x0, -1

    PC: 00    Steps: 10

    x1: 0006, x4: 0000, x7: 0000
    x2: 0001, x5: 0000, x8: 0000
    x3: 0000, x6: 0000, x9: 0000

    Inputs:  2, 8, 6, 1
    Outputs: 8, 6

The `maximizer` program looks at pairs of inputs and writes out the larger value to output. The first pair is 2 and 8 which produce an 8 on the output, while the second pair is 6 and 1 which produce a 6 on the output.

# Supported Instructions
All instructions are encoded as 4-digit decimal number where the first number indicates the opcode (the operation to perform) and the rest encode the operands (arguments to instruction). In theory this should give only 10 unique instructions but Calcutron-33 has a number of _pseudo instructions_ which is assembly code mnemonics which translates into one of the base instructions.

In the description whenever you read `Rd`, `Ra` or `Rb` then that  refers to a register from `x0`, `x1` to `x9`. Whenever you see a `k` that refers to a constant value. Instructions which have two register arguments in addition to the constant will only allow small constant values in the range -5 to 4. Those with only an `Rd` register operand will take `k` values in the trange -50 to 99.

## Arithmetic Operations
Typically you perform an operation with two source registers `Ra` and `Rb` and store the result in `Rd`.

The shift instruction `SHFT` is special in that it affects two registers `Rd` and `Ra`. When `k > 0` it multiplies `Ra` with `10^k`. Digits outside the range 0 to 9999 will be pushed over to `Rd`. For `k < 0` we get divisions instead. 

- `ADD Rd, Ra, Rb` - ADD registers
- `ADDI Rd, k` - ADD Immediate
- `SUB Rd, Ra, Rb` - SUBtract registers
- `SHFT Rd, Ra, k` - SHiFT digits left, or right if k is negative

## Load and Store Operations
The instructions combine a register `Ra` and a constant `k` to form a memory address, which we either read from or write to.

- `LOAD Rd, Ra, k` - LOAD
- `MOVE Rd, k` - MOVE value k to register `Rd`
- `STOR Rd, Ra, k` - Store to memory

## Jumping and Branching
Jump instruction are unconditional while branch instructions are conditional. `JMP` will form destination address with register `Rd` and constant `k`. The register serves a dual purpose in that the return address `PC + 1` is stored in `Rd` after a jump is executed. This way you can use `JMP` to return to the instruction after a call site.

The `BGT` and `BEQ` instructions compare registers `Ra` and `Rb` and depending on result jumps to a relative address  tiven by constant `k`. You can only jump 5 instructions back or 4 instructions forward, meaning `k` must be in the range -5 to 4. Using `k = 0` turns the branch into a `HLT` instruction because it would produce an infinite loop.

- `JMP  Rd,k` - JuMP to address
- `BGT  Ra, Rb, k` - Branch if Greater Than
- `BEQ  Ra, Rb, k` - Branch if EQual

## Pseudo Instructions
These instructions are all just shorthands for other instructions. For instance `INC Rd` is just short for `ADDI Rd, 1` and `COPY Rd, Ra` is short for `ADD Rd, x0, Ra`. Remember register `x0` is always zero.

- `DEC  Rd` - DECrement
- `INC  Rd` - INCrement
- `SUBI Rd, k`  - SUBtract Immediate

- `BLT  Ra, Rb, k` - Branch if Less Than
- `CLR  Rd` - CLeaR register
- `COPY Rd, Ra`  - COPY from one register to another

- `NOP` - No Operation
- `HLT` - HaLT execution
- `INP Rd` - INput instruction. Reads a number form input.
- `OUT Rd` - OUTput instruction. Writes number to output.

# History and Other Implementations
An early version of the idea of this decimal RISC based computer was described in this old article [this medium article](https://medium.com/@Jernfrost/decimal-risc-cpu-a13968922812). The instruction-set described in this article is now outdated.

The original article served as a based for an early prototype written in the Julia programming language. To explore the Zig language by implementing something of moderate complexity I made a variant called Zacktron-33 in Zig.

Remaking the whole thing in Go is both to get a sense of how Go compares with Zig but also because I believe Go is very well suited for this kind of project. I want to be able to easily distribute binaries that can run any macOS, Linux and Windows and that is possible with the good Go cross-compile support.

# Current Status
We got all the programs I desired to: assembler, disassembler, simulator and debugger. All the core stuff I would want to test and teach about assembly programming works. Yet one could say all of these programs are still a bit buggy or unpolished.

Addition and subtraction deal with negative numbers while shift and  conditional branching operate on numbers as if they were unsigned. Registers work on 4 digit numbers but there are only 100 memory locations numbered from 0 to 99.

Documenting properly the instruction set and their usage is still lacking.

