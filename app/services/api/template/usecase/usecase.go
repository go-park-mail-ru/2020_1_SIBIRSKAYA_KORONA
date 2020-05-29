package usecase

import (
	"strings"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/spf13/viper"

	// нужны почти все сущности
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/template"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user"
)

type TemplateUsecase struct {
	boardRepo   board.Repository
	columnRepo  column.Repository
	taskRepo    task.Repository
	commentRepo comment.Repository
	labelRepo   label.Repository
	chlistRepo  checklist.Repository
	userRepo    user.Repository

	templateCache map[string]*viper.Viper
}

func CreateUseCase(labelRepo_ label.Repository, columnRepo_ column.Repository, taskRepo_ task.Repository, commentRepo_ comment.Repository,
	chlistRepo_ checklist.Repository, boardRepo_ board.Repository, userRepo_ user.Repository, templateCache_ map[string]*viper.Viper) template.Usecase {
	return &TemplateUsecase{
		labelRepo:     labelRepo_,
		columnRepo:    columnRepo_,
		taskRepo:      taskRepo_,
		commentRepo:   commentRepo_,
		chlistRepo:    chlistRepo_,
		boardRepo:     boardRepo_,
		userRepo:      userRepo_,
		templateCache: templateCache_,
	}
}

// Здесь все плохо и это нужно будет переделать
func (tmplUsecase *TemplateUsecase) Create(uid uint, tmpl *models.Template) (*models.Board, error) {
	reader, exist := tmplUsecase.templateCache[tmpl.Variant]
	if !exist {
		return nil, errors.ErrDbBadOperation
	}

	board := ParseBoard(reader)
	board.InviteLink = tmplUsecase.boardRepo.GenerateInviteLink(32)

	_, err := tmplUsecase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	err = tmplUsecase.boardRepo.Create(uid, board)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	//board.Admins = append(board.Admins, *usr)

	labelsMap, err := tmplUsecase.CreateLabels(reader, board.ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	err = tmplUsecase.CreateColumnsAndTask(reader, board.ID, labelsMap)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return board, nil
}

func ParseBoard(reader *viper.Viper) *models.Board {
	boardName := reader.GetString("board.name")
	return &models.Board{Name: boardName}
}

func (tmplUsecase *TemplateUsecase) CreateLabels(reader *viper.Viper, bid uint) (map[string]uint, error) {
	labelsMap := reader.GetStringMapString("labels")
	labelNameToID := make(map[string]uint, 0)

	for labelName, labelColor := range labelsMap {
		label := models.Label{Name: labelName, Color: labelColor, Bid: bid}
		err := tmplUsecase.labelRepo.Create(&label)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		labelNameToID[label.Name] = label.ID
	}

	return labelNameToID, nil
}

// Один большой костыль, переделать формат файла и подобрать отдельный парсер для yaml
func (tmplUsecase *TemplateUsecase) CreateColumnsAndTask(reader *viper.Viper, bid uint, labelsMap map[string]uint) error {
	columns := reader.GetStringMap("columns")

	for columnName, tasksNode := range columns {
		var column models.Column
		column.Bid = bid
		column.Name = columnName

		tasks := tasksNode.(map[string]interface{})

		pos := tasks["position"].([]interface{})
		column.Pos = pos[0].(float64)

		err := tmplUsecase.columnRepo.Create(&column)
		if err != nil {
			logger.Error(err)
			return err
		}

		for taskName, taskInter := range tasks {
			// Костыль
			if taskName == "position" {
				continue
			}

			var task models.Task
			task.Name = taskName
			task.Cid = column.ID
			taskNode := taskInter.(map[string]interface{})
			taskPosition := taskNode["position"].([]interface{})
			task.Pos = taskPosition[0].(float64)

			err := tmplUsecase.taskRepo.Create(&task)
			if err != nil {
				logger.Error(err)
				return err
			}
			labels := taskNode["labels"].([]interface{})
			for _, labelName := range labels {
				err := tmplUsecase.labelRepo.AddLabelOnTask(labelsMap[strings.ToLower(labelName.(string))], task.ID)
				if err != nil {
					logger.Error(err)
					return err
				}
			}
		}

	}

	return nil
}
