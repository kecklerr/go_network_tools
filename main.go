package main

import (
	"fmt"
	"github.com/sparrc/go-ping"
	"html/template"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"
)

type Message struct {
	Body string
	Pingfirst string
	Pingfinal string
}

func Index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()  // parse arguments, you have to call this by yourself
	fmt.Println("method:", r.Method) //get request method
	fmt.Println(r.Form)  // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	if r.Method == "GET" {
		output := Message{}
		DisplayTmpl(w, output, "template/index.gtpl")
		//t, _ := template.ParseFiles("template/index.gtpl")
		//t.Execute(w, nil)
	} else {
		r.ParseForm()
		fmt.Println("host:", r.Form["host"])
		fmt.Println("port:", r.Form["port"])
	}
}

func LookupHost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	output := Message{}
	myhost := fmt.Sprintf("%v", r.Form["host"][0])
	ipaddress, err := net.LookupHost(myhost)
	if err != nil {
		fmt.Println("ERROR: %s", err)
		output = Message{Body: err.Error()}
		DisplayTmpl(w, output, "template/nslookup.gtpl")
	} else {
		output = Message{Body: ipaddress[0]}
		DisplayTmpl(w, output, "template/nslookup.gtpl")
	}

}

func Ping(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	output := Message{}
	myhost := fmt.Sprintf("%v", r.Form["host"][0])
	pinger, err := ping.NewPinger(myhost)
	if err != nil {
		fmt.Println("ERROR: %s", err.Error())
		myerror := Message{Body: err.Error()}
		DisplayTmpl(w, myerror, "template/ping.gtpl")
	} else {

		pinger.Count = 3

		pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}
		pinger.OnFinish = func(stats *ping.Statistics) {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}

		fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
		finaloutput := fmt.Sprintf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

		pinger.Run() // blocks until finished
		stats := ping.Statistics{}
		pingfirst := fmt.Sprintf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		pingfinal := fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		output = Message{Body: finaloutput, Pingfirst: pingfirst, Pingfinal: pingfinal}

		DisplayTmpl(w, output, "template/ping.gtpl")
	}
}

func Telnet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mymessage := Message{}
	myhost := fmt.Sprintf("%v",r.Form["host"][0])
	myport := fmt.Sprintf("%v",r.Form["port"][0])

	match, _ := regexp.MatchString("[\\w/\\-?=%.]+\\.[\\w/\\-?=%.]+", myhost)
	if match != false && len(r.Form["port"][0])!=0 {
		fmt.Println("It matched")
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(myhost,myport), time.Second)
		if err != nil {
			fmt.Println("could not connect to server: ", err)
			mymessage = Message{Body: "No Connection - Error"}
			DisplayTmpl(w, mymessage, "template/telnet.gtpl")
		}
		if conn != nil {
			defer conn.Close()
			mymessage = Message{Body: "Good Connection"}
			DisplayTmpl(w, mymessage, "template/telnet.gtpl")
		}

	} else {
		fmt.Println("It did NOT match")
		mymessage = Message{Body: "Either the Host or port is not correct"}
		DisplayTmpl(w, mymessage, "template/telnet.gtpl")
	}
}

func DisplayTmpl(w http.ResponseWriter, message Message, htmltpl string){
	fmt.Println(message)
	tmpl, err := template.ParseFiles(htmltpl)
	if err == nil {
		errs := tmpl.Execute(w, message)
		if errs != nil {
			fmt.Println("Error in template: ", errs)
		}
	} else {
		fmt.Println("Eror in displaytmpl: ", err)
	}
}

func main() {
	http.HandleFunc("/", Index) // set router
	http.HandleFunc("/telnet", Telnet)
	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/nslookup", LookupHost)
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
