package service

import (
	"sort"
	"strings"

	"abasithdev.github.io/internal-cs-center-backend/internal/domain"
	"abasithdev.github.io/internal-cs-center-backend/internal/errors"
	"abasithdev.github.io/internal-cs-center-backend/internal/storage"
)

type PaymentService struct {
	store *storage.MemoryStore
}

type ListRequest struct {
	Page    int
	Size    int
	Status  string
	Search  string
	SortBy  string
	OrderBy string
}

type ListResult struct {
	Total      int               `json:"total"`
	Size       int               `json:"size"`
	Page       int               `json:"page"`
	TotalPages int               `json:"total_pages"`
	Data       []*domain.Payment `json:"data"`
}

func NewPaymentService(store *storage.MemoryStore) *PaymentService {
	return &PaymentService{store: store}
}

func (payment *PaymentService) GetTotalByFilter(request ListRequest) int {
	filtered := payment.getListPayment(request)
	return len(filtered)
}

func (payment *PaymentService) GetStatusSummary() (int, int, int) {
	all := payment.store.GetPaymentList()
	completed, process, failed := 0, 0, 0
	for _, paymentData := range all {
		switch paymentData.Status {
		case "completed":
			completed++
		case "processing":
			process++
		case "failed":
			failed++
		}
	}

	return completed, process, failed
}

func (payment *PaymentService) GetList(request ListRequest) ListResult {
	filtered := payment.getListPayment(request)

	switch request.SortBy {
	case "amount":
		sort.Slice(filtered, func(i, j int) bool {
			if request.OrderBy == "asc" {
				return filtered[i].Amount < filtered[j].Amount
			}

			return filtered[i].Amount > filtered[j].Amount
		})
	default:
		sort.Slice(filtered, func(i, j int) bool {
			if request.OrderBy == "asc" {
				return filtered[i].Date.Before(filtered[j].Date)
			}
			return filtered[i].Date.After(filtered[j].Date)
		})
	}

	if request.Size <= 0 {
		request.Size = 10
	}

	if request.Page <= 0 {
		request.Page = 1
	}

	start := (request.Page - 1) * request.Size
	total := len(filtered)
	if start > total {
		start = total
	}

	end := start + request.Size
	if end > total {
		end = total
	}

	perItems := filtered[start:end]
	totalPage := 0
	if total > 0 {
		totalPage = (total + request.Size - 1) / request.Size
	}

	return ListResult{
		Total:      total,
		Size:       request.Size,
		Page:       request.Page,
		TotalPages: totalPage,
		Data:       perItems,
	}
}

func (payment *PaymentService) Review(paymentID string) error {
	paymentResult, ok := payment.store.GetPaymentById(paymentID)
	if !ok {
		return errors.NewNotFoundError("paymentId: " + paymentID)
	}

	paymentResult.Reviewed = true
	payment.store.UpdatePayment(paymentResult)
	return nil
}

// private
func (payment *PaymentService) getListPayment(request ListRequest) []*domain.Payment {

	all := payment.store.GetPaymentList()
	filtered := []*domain.Payment{}
	for _, paymentData := range all {
		if request.Status != "" && request.Status != paymentData.Status {
			continue
		}
		if request.Search != "" && !strings.Contains(paymentData.ID, request.Search) {
			continue
		}

		filtered = append(filtered, paymentData)
	}

	return filtered
}
