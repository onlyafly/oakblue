# Oakblue Virtual Machine

Oakblue is a virtual machine and assembler for that virtual machine, as a learning project.

The program started off as an implementation of the LC3 abstract machine.

## Using

## Developing the interpreter/compiler

Architecture of assembler:

ASM file -> scanner -----> parser -----> analyzer -----> emitter ----------> BIN file
                    tokens         CST             AST            bytecode

Architecture of VM:

BIN file -> executor

## Other

To lint the project, use [golangci-lint](https://github.com/golangci/golangci-lint).

The project is laid out according to these guidelines:
[Golang Standards Project Layout](https://github.com/golang-standards/project-layout)
