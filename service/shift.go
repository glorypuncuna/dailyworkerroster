package service

import (
	"context"
	"dailyworkerroster/model"
	"dailyworkerroster/repository"
	"fmt"
	"log"

	"github.com/spf13/cast"
)

type ShiftServiceItf interface {
	// // Worker
	GetAssignedShifts(ctx context.Context) (*model.ListShiftDetail, error)
	GetAvailableShifts(ctx context.Context, workerID int64) ([]*model.ShiftStatus, error)
	RequestShift(ctx context.Context, shiftID, workerID int64) error
	GetAllRequestedShift(ctx context.Context, workerID int64) ([]*model.ShiftStatus, error)

	// // Admin
	CreateShift(ctx context.Context, shift *model.Shift) (int64, error)
	UpdateShift(ctx context.Context, shift *model.Shift) error
	DeleteShift(ctx context.Context, shiftID int64) error
	GetAllShiftRequests(ctx context.Context, queryParam model.WorkerShiftDetailQuery) ([]model.WorkerShiftDetail, error)
	ApproveShiftRequest(ctx context.Context, shiftID, workerID int64) error
	RejectShiftRequest(ctx context.Context, shiftID, workerID int64) error
	GetShiftsByDay(ctx context.Context, date string) ([]*model.ShiftStatus, error)
}

type ShiftService struct {
	ShiftRepo       repository.ShiftRepoItf
	WorkerShiftRepo repository.WorkerShiftRepoItf
}

func NewShiftService(
	shiftRepo repository.ShiftRepoItf,
	workerShiftRepo repository.WorkerShiftRepoItf) ShiftServiceItf {
	return &ShiftService{
		ShiftRepo:       shiftRepo,
		WorkerShiftRepo: workerShiftRepo,
	}
}

func (s *ShiftService) GetAssignedShifts(ctx context.Context) (*model.ListShiftDetail, error) {
	funcName := "/service/shift/GetAssignedShifts"

	listShiftResp := &model.ListShiftDetail{
		UserAccountID: cast.ToInt64(ctx.Value("user_account_id")),
		Name:          cast.ToString(ctx.Value("name")),
	}

	shiftList := make([]model.WorkerShiftDetail, 0)

	// Call the repository method to get the assigned shifts
	status := model.WORKER_SHIFT_APPROVED
	workerShift, err := s.WorkerShiftRepo.GetWorkerShiftListByFilter(&listShiftResp.UserAccountID, &status)
	if err != nil {
		log.Printf("%s: GetWorkerShiftByAccountID error: %v", funcName, err)
		return nil, err
	}

	shiftIDs := make([]int64, 0)
	workerShiftMap := make(map[int64]model.WorkerShift)
	for _, shift := range workerShift {
		shiftIDs = append(shiftIDs, shift.ShiftID)
		workerShiftMap[shift.ShiftID] = shift
	}

	shift, err := s.ShiftRepo.GetShiftsByIDs(shiftIDs)
	if err != nil {
		log.Printf("%s: GetShiftsByIDs error: %v", funcName, err)
		return nil, err
	}

	for _, s := range shift {
		shiftDetail := &model.WorkerShiftDetail{
			ShiftID:        s.ID,
			Date:           s.Date,
			StartTime:      s.StartTime,
			EndTime:        s.EndTime,
			RoleAssignment: s.RoleAssignment,
			Location:       s.Location,
			IsAvailable:    s.IsAvailable,
		}

		if workerShift, ok := workerShiftMap[s.ID]; ok {
			shiftDetail.ID = workerShift.ID
			shiftDetail.ApprovedBy = workerShift.ApprovedBy
			shiftDetail.Status = workerShift.Status
		}

		shiftList = append(shiftList, *shiftDetail)
	}

	listShiftResp.ShiftDetails = shiftList

	return listShiftResp, nil
}

func (s *ShiftService) GetAvailableShifts(ctx context.Context, workerID int64) ([]*model.ShiftStatus, error) {
	funcName := "/service/shift/GetAvailableShifts"

	var availableShiftStatus []*model.ShiftStatus

	isAvailable := true
	availableShift, err := s.ShiftRepo.GetListShifts(model.ShiftListQuery{
		IsAvailable: &isAvailable,
	})
	if err != nil {
		log.Printf("%s: GetListShifts error: %v", funcName, err)
		return nil, err
	}

	statusPending := model.WORKER_SHIFT_PENDING
	statusApproved := model.WORKER_SHIFT_APPROVED
	statusRejected := model.WORKER_SHIFT_REJECTED

	var shiftStatusMap = make(map[int64]string)
	shiftPending, err := s.WorkerShiftRepo.GetWorkerShiftListByFilter(&workerID, &statusPending)
	if err != nil {
		log.Printf("%s: GetWorkerShiftListByFilter error: %v", funcName, err)
		return nil, err
	}

	for _, shift := range shiftPending {
		shiftStatusMap[shift.ShiftID] = shift.Status
	}

	shiftApproved, err := s.WorkerShiftRepo.GetWorkerShiftListByFilter(&workerID, &statusApproved)
	if err != nil {
		log.Printf("%s: GetWorkerShiftListByFilter error: %v", funcName, err)
		return nil, err
	}

	for _, shift := range shiftApproved {
		shiftStatusMap[shift.ShiftID] = shift.Status
	}

	shiftRejected, err := s.WorkerShiftRepo.GetWorkerShiftListByFilter(&workerID, &statusRejected)
	if err != nil {
		log.Printf("%s: GetWorkerShiftListByFilter error: %v", funcName, err)
		return nil, err
	}

	for _, shift := range shiftRejected {
		shiftStatusMap[shift.ShiftID] = shift.Status
	}

	for _, shift := range availableShift {
		shiftStatus := &model.ShiftStatus{
			ID:             shift.ID,
			Date:           shift.Date,
			StartTime:      shift.StartTime,
			EndTime:        shift.EndTime,
			RoleAssignment: shift.RoleAssignment,
			Location:       shift.Location,
			IsAvailable:    shift.IsAvailable,
		}

		if status, ok := shiftStatusMap[shift.ID]; ok {
			shiftStatus.StatusWorker = status
		} else {
			shiftStatus.StatusWorker = ""
		}

		availableShiftStatus = append(availableShiftStatus, shiftStatus)
	}

	return availableShiftStatus, nil
}

func (s *ShiftService) RequestShift(ctx context.Context, shiftID, workerID int64) error {
	funcName := "/service/shift/RequestShift"

	shift, err := s.ShiftRepo.GetShiftByID(shiftID)
	if err != nil {
		log.Printf("%s: GetShiftByID error: %v", funcName, err)
		return err
	}
	if !shift.IsAvailable {
		log.Printf("%s: Shift is not available", funcName)
		return fmt.Errorf("shift is not available")
	}

	statuses := []string{model.WORKER_SHIFT_PENDING, model.WORKER_SHIFT_APPROVED}
	for _, status := range statuses {
		shifts, err := s.WorkerShiftRepo.GetWorkerShiftListByFilter(&workerID, &status)
		if err != nil {
			log.Printf("%s: GetWorkerShiftListByFilter error: %v", funcName, err)
			return err
		}
		for _, ws := range shifts {
			if ws.ShiftID == shiftID {
				log.Printf("%s: Already requested or assigned to this shift", funcName)
				return fmt.Errorf("already requested or assigned to this shift")
			}
		}
	}

	status := model.WORKER_SHIFT_APPROVED
	approvedShifts, err := s.WorkerShiftRepo.GetWorkerShiftListByFilter(&workerID, &status)
	if err != nil {
		log.Printf("%s: GetWorkerShiftListByFilter error: %v", funcName, err)
		return err
	}
	for _, ws := range approvedShifts {
		otherShift, err := s.ShiftRepo.GetShiftByID(ws.ShiftID)
		if err != nil {
			continue
		}
		if otherShift.Date == shift.Date {
			// Overlap if times intersect
			if shift.StartTime < otherShift.EndTime && shift.EndTime > otherShift.StartTime {
				log.Printf("%s: Overlapping shift on this day", funcName)
				return fmt.Errorf("overlapping shift on this day")
			}

			return fmt.Errorf("already has a shift on this day")
		}
	}

	hasShiftOnDay, shiftsThisWeek, err := s.WorkerShiftRepo.CheckWorkerShiftLimits(workerID, shift.Date)
	if err != nil {
		log.Printf("%s: CheckWorkerShiftLimits error: %v", funcName, err)
		return err
	}
	if hasShiftOnDay {
		return fmt.Errorf("already has a shift on this day")
	}
	if shiftsThisWeek >= 5 {
		return fmt.Errorf("already has 5 shifts this week")
	}

	ws := &model.WorkerShift{
		ShiftID:       shiftID,
		UserAccountID: workerID,
		Status:        model.WORKER_SHIFT_PENDING,
	}
	_, err = s.WorkerShiftRepo.CreateWorkerShift(ws)
	if err != nil {
		log.Printf("%s: CreateWorkerShift error: %v", funcName, err)
		return err
	}

	return nil
}

func (s *ShiftService) GetAllRequestedShift(ctx context.Context, workerID int64) ([]*model.ShiftStatus, error) {
	funcName := "/service/shift/GetAllRequestedShift"

	var requestedShiftStatus []*model.ShiftStatus

	workerShift, err := s.WorkerShiftRepo.GetWorkerShiftListByFilter(&workerID, nil)
	if err != nil {
		log.Printf("%s: GetWorkerShiftListByFilter error: %v", funcName, err)
		return nil, err
	}

	shiftIDs := make([]int64, 0)
	for _, shift := range workerShift {
		shiftIDs = append(shiftIDs, shift.ShiftID)
	}

	shift, err := s.ShiftRepo.GetShiftsByIDs(shiftIDs)
	if err != nil {
		log.Printf("%s: GetShiftsByIDs error: %v", funcName, err)
		return nil, err
	}

	for _, s := range shift {
		shiftStatus := &model.ShiftStatus{
			ID:             s.ID,
			Date:           s.Date,
			StartTime:      s.StartTime,
			EndTime:        s.EndTime,
			RoleAssignment: s.RoleAssignment,
			Location:       s.Location,
			IsAvailable:    s.IsAvailable,
		}

		for _, ws := range workerShift {
			if ws.ShiftID == s.ID {
				shiftStatus.StatusWorker = ws.Status
				break
			}
		}

		requestedShiftStatus = append(requestedShiftStatus, shiftStatus)
	}

	return requestedShiftStatus, nil
}

func (s *ShiftService) CreateShift(ctx context.Context, shift *model.Shift) (int64, error) {
	funcName := "/service/shift/CreateShift"

	shiftID, err := s.ShiftRepo.CreateShift(shift)
	if err != nil {
		log.Printf("%s: CreateShift error: %v", funcName, err)
		return 0, err
	}
	return shiftID, nil
}

func (s *ShiftService) UpdateShift(ctx context.Context, shift *model.Shift) error {
	funcName := "/service/shift/UpdateShift"

	err := s.ShiftRepo.UpdateShiftByID(shift)
	if err != nil {
		log.Printf("%s: UpdateShift error: %v", funcName, err)
	}

	return err
}

func (s *ShiftService) DeleteShift(ctx context.Context, shiftID int64) error {
	funcName := "/service/shift/DeleteShift"

	err := s.ShiftRepo.DeleteShiftByID(shiftID)
	if err != nil {
		log.Printf("%s: DeleteShift error: %v", funcName, err)
	}
	return err
}

func (s *ShiftService) GetAllShiftRequests(ctx context.Context, queryParam model.WorkerShiftDetailQuery) ([]model.WorkerShiftDetail, error) {
	funcName := "/service/shift/GetAllShiftRequests"

	workerShift, err := s.WorkerShiftRepo.GetWorkerShiftDetailListByFilter(&queryParam)
	if err != nil {
		log.Printf("%s: GetAllWorkerShift error: %v", funcName, err)
		return nil, err
	}

	return workerShift, nil
}

func (s *ShiftService) ApproveShiftRequest(ctx context.Context, shiftID, workerID int64) error {
	funcName := "/service/shift/ApproveShiftRequest"

	shift, err := s.ShiftRepo.GetShiftByID(shiftID)
	if err != nil {
		log.Printf("%s: GetShiftByID error: %v", funcName, err)
		return err
	}

	shift.IsAvailable = false
	if err := s.ShiftRepo.UpdateShiftByID(shift); err != nil {
		log.Printf("%s: UpdateShiftByID error: %v", funcName, err)
		return err
	}

	workerShifts, err := s.WorkerShiftRepo.ListWorkerShiftsByShift(shiftID)
	if err != nil {
		log.Printf("%s: ListWorkerShiftsByShift error: %v", funcName, err)
		return err
	}

	for _, ws := range workerShifts {
		if ws.UserAccountID == workerID {
			err := s.WorkerShiftRepo.UpdatesWorkerShiftStatus(ws.ID, model.WORKER_SHIFT_APPROVED, nil)
			if err != nil {
				log.Printf("%s: Approve error for wsID %d: %v", funcName, ws.ID, err)
				return err
			}
		} else {
			err := s.WorkerShiftRepo.UpdatesWorkerShiftStatus(ws.ID, model.WORKER_SHIFT_REJECTED, nil)
			if err != nil {
				log.Printf("%s: Reject error for wsID %d: %v", funcName, ws.ID, err)
				return err
			}
		}
	}

	return nil
}

func (s *ShiftService) RejectShiftRequest(ctx context.Context, shiftID, workerID int64) error {
	funcName := "/service/shift/RejectShiftRequest"

	workerShifts, err := s.WorkerShiftRepo.ListWorkerShiftsByShift(shiftID)
	if err != nil {
		log.Printf("%s: ListWorkerShiftsByShift error: %v", funcName, err)
		return err
	}

	for _, ws := range workerShifts {
		if ws.UserAccountID == workerID {
			err := s.WorkerShiftRepo.UpdatesWorkerShiftStatus(ws.ID, model.WORKER_SHIFT_REJECTED, nil)
			if err != nil {
				log.Printf("%s: Reject error for wsID %d: %v", funcName, ws.ID, err)
				return err
			}
			break
		}
	}

	return nil
}

func (s *ShiftService) GetShiftsByDay(ctx context.Context, date string) ([]*model.ShiftStatus, error) {
	funcName := "/service/shift/GetShiftsByDay"

	shifts, err := s.ShiftRepo.GetListShifts(model.ShiftListQuery{
		Date: date,
	})
	if err != nil {
		log.Printf("%s: GetListShifts error: %v", funcName, err)
		return nil, err
	}

	var result []*model.ShiftStatus
	for _, shift := range shifts {
		workerShifts, err := s.WorkerShiftRepo.ListWorkerShiftsByShift(shift.ID)
		if err != nil {
			log.Printf("%s: ListWorkerShiftsByShift error: %v", funcName, err)
			return nil, err
		}

		statusWorker := ""
		for _, ws := range workerShifts {
			if ws.Status == model.WORKER_SHIFT_APPROVED {
				statusWorker = model.WORKER_SHIFT_APPROVED
				break
			}
		}

		shiftStatus := &model.ShiftStatus{
			ID:             shift.ID,
			Date:           shift.Date,
			StartTime:      shift.StartTime,
			EndTime:        shift.EndTime,
			RoleAssignment: shift.RoleAssignment,
			Location:       shift.Location,
			IsAvailable:    shift.IsAvailable,
			StatusWorker:   statusWorker,
		}
		result = append(result, shiftStatus)
	}

	return result, nil
}
