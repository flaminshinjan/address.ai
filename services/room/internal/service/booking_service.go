package service

import (
	"errors"
	"math"
	"time"

	"github.com/flaminshinjan/address.ai/services/room/internal/model"
	"github.com/flaminshinjan/address.ai/services/room/internal/repository"
)

// BookingService handles business logic for bookings
type BookingService struct {
	bookingRepo *repository.BookingRepository
	roomRepo    *repository.RoomRepository
}

// NewBookingService creates a new BookingService
func NewBookingService(bookingRepo *repository.BookingRepository, roomRepo *repository.RoomRepository) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
	}
}

// CreateBooking creates a new booking
func (s *BookingService) CreateBooking(userID string, req model.BookingRequest) (model.BookingResponse, error) {
	// Validate dates
	if req.StartDate.After(req.EndDate) {
		return model.BookingResponse{}, errors.New("start date must be before end date")
	}

	if req.StartDate.Before(time.Now()) {
		return model.BookingResponse{}, errors.New("start date must be in the future")
	}

	// Check if room exists
	room, err := s.roomRepo.GetByID(req.RoomID)
	if err != nil {
		return model.BookingResponse{}, err
	}

	// Check if room is available
	if room.Status != "available" {
		return model.BookingResponse{}, errors.New("room is not available")
	}

	// Check if room is available for the given dates
	available, err := s.bookingRepo.CheckRoomAvailability(req.RoomID, req.StartDate, req.EndDate)
	if err != nil {
		return model.BookingResponse{}, err
	}

	if !available {
		return model.BookingResponse{}, errors.New("room is not available for the given dates")
	}

	// Calculate total price
	days := int(math.Ceil(req.EndDate.Sub(req.StartDate).Hours() / 24))
	totalPrice := room.PricePerDay * float64(days)

	// Create booking
	booking := model.Booking{
		RoomID:     req.RoomID,
		UserID:     userID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		TotalPrice: totalPrice,
		Status:     "confirmed",
	}

	createdBooking, err := s.bookingRepo.Create(booking)
	if err != nil {
		return model.BookingResponse{}, err
	}

	// Create response
	response := model.BookingResponse{
		ID:         createdBooking.ID,
		Room:       room.ToResponse(),
		UserID:     createdBooking.UserID,
		StartDate:  createdBooking.StartDate,
		EndDate:    createdBooking.EndDate,
		TotalPrice: createdBooking.TotalPrice,
		Status:     createdBooking.Status,
		CreatedAt:  createdBooking.CreatedAt,
	}

	return response, nil
}

// GetBookingByID gets a booking by ID
func (s *BookingService) GetBookingByID(id string) (model.BookingResponse, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		return model.BookingResponse{}, err
	}

	room, err := s.roomRepo.GetByID(booking.RoomID)
	if err != nil {
		return model.BookingResponse{}, err
	}

	response := model.BookingResponse{
		ID:         booking.ID,
		Room:       room.ToResponse(),
		UserID:     booking.UserID,
		StartDate:  booking.StartDate,
		EndDate:    booking.EndDate,
		TotalPrice: booking.TotalPrice,
		Status:     booking.Status,
		CreatedAt:  booking.CreatedAt,
	}

	return response, nil
}

// GetBookingsByUserID gets bookings by user ID
func (s *BookingService) GetBookingsByUserID(userID string, limit, offset int) ([]model.BookingResponse, error) {
	bookings, err := s.bookingRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []model.BookingResponse
	for _, booking := range bookings {
		room, err := s.roomRepo.GetByID(booking.RoomID)
		if err != nil {
			return nil, err
		}

		response := model.BookingResponse{
			ID:         booking.ID,
			Room:       room.ToResponse(),
			UserID:     booking.UserID,
			StartDate:  booking.StartDate,
			EndDate:    booking.EndDate,
			TotalPrice: booking.TotalPrice,
			Status:     booking.Status,
			CreatedAt:  booking.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// CancelBooking cancels a booking
func (s *BookingService) CancelBooking(id string) error {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		return err
	}

	if booking.Status != "confirmed" {
		return errors.New("booking is not in a confirmed state")
	}

	if booking.StartDate.Before(time.Now()) {
		return errors.New("cannot cancel a booking that has already started")
	}

	return s.bookingRepo.UpdateStatus(id, "cancelled")
}

// ListBookings lists all bookings
func (s *BookingService) ListBookings(limit, offset int) ([]model.BookingResponse, error) {
	bookings, err := s.bookingRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []model.BookingResponse
	for _, booking := range bookings {
		room, err := s.roomRepo.GetByID(booking.RoomID)
		if err != nil {
			return nil, err
		}

		response := model.BookingResponse{
			ID:         booking.ID,
			Room:       room.ToResponse(),
			UserID:     booking.UserID,
			StartDate:  booking.StartDate,
			EndDate:    booking.EndDate,
			TotalPrice: booking.TotalPrice,
			Status:     booking.Status,
			CreatedAt:  booking.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, nil
}
