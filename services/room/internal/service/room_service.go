package service

import (
	"errors"
	"time"

	"github.com/flaminshinjan/address.ai/services/room/internal/model"
	"github.com/flaminshinjan/address.ai/services/room/internal/repository"
)

// RoomService handles business logic for rooms
type RoomService struct {
	roomRepo    *repository.RoomRepository
	bookingRepo *repository.BookingRepository
}

// NewRoomService creates a new RoomService
func NewRoomService(roomRepo *repository.RoomRepository, bookingRepo *repository.BookingRepository) *RoomService {
	return &RoomService{
		roomRepo:    roomRepo,
		bookingRepo: bookingRepo,
	}
}

// CreateRoom creates a new room
func (s *RoomService) CreateRoom(room model.Room) (model.RoomResponse, error) {
	// Check if room number already exists
	_, err := s.roomRepo.GetByNumber(room.Number)
	if err == nil {
		return model.RoomResponse{}, errors.New("room number already exists")
	}

	// Create room
	createdRoom, err := s.roomRepo.Create(room)
	if err != nil {
		return model.RoomResponse{}, err
	}

	return createdRoom.ToResponse(), nil
}

// GetRoomByID gets a room by ID
func (s *RoomService) GetRoomByID(id string) (model.RoomResponse, error) {
	room, err := s.roomRepo.GetByID(id)
	if err != nil {
		return model.RoomResponse{}, err
	}

	return room.ToResponse(), nil
}

// UpdateRoom updates a room
func (s *RoomService) UpdateRoom(id string, room model.Room) (model.RoomResponse, error) {
	// Check if room exists
	existingRoom, err := s.roomRepo.GetByID(id)
	if err != nil {
		return model.RoomResponse{}, err
	}

	// Check if room number already exists (if changing number)
	if room.Number != existingRoom.Number {
		_, err := s.roomRepo.GetByNumber(room.Number)
		if err == nil {
			return model.RoomResponse{}, errors.New("room number already exists")
		}
	}

	// Update fields
	existingRoom.Number = room.Number
	existingRoom.Type = room.Type
	existingRoom.Floor = room.Floor
	existingRoom.Description = room.Description
	existingRoom.Capacity = room.Capacity
	existingRoom.PricePerDay = room.PricePerDay
	existingRoom.Status = room.Status
	existingRoom.UpdatedAt = time.Now()

	// Save room
	updatedRoom, err := s.roomRepo.Update(existingRoom)
	if err != nil {
		return model.RoomResponse{}, err
	}

	return updatedRoom.ToResponse(), nil
}

// DeleteRoom deletes a room
func (s *RoomService) DeleteRoom(id string) error {
	// Check if room has bookings
	bookings, err := s.bookingRepo.GetByRoomID(id, 1, 0)
	if err != nil {
		return err
	}

	if len(bookings) > 0 {
		return errors.New("cannot delete room with bookings")
	}

	return s.roomRepo.Delete(id)
}

// ListRooms lists all rooms
func (s *RoomService) ListRooms(limit, offset int) ([]model.RoomResponse, error) {
	rooms, err := s.roomRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	var roomResponses []model.RoomResponse
	for _, room := range rooms {
		roomResponses = append(roomResponses, room.ToResponse())
	}

	return roomResponses, nil
}

// ListAvailableRooms lists all available rooms for the given dates
func (s *RoomService) ListAvailableRooms(startDate, endDate time.Time, limit, offset int) ([]model.RoomResponse, error) {
	// Validate dates
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}

	if startDate.Before(time.Now()) {
		return nil, errors.New("start date must be in the future")
	}

	rooms, err := s.roomRepo.ListAvailable(startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}

	var roomResponses []model.RoomResponse
	for _, room := range rooms {
		roomResponses = append(roomResponses, room.ToResponse())
	}

	return roomResponses, nil
}

// UpdateRoomStatus updates a room's status
func (s *RoomService) UpdateRoomStatus(id, status string) error {
	// Check if room exists
	_, err := s.roomRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Validate status
	validStatuses := map[string]bool{
		"available":   true,
		"occupied":    true,
		"maintenance": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	return s.roomRepo.UpdateStatus(id, status)
}
