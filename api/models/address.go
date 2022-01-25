package models

import (
	"errors"
	"html"    //Escapestring 
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Address struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`    //GORM prefer convention over configuration, by default, GORM uses ID as primary key,
	Address   string    `gorm:"size:100;not null;unique" json:"address"`//pluralize struct name to snake_cases as table name, snake_case as column name,
	Phone     string    `gorm:"size:100;not null;" json:"phone"`        //and uses CreatedAt, UpdatedAt to track creating/updating time
	Author    User      `json:"author"`
	AuthorID  uint32    `sql:"type:int REFERENCES users(id)" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Address) Prepare() {                                       // *Address(pointer recevier)
	p.ID = 0
	p.Address = html.EscapeString(strings.TrimSpace(p.Address))
	p.Phone = html.EscapeString(strings.TrimSpace(p.Phone))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}
func (p *Address) Validate() error {

	if p.Address == "" {
		return errors.New("Required Address")
	}
	if p.Phone == "" {
		return errors.New("Required Phone")
	}
	if p.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}
func (p *Address) Saveaddress(db *gorm.DB) (*Address, error) {
	var err error
	err = db.Debug().Model(&Address{}).Create(&p).Error
	if err != nil {
		return &Address{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Address{}, err
		}
	}
	return p, nil
}
func (p *Address) FindAllAddress(db *gorm.DB) (*[]Address, error) {
	var err error
	addr := []Address{}
	err = db.Debug().Model(&Address{}).Limit(100).Find(&addr).Error // error handling
	if err != nil {
		return &[]Address{}, err
	}
	if len(addr) > 0 {
		for i, _ := range addr {
			err := db.Debug().Model(&User{}).Where("id = ?", addr[i].AuthorID).Take(&addr[i].Author).Error
			if err != nil {
				return &[]Address{}, err
			}
		}
	}
	return &addr, nil
}
func (p *Address) FindAddressByID(db *gorm.DB, pid uint64) (*Address, error) {
	var err error
	err = db.Debug().Model(&Address{}).Where("id = ?", pid).Take(&p).Error //error handling
	if err != nil {
		return &Address{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Address{}, err
		}
	}
	return p, nil
}

func (p *Address) UpdateAddressPost(db *gorm.DB) (*Address, error) {

	var err error
	err = db.Debug().Model(&Address{}).Where("id = ?", p.ID).Updates(Address{Address: p.Address, Phone: p.Phone, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Address{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Address{}, err
		}
	}
	return p, nil
}

func (p *Address) DeleteAddress(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Address{}).Where("id = ? and author_id = ?", pid, uid).Take(&Address{}).Delete(&Address{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Address not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
