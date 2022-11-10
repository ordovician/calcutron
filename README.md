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
We got all the programs I desired to: assembler, disassembler, simulator and debugger. All the core stuff I would want to test and teach about assembly programming works. Yet one could say all of these programs are still a bit buggy or unpolished.

Addition and subtraction deal with negative numbers while shift and  conditional branching operate on numbers as if they were unsigned. Registers work on 4 digit numbers but there are only 100 memory locations numbered from 0 to 99.

Documenting properly the instruction set and their usage is still lacking.

