# Golang Wget

![Wget Logo](https://images.ctfassets.net/aoyx73g9h2pg/7z9SIh6Z6kTyYWKqIlg4VG/56519fa840f130b4ee79d109ecfef9d7/Wget-Linux-Diagram.jpg)

## Overview

`golang-wget` is a simple command-line tool written in Golang that mimics the functionality of the `wget` command with a subset of its commonly used flags. This tool allows you to download files from the internet using the specified URL and additional options.

## Features

- **-B, --base**: Set base URL for the relative links in the input file.
- **-i, --input-file**: Read URLs from a local or remote file.
- **--limit-rate**: Limit the download speed, specified in bytes per second.
- **--mirror**: Enable mirroring of the directory structure.
- **-P, --directory-prefix**: Set directory prefix to save the files.
- **-O, --output-document**: Specify the output file name.

## Installation

To install `golang-wget`, you need to have Golang installed on your system. If Golang is not installed, please visit [Golang's official website](https://golang.org/) for instructions on how to install it.

Once Golang is installed, run the following command to install `golang-wget`:

```
go get -u github.com/your-username/golang-wget

Usage



golang-wget [options] [URL]

Options

    -B, --base: Set the base URL for relative links.
    -i, --input-file: Read URLs from a local or remote file.
    --limit-rate: Limit download speed (bytes per second).
    --mirror: Mirror the directory structure.
    -P, --directory-prefix: Set the directory prefix to save files.
    -O, --output-document: Specify the output file name.

Examples

Download a file from a URL:


golang-wget https://example.com/file.txt

Download files listed in a file:



golang-wget -i urls.txt

Limit download speed to 1MB/s:



golang-wget --limit-rate=300k https://example.com/large-file.zip

Mirror a website:



golang-wget --mirror https://example.com/

Save downloaded files to a specific directory:



golang-wget -P /path/to/save https://example.com/file.txt

Specify output file name:



golang-wget -O my-file.txt https://example.com/download
