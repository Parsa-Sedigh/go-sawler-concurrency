package main

import (
	"bytes"
	"fmt"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"sync"
	"time"
)

// This type describes the actual mail server
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string // username to login to the mail server to send the email
	Password    string
	Encryption  string // mail server encryption
	FromAddress string // the default "from" address(when you send email it has to come from someone!)
	FromName    string
	Wait        *sync.WaitGroup

	// this is where we're gonna send any email we want sent in the background
	MailerChan chan Message
	ErrorChan  chan error
	DoneChan   chan bool
}

// describes an actual message
type Message struct {
	From          string // in case you wanna overwrite the default FromAddress
	FromName      string
	To            string // holds the email address we're sending to
	Subject       string
	Attachments   []string // would contain full path names to any file we want to attach to the email message
	AttachmentMap map[string]string
	Data          any            // body of the message
	DataMap       map[string]any // just a convenient of getting data to the actual template we'll be rendering
	Template      string
}

// a function to listen for messages on the mailer chan. This function is gonna run in the background
func (app *Config) listenForMail() {
	for {
		select {
		case msg := <-app.Mailer.MailerChan:
			go app.Mailer.sendMail(msg, app.Mailer.ErrorChan)

			// we tried to send email and sth went wrong:
		case err := <-app.Mailer.ErrorChan:
			// you might want to notify someone on slack or ...
			app.ErrorLog.Println(err)
		case <-app.Mailer.DoneChan:
			// stop processing email in the background(by returning, we quit this goroutine)
			return
		}
	}
}

func (m *Mail) sendMail(msg Message, errorChan chan error) {
	defer m.Wait.Done()

	// check to see if they specified a custom template
	if msg.Template == "" {
		// by setting this, it's gonna capture both mail.gohtml and mail.plain.gohtml
		msg.Template = "mail"
	}

	// did we specify a custom from address in this message?
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// tutor is not sure we have to odo this,but it's not gonna hurt to do this!
	if msg.AttachmentMap == nil {
		msg.AttachmentMap = make(map[string]string) // just make a initialized map
	}

	// send information to the mail templates
	//data := map[string]any{
	//	"message": msg.Data,
	//}

	if len(msg.DataMap) == 0 {
		msg.DataMap = make(map[string]any) // just make a initialized map
	}

	// do not overwrite everything with the line below:
	//msg.DataMap = data
	// instead use this:
	msg.DataMap["message"] = msg.Data

	// build two versions of the message:
	// build html mail
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		errorChan <- err
	}

	// build plain text mail
	plainMessage, err := m.buildPlainTextMessage(msg)

	// create an SMTP client. Sth iin mmy code that will connect to the mail server:
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)

	// we're not gonna be sending mail every second, so set this to false
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		errorChan <- err
	}

	// create an email message:
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	// by default, the body of our message will be the plain text:
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	if len(msg.AttachmentMap) != 0 {
		for key, value := range msg.AttachmentMap {
			/* first args is the name oof the file stored somewhere and second arg is what we'll call it when we add the attachment. So we can overwrite the
			stored file name when attaching it to the email. This allows us to have a user-friendly name when attaching generated files with weird names
			to the email. */
			email.AddAttachment(value, key)
		}
	}

	// let's send the email!:
	err = email.Send(smtpClient)
	if err != nil {
		errorChan <- err
	}
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	// build the name of the file using fmt.Sprintf() to get the full path to the file we want to use to render this email
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.html.gohtml", msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// now that we arrived here, it means we have a formatted message or sth that we can use as a formatted message, we still have a bit of work to do on it
	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.plain.gohtml", msg.Template)

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

// we're making css more acceptable to various email clients
func (*Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	// prem stands for premailer
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(e string) mail.Encryption {
	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
