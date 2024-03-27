package testrepohandlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vvjke314/mkc-backend/internal/pkg/db"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

type ApplicationTest struct {
	ctx context.Context
	r   *db.Repo
	//srv *service.Service
}

func NewApplicationTest() *ApplicationTest {
	return &ApplicationTest{}
}

// Init
// Инициализирует тест-сервис
func (app *ApplicationTest) Init() error {
	app.ctx = context.Background()
	app.r = db.NewRepo()
	err := app.r.Init()
	if err != nil {
		return fmt.Errorf("[db.Init]: Can't initialize to database: %w", err)
	}

	//TO-DO: SERVICE INIT
	return nil
}

// Run
// Запускает тест-сервис с симуляцией работы сервиса
func (app *ApplicationTest) Run() error {
	// подключение к бд
	err := app.r.Connect()
	if err != nil {
		return fmt.Errorf("[db.Connect]: Can't connect to database: %w", err)
	}
	defer app.r.Close()

	// 1. Регистрация пользователя
	customer := ds.Customer{
		Id:         uuid.New(),
		FirstName:  "Vladimir",
		SecondName: "Abramov",
		Login:      "vvjkee",
		Password:   "bufybuff2002",
		Email:      "vvjkee@mail.ru",
		Type:       0,
	}
	err = app.r.SignUpCustomer(customer)
	if err != nil {
		return fmt.Errorf("[db.SignUpCustomer]: Can't signup customer: %w", err)
	}

	// 2. Регистрация администратора
	administrator := ds.Administrator{
		Id:       uuid.New(),
		Name:     "Polina",
		Email:    "polina.andronova@mail.ru",
		Password: "lyblyuVovu",
	}
	err = app.r.SignUpAdministrator(administrator)
	if err != nil {
		return fmt.Errorf("[db.SignUpAdministrator]: Can't signup admin: %w", err)
	}

	// 3. Создание пользователем проекта
	project := ds.Project{
		Id:           uuid.New(),
		OwnerId:      customer.Id,
		Name:         "NewProject",
		CreationDate: time.Now(),
	}
	err = app.r.CreateProject(project)
	if err != nil {
		return fmt.Errorf("[db.CreateProject]: Can't create project: %w", err)
	}

	// 3.1. Удаление пользователем проекта
	// err = app.r.DeleteProject(project.Id.String())
	// if err != nil {
	// 	return fmt.Errorf("[db.DeleteProject]: Can't delete project: %w", err)
	// }

	// 4. Добавления файла в проект
	file := ds.File{
		Id:             uuid.New(),
		ProjectId:      project.Id,
		Filename:       "New File",
		Extension:      "txt",
		Size:           200,
		FilePath:       project.Id.String(),
		UpdateDatetime: time.Now(),
	}
	err = app.r.CreateFile(file)
	if err != nil {
		return fmt.Errorf("[db.CreateFile]: Can't create file: %w", err)
	}

	// 5. Назначение администратора
	err = app.r.SetAdministrator(administrator.Id.String(), project.Id.String())
	if err != nil {
		return fmt.Errorf("[db.SetAdministator]: Can't set administrator to project: %w", err)
	}

	// 6. Назначение нового администратора
	admin2 := ds.Administrator{
		Id:       uuid.New(),
		Name:     "Miwa",
		Email:    "miwamiwa",
		Password: "lyblyuVovu",
	}
	err = app.r.SignUpAdministrator(admin2)
	if err != nil {
		return fmt.Errorf("[db.SignUpAdministrator]: Can't signup admin: %w", err)
	}

	err = app.r.SetAdministrator(admin2.Id.String(), project.Id.String())
	if err != nil {
		return fmt.Errorf("[db.SetAdministator]: Can't set administrator to project: %w", err)
	}

	// 7. Повышение статуса клиента
	err = app.r.UpgradeCustomerStatus(customer.Id.String(), 1)
	if err != nil {
		return fmt.Errorf("[db.UpgradeCustomerStatus]: Can't upgrade user status: %w", err)
	}

	// 8. Удаление файла из проекта
	file2 := ds.File{
		Id:             uuid.New(),
		ProjectId:      project.Id,
		Filename:       "Newest File",
		Extension:      "txt",
		Size:           150,
		FilePath:       project.Id.String(),
		UpdateDatetime: time.Now(),
	}
	err = app.r.CreateFile(file2)
	if err != nil {
		return fmt.Errorf("[db.CreateFile] %w", err)
	}
	err = app.r.DeleteFile(file2.Id.String())
	if err != nil {
		return fmt.Errorf("[db.DeleteFile] %w", err)
	}

	// 9. Изменение имени файла
	err = app.r.UpdateFileName(file.Id.String(), "Updated filename")
	if err != nil {
		return fmt.Errorf("[db.UpdateFileName] %w", err)
	}

	// 10. Считывание файла из БД
	var f ds.File
	err = app.r.GetFileById(file.Id.String(), &f)
	if err != nil {
		return fmt.Errorf("[db.GetFileById] %w", err)
	}

	// 11. Получение всех файлов в проекте
	file3 := ds.File{
		Id:             uuid.New(),
		ProjectId:      project.Id,
		Filename:       "File-3",
		Extension:      "txt",
		Size:           150,
		FilePath:       project.Id.String(),
		UpdateDatetime: time.Now(),
	}
	file4 := ds.File{
		Id:             uuid.New(),
		ProjectId:      project.Id,
		Filename:       "File-4",
		Extension:      "txt",
		Size:           225,
		FilePath:       project.Id.String(),
		UpdateDatetime: time.Now(),
	}
	err = app.r.CreateFile(file3)
	if err != nil {
		return fmt.Errorf("[db.CreateFile] %w", err)
	}
	err = app.r.CreateFile(file4)
	if err != nil {
		return fmt.Errorf("[db.CreateFile] %w", err)
	}
	// var files []ds.File
	// files, err = app.r.GetFiles(project.Id.String())
	// if err != nil {
	// 	return fmt.Errorf("[db.GetFiles] %w", err)
	// }

	// 12. Создание заметки в проекте
	note := ds.Note{
		Id:             uuid.New(),
		ProjectId:      project.Id,
		Title:          "first-note",
		Content:        "",
		UpdateDatetime: time.Now(),
		Deadline:       time.Date(2024, time.March, 23, 12, 50, 0, 0, time.Local),
	}
	err = app.r.CreateNote(note)
	if err != nil {
		return fmt.Errorf("[db.CreateNote] %w", err)
	}

	// 13. Удаление заметки из проекта
	// err = app.r.DeleteNote(note.Id.String())
	// if err != nil {
	// 	return fmt.Errorf("[db.DeleteNote] %w", err)
	// }

	// 14. Изменение имени заметки
	err = app.r.UpdateNoteName(note.Id.String(), "New-note-Name")
	if err != nil {
		return fmt.Errorf("[db.UpdateNoteName] %w", err)
	}

	// 15. Удаление всего проекта
	// err = app.r.DeleteProject(project.Id.String())
	// if err != nil {
	// 	return fmt.Errorf("[db.UpdateNoteName] %w", err)
	// }

	return nil
}