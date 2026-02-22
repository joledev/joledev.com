package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joledev/api-scheduler/models"
)

type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
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

func getAPIBaseURL() string {
	url := os.Getenv("API_BASE_URL")
	if url == "" {
		url = "http://localhost:8082"
	}
	return strings.TrimRight(url, "/")
}

func formatDate(date, lang string) string {
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return date
	}

	months := map[string]map[string]string{
		"01": {"es": "enero", "en": "January"},
		"02": {"es": "febrero", "en": "February"},
		"03": {"es": "marzo", "en": "March"},
		"04": {"es": "abril", "en": "April"},
		"05": {"es": "mayo", "en": "May"},
		"06": {"es": "junio", "en": "June"},
		"07": {"es": "julio", "en": "July"},
		"08": {"es": "agosto", "en": "August"},
		"09": {"es": "septiembre", "en": "September"},
		"10": {"es": "octubre", "en": "October"},
		"11": {"es": "noviembre", "en": "November"},
		"12": {"es": "diciembre", "en": "December"},
	}

	month := months[parts[1]][lang]
	day := strings.TrimLeft(parts[2], "0")

	if lang == "en" {
		return fmt.Sprintf("%s %s, %s", month, day, parts[0])
	}
	return fmt.Sprintf("%s de %s, %s", day, month, parts[0])
}

func formatTime(t string) string {
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return t
	}
	hour := parts[0]
	min := parts[1]

	h := 0
	fmt.Sscanf(hour, "%d", &h)

	if h == 0 {
		return fmt.Sprintf("12:%s AM", min)
	} else if h < 12 {
		return fmt.Sprintf("%d:%s AM", h, min)
	} else if h == 12 {
		return fmt.Sprintf("12:%s PM", min)
	}
	return fmt.Sprintf("%d:%s PM", h-12, min)
}

func meetingTypeLabel(mt, lang string) string {
	if mt == "presencial" {
		if lang == "en" {
			return "In-person"
		}
		return "Presencial"
	}
	if lang == "en" {
		return "Video call"
	}
	return "Videollamada"
}

// SendAdminPendingNotification sends an email to the admin when a new booking request comes in.
// Includes Confirm and Reject buttons with secure token links.
func SendAdminPendingNotification(b *models.Booking) error {
	contactEmail := os.Getenv("CONTACT_EMAIL")
	if contactEmail == "" {
		contactEmail = "joel@joledev.com"
	}

	baseURL := getAPIBaseURL()
	confirmURL := fmt.Sprintf("%s/scheduler/bookings/confirm?token=%s", baseURL, b.ConfirmToken)
	rejectURL := fmt.Sprintf("%s/scheduler/bookings/reject?token=%s", baseURL, b.RejectToken)

	dateStr := formatDate(b.Date, "es")
	timeStr := fmt.Sprintf("%s - %s", formatTime(b.StartTime), formatTime(b.EndTime))
	mtLabel := meetingTypeLabel(b.MeetingType, "es")

	addressLine := ""
	if b.MeetingType == "presencial" && b.ClientAddress != "" {
		addressLine = fmt.Sprintf(`<p><strong>Dirección:</strong> %s</p>`, b.ClientAddress)
	}

	notesLine := ""
	if b.Notes != "" {
		notesLine = fmt.Sprintf(`<p><strong>Notas:</strong><br>"%s"</p>`, b.Notes)
	}

	tzLine := ""
	if b.ClientTimezone != "" {
		tzLine = fmt.Sprintf(`<p><strong>Zona horaria del cliente:</strong> %s</p>`, b.ClientTimezone)
	}

	subject := fmt.Sprintf("Nueva solicitud de reunión - %s", b.BookingID)
	html := fmt.Sprintf(`<h2>Nueva solicitud de reunión: %s</h2>
<p><strong>Estado:</strong> <span style="color:#f59e0b;font-weight:bold">Pendiente de confirmación</span></p>
<p><strong>Cliente:</strong> %s<br>
<strong>Email:</strong> %s<br>
<strong>Teléfono:</strong> %s<br>
<strong>Empresa:</strong> %s</p>
<p><strong>Tipo:</strong> %s<br>
<strong>Fecha:</strong> %s<br>
<strong>Hora:</strong> %s</p>
%s
%s
%s
<div style="margin-top:1.5rem">
<a href="%s" style="display:inline-block;padding:12px 24px;background:#22c55e;color:#fff;text-decoration:none;border-radius:8px;font-weight:bold;margin-right:12px">Confirmar</a>
<a href="%s" style="display:inline-block;padding:12px 24px;background:#ef4444;color:#fff;text-decoration:none;border-radius:8px;font-weight:bold">Rechazar</a>
</div>`,
		b.BookingID, b.ClientName, b.ClientEmail, b.ClientPhone,
		b.ClientCompany, mtLabel, dateStr, timeStr, tzLine, addressLine, notesLine,
		confirmURL, rejectURL)

	return sendEmail(contactEmail, subject, html)
}

// SendClientPendingNotification notifies the client that their request was received and is pending.
func SendClientPendingNotification(b *models.Booking) error {
	lang := b.Lang
	if lang == "" {
		lang = "es"
	}

	dateStr := formatDate(b.Date, lang)
	timeStr := fmt.Sprintf("%s - %s", formatTime(b.StartTime), formatTime(b.EndTime))
	mtLabel := meetingTypeLabel(b.MeetingType, lang)

	var subject, html string

	if lang == "en" {
		subject = fmt.Sprintf("Meeting request received - JoleDev - %s", b.BookingID)
		html = fmt.Sprintf(`<p>Hi %s,</p>
<p>Your meeting request has been received and is <strong>pending confirmation</strong>.</p>
<p>📅 <strong>Date:</strong> %s<br>
🕐 <strong>Time:</strong> %s<br>
📍 <strong>Type:</strong> %s</p>
<p>I'll review your request and you'll receive another email once it's confirmed or if there's any issue.</p>
<p>Best regards,<br>Joel López Verdugo<br>JoleDev</p>`,
			b.ClientName, dateStr, timeStr, mtLabel)
	} else {
		subject = fmt.Sprintf("Solicitud de reunión recibida - JoleDev - %s", b.BookingID)
		html = fmt.Sprintf(`<p>Hola %s,</p>
<p>Tu solicitud de reunión ha sido recibida y está <strong>pendiente de confirmación</strong>.</p>
<p>📅 <strong>Fecha:</strong> %s<br>
🕐 <strong>Hora:</strong> %s<br>
📍 <strong>Tipo:</strong> %s</p>
<p>Revisaré tu solicitud y recibirás otro correo cuando sea confirmada o si hay algún inconveniente.</p>
<p>Saludos,<br>Joel López Verdugo<br>JoleDev</p>`,
			b.ClientName, dateStr, timeStr, mtLabel)
	}

	return sendEmail(b.ClientEmail, subject, html)
}

// SendBookingConfirmation sends a confirmation email to the client when the admin approves.
func SendBookingConfirmation(b *models.Booking) error {
	lang := b.Lang
	if lang == "" {
		lang = "es"
	}

	dateStr := formatDate(b.Date, lang)
	timeStr := fmt.Sprintf("%s - %s", formatTime(b.StartTime), formatTime(b.EndTime))
	mtLabel := meetingTypeLabel(b.MeetingType, lang)

	var subject, html string

	if lang == "en" {
		subject = fmt.Sprintf("Meeting confirmed - JoleDev - %s", b.BookingID)
		locationLine := ""
		if b.MeetingType == "presencial" && b.ClientAddress != "" {
			locationLine = fmt.Sprintf(`<p>📍 <strong>Address:</strong> %s</p>
<p>I'll be at your office at the indicated time.</p>`, b.ClientAddress)
		} else {
			locationLine = `<p>I'll send you the video call link by email before the meeting.</p>`
		}

		html = fmt.Sprintf(`<p>Hi %s,</p>
<p>Your meeting has been <strong style="color:#22c55e">confirmed</strong>!</p>
<p>📅 <strong>Date:</strong> %s<br>
🕐 <strong>Time:</strong> %s<br>
📍 <strong>Type:</strong> %s</p>
%s
<p>If you need to reschedule, contact me at joel@joledev.com.</p>
<p>Best regards,<br>Joel López Verdugo<br>JoleDev</p>`,
			b.ClientName, dateStr, timeStr, mtLabel, locationLine)
	} else {
		subject = fmt.Sprintf("Reunión confirmada - JoleDev - %s", b.BookingID)
		locationLine := ""
		if b.MeetingType == "presencial" && b.ClientAddress != "" {
			locationLine = fmt.Sprintf(`<p>📌 <strong>Dirección:</strong> %s</p>
<p>Me presentaré en tu oficina a la hora indicada.</p>`, b.ClientAddress)
		} else {
			locationLine = `<p>Te enviaré el link de la videollamada por email antes de la reunión.</p>`
		}

		html = fmt.Sprintf(`<p>Hola %s,</p>
<p>Tu reunión ha sido <strong style="color:#22c55e">confirmada</strong>!</p>
<p>📅 <strong>Fecha:</strong> %s<br>
🕐 <strong>Hora:</strong> %s<br>
📍 <strong>Tipo:</strong> %s</p>
%s
<p>Si necesitas reprogramar, contáctame a joel@joledev.com.</p>
<p>Saludos,<br>Joel López Verdugo<br>JoleDev</p>`,
			b.ClientName, dateStr, timeStr, mtLabel, locationLine)
	}

	return sendEmail(b.ClientEmail, subject, html)
}

// SendBookingRejection notifies the client that their booking was not approved.
func SendBookingRejection(b *models.Booking) error {
	lang := b.Lang
	if lang == "" {
		lang = "es"
	}

	dateStr := formatDate(b.Date, lang)
	timeStr := fmt.Sprintf("%s - %s", formatTime(b.StartTime), formatTime(b.EndTime))

	var subject, html string

	if lang == "en" {
		subject = fmt.Sprintf("Meeting request not available - JoleDev - %s", b.BookingID)
		html = fmt.Sprintf(`<p>Hi %s,</p>
<p>Unfortunately, I wasn't able to confirm your meeting scheduled for <strong>%s</strong> at <strong>%s</strong>.</p>
<p>This could be due to a scheduling conflict. Please feel free to select a different time at <a href="https://joledev.com/en/schedule">joledev.com/en/schedule</a>.</p>
<p>Sorry for the inconvenience!</p>
<p>Best regards,<br>Joel López Verdugo<br>JoleDev</p>`,
			b.ClientName, dateStr, timeStr)
	} else {
		subject = fmt.Sprintf("Solicitud de reunión no disponible - JoleDev - %s", b.BookingID)
		html = fmt.Sprintf(`<p>Hola %s,</p>
<p>Lamentablemente no pude confirmar tu reunión programada para el <strong>%s</strong> a las <strong>%s</strong>.</p>
<p>Esto puede deberse a un conflicto de horario. Por favor selecciona otro horario en <a href="https://joledev.com/es/agendar">joledev.com/es/agendar</a>.</p>
<p>Disculpa las molestias.</p>
<p>Saludos,<br>Joel López Verdugo<br>JoleDev</p>`,
			b.ClientName, dateStr, timeStr)
	}

	return sendEmail(b.ClientEmail, subject, html)
}

// SendBookingCancellation notifies the client that their booking was cancelled.
func SendBookingCancellation(b *models.Booking) error {
	lang := b.Lang
	if lang == "" {
		lang = "es"
	}

	dateStr := formatDate(b.Date, lang)
	timeStr := fmt.Sprintf("%s - %s", formatTime(b.StartTime), formatTime(b.EndTime))

	var subject, html string

	if lang == "en" {
		subject = fmt.Sprintf("Meeting cancelled - JoleDev - %s", b.BookingID)
		html = fmt.Sprintf(`<p>Hi %s,</p>
<p>Your meeting scheduled for <strong>%s</strong> at <strong>%s</strong> has been cancelled.</p>
<p>If you'd like to reschedule, visit <a href="https://joledev.com/en/schedule">joledev.com/en/schedule</a>.</p>
<p>Best regards,<br>Joel López Verdugo<br>JoleDev</p>`,
			b.ClientName, dateStr, timeStr)
	} else {
		subject = fmt.Sprintf("Reunión cancelada - JoleDev - %s", b.BookingID)
		html = fmt.Sprintf(`<p>Hola %s,</p>
<p>Tu reunión programada para el <strong>%s</strong> a las <strong>%s</strong> ha sido cancelada.</p>
<p>Si deseas reagendar, visita <a href="https://joledev.com/es/agendar">joledev.com/es/agendar</a>.</p>
<p>Saludos,<br>Joel López Verdugo<br>JoleDev</p>`,
			b.ClientName, dateStr, timeStr)
	}

	return sendEmail(b.ClientEmail, subject, html)
}
