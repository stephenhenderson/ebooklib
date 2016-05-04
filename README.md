# A simple ebook library in Golang

A basic webapp for managing a library of ebooks. Ebooks have meta information 
(title, authors, etc.) and an arbitrary collection of files (the ebooks 
themselves + code sample zips, etc.). 

## Running
The library takes a single "-config" parameter pointing to a json config
containing details of where to store the library, the port to listen on, etc.
See [config_example.json](config_example.json) for details. 

## TODO
* Images for books
* Date updated for books
* Search/filtering
* CSS
* Authentication
* More tests for webservice (form handling, etc.)