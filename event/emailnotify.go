package event

import (
	"fmt"
	"log"
	"net/smtp"
	"tiim/go-comment-api/model"

	"github.com/jordan-wright/email"
)

type EmailNotify struct {
	From               string
	To                 string
	Subject            string
	Username           string
	Password           string
	SmtpHost           string
	SmtpPort           string
	CommentToUrlMapper func(model.Comment) string
}

func (n *EmailNotify) OnNewComment(c *model.Comment) (bool, error) {
	go n.doSendEmail(c)
	return true, nil
}

func (n *EmailNotify) doSendEmail(c *model.Comment) {

	log.Printf("sending notification email from %s to %s", n.From, n.To)

	page := c.Page
	if n.CommentToUrlMapper != nil {
		page = n.CommentToUrlMapper(*c)
	}

	e := email.NewEmail()
	e.From = n.From
	e.To = []string{n.To}
	e.Subject = n.Subject
	e.Text = []byte(fmt.Sprintf("New Comment\nid:\t%s\nfrom:\t%s <%s>\npage:\t%s\n\n%s", c.Id, c.Name, c.Email, page, c.Content))

	log.Printf("sending mail: %s:%s user:%s", n.SmtpHost, n.SmtpPort, n.Username)
	err := e.Send(n.SmtpHost+":"+n.SmtpPort, smtp.PlainAuth("", n.Username, n.Password, n.SmtpHost))

	if err != nil {
		log.Println("error sending notification email:", err)
	} else {
		log.Println("notification email sent")
	}
}

func (n *EmailNotify) OnDeleteComment(c *model.Comment) (bool, error) {
	return true, nil
}

func (n *EmailNotify) Name() string {
	return "EmailNotify"
}
