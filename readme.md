![](https://github.com/jcorry/nmap-scan-api/workflows/Go%20build/badge.svg)
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
![](https://imgflip.com/i/3am9pr)

### Run it
1. Build the image

    `docker build -t nmap-api .`

2. Run the image in a container

    `docker run -it -e "PORT=8080" -p 8080:8080 nmap-api`

3. Make requests!

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
ERROR_MESSAGE string
```

## **GET** /api/v1/nmap
## **GET** /nmap (HTML representation)

Get a paginated list of nmap hosts

### Request

- start
  - In: query
  - Valid: int
  - Matches: The starting index of records to retrieve

- length
  - In: query
  - Valid: int (max: 1000)
  - Matches: The number of records to retrieve

Example: http://localhost:4000/api/v1/nmap?start=0&length=400

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

I added a file import table to log the importation of files with a unique hash per file. This way, I can prevent duplicate
data imports. I'm filing an issue to come back and handle the case where the file was accepted on upload and its hash saved
to the imports table but the batch insert of hosts fails, rollsback and leaves the hash in the imports table. This case would
prevent the user from ever successfully uploading that file.

I've never seen nmap data before today but found a Go package used for parsing nmap XML files. The author clearly knows 
way more than me about nmap data structures so I am going to use their parser and model my app around their structs.

[parser](https://godoc.org/github.com/tomsteele/go-nmap)

There's a ton of data that is available in nmap results, for our demo project we will focus on only a subset of that data. 
My API response as designed above will inform the SQL structure and the internal structs that query data is scanned to.
Much to my dismay, the structure of the XML does not really map well to a relational DB. Ports and Hostnames are (I think)
related to Addresses, but there's no relationship between these depicted in the XML structure or its parsed equivalent.

One mistake I made was in the `sqlite HostRepo.list()` function. The way I handled scanning row results into structs and
then storing those structs in maps so they can be keyed to the host prevents DB based sorting. This should be refactored
and can be remedied by using a `[]*models.Host` as the parent structure and eliminating the `hostMap`. 

#### SQL
![sql schema](https://docs.google.com/drawings/d/e/2PACX-1vQM_B93_LE8tp0kMWfel9LPAaOtlSLgKrqUsvxNN5B6HJIz0s92p91tNwnQCx1D6CYmH0ir8VGl9hVQ/pub?w=960&h=720)

### UI
The UI is a simple HTML grid supplied by Bootstrap 4. The HTML is generated from templates using Go's `html/template` package
which is sparse, but adequate for very basic display.