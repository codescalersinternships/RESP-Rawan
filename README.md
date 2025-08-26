# RESP file Parser

This repository implements a RESP file parser

## In this README 👇

- [Features](#features)
- [Usage](#usage)

## Features
 - `NewParser()`  returns a parser object to call our APIs
 - `Parse()` takes the RESP file lines and parses them into the calling parser object 
 - `LoadFromFile()` takes a path to an input RESP file and parses it to the calling parser object

## Usage

1. 
    ```go
        import github.com/codescalersinternships/RESP-Rawan
    ```

2. Get a new parser first
    ```go
        p := NewParser()
    ```

3. Example usage:
    ```go
        err := p.LoadFromFile(filePath)
			if err != nil {
				return err
		}
    ```


