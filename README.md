# golang-web-api-graphql-sqlserver-ssl-mutation

## Description
Using an Apollo server, built-in, query database
data of 10 country objects with internet codes.
Multiple filters applied to `Countries` field only.
Added a mutation query that creates a new country
object.

Sql server uses self-signed ssl.

## Tech stack
- bash
- golang 1.13

## Docker stack
- alpine:edge
- ubuntu:latest
- mcr.microsoft.com/mssql/server:2017-latest-ubuntu

## To run
`sudo ./install.sh -u`

## To stop
`sudo ./install.sh -d`

## To see help
`sudo ./install.sh -h`

## Credits
- https://www.worldstandards.eu/other/tlds/
- https://tutorialedge.net/golang/go-graphql-beginners-tutorial/
