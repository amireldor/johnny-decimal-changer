# Johnny Decimal Changer

A CLI tool for managing folders using the [Johnny Decimal](https://johnnydecimal.com/) system. This tool helps you rename and renumber folders while maintaining the Johnny Decimal format (xx.yy).

I built this tool to help me reorganize and manage my Johnny Decimal based folders. To be honest, I am currently trying the PARA method but I still like the folder sorting and indexing in Johnny Decimal, so I use a hybrid approach.

## Features

- **Rename Folders**: Change the prefix of folders (e.g., from `10.xx` to `20.xx`) while preserving the decimal parts and folder names
- **Sequential Renumbering**: Renumber folders to be sequential starting from a specified number
- **Gap Handling**: Automatically handles gaps in folder numbers when renumbering
- **Safe Operations**:
  - Dry run mode to preview changes before making them
  - Skips folders that already have the correct name
  - Prevents renumbering that would exceed the maximum for the specified number of digits
- **Flexible Numbering**:
  - Configurable number of digits after the decimal point (e.g., xx.99, xx.999, xx.9999)

## Usage

```bash
# Change prefix from 10 to 20
./johnny-decimal-changer -from 10 -to 20

# Renumber folders sequentially starting from 5
./johnny-decimal-changer -from 20 -start 5

# Preview changes without making them
./johnny-decimal-changer -from 10 -to 20 -dry-run
```

## Command Line Options

- `-from`: The current prefix of the folders to rename (required)
- `-to`: The new prefix to use (optional if -start is specified)
- `-start`: Start renumbering from this number (optional)
- `-dry-run`: Preview changes without making them (optional)
- `-digits`: Number of digits after the decimal point (default: 2)

## Examples

### Renaming with a New Prefix
```bash
# Before:
10.01 Projects
10.02 Documents
10.03 Archive

# Command:
./johnny-decimal-changer -from 10 -to 20

# After:
20.01 Projects
20.02 Documents
20.03 Archive
```

### Sequential Renumbering
```bash
# Before:
20.01 Projects
20.02 Archive
20.10 Documents

# Command:
./johnny-decimal-changer -from 20 -start 5

# After:
20.05 Projects
20.06 Archive
20.07 Documents
```

### Using More Digits
```bash
# Before:
20.01 Projects
20.02 Archive
20.10 Documents

# Command:
./johnny-decimal-changer -from 20 -to 90 -digits 4

# After:
90.0001 Projects
90.0002 Archive
90.0003 Documents
```

### Dry Run Mode
```bash
# Command:
./johnny-decimal-changer -from 10 -to 20 -dry-run

# Output:
Would rename: 10.01 Projects -> 20.01 Projects
Would rename: 10.02 Documents -> 20.02 Documents
Would rename: 10.03 Archive -> 20.03 Archive
```

## Safety Features

1. **Dry Run Mode**: Use `-dry-run` to preview changes before making them
2. **Skip Existing**: Folders that already have the correct name are skipped
3. **Format Protection**: Prevents renumbering that would exceed the maximum for the specified number of digits (e.g., xx.99 for 2 digits, xx.9999 for 4 digits)
4. **Validation**: Ensures all folder names follow the Johnny Decimal format

## Installation

```bash
go build
```

## Development

Run tests:
```bash
go test -v
```
