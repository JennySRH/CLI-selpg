package main

import (
	"bufio"
	"fmt"
	"errors"
	"strings"
	"log"
	"os"
	"math"
	"os/exec"
	"github.com/spf13/pflag"
)

const INT_MAX = 99999999


var start_page *int
var end_page *int
var page_len *int
var page_type *bool
var print_dest *string

func init() {
	start_page = pflag.IntP("start_page", "s", -1, "The start page.")
	end_page   = pflag.IntP("end_page", "e", -1, "The end page.")
	page_len   = pflag.IntP("page_len", "l", 72, "Default value, can be overriden by \"-l number\" on command line ")
	page_type  = pflag.BoolP("page_type", "f", false,  " for lines-delimited, f for form-feed-delimited.")
	print_dest = pflag.StringP("print_dest", "d", "", "The destination of the printer.")
}

func main() {
	pflag.Parse()
	err := process_args()
	if err!=nil {
		log.Println(err)
	}
	err = process_input()
	if err!=nil {
		log.Println(err)
	}
}

// It helps to handle the error.
func handleError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

/*================================= process_args() ================*/
func process_args() error {
	// check mandatory opts
	if *start_page == -1 || *end_page == -1 {
		return errors.New("not enough arguments")
	}
	// validity check of start page and end page
	if *start_page < 1 || *start_page > INT_MAX-1 {
		return errors.New("invalid start page")
	}
	if *end_page < 1 || *end_page > INT_MAX-1 || *end_page < *start_page {
		return errors.New("invalid end page")
	}
	// check mutually exclusive options
	if *page_type == true && *page_len != 72 {
		return errors.New("-l and -f are mutually exclusive options")
	}
	// validity check of page length
	if *page_len < 1 || *page_len > INT_MAX-1 {
		return errors.New("invalid page length")
	}
	// check other args
	if pflag.NArg() > 1 {
		return errors.New("invalid infilename")
	}
	return nil;
}

/*================================= process_input() ===============*/

func process_input() error {
	var line string
	line_cnt := 0
	// io process
	if pflag.NArg() == 0{
		// from stdin
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			if line_cnt > 0 {
				line += "\n"
			}
			line += input.Text()
			line_cnt++
		}
	} else {
		// from file
		infilename,err := os.Open(pflag.Arg(0))
		if(err != nil) {
			return err
		}
		input := bufio.NewScanner(infilename)
		for input.Scan() {
			if line_cnt > 0 {
				line += "\n"
			}
			line += input.Text()
			line_cnt++
		}
	}

	page_ctr := "\n"
	if(*page_type == true) {
		page_ctr = "\f"
	} 

	arr := strings.Split(line,page_ctr);

	pages := math.Ceil(float64(len(arr))/float64(*page_len))
	// check start page and end page	
	if *end_page >int(pages) || *start_page > int(pages) {
		return errors.New("invalid end page or start page : out of range")
	}


	var sub string;
	for i:= (*page_len) *(*start_page-1); i < (*page_len) *(*end_page)    && i < len(arr); i++ {
		if i!= (*page_len) *(*start_page-1){
			sub += page_ctr;
		}
		sub += arr[i];
	}
	if *print_dest =="" {
		fmt.Printf(sub)
	}  else {
		cmd := exec.Command("lp", "-d"+*print_dest)
		cmd.Stdin = strings.NewReader(sub);
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
