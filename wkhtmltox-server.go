package main

import (
    "errors"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func parse_int_option(request *http.Request, key string) (int, error) {
    if value := request.URL.Query().Get(key); value != "" {
        valueStrConv, err := strconv.Atoi(value)
        if err == nil {
            return valueStrConv, nil
        } else {
            return -1, errors.New("invalid integer option value provided")
        }
    }
    return -1, nil
}

func parse_bool_option(request *http.Request, key string) (bool) {
    value := request.URL.Query().Get(key)
    return value == "1"
}

type pdf_opts struct {
    grayscale bool
    lowquality bool
    orientation string
    forms bool
    images bool
    javascript bool
    pagesize string
    title string
    imagedpi *int
    imagequality *int
}

func parse_pdf_options(request *http.Request) (pdf_opts, error) {
    opts := pdf_opts{}

    opts.grayscale = parse_bool_option(request, "grayscale")
    opts.lowquality = parse_bool_option(request, "lowquality")
    opts.forms = parse_bool_option(request, "forms")
    opts.images = !parse_bool_option(request, "noimages")
    opts.javascript = !parse_bool_option(request, "nojavascript")
    if imagedpi, err := parse_int_option(request, "imagedpi"); imagedpi > 0 {
        opts.imagedpi = &imagedpi
    } else if err != nil {
        return opts, errors.New("invalid imagedpi value provided")
    }
    if imagequality, err := parse_int_option(request, "imagequality"); imagequality > 0 {
        opts.imagequality = &imagequality
    } else if err != nil {
        return opts, errors.New("invalid imagequality value provided")
    }

    orientation := request.URL.Query().Get("orientation")
    if orientation == "P" {
        opts.orientation = "Portrait"
    } else if orientation == "L" {
        opts.orientation = "Landscape"
    } else if orientation == "" {
        opts.orientation = "Portrait"
    } else {
        return opts, errors.New("invalid orientation value provided")
    }

    opts.pagesize = request.URL.Query().Get("pagesize")
    if opts.pagesize == "" {
        opts.pagesize = "A4"
    }

    opts.title = request.URL.Query().Get("title")

    return opts, nil
}

func prepare_pdf_args(opts pdf_opts) ([]string) {
    args := []string{"--encoding", "utf-8"}

    if opts.grayscale {
        args = append(args, "--grayscale")
    }
    if opts.lowquality {
        args = append(args, "--lowquality")
    }
    if opts.forms {
        args = append(args, "--enable-forms")
    }
    if opts.images == false {
        args = append(args, "--no-images")
    }
    if opts.javascript == false {
        args = append(args, "--disable-javascript")
    }
    args = append(args, "--orientation", opts.orientation)
    args = append(args, "--page-size", opts.pagesize)
    if opts.title != "" {
        args = append(args, "--title", opts.title)
    }
    if opts.imagedpi != nil {
        args = append(args, "--image-dpi", strconv.Itoa(*opts.imagedpi))
    }
    if opts.imagequality != nil {
        args = append(args, "--image-quality", strconv.Itoa(*opts.imagequality))
    }

    // Tell wkhtmltopdf to accept input from STDIN, and output to STDOUT
    args = append(args, "-", "-")
    return args
}

func run_wkhtmltopdf(response http.ResponseWriter, html string, args []string) {
    binPath := os.Getenv("WKHTMLTOPDF_PATH")
    if binPath == "" {
        // The exec module will use PATH if it needs to, so this is fine
        binPath = "wkhtmltopdf"
    }
    cmd := exec.Command(binPath, args...)
    stdin, err := cmd.StdinPipe()
    if err != nil {
        response.WriteHeader(500)
        fmt.Println("An error occurred: ", err)
        return
    }
    defer stdin.Close()
    cmd.Stdout = response

    fmt.Printf("Rendering PDF. Arguments: '%s' ... ", strings.Join(args, " "))
    if err = cmd.Start(); err != nil {
        response.WriteHeader(500)
        fmt.Println("An error occurred: ", err)
        return
    }

    response.Header().Set("Content-type", "application/pdf")
    fmt.Fprintf(stdin, "%s", html)
    stdin.Close()
    cmd.Wait()
    fmt.Println("done")
}

func handle_pdf(response http.ResponseWriter, request *http.Request) {
    if request.Method != "POST" {
        response.WriteHeader(405)
        response.Header().Set("Allow", "POST")
        return
    }

    bhtml, err := ioutil.ReadAll(request.Body)
    if err != nil {
        response.WriteHeader(400)
        return
    }
    html := fmt.Sprintf("%s", bhtml)
    if html == "" {
        response.WriteHeader(400)
        return
    }

    opts, err := parse_pdf_options(request)
    if err != nil {
        response.WriteHeader(400)
        fmt.Fprintf(response, err.Error())
        return
    }
    args := prepare_pdf_args(opts)
    run_wkhtmltopdf(response, html, args)
}

type image_opts struct {
    cropheight *int
    cropwidth *int
    cropx *int
    cropy *int
    disablesmartwidth bool
    images bool
    javascript bool
    format string
    height *int
    quality *int
    width *int
}

func parse_image_options(request *http.Request) (image_opts, error) {
    opts := image_opts{}

    if cropheight, err := parse_int_option(request, "cropheight"); cropheight > 0 {
        opts.cropheight = &cropheight
    } else if err != nil {
        return opts, errors.New("invalid cropheight value provided")
    }
    if cropwidth, err := parse_int_option(request, "cropwidth"); cropwidth > 0 {
        opts.cropwidth = &cropwidth
    } else if err != nil {
        return opts, errors.New("invalid cropwidth value provided")
    }
    if cropx , err := parse_int_option(request, "cropx"); cropx > 0 {
        opts.cropx = &cropx
    } else if err != nil {
        return opts, errors.New("invalid cropx value provided")
    }
    if cropy, err := parse_int_option(request, "cropy"); cropy > 0 {
        opts.cropy = &cropy
    } else if err != nil {
        return opts, errors.New("invalid cropy value provided")
    }
    quality, err := parse_int_option(request, "quality")
    if quality > 0 {
        if quality <= 100 {
            opts.quality = &quality
        } else {
            return opts, errors.New("invalid quality value provided")
        }
    } else if err != nil {
        return opts, errors.New("invalid quality value provided")
    }
    if height, err := parse_int_option(request, "height"); height > 0 {
        opts.height = &height
    } else if err != nil {
        return opts, errors.New("invalid height value provided")
    }
    if width, err := parse_int_option(request, "width"); width > 0 {
        opts.width = &width
    } else if err != nil {
        return opts, errors.New("invalid width value provided")
    }
    opts.disablesmartwidth = parse_bool_option(request, "disablesmartwidth")
    opts.images = !parse_bool_option(request, "noimages")
    opts.javascript = !parse_bool_option(request, "nojavascript")

    if format := request.URL.Query().Get("format"); format != "" {
        if format == "png" {
            opts.format = "png"
        } else if format == "jpg" {
            opts.format = "jpg"
        } else {
            return opts, errors.New("invalid format provided")
        }
    } else {
        opts.format = "png"
    }

    return opts, nil
}

func prepare_image_args(opts image_opts) ([]string) {
    args := []string{"--encoding", "utf-8"}

    if opts.cropheight != nil {
        args = append(args, "--crop-h", strconv.Itoa(*opts.cropheight))
    }
    if opts.cropwidth != nil {
        args = append(args, "--crop-w", strconv.Itoa(*opts.cropwidth))
    }
    if opts.cropx != nil {
        args = append(args, "--crop-x", strconv.Itoa(*opts.cropx))
    }
    if opts.cropy != nil {
        args = append(args, "--crop-y", strconv.Itoa(*opts.cropy))
    }
    if opts.disablesmartwidth {
        args = append(args, "--disable-smart-width")
    }
    if opts.images == false {
        args = append(args, "--no-images")
    }
    if opts.javascript == false {
        args = append(args, "--disable-javascript")
    }
    args = append(args, "--format", opts.format)
    if opts.height != nil {
        args = append(args, "--height", strconv.Itoa(*opts.height))
    }
    if opts.width != nil {
        args = append(args, "--width", strconv.Itoa(*opts.width))
    }
    if opts.quality != nil {
        args = append(args, "--quality", strconv.Itoa(*opts.quality))
    }

    // Tell wkhtmltoimage to accept input from STDIN, and output to STDOUT
    args = append(args, "-", "-")
    return args
}

func run_wkhtmltoimage(response http.ResponseWriter, html string, format string, args []string) {
    binPath := os.Getenv("WKHTMLTOIMAGE_PATH")
    if binPath == "" {
        // The exec module will use PATH if it needs to, so this is fine as a fallback
        binPath = "wkhtmltoimage"
    }
    cmd := exec.Command(binPath, args...)
    stdin, err := cmd.StdinPipe()
    if err != nil {
        response.WriteHeader(500)
        fmt.Println("An error occurred: ", err)
        return
    }
    defer stdin.Close()
    cmd.Stdout = response

    fmt.Printf("Rendering image. Arguments: '%s' ... ", strings.Join(args, " "))
    if err = cmd.Start(); err != nil {
        response.WriteHeader(500)
        fmt.Println("An error occurred: ", err)
        return
    }

    if format == "png" {
        response.Header().Set("Content-type", "image/png")
    } else if format == "jpg" {
        response.Header().Set("Content-type", "image/jpeg")
    } else {
        response.WriteHeader(500)
        fmt.Println("Invalid format")
        return
    }

    fmt.Fprintf(stdin, "%s", html)
    stdin.Close()
    cmd.Wait()
    fmt.Println("done")
}

func handle_image(response http.ResponseWriter, request *http.Request) {
    if request.Method != "POST" {
        response.WriteHeader(405)
        response.Header().Set("Allow", "POST")
        return
    }

    bhtml, err := ioutil.ReadAll(request.Body)
    if err != nil {
        response.WriteHeader(400)
        return
    }
    html := fmt.Sprintf("%s", bhtml)
    if html == "" {
        response.WriteHeader(400)
        return
    }

    opts, err := parse_image_options(request)
    if err != nil {
        response.WriteHeader(400)
        fmt.Fprintf(response, err.Error())
        return
    }
    args := prepare_image_args(opts)
    run_wkhtmltoimage(response, html, opts.format, args)
}

func handle_license(response http.ResponseWriter, request *http.Request) {
    if request.Method != "GET" {
        response.WriteHeader(405)
        response.Header().Set("Allow", "GET")
        return
    }

    binPath := os.Getenv("WKHTMLTOPDF_PATH")
    if binPath == "" {
        // The exec module will use PATH if it needs to, so this is fine
        binPath = "wkhtmltopdf"
    }
    cmd := exec.Command(binPath, "--license")
    cmd.Stdout = response
    fmt.Printf("Printing license... ")
    response.Header().Set("Content-type", "text/plain")
    if err := cmd.Start(); err != nil {
        response.WriteHeader(500)
        fmt.Println("An error occurred: ", err)
        return
    }

    cmd.Wait()
    fmt.Println("done")
}

func handle_version(response http.ResponseWriter, request *http.Request) {
    if request.Method != "GET" {
        response.WriteHeader(405)
        response.Header().Set("Allow", "GET")
        return
    }

    binPath := os.Getenv("WKHTMLTOPDF_PATH")
    if binPath == "" {
        // The exec module will use PATH if it needs to, so this is fine
        binPath = "wkhtmltopdf"
    }
    cmd := exec.Command(binPath, "--version")
    cmd.Stdout = response
    fmt.Printf("Printing version... ")
    response.Header().Set("Content-type", "text/plain")
    if err := cmd.Start(); err != nil {
        response.WriteHeader(500)
        fmt.Println("An error occurred: ", err)
        return
    }

    cmd.Wait()
    fmt.Println("done")
}

func main() {
    var portNumber int
    flag.IntVar(&portNumber, "port", 0, "Port to listen on")
    flag.Parse()
    if portNumber == 0 {
        portNumberStr := os.Getenv("WKHTMLTOX_PORT")
        if portNumberStr == "" {
            portNumber = 8000 // The default port number
        } else {
            portNumberStrConv, err := strconv.Atoi(portNumberStr)
            if err == nil {
                portNumber = portNumberStrConv
            } else {
                fmt.Println(err)
                os.Exit(2)
            }
        }
    }

    http.HandleFunc("/pdf", handle_pdf)
    http.HandleFunc("/image", handle_image)
    http.HandleFunc("/license", handle_license)
    http.HandleFunc("/version", handle_version)
    fmt.Println(fmt.Sprintf("Starting webserver on port %d", portNumber))
    http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil)
}
