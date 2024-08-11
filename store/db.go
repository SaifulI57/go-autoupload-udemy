package store

import (
	"errors"
	"fmt"
	"time"

	log "github.com/SaifulI57/uploader-udemy/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UploaderData struct {
	gorm.Model
	UID                int64     `gorm:"column:uid;unique"`
	UdemyUrl           string    `gorm:"column:udemy_url;unique"`
	Description        string    `gorm:"column:description"`
	HashUrl            string    `gorm:"column:hash_url:unique"`
	UploadedToFacebook bool      `gorm:"column:uploaded_to_facebook"`
	UploadedToBot      bool      `gorm:"column:uploaded_to_bot"`
	UploadedToWebsite  bool      `gorm:"column:uploaded_to_website"`
	UploadedToWhatsapp bool      `gorm:"column:uploaded_to_whatsapp"`
	ExpiredIn          time.Time `gorm:"column:expired_in"`
}

type ListTelegramPost struct {
	gorm.Model
	UID           int64     `gorm:"column:uid;unique"`
	Today         time.Time `gorm:"column:today"`
	StartPostID   int64     `gorm:"column:start_post_id"`
	ListPostID    []int64   `gorm:"column:list_post_id"`
	CurrentPostID int64     `gorm:"column:current_post_id"`
}

type ListLinkUdemy struct {
	gorm.Model
	UID       int64     `gorm:"column:uid;unique"`
	Today     time.Time `gorm:"column:today"`
	ListUdemy []string  `gorm:"column:list_udemy"`
}

type SubscribeBotDiscord struct {
	gorm.Model
	UID       int64  `gorm:"column:uid;unique"`
	ServerID  string `gorm:"column:server_id;unique"`
	ChannelID string `gorm:"column:channel_id;unique"`
}

type ShortLink struct {
	gorm.Model
	UID      int64        `gorm:"column:uid;unique"`
	Short    string       `gorm:"column:short_link"`
	LongLink string       `gorm:"column:long_link"`
	Coupon   UploaderData `gorm:"column:coupon"`
}

type UploadDB struct{}
type ShortDB struct{}
type SubscribeDB struct{}

var Store *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	// dsn := os.Getenv("dsn")
	dsn := "host=localhost user=postgres password=12345678 dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Failed to ConnectDB: %s", err))
		return nil, err
	}

	Store = db

	err = Store.AutoMigrate(&UploaderData{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (d *UploadDB) CreateDBUploader(db *gorm.DB, upload *UploaderData) error {
	first := db.First(upload)
	if first == nil {
		log.Logger.Error("duplicate upload")
		return errors.New("duplicate upload")
	}
	if err := db.Create(upload).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Logger.Error("Error creating upload: Duplicate entry")
			return err
		}
		log.Logger.Error(fmt.Sprintf("Failed to CreateDB: %s", err))
		return err
	}
	log.Logger.Info(fmt.Sprintf("Successfully created uploader data with ID: %d", upload.ID))

	return nil
}

func (d *UploadDB) UpdateDBUploader(db *gorm.DB, upload *UploaderData) error {
	if err := db.Model(&UploaderData{}).Updates(&upload).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("Failed to update Data uploader: %v", err))
		return err
	}
	log.Logger.Info("Successfully updated Data uploader")

	return nil
}

func (d *UploadDB) BatchUploadDBUploader(db *gorm.DB, upload []*UploaderData) error {
	for _, up := range upload {
		first := db.First(up)
		log.Logger.Info(fmt.Sprintf("Uploading UID: %d", up.UID))
		if first == nil {
			log.Logger.Error("duplicate upload")
			return errors.New("duplicate upload")
		}
	}

	if err := db.CreateInBatches(upload, 5).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Logger.Error("Error creating upload: Duplicate entry")
			return err
		}
		log.Logger.Error(fmt.Sprintf("Failed to CreateDB: %s", err))
		return err
	}
	log.Logger.Info(fmt.Sprintf("Successfully created uploader data with: %v", upload))

	return nil
}

func (d *UploadDB) GetAllSDateDBUploader(db *gorm.DB, dates time.Time) ([]UploaderData, error) {
	var uploads []UploaderData
	date := "2006-01-02"
	if err := db.Where("DATE(created_at) = ?", dates.Format(date)).Find(&uploads).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("Failed to query by date: %s", err))
		return nil, err
	}
	log.Logger.Info(fmt.Sprintf("Successfully retrieved all upload %v", uploads))
	return uploads, nil
}

func (d *UploadDB) GetLastTwoDaysDBUploader(db *gorm.DB) ([]UploaderData, error) {
	var uploads []UploaderData
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")

	if err := db.Where("DATE(created_at) BETWEEN ? AND ?", twoDaysAgo, today).Find(&uploads).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("Failed to query by date range (2 days ago to today): %s", err))
		return nil, err
	}
	log.Logger.Info(fmt.Sprintf("Successfully retrieved uploads from 2 days ago to today: %v", uploads))
	return uploads, nil
}

func (d *UploadDB) GetLastThreeDaysDBUploader(db *gorm.DB) ([]UploaderData, error) {
	var uploads []UploaderData
	threeDaysAgo := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")

	if err := db.Where("DATE(created_at) BETWEEN ? AND ?", threeDaysAgo, today).Find(&uploads).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("Failed to query by date range (3 days ago to today): %s", err))
		return nil, err
	}
	log.Logger.Info(fmt.Sprintf("Successfully retrieved uploads from 3 days ago to today: %v", uploads))
	return uploads, nil
}

func (d *UploadDB) GetAllDBUploader(db *gorm.DB) ([]UploaderData, error) {
	var uploads []UploaderData
	res := db.Find(&uploads)
	if err := res.Error; err != nil {
		log.Logger.Error("Error getting uploader")
		return nil, err
	}
	log.Logger.Info(fmt.Sprintf("Successfully get all uploaded data: %v", uploads))
	return uploads, nil

}

func (d *UploadDB) DeleteDBUploader(db *gorm.DB, upload *UploaderData) error {
	if err := db.Unscoped().Delete(&UploaderData{}, upload.ID).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("Error deleting uploader data: %v", err))
		return err
	}
	log.Logger.Info("Successfully deleted uploader data")
	return nil
}

func (s *ShortDB) CreateDBShort(db *gorm.DB, short *ShortLink) error {
	if err := db.Create(short).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("Error creating shortlink: %v", err))
		return err
	}
	log.Logger.Info("Successfully created short")
	return nil
}

func (s *ShortDB) UpdateDBShort(db *gorm.DB, short *ShortLink) error {
	if err := db.Model(&ShortLink{}).Updates(short).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("Error updating shortlink: %v", err))
		return err
	}
	log.Logger.Info("Updated shortlink successfully")
	return nil
}

func (s *ShortDB) DeleteShortLink(db *gorm.DB, short *ShortLink) error {
	if err := db.Unscoped().Delete(&short).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("delete short link failed %v", err))
		return err
	}
	log.Logger.Info("delete short link succeeded")
	return nil
}

func (sub *SubscribeDB) CreateDBSubscription(db *gorm.DB, subs *SubscribeBotDiscord) error {
	if err := db.Create(subs).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("create subscription failed %v", err))
		return err
	}
	log.Logger.Info("create subscription succeeded")
	return nil
}

func (sub *ShortDB) UpdateDBSubscription(db *gorm.DB, subs *SubscribeBotDiscord) error {
	if err := db.Model(&SubscribeBotDiscord{}).Updates(subs).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("Error updating shortlink: %v", err))
		return err
	}
	log.Logger.Info("Updated shortlink successfully")
	return nil
}

func (sub *ShortDB) DeleteDBSubscription(db *gorm.DB, subs *SubscribeBotDiscord) error {
	if err := db.Unscoped().Delete(&subs).Error; err != nil {
		log.Logger.Error(fmt.Sprintf("delete subscription failed %v", err))
		return err
	}
	log.Logger.Info("delete subscription succeeded")
	return nil
}
