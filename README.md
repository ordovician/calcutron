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

# History and Other Implementations
I described the computer, or more specifically the CPU on this computer in more detail in [this medium article](https://medium.com/@Jernfrost/decimal-risc-cpu-a13968922812).

That led to an early implementation in the Julia programming language. To explore the Zig language by implementing something of moderate complexity I made a variant called Zacktron-33 in Zig.

Remaking the whole thing in Go is both to get a sense of how Go compares with Zig but also because I believe Go is very well suited for this kind of project. I want to be able to easily distribute binaries that can run any macOS, Linux and Windows and that is possible with the good Go cross-compile support.

# Current Status
We got all the programs I desired to: assembler, disassembler, simulator and debugger. All the core stuff I would want to test and teach about assembly programming works. Yet one could say all of these programs are still a bit buggy or unpolished. Handling negative numbers is not working great. The example programs however don't require any negative number usage.

There are still cases where the error messages and feedback to the user could be a lot better.

Currently registers work with 2 digit decimal numbers, but I have been thinking that the more sensible choice is 4 digit decimal numbers because that is the size of each memory location.


# Refactoring Ideas
The whole project is built to a large degree like a prototype. There was no thoughtful upfront design. However because this is the third iteration of the assembler there has been some evolution of ideas. However, as complexity has grown I have noticed a few key problems.

There are switch-case statement on instuctions such as ADD, SUB and LD a large number of places. It means that adding or modifying and instruction requires making code changes in numerous locations. That is not an ideal solution.

Ideally each instruction should be self contained and know how to do the following:

- Disassemble machine code into assembly code
- Assemble source code into machine code
- Execute instruction on virtual CPU
- Visualize instruction with colors or without
- Know what info is most relevant to show after instruction is run in debugger. E.g. ADD and SUB benefit from showing registers, while a branch benefit from showing the program counter (PC).

Several instructions are strongly related which means having a single type representing each instruction may not make sense.  Instead one can bundle instructions such as ADD and SUB because they both work in 3 registers. While SUBI, LSH, RSH and ADDI can be bundled as they all use a constant (immediate value) as 3rd operand.

Likewises LD, ST, BRZ, BGT and BRA all take a register and address as operands which make them suitable to be used in a bundle.

A similar idea can be used for commands within the debugger. Each command needs to have support for the following:

- Handle command completion
- Execute when selected. Will likely require access to CPU/Computer object
- Help info. Explanation to user of what the command actually does.
