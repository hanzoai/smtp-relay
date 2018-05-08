package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/textproto"
	"os"

	"github.com/jordan-wright/email"
)

const authUser = os.Getenv("SMTP_RELAY_USER")
const authPass = os.Getenv("SMTP_RELAY_PASS")

// Unsupported
type attachment struct {
	buf         io.Reader
	filename    string
	contentType string
}

type request struct {
	// SMTP Server configuration
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`

	// Email to send
	FromEmail string `json:"fromEmail"`
	FromName  string `json:"fromName"`

	To          []string `json:"to"`
	Cc          []string `json:"cc"`
	Bcc         []string `json:"bcc"`
	ReadReceipt []string `json:"readReceipt"`

	Subject string `json:"subject"`

	Html string `json:"html"`
	Text string `json:"text"`

	// TODO: Support attachments
	// Attachments []attachment `json:"attachments"`
}

func formatFrom(req request) string {
	from := ""
	if req.FromName != "" {
		from = req.FromName
	}

	if from == "" {
		from = req.FromEmail
	} else {
		from = from + "<" + req.FromEmail + ">"
	}
	return from
}

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "ok\n")
	})

	router.POST("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Ensure only authenticated requests are allowed to do relay
		fmt.Println("authorizing user")
		user, pass, _ := r.BasicAuth()
		if user != authUser || pass != authPass {
			fmt.Printf("Unauthorized: %s %s\n", user, pass)
			http.Error(w, "Unauthorized.", 401)
			return
		}

		// Decode request body
		req := request{}
		reqBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Problem reading request body: %v", err)
			http.Error(w, "Problem reading request body", 400)
			return
		}

		err = json.Unmarshal(reqBytes, &req)
		if err != nil {
			fmt.Printf("Problem unmarshalling JSONy: %v", err)
			http.Error(w, "Problem unmarshalling JSON", 400)
			return
		}

		fmt.Printf("Request: %+v\n", req)

		e := &email.Email{
			To:          req.To,
			Bcc:         req.Bcc,
			Cc:          req.Cc,
			ReadReceipt: req.ReadReceipt,
			From:        formatFrom(req),
			Subject:     req.Subject,
			Text:        []byte(req.Text),
			HTML:        []byte(req.Html),
			Headers:     textproto.MIMEHeader{},
		}

		auth := smtp.PlainAuth("", req.Username, req.Password, req.Host)

		fmt.Println("sending email")
		err = e.Send(req.Host+":"+req.Port, auth)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Client error", 500)
			return
		}

		fmt.Fprint(w, "ok\n")
	})

	fmt.Println("listening on port 443")
	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
