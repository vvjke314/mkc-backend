package db

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

func (r *Repo) ProccessNotes(durationBefore time.Duration, query string) error {
	// Определение временной метки для сравнения
	compareTime := time.Now().Add(durationBefore)
	// Запрос данных из базы данных
	notes := []ds.Note{}
	rows, err := r.pool.Query(r.ctx, query, compareTime)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Query] can't exec query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n ds.Note
		if err := rows.Scan(&n.Id, &n.ProjectId, &n.Title, &n.Content, &n.UpdateDatetime, &n.Deadline, &n.Overdue); err != nil {
			return fmt.Errorf("[pgx.Rows.Scan] can't scan data: %w", err)
		}
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("[pgx.Rows.Err] error occured when iterating data: %w", err)
	}

	for _, note := range notes {
		if time.Now().Add(+3*time.Hour).After(note.Deadline) && note.Overdue == 0 {
			err = r.NoteOverdue(note.Id.String())
			if err != nil {
				return fmt.Errorf("[repo.NoteOverdue]: %w", err)
			}
		}
		customers, err := r.GetParticipants(note.ProjectId.String())
		if err != nil {
			return fmt.Errorf("[repo.Getparticipants]: %w", err)
		}
		var emails []string
		for _, cust := range customers {
			emails = append(emails, cust.Email)
		}
		err = sendEmailNotification(note, emails)
		if err != nil {
			return fmt.Errorf("[sendEmailNotification]: %w", err)
		}
	}

	return nil
}

func sendEmailNotification(note ds.Note, emails []string) error {
	// Настройки SMTP сервера
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := "mknnotifyer@gmail.com"
	senderPassword := "uxluuqehymoyqjnx"

	// Формирование сообщения
	message := "Subject: Крайний срок заметки\n\n"
	message += "Здравствуйте,\n\n"
	message += "Это уведомление говорит о приближении дедлайна заметки:\n\n"
	message += "Название: " + note.Title + "\n"
	message += "Дедлайн: " + note.Deadline.Format("2006-01-02 15:04:05") + "\n\n"
	message += "Наилучших пожеланий,\n MKC team"

	// Аутентификация на SMTP сервере и отправка сообщения
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, emails, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
