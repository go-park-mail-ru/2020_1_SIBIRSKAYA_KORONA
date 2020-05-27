package usecase

import (
	"os"
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
}

func CreateUseCase(labelRepo_ label.Repository, columnRepo_ column.Repository, taskRepo_ task.Repository, commentRepo_ comment.Repository,
	chlistRepo_ checklist.Repository, boardRepo_ board.Repository, userRepo_ user.Repository) template.Usecase {
	return &TemplateUsecase{
		labelRepo:   labelRepo_,
		columnRepo:  columnRepo_,
		taskRepo:    taskRepo_,
		commentRepo: commentRepo_,
		chlistRepo:  chlistRepo_,
		boardRepo:   boardRepo_,
		userRepo:    userRepo_,
	}
}

// Здесь все плохо и это нужно будет переделать

func (tmplUsecase *TemplateUsecase) Create(uid uint, tmpl *models.Template) error {
	_ = ReadTemplateByVariant(tmpl)
	board := ParseBoard()

	_, err := tmplUsecase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = tmplUsecase.boardRepo.Create(uid, board)
	if err != nil {
		logger.Error(err)
		return err
	}

	//board.Admins = append(board.Admins, *usr)

	labelsMap, err := tmplUsecase.CreateLabels(board.ID)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = tmplUsecase.CreateColumnsAndTask(board.ID, labelsMap)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// читаем файл
func ReadTemplateByVariant(tmpl *models.Template) error {
	tmplPath, exists := os.LookupEnv("BOARD_TEMPLATES_PATH")
	if !exists {
		logger.Fatal("BOARD_TEMPLATES_PATH environment variable not exist")
	}

	switch tmpl.Variant {
	case "week_plan":
		{
			viper.SetConfigName("week_plan")
			viper.AddConfigPath(tmplPath)
			err := viper.MergeInConfig()
			if err != nil {
				logger.Error(err)
			}
		}
	default:
		{
			return errors.ErrDbBadOperation
		}
	}

	return nil
}

func ParseBoard() *models.Board {
	boardName := viper.GetString("board.name")
	return &models.Board{Name: boardName}
}

func (tmplUsecase *TemplateUsecase) CreateLabels(bid uint) (map[string]uint, error) {
	labelsMap := viper.GetStringMapString("labels")
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

// Один большой костыль
func (tmplUsecase *TemplateUsecase) CreateColumnsAndTask(bid uint, labelsMap map[string]uint) error {
	columns := viper.GetStringMap("columns")

	for columnName, tasksNode := range columns {
		var column models.Column
		column.Bid = bid
		column.Name = columnName

		tasks := tasksNode.(map[string]interface{})

		pos := tasks["position"].([]interface{})
		column.Pos = pos[0].(float64)
		delete(tasks, "position")

		err := tmplUsecase.columnRepo.Create(&column)
		if err != nil {
			logger.Error(err)
			return err
		}

		for taskName, taskInter := range tasks {
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
