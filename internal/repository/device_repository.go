package repository

import (
	"context"
	"strings"
	"time"

	"ancianato-backend/internal/domain"
	"ancianato-backend/internal/infrastructure/cache"
	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) domain.DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(ctx context.Context, device *domain.Wearable) error {
	return r.db.WithContext(ctx).Create(device).Error
}

func (r *deviceRepository) GetByID(ctx context.Context, id string) (*domain.Wearable, error) {
	cacheKey := "device:" + id
	if cachedData, found := cache.AppCache.Get(cacheKey); found {
		return cachedData.(*domain.Wearable), nil
	}

	var device domain.Wearable
	err := r.db.WithContext(ctx).First(&device, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheKey, &device, 5*time.Minute)
	return &device, nil
}

func (r *deviceRepository) GetByMacAddress(ctx context.Context, mac string) (*domain.Wearable, error) {
	// Normalizar: quitar dos puntos y guiones
	cleanMac := strings.ReplaceAll(strings.ReplaceAll(mac, ":", ""), "-", "")
	
	// Crear versión con dos puntos
	var colonMac string
	if len(cleanMac) == 12 {
		var parts []string
		for i := 0; i < 12; i += 2 {
			parts = append(parts, cleanMac[i:i+2])
		}
		colonMac = strings.Join(parts, ":")
	} else {
		colonMac = cleanMac
	}

	var device domain.Wearable
	err := r.db.WithContext(ctx).Where(
		"LOWER(mac_address) = LOWER(?) OR LOWER(mac_address) = LOWER(?)", 
		cleanMac, 
		colonMac,
	).First(&device).Error
	
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) GetAll(ctx context.Context) ([]domain.Wearable, error) {
	cacheKey := "device:ALL"
	if cachedData, found := cache.AppCache.Get(cacheKey); found {
		return cachedData.([]domain.Wearable), nil
	}

	var devices []domain.Wearable
	err := r.db.WithContext(ctx).Find(&devices).Error
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheKey, devices, 5*time.Minute)
	return devices, nil
}

func (r *deviceRepository) Update(ctx context.Context, device *domain.Wearable) error {
	err := r.db.WithContext(ctx).Save(device).Error
	if err == nil {
		cache.AppCache.Delete("device:" + device.ID)
		cache.AppCache.Delete("device:ALL")
	}
	return err
}

func (r *deviceRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Delete(&domain.Wearable{}, "id = ?", id).Error
	if err == nil {
		cache.AppCache.Delete("device:" + id)
		cache.AppCache.Delete("device:ALL")
	}
	return err
}
