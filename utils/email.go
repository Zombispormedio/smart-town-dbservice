package utils

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/sendgrid/sendgrid-go"
)

func SendInvitation(code string, email string) *RequestError {
	var Error *RequestError

	sendgridKey := os.Getenv("SENDGRID_API_KEY")

	if sendgridKey == "" {

		log.Error("Environment variable SENDGRID_API_KEY is undefined. Did you forget to source sendgrid.env?")

		return BadRequestError("Internal Error Sending Invitation")
	}

	sg := sendgrid.NewSendGridClientWithApiKey(sendgridKey)

	message := sendgrid.NewMail()

	message.AddTo(email)
	message.AddToName("Someone Special")
	message.SetSubject("Smart Town Invitation")

	template := "<html><body><h1>Congratulations! You have just received a invitation to participate in Smart Town as Admin</h1><p>Accept your <a href='%s'>invitation</a></p></body></html>"
	link := os.Getenv("SMART_INVITATION_HOST") + code

	result := fmt.Sprintf(template, link)
	message.SetHTML(result)
	message.SetFromName("Smart Town Dev Team")
	message.SetFrom("zombispormedio007@gmail.com")

	r := sg.Send(message)

	if r != nil {
		log.WithFields(log.Fields{
			"message": r.Error(),
		}).Error("InvitationSendError")
		Error = BadRequestError("Internal Error Sending Invitation")
	}

	return Error
}
