# nmap-scan API

## Problem

* Build a REST API to import a single nmap scan result.
    * The API should accept a single nmap scan file.
    * Detail which you chose and why.
* Ingest the nmap result into a sqlite database.
* Build a UI that allows a user to view the results of nmap scans by IP address.

In the above challenge, feel free to use any language. For your submission, package source code and files into a zip file. Include a README.md outlining:

1. Precise instructions to launch and use your submission including runtimes and versions.
2. Any assumptions you made.
3. Additional thoughts about the project.

## Solution

# HTTP Endpoints

## **POST** /api/v1/nmap/upload

Upload an nmap XML file to be parsed and saved to the DB.

## Request

## Response
201
```
created
```
400
```
{
    "error": ERROR_MESSAGE
}
```

## **GET** /api/v1/nmap

Get a list of nmap hosts

### Request

- addr
  - In: query
  - Valid: string
  - Matches: IP Address of Host to return data for

- start
  - In: query
  - Valid: int
  - Matches: The starting index of records to retrieve

- length
  - In: query
  - Valid: int (max: 1000)
  - Matches: The number of records to retrieve

Example: http://localhost:4000/api/v1/nmap?addr=192.168.1.1&start=0&length=400

### Response
```
{
    "links": {
        "self": {
            "href": "/api/v1/nmap"
        }
    },
    "meta": {
        "start": 0,
        "length": 100,
        "total": 492
    },
    "items": [
        {
            "starttime": timestamp,
            "endtime": timestamp,
            "comment": string
            "status": string,
            "hostnames": [
                "name": string,
                "type": string,
            ],
            "addresses": [
                "addr": string,
                "addrtype": string
            ],
            "ports": [
                "protocol": string,
                "portid": int,
                "owner": string,
                "service": string,
            ]
        }
    ]
}
```

### Data Structures

nmap data must be mapped to internal structs and SQL tables.

I've never seen nmap data before today but found a Go package used for parsing nmap XML files. The author clearly knows way more than me about nmap data structures so I am going to use their parser and model my app around their structs.

[parser](https://godoc.org/github.com/tomsteele/go-nmap)

There's a ton of data that is available in nmap results, for our demo project we will focus on only a subset of that data. My API response as designed above will inform the SQL structure and the internal structs that query data is scanned to.

#### SQL
![sql schema](https://docs.google.com/drawings/d/e/2PACX-1vQM_B93_LE8tp0kMWfel9LPAaOtlSLgKrqUsvxNN5B6HJIz0s92p91tNwnQCx1D6CYmH0ir8VGl9hVQ/pub?w=960&h=720)

### UI