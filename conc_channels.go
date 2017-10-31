/******************************************************************************

 conc_channels.go

 Playing around with managing GO-routines (concurrency) and channels.
 This app spawns some workers who send their work output back to main()
 on a data-channel. The workers signal their completion of work on a
 separate control-channel. All worker data is consolidated onto a fan-in
 channel to be formatted and printed to stdout.

 TODO: Would be ideal to put all the random-string generator functions in a
       separate pakage (i.e. random.go) to be imported ("Do one thing and do
       it well").

 bfanselow 2017-10-31

******************************************************************************/
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//-------------------------------------------------------------------------------
const (
	// change this to make the whole thing more/less impressive
	Num_workers = 3

	// Strings used in random-string generators below
	alphanumset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	alphaset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphacapset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numset      = "0123456789"
)

//-------------------------------------------------------------------------------
// Generate new random seed on each invocation of app
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

//-------------------------------------------------------------------------------
// Create random string based on input characeter set
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

//-------------------------------------------------------------------------------
// Create random string of letters (upper and lower case)
func RandomAlphaString(length int) string {
	return StringWithCharset(length, alphaset)
}

//-------------------------------------------------------------------------------
// Create random string of alpha-numeric chaacters
func RandomAlphaNumString(length int) string {
	return StringWithCharset(length, alphanumset)
}

//-------------------------------------------------------------------------------
// Create random string of CAPITAL letters (A-Z)
func RandomAlphaCapString(length int) string {
	return StringWithCharset(length, alphacapset)
}

//-------------------------------------------------------------------------------
// Create random string of integers
func RandomNumString(length int) string {
	return StringWithCharset(length, numset)
}

//-------------------------------------------------------------------------------
// Create random integer within range of input min/max
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

//-------------------------------------------------------------------------------
// Format input string, prepend with timestamp and output to STDOUT
func ts_stdout(log_msg string) {
	loc, _ := time.LoadLocation("UTC")
	ts_now := time.Now().In(loc)
	fmt.Printf("%s: %s\n", ts_now, log_msg)
}

//-------------------------------------------------------------------------------
// Worker: do a specified number of  random-string operations sleeping a specified
// amount between each.  Put output of each operation onto the data-channel.
// Once work is finished, signal completion of work on the control-channel.
func worker(data_ch chan string, control_ch chan string, d time.Duration, wname string, limit int) {
	time.Sleep(2 * d)
	s := int(2 * d / 1000000)
	ts_stdout(fmt.Sprintf("Worker [%s] is starting after sleeping (%d) ms. Performing %d operations...", wname, s, limit))
	for i := 0; i <= limit; i++ {
		random_str := RandomNumString(10)

		// worker output. Put this on the data-channel
		ret_str := fmt.Sprintf("[id:%s]: msg-%s: (%s)", wname, strconv.Itoa(i), random_str)
		data_ch <- ret_str
		time.Sleep(d)
	}
	// worker now signals (on control channel) this it is done
	ts_stdout(fmt.Sprintf("Worker [%s] has completed and is sending its ID on control channel...", wname))
	control_ch <- wname
}

//-------------------------------------------------------------------------------
// "fan-in" channel reader
func reader(fi_ch chan string) {
	for msg := range fi_ch {
		fmt.Printf("  > %s\n", msg)
	}
	ts_stdout("FanIn Channel Reader is done reading from channel")
}

//-------------------------------------------------------------------------------
func main() {
	var N_comp = 0

	// create channels
	data_channel := make(chan string)
	control_channel := make(chan string)
	fanIn_channel := make(chan string)

	// Create Num_workers each with random worker-ids, sleep-durations, and num-operations
	ts_stdout(fmt.Sprintf("Main() creating %d workers...", Num_workers))
	for i := 0; i <= Num_workers; i++ {
		worker_id := RandomAlphaCapString(3)
		limit := RandomInt(20, 50)
		ms_coef := RandomInt(200, 500)
		var duration_ms = time.Duration(ms_coef) * time.Millisecond
		go worker(data_channel, control_channel, duration_ms, worker_id, limit)
	}

	go reader(fanIn_channel)

	// FAN-IN: put all data returned from worker channels and put onto single fanIn channel for printing
	ts_stdout("Main() Listening on data channel for worker data")
loop:
	for {
		select {
		case str := <-data_channel:
			fanIn_channel <- str
		case sname := <-control_channel:
			if sname != "" {
				N_comp++
				ts_stdout(fmt.Sprintf("Main() has detected complete signal from worker [%s]. Total-comp=(%d)", sname, N_comp))
				if N_comp == Num_workers {
					ts_stdout("Main() has determined that ALL workers are done sending on data-channel")
					break loop
				}
			}
		}
	}
}
