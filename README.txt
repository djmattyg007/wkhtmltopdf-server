Wkhtmltopdf as a Service

This is a simple program to run a webserver that responds to all POST requests
that provide HTML in the request body with a fully-rendered PDF for that HTML.

By default, the program listens on port 8000. It takes all of the content in
the the request body and passes it to wkhtmltopdf via STDIN. The STDOUT of
wkhtmltopdf is then sent back as the response body.

Assuming there is a file named test.html full of HTML, here is an example:

curl -s -X POST -d "$(cat test.html)" -H 'Content-type: text/html' http://localhost:8000/pdf

This will return a fully-formed PDF document with the HTML rendered by
wkhtmltopdf.

You can configure some of the options that wkhtmltopdf accepts by passing
query string parameters.

- grayscale: set to "1" to pass the --grayscale argument.
- lowquality: set to "1" to pass the --lowquality argument.
- forms: set to "1" to pass the --enable-forms argument.
- noimages: set to "1" to pass the --no-images argument.
- nojavascript: set to "1" to pass the --disable-javascript argument.
- orientation: set to "P" for Portrait or "L" for "Landscape". Passed with
  the --orientation argument. Defaults to "P".
- pagesize: set to any page size accepted by wkhtmltopdf. Passed with the
  --page-size argument. Defaults to "A4".
- title: set the title of the PDF. Passed with the --title argument.
- imagedpi: set the DPI of images rendered by wkhtmltopdf. Passed with the
  --image-dpi argument.
- imagequality: set the quality of the images rendered by wkhtmltopdf. Passed
  with the --image-quality argument.

An example with some query parameters passed:

curl -s -X POST -d "$(cat test.html)" -H 'Content-type: text/html' 'http://localhost:8000/pdf?orientation=L&pagesize=A3&noimages=1'

If instead, you want an image, you can POST to the /image endpoint:

curl -s -X POST -d "$(cat test.html)" -H 'Content-type: text/html' http://localhost:8000/image

This will return a fully-formed PNG image file with the HTML rendered by wkhtmltoimage.

You can also configure some of the options that wkhtmltoimage accepts by
passing query string parameters.

- cropheight: set the height for cropping. Passed with the --crop-h argument.
- cropwidth: set the width for cropping. Passed with the --crop-w argument.
- cropx: set the x coordinate for cropping. Passed with the --crop-x argument.
- cropy: set the y coordinate for cropping. Passed with the --crop-y argument.
- height: set the screen height. Passed with the --height argument.
- width: set the screen width. Passed with the --width argument.
- disablesmartwidth: set to "1" to pass the --disable-smart-width argument.
- quality: a percentage value controlling the final quality of the image.
  Passed with the --quality argument.
- format: set to "png" for a PNG file, "jpg" for a JPEG" file. Passed with the
  --format argument. Defaults to "png".
- noimages: set to "1" to pass the --no-images argument.
- nojavascript: set to "1" to pass the --disable-javascript argument.

An example with some query parameters passed:

curl -s -X POST -d "$(cat test.html)" -H 'Content-type: text/html' 'http://localhost:8000/image?format=jpg&disablesmartwidth=1'

To obtain the license text from wkhtmltopdf, you would make a GET request to
the /license endpoint:

curl -s http://localhost:8000/license

To obtain the wkhtmltopdf version information, you would make a GET request to
the /version endpoint:

curl -s http://localhost:8000/version

To configure the port number that is used, you can either use the --port CLI
argument, or you can set the WKHTMLTOX_PORT environment variable.

This software is released into the public domain without any warranty.
