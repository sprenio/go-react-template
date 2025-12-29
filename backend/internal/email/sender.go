package email

import (
	"backend/assets"
	"backend/config"
	"backend/locale"
	"backend/pkg/logger"
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	mail "gopkg.in/mail.v2"
)

type embedImage struct {
	ContentID   string
	ContentType string
	FilePath    string
	FileName    string
}

type EmailSender struct {
	smtpHost    string
	smtpPort    int
	username    string
	password    string
	from        string
	embedImages []embedImage
}

//go:embed templates/*.html
var templateFiles embed.FS

func GetEmailSender() *EmailSender {

	cfg, err := config.GetConfig()
	if err != nil {
		return nil
	}
	from := cfg.AppName + " <" + cfg.Email.From + ">"
	return newEmailSender(cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.Username, cfg.Email.Password, from)
}

func newEmailSender(cfgHost, cfgPort, user, pass, from string) *EmailSender {

	var intPort int
	fmt.Sscanf(cfgPort, "%d", &intPort)
	return &EmailSender{
		smtpHost: cfgHost,
		smtpPort: intPort,
		username: user,
		password: pass,
		from:     from,
	}
}

func (es *EmailSender) AddEmbedImage(contentID, contentType, filePath, fileName string) {
	es.embedImages = append(es.embedImages, embedImage{
		ContentID:   contentID,
		ContentType: contentType,
		FilePath:    filePath,
		FileName:    fileName,
	})
}

func (es *EmailSender) SendHtmlEmail(to, subject, htmlContent, PageTitle, pageHeader, HeadExtra string) error {
	var htmlBody bytes.Buffer

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	m := mail.NewMessage()
	m.SetHeader("From", es.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	tmpl, err := template.ParseFS(templateFiles, "templates/base.html")
	if err != nil {
		return errors.Wrap(err, "parse base template")
	}

	randGen := rand.New(rand.NewSource(time.Now().UnixNano())) // lokalny generator
	num := randGen.Intn(900000) + 100000                       // losowa liczba 100000–999999
	logoCID := fmt.Sprintf("logoCID%d", num)

	err = tmpl.Execute(&htmlBody, map[string]interface{}{
		"PageTitle": PageTitle,
		"AppName":   cfg.AppName,
		"logoCID":   logoCID,
		"HeadExtra": HeadExtra,
		"Header":    pageHeader,
		"Content":   template.HTML(htmlContent),
		"Year":      time.Now().Year(),
	})
	if err != nil {
		return errors.Wrap(err, "execute base template")
	}
	logoPath, err := es.AddEmbeddedImageFromBytes(logoCID, "image/png", "logo.png", assets.Logo)
	if err == nil {
		logger.Info("Embedding logo image from embedded FS")
		defer os.Remove(logoPath)
	}

	logger.Info("Sending email to %s with subject: %s", to, subject)
	m.SetBody("text/html", htmlBody.String())
	//logger.Info("Email content:\n%s", htmlBody.String())
	for _, img := range es.embedImages {
		myMap := map[string][]string{
			"Content-ID":          {"<" + img.ContentID + ">"},
			"Content-Disposition": {"inline; filename=\"" + img.FileName + "\""},
		}
		m.Embed(img.FilePath, mail.SetHeader(myMap))
	}

	d := mail.NewDialer(es.smtpHost, es.smtpPort, es.username, es.password)

	//var buf bytes.Buffer
	//if _, err := m.WriteTo(&buf); err != nil {
	//	logger.Info("Błąd przy generowaniu maila: %v", err)
	//} else {
	//	logger.Info("Cała wiadomość do wysłania:\n%s", buf.String())
	//}
	return d.DialAndSend(m)
}

func (es *EmailSender) SendWelcomeEmail(to, userName, langCode, confirmationLink string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	logger.Info("Preparing to send welcome email to %s in language %s", to, langCode)
	loc := locale.GetNewLocalizer(langCode)
	// Fallback do EN jeśli brak tłumaczenia

	// Wygeneruj HTML z szablonu
	tmpl, err := template.ParseFS(templateFiles, "templates/welcome.html")
	if err != nil {
		return errors.Wrap(err, "parse welcome email template")
	}

	var htmlContent bytes.Buffer
	err = tmpl.Execute(&htmlContent, map[string]interface{}{
		"Hello":            template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "general.hello_user", TemplateData: map[string]string{"UserName": userName}})),
		"ThankYouAndClick": template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "welcome.thank_you_and_click", TemplateData: map[string]string{"AppName": cfg.AppName}})),
		"ConfirmationLink": confirmationLink,
		"ConfirmEmail":     template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "welcome.confirm_email"})),
		"IfButtonFails":    template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "general.if_button_fails"})),
		"LinkExpiryInfo":   template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "welcome.link_expiry_info", TemplateData: map[string]int{"ExpiryDays": cfg.Register.ExpirationDays}, PluralCount: cfg.Register.ExpirationDays})),
		"IfNotYou":         template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "welcome.if_not_you"})),
		"BestRegards":      template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "general.best_regards", TemplateData: map[string]string{"AppName": cfg.AppName}})),
	})

	if err != nil {
		return errors.Wrap(err, "execute welcome email template")
	}

	subject := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "welcome.subject", TemplateData: map[string]string{"AppName": cfg.AppName}})
	pageTitle := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "welcome.page_title"})
	pageHeader := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "welcome.page_header", TemplateData: map[string]string{"AppName": cfg.AppName}})
	if err := es.SendHtmlEmail(to, subject, htmlContent.String(), pageTitle, pageHeader, ""); err != nil {
		return errors.Wrap(err, "send email")
	}
	return nil
}

func (es *EmailSender) SendEmailChangeEmail(to, userName, confirmationLink string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}
	loc := locale.GetNewLocalizer("pl")
	// Fallback do EN jeśli brak tłumaczenia

	// Wygeneruj HTML z szablonu
	tmpl, err := template.ParseFS(templateFiles, "templates/email_change.html")
	if err != nil {
		return errors.Wrap(err, "parse email change email template")
	}

	var htmlContent bytes.Buffer
	err = tmpl.Execute(&htmlContent, map[string]interface{}{
		"Hello":                   template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "general.hello_user", TemplateData: map[string]string{"UserName": userName}})),
		"YouRequestedEmailChange": template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_change.you_requested_email_change"})),
		"ConfirmationLink":        confirmationLink,
		"PleaseConfirmNewEmail":   template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_change.please_confirm_new_email"})),
		"ConfirmNewEmail":         template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_change.confirm_new_email"})),
		"IfButtonFails":           template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "general.if_button_fails"})),
		"LinkExpiryInfo":          template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_change.link_expiry_info", TemplateData: map[string]int{"ExpiryDays": cfg.Register.ExpirationDays}, PluralCount: cfg.Register.ExpirationDays})),
		"IfNotYouEmailChange":     template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_change.if_not_you_email_change"})),
		"BestRegards":             template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "general.best_regards", TemplateData: map[string]string{"AppName": cfg.AppName}})),
	})

	if err != nil {
		return errors.Wrap(err, "execute welcome email template")
	}

	subject := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_change.subject", TemplateData: map[string]string{"AppName": cfg.AppName}})
	pageTitle := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_change.page_title"})
	pageHeader := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_change.page_header", TemplateData: map[string]string{"AppName": cfg.AppName}})
	if err := es.SendHtmlEmail(to, subject, htmlContent.String(), pageTitle, pageHeader, ""); err != nil {
		return errors.Wrap(err, "send email")
	}
	return nil
}

func (es *EmailSender) SendPasswordResetEmail(to, userName, resetLink string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}
	loc := locale.GetNewLocalizer("pl")

	// Wygeneruj HTML z szablonu
	tmpl, err := template.ParseFS(templateFiles, "templates/password_reset.html")
	if err != nil {
		return errors.Wrap(err, "parse password reset email template")
	}

	var htmlContent bytes.Buffer
	err = tmpl.Execute(&htmlContent, map[string]interface{}{
		"Hello":          template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "general.hello_user", TemplateData: map[string]string{"UserName": userName}})),
		"YouRequested":   template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset.you_requested"})),
		"ResetLink":      resetLink,
		"PleaseReset":    template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset.please_reset"})),
		"ResetPassword":  template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset.reset_password"})),
		"IfButtonFails":  template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "general.if_button_fails"})),
		"LinkExpiryInfo": template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset.link_expiry_info", TemplateData: map[string]int{"ExpiryDays": cfg.ResetPassword.ExpirationDays}, PluralCount: cfg.ResetPassword.ExpirationDays})),
		"IfNotYou":       template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset.if_not_you"})),
		"BestRegards":    template.HTML(loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "general.best_regards", TemplateData: map[string]string{"AppName": cfg.AppName}})),
	})

	if err != nil {
		return errors.Wrap(err, "execute password reset email template")
	}

	subject := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset.subject", TemplateData: map[string]string{"AppName": cfg.AppName}})
	pageTitle := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset.page_title"})
	pageHeader := loc.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset.page_header", TemplateData: map[string]string{"AppName": cfg.AppName}})
	if err := es.SendHtmlEmail(to, subject, htmlContent.String(), pageTitle, pageHeader, ""); err != nil {
		return errors.Wrap(err, "send email")
	}
	return nil
}
func (es *EmailSender) AddEmbeddedImageFromBytes(contentID, contentType, fileName string, data []byte) (string, error) {
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%d", fileName, time.Now().UnixNano()))
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return "", err
	}
	es.AddEmbedImage(contentID, contentType, tempPath, fileName)
	return tempPath, nil
}
