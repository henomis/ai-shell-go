# ai-shell-go
This is a simple AI shell helper written in GO. It use OpenAI API to generate a plausible command from a given prompt.
As soon as the command is generated, the user can choose to execute it or revise it adding more context.

## Installation
Be sure to have a working Go environment, then clone the repository and run the following command:

```
$ make 
```

## Usage

```
$ ./bin/ai-shell-go print first 3 lines of each file in a directory
```

## Output

```
Here is your command line:

$ head -n 3 *
--
This command uses the `head` utility to print the first 3 lines of each file in the current directory (`*` is a wildcard that matches all files in the directory). The `-n 3` flag specifies that it should print only the first 3 lines.

[E]xecute, [R]evise, [Q]uit? > 
```