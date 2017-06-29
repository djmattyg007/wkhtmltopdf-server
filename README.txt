Wkhtmltopdf as a Service

This is a simple program to run a webserver that responds to all POST requests
that provide HTML in the request body with a fully-rendered PDF for that HTML.

By default, the program listens on port 8000. It takes all of the content in
the the request body and passes it to wkhtmltopdf via STDIN. The STDOUT of
wkhtmltopdf is then sent back as the response body.

Assuming there is a file named test.html full of HTML, here is an example:

curl -s -X POST -d "$(cat test.html)" -H 'Content-type: text/html' http://localhost:8000/pdf

This will return a fully-formed PDF document with the graph rendered by
wkhtmltopdf.

You can configure some of the options that wkhtmltopdf accepts by passing
query string parameters.

- grayscale: set to "1" to pass the --grayscale argument.
- lowquality: set to "1" to pass the --lowquality argument.
- orientation: set to "P" for Portrait or "L" for "Landscape". Passed with
  the --orientation argument. Defaults to "P".
- pagesize: set to any page size accepted by wkhtmltopdf. Passed with the
  --page-size argument. Defaults to "A4".
- title: set the title of the PDF. Passed with the --title argument.

An example with some query parameters passed:

curl -s -X POST -d "$(cat test.html)" -H 'Content-type: text/html' 'http://localhost:8000/pdf?orientation=L&pagesize=A3'

To configure the port number that is used, you can either use the --port CLI
argument, or you can set the WKHTMLTOX_PORT environment variable.

This software is released into the public domain without any warranty.
