# EANBaker

A Go application that generates PDF files containing EAN barcode labels from CSV or Excel data. EANBaker features both a graphical user interface and command-line interface for flexible usage.

## Features

- **Dual Interface**: Choose between GUI mode for easy interaction or CLI mode for automation
- **Multiple Input Formats**: Supports both CSV and Excel files
- **Barcode Generation**: Creates EAN barcodes with accompanying text labels
- **Customizable Layout**: Generate multiple copies of each barcode
- **Flexible Configuration**: Configurable column headers and CSV separators
- **Cross-Platform**: Runs on Windows, macOS, and Linux

## Installation

### Download Pre-built Binary

Download the latest release for your platform from the [GitHub Releases](https://github.com/Fanteria/EANBaker/releases).

### Build from Source

If you prefer to build from source or need a custom build:

#### Prerequisites

- Go 1.19 or later
- Git (for cloning the repository)

```bash
git clone https://github.com/Fanteria/EANBaker.git
cd EANBaker
go build -o eanbaker
```

## Usage

### GUI Mode (Default)

Launch the application without any arguments to start the graphical interface:

```bash
./eanbaker
```

#### GUI Features:

- **File Selection**: Click "Choose file" to select your CSV or Excel file
- **Column Headers**: Specify the column names for text and EAN data
- **Options Page**: Configure advanced settings like CSV separator, PDF output path, and barcode repetition

### Command Line Mode

Use command-line flags for automated processing:

```bash
./eanbaker -csv data.csv -text-header "Product Name" -ean-header "Barcode"
```

#### Command Line Options:

| Flag              | Default            | Description                                           |
| ----------------- | ------------------ | ----------------------------------------------------- |
| `-csv`            | `data.csv`         | Path to input CSV file                                |
| `-pdf`            | (_CSV file name_)  | Output PDF file path                                  |
| `-text-header`    | `Material Number`  | Column header for text labels (case-insensitive)      |
| `-ean-header`     | `ean`              | Column header for EAN codes (case-insensitive)        |
| `-times-header`   | `""`               | Column containing repetition counts for each EAN code | 
| `-times-each-ean` | `1`                | Number of copies per barcode                          |
| `-csv-separator`  | `,`                | CSV column separator character                        |

#### Examples:

Generate PDF with default settings:

```bash
./eanbaker -csv products.csv
```

Custom headers and output path:

```bash
./eanbaker -csv inventory.csv -text-header "Item Name" -ean-header "SKU" -pdf labels.pdf
```

Multiple copies with semicolon separator:

```bash
./eanbaker -csv data.csv -csv-separator ";" -times-each-ean 3
```

## Configuration

EANBaker automatically saves your settings to `.EANBaker.json` in the current directory. This hidden file stores:

- Column header preferences
- CSV separator settings
- PDF output path preferences
- Barcode repetition settings

