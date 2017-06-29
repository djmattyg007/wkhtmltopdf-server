package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "strconv"
)

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

    fmt.Printf("Rendering PDF... ")
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

func prepare_args(grayscale bool, lowquality bool, orientation string, pagesize string, title string) ([]string) {
    args := []string{}
    if grayscale {
        args = append(args, "--grayscale")
    }
    if lowquality {
        args = append(args, "--lowquality")
    }
    args = append(args, "--orientation", orientation)
    args = append(args, "--page-size", pagesize)
    if title != "" {
        args = append(args, "--title", title)
    }

    // Tell wkhtmltopdf to accept input from STDIN, and output to STDOUT
    args = append(args, "-", "-")
    return args
}

// TODO: Move parsing of query string to separate function.
// Have it return a struct and an err variable, that if not nil,
// indicates an error occurred.
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
    html:= fmt.Sprintf("%s", bhtml)
    if html == "" {
        response.WriteHeader(400)
        return
    }

    grayscale := request.URL.Query().Get("grayscale")
    grayscaleOpt := grayscale == "1"

    lowquality := request.URL.Query().Get("lowquality")
    lowqualityOpt := lowquality == "1"

    var orientationOpt string
    orientation := request.URL.Query().Get("orientation")
    if orientation == "P" {
        orientationOpt = "Portrait"
    } else if orientation == "L" {
        orientationOpt = "Landscape"
    } else if orientation == "" {
        orientationOpt = "Portrait"
    } else {
        response.WriteHeader(400)
        return
    }

    pagesizeOpt := request.URL.Query().Get("pagesize")
    if pagesizeOpt == "" {
        pagesizeOpt = "A4"
    }

    titleOpt := request.URL.Query().Get("title")

    args := prepare_args(grayscaleOpt, lowqualityOpt, orientationOpt, pagesizeOpt, titleOpt)
    run_wkhtmltopdf(response, html, args)
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
    //http.HandleFunc("/image", handle_image)
    fmt.Println(fmt.Sprintf("Starting webserver on port %d", portNumber))
    http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil)
}
