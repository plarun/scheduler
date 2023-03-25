package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/eventserver/db"
	"github.com/plarun/scheduler/internal/eventserver/db/query"
	"github.com/plarun/scheduler/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newTaskEntity(vt *proto.ValidatedTaskEntity) *task.TaskEntity {
	tsk := task.NewTaskEntity(vt.Name)

	if empty := false; vt.Command.Flag != proto.NullableFlag_NotAvailable {
		if vt.Command.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldCommand(vt.Command.Value, empty)
	}

	if empty := false; vt.Condition.Flag != proto.NullableFlag_NotAvailable {
		if vt.Condition.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldCondition(vt.Condition.Value, empty)
	}
	//
	if empty := false; vt.ErrLogFile.Flag != proto.NullableFlag_NotAvailable {
		if vt.ErrLogFile.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldErrLogFile(vt.ErrLogFile.Value, empty)
	}

	if empty := false; vt.Label.Flag != proto.NullableFlag_NotAvailable {
		if vt.Label.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldLabel(vt.Label.Value, empty)
	}

	if empty := false; vt.Machine.Flag != proto.NullableFlag_NotAvailable {
		if vt.Machine.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldMachine(vt.Machine.Value, empty)
	}

	if empty := false; vt.OutLogFile.Flag != proto.NullableFlag_NotAvailable {
		if vt.OutLogFile.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldOutLogFile(vt.OutLogFile.Value, empty)
	}

	if empty := false; vt.Parent.Flag != proto.NullableFlag_NotAvailable {
		if vt.Parent.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldParent(vt.Parent.Value, empty)
	}

	if empty := false; vt.Priority.Flag != proto.NullableFlag_NotAvailable {
		if vt.Priority.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldPriority(vt.Priority.Value, empty)
	}

	if empty := false; vt.Profile.Flag != proto.NullableFlag_NotAvailable {
		if vt.Profile.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldProfile(vt.Profile.Value, empty)
	}

	if empty := false; vt.RunDays.Flag != proto.NullableFlag_NotAvailable {
		if vt.RunDays.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldRunDays(vt.RunDays.Value, empty)
	}

	if empty := false; vt.RunWindow.Flag != proto.NullableFlag_NotAvailable {
		if vt.RunWindow.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldRunWindow(vt.RunWindow.Value.Start, vt.RunWindow.Value.End, empty)
	}

	if empty := false; vt.StartMins.Flag != proto.NullableFlag_NotAvailable {
		mins := make([]uint8, 0)
		if vt.StartMins.Flag == proto.NullableFlag_Empty {
			empty = true
		} else {
			for _, m := range vt.StartMins.Value {
				mins = append(mins, uint8(m))
			}
		}
		tsk.SetFieldStartMins(mins, empty)
	}

	if empty := false; vt.StartTimes.Flag != proto.NullableFlag_NotAvailable {
		if vt.StartTimes.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldStartTimes(vt.StartTimes.Value, empty)
	}

	if empty := false; vt.Type.Flag != proto.NullableFlag_NotAvailable {
		if vt.Type.Flag == proto.NullableFlag_Empty {
			empty = true
		}
		tsk.SetFieldType(vt.Type.Value, empty)
	}

	return tsk
}

func ProcessTaskActions(ctx context.Context, tasks []*proto.ValidatedTaskEntity) error {
	// database transaction helps to execute list
	// of queries after all queries are success
	// then commits otherwise rollbacks
	database := db.GetDatabase()
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error: ProcessTaskActions: unable to create db transaction: %v", err)
	}

	log.Println("Starting to process task actions. New sql transaction begins")

	for _, t := range tasks {
		tsk := newTaskEntity(t)

		err := func(action task.Action) error {
			switch task.Action(t.Action) {
			case task.ActionInsert:
				if err = query.InsertTask(tx, tsk); err != nil {
					return err
				}
			case task.ActionUpdate:
				if err = query.UpdateTask(tx, tsk); err != nil {
					return err
				}
			case task.ActionDelete:
				if err = query.DeleteTask(tx, tsk); err != nil {
					return err
				}
			}
			return nil
		}(task.Action(t.Action))

		if err != nil {
			log.Printf("Error: ProcessTaskActions: %v", err)

			// rollback the transaction
			if err := tx.Rollback(); err != nil {
				log.Printf("Error: ProcessTaskActions: failed to rollback sql transaction: %v", err)
				return status.Error(codes.Internal, "internal error")
			}

			log.Printf("Error: ProcessTaskActions: successfully rollbacked sql transaction: %v", err)

			// failed to process JIL
			for errors.Unwrap(err) != nil {
				err = errors.Unwrap(err)
			}
			return err
		}

	}

	// commit the transaction
	if err := tx.Commit(); err != nil {
		// logger.Log.WithFields(logrus.Fields{"prefix": "JilService"}).Errorf("Submit: failed to commit: %v", err)
		log.Printf("Error: ProcessTaskActions: failed to commit sql transaction: %v", err)
		return status.Error(codes.Internal, "internal error")
	}

	return nil
}

func GetTaskDefinition(ctx context.Context, name string) (*proto.TaskDefinition, error) {
	if res, err := query.GetTaskDetails(name); err != nil {
		return nil, fmt.Errorf("GetTaskDefinition: %w", err)
	} else {
		return res, nil
	}
}

func GetTaskLatestStatus(ctx context.Context, name string) (*proto.TaskRunStatus, error) {
	if res, err := query.GetLatestStatus(name); err != nil {
		return nil, fmt.Errorf("GetTaskLatestStatus: %w", err)
	} else {
		return res, nil
	}
}

func GetTaskRuns(ctx context.Context, name string, n int32, date string) ([]*proto.TaskRunStatus, error) {
	if res, err := query.GetRuns(name, n, date); err != nil {
		return nil, fmt.Errorf("GetTaskLatestStatus: %w", err)
	} else {
		return res, nil
	}
}
