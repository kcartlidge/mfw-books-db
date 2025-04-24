# MFW Books DB

A Go tool that fetches book information from Google Books and Open Library APIs using ISBNs, maintaining a catalogue in JSON form.

- Copyright 2025 K Cartlidge - [AGPL license](./LICENSE.txt)

ISBNs can be scanned in via mobile apps like [Alfa ISBN Scanner](https://www.alfaebooks.com/help/isbn_scanner) which lets you gather ISBNs then share them to an email as a text file suitable for input here.

## Contents

- [Features](#features)
- [API Rate Limits](#api-rate-limits)
- [Program Rate Limiting](#program-rate-limiting)
- [Usage](#usage)
- [File Formats](#file-formats)
- [Error Handling](#error-handling)
- [Producing New Builds](#producing-new-builds)

## Features

- Fetches book metadata from both Google Books and Open Library APIs
- Combines data from both sources for comprehensive information
- Handles rate limits and retries automatically
- Maintains a collection in JSON format
- Skips ISBNs already successfully fetched previously
- Sanitizes JSON file by removing books without ISBNs

## API Rate Limits

### Google Books API
- Free tier: 1,000 requests per day
- No authentication required for basic queries
- No per-second rate limit specified

### Open Library API
- No official rate limit documented
- Recommended limits:
    - Maximum 20 requests per second
    - Cache results when possible
    - Use appropriate delays between requests

## Program Rate Limiting

The program implements the following rate limiting:
- 333ms delay between requests (3 requests per second)
- Up to 3 retries for failed requests
- 2-second delay between retries
- Automatic handling of rate limit errors (HTTP 429)

## Usage

There are builds for Mac, Windows, and Linux in the [builds](./cmd/builds) folder.
Download the relevant one and place it somewhere reachable in your `PATH`.
There are *no* external dependencies to install.

1. Create a folder for your book data
2. Create a file named `isbns.txt` in that folder
3. Add ISBNs to the file, one per line
4. Run your downloaded build
    ```bash
    cd <wherever>
    ./mfw-books-db <folder_path>
    ```
5. Results will be written/updated in `books.json` in the specified folder
6. Repeat from step 3 to add more ISBNS (duplicates are ignored)

## File Formats

`isbns.txt`

    9781841493138
    9781841493145
    9781841493152

`books.json`

``` json
[
    {
        "isbn": "9780356517186",
        "title": "Some Desperate Glory",
        "authors": "Emily Tesh",
        "authorSort": "Tesh, Emily",
        "series": null,
        "sequence": null,
        "genres": [ "Science Fiction" ],
        "link": "https://books.google.com/books?id=fMSLzgEACAAJ\u0026dq=isbn:9780356517186",
        "isException": false,
        "exceptionReason": "",
        "modifiedUtc": "2025-04-23T18:07:40.615828Z",
        "status": "U - Unread",
        "rating": null,
        "notes": "",
        "statusIcon": "U",
        "seriesSort": ""
    },
    ...
]
```

## Error Handling

The program includes error handling for:
- Rate limit errors
- Network errors
- Missing or invalid ISBNs
- API response errors

Failed requests will still be added to the JSON file but with the exception flagged. The program will continue processing remaining ISBNs.

If an errored ISBN is later included in another search it's presence in the JSON file will be ignored and another attempt made.

## Producing New Builds

There are 3 scripts for producing builds, one each for Mac, Windows, and Linux.
*Run the one that relates to your system, not your build target.*

When run, those scripts will generate new builds for all platforms and place them in the `cmd/builds` folder.

Be sure to run the scripts from within the `cmd` folder.

    cd cmd
    ./scripts/macos.sh
    ./scripts/linux.sh
    scripts\windows.bat

For Mac and Linux you may need to `chmod a+x` to make them executable first.
