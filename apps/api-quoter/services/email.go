package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joledev/api-quoter/models"
)

type resendPayload struct {
	From    string `json:"from"`
	To      []string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
}

func sendEmail(to, subject, html string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY not set")
	}

	payload := resendPayload{
		From:    "JoleDev <noreply@joledev.com>",
		To:      []string{to},
		Subject: subject,
		HTML:    html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend API returned status %d", resp.StatusCode)
	}

	return nil
}

var planLabels = map[string]map[string]string{
	"fullPayment":   {"es": "Pago completo (-10%)", "en": "Full payment (-10%)"},
	"splitPayment":  {"es": "50% inicio / 50% entrega", "en": "50% upfront / 50% delivery"},
	"msi3":          {"es": "3 meses sin intereses", "en": "3 interest-free installments"},
	"msi6":          {"es": "6 meses sin intereses", "en": "6 interest-free installments"},
	"financing12":   {"es": "Financiamiento a 12 meses", "en": "12-month financing"},
	"saasMonthly":   {"es": "SaaS mensual", "en": "Monthly SaaS"},
	"annualLicense": {"es": "Licencia anual + servidor", "en": "Annual license + server"},
	"timeRetainer":  {"es": "Retainer por horas", "en": "Hourly retainer"},
	"payroll":       {"es": "Equivalente a nómina", "en": "Payroll-style"},
}

func getPlanLabel(key, lang string) string {
	if labels, ok := planLabels[key]; ok {
		if label, ok := labels[lang]; ok {
			return label
		}
	}
	return key
}

func formatSourceCode(include bool, lang string) string {
	if include {
		if lang == "en" {
			return "Yes"
		}
		return "Sí"
	}
	if lang == "en" {
		return "No"
	}
	return "No"
}

func formatCurrency(amount int, currency string) string {
	if currency == "USD" {
		return fmt.Sprintf("$%d USD", amount)
	}
	return fmt.Sprintf("$%d MXN", amount)
}

func SendQuoteNotification(q *models.QuoteRequest, quoteID string) error {
	contactEmail := os.Getenv("CONTACT_EMAIL")
	if contactEmail == "" {
		contactEmail = "joel@joledev.com"
	}

	projectTypes := strings.Join(q.ProjectTypes, ", ")
	features := strings.Join(q.Features, ", ")
	estimate := fmt.Sprintf("%s — %s", formatCurrency(q.EstimatedMin, q.Currency), formatCurrency(q.EstimatedMax, q.Currency))

	subject := fmt.Sprintf("Nueva cotización - %s - %s", q.Contact.Company, quoteID)
	if q.Contact.Company == "" {
		subject = fmt.Sprintf("Nueva cotización - %s - %s", q.Contact.Name, quoteID)
	}

	html := fmt.Sprintf(`<h2>Nueva cotización recibida: %s</h2>
<p><strong>Cliente:</strong> %s<br>
<strong>Email:</strong> %s<br>
<strong>Teléfono:</strong> %s<br>
<strong>Empresa:</strong> %s</p>
<p><strong>Proyectos:</strong> %s<br>
<strong>Funcionalidades:</strong> %s</p>
<p><strong>Tamaño:</strong> %s<br>
<strong>Estado:</strong> %s<br>
<strong>Plazo:</strong> %s<br>
<strong>Moneda:</strong> %s</p>
<p><strong>Presupuesto estimado:</strong> %s</p>
<p><strong>Plan de pago:</strong> %s<br>
<strong>Código fuente:</strong> %s</p>
<p><strong>Notas:</strong><br>%s</p>`,
		quoteID, q.Contact.Name, q.Contact.Email, q.Contact.Phone,
		q.Contact.Company, projectTypes, features,
		q.BusinessSize, q.CurrentState, q.Timeline, q.Currency,
		estimate, getPlanLabel(q.PaymentPlan, "es"), formatSourceCode(q.IncludeSourceCode, "es"),
		q.Contact.Notes)

	return sendEmail(contactEmail, subject, html)
}

func SendQuoteConfirmation(q *models.QuoteRequest, quoteID string) error {
	estimate := fmt.Sprintf("%s — %s", formatCurrency(q.EstimatedMin, q.Currency), formatCurrency(q.EstimatedMax, q.Currency))

	var subject, html string

	if q.Lang == "en" {
		subject = fmt.Sprintf("Your JoleDev quote - %s", quoteID)
		html = fmt.Sprintf(`<p>Hi %s,</p>
<p>Thank you for your interest in my services. I've received your quote request and will review it in detail.</p>
<p>I'll contact you within the next 24 hours to discuss your project and prepare a personalized proposal.</p>
<p><strong>Summary:</strong><br>
Projects: %s<br>
Estimated budget: %s<br>
Payment plan: %s</p>
<p>If you have any questions, feel free to reach out at joel@joledev.com.</p>
<p>Best regards,<br>Joel López Verdugo<br>JoleDev — Technology tailored to your business</p>`,
			q.Contact.Name, strings.Join(q.ProjectTypes, ", "), estimate,
			getPlanLabel(q.PaymentPlan, "en"))
	} else {
		subject = fmt.Sprintf("Tu cotización JoleDev - %s", quoteID)
		html = fmt.Sprintf(`<p>Hola %s,</p>
<p>Gracias por tu interés en mis servicios. He recibido tu solicitud de cotización y la revisaré en detalle.</p>
<p>Me pondré en contacto contigo en las próximas 24 horas para discutir tu proyecto y preparar una propuesta personalizada.</p>
<p><strong>Resumen:</strong><br>
Proyectos: %s<br>
Presupuesto estimado: %s<br>
Plan seleccionado: %s</p>
<p>Si tienes alguna pregunta, escríbeme a joel@joledev.com.</p>
<p>Saludos,<br>Joel López Verdugo<br>JoleDev — Desarrollo a la medida de tu negocio</p>`,
			q.Contact.Name, strings.Join(q.ProjectTypes, ", "), estimate,
			getPlanLabel(q.PaymentPlan, "es"))
	}

	return sendEmail(q.Contact.Email, subject, html)
}
