![](https://github.com/jcorry/nmap-scan-api/workflows/Go/badge.svg)
# nmap-scan API

A simple file parser for nmap XML data.


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
![](https://i.imgflip.com/3am9pr.jpg)

One of the things I love the most about software creation is the great fact that no matter how far you progress into the 
disciplines, you've merely scratched the surface. There is always another layer that can be added to make a program 
better.

I love building things and I really want to build them well. Perfection is probably not possible, but striving to become 
better at our craft and seeking better ways of doing things can be a part of every day and every task. 
[Kaizen](https://en.wikipedia.org/wiki/Kaizen), improvement.

The biggest assumption I made was that the purpose of the exercise is to demonstrate principles, values, patterns and 
my ideas about how software should be built. I tried to do that faithfully. I have interviewed many software engineering 
candidates and know what I look for, I have tried to deliver indications of my values.

Some guiding principles that I tried to convey here:
1. Separation of concerns
2. Security should be considered
3. Documentation should be adequate
4. Unit/Integration test should back the code
5. CI protects the codebase
6. Issues should be tracked and changes to code should reference the issue that required the change
7. Source control should be used effectively and should tell a story about the development of a project

I also made some compromises for the sake of expediency:
1. There's some duplication in test setup between `/cmd` and `/pkg`. I normally would use DB mocks in my http handler tests
but am experimenting with this idea that better tests use an actual DB instead of a mock. I'm not sure how I feel about 
this but wanted to try it.

2. I didn't have time to write Swagger/Openapi docs for the REST endpoints, so I just loosely documented them here in 
the README.

3. I would probably revise the data model if I were to do this over again. Once I got into implementing the UI (last!)
I found some characteristics of the data models that I didn't like.

### Run it!
The most reliable way to distribute the application is as a docker image. You will need docker installed on your 
machine to run the Dockerfile. If you do not have Docker installed on your machine, you may be able to build and run
the Go binary, but it will depend on your having SQLite installed. "It works on my machine", but to make sure you can run 
it too, I've packaged the app in a docker image.

1. Clone the repo
    `git clone git@github.com:jcorry/nmap-scan-api.git`

2. Build the image

    `docker build -t nmap-api .`

3. Run the image in a container

    `docker run -it -e "PORT=8080" -p 8080:8080 nmap-api`

4. Make requests!
[My Postman Collection](https://www.getpostman.com/collections/50372469fa0e3f090a47)

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
## **GET** /nmap/list (HTML representation)

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

Example: http://localhost:8080/api/v1/nmap?start=0&length=400

### Response
```
{
    "meta": {
        "start": 0,
        "length": 100,
        "total": 492
    },
    "items": [
        {
            "fileid": string,
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
        },
        ...
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
way more than me about nmap data structures so I am going to use their parser and model my app around their structs. I 
used the XML data format from the available files because it was the first I found suitable tools for parsing it.

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

The list view also makes a JSON representation available. An API URL is provided. If the `Content-Type` request header
is `application/json` the handler will respond with JSON.