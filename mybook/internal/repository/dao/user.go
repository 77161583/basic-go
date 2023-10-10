package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUseDuplicate = errors.New("邮箱冲突")
	ErrUserNotFund  = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	//第二种写法
	//err := dao.db.WithContext(ctx).Where(&u, "email = ?", email).Error
	return u, err
}

func (dao *UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	//第二种写法
	//err := dao.db.WithContext(ctx).Where(&u, "email = ?", email).Error
	return u, err
}

func (dao *UserDao) FindById(ctx context.Context, userId int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", userId).First(&u).Error
	return u, err
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	//存更新时间
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//邮箱冲突 or 手机号码冲突
			return ErrUseDuplicate
		}
	}
	return err
}

func (dao *UserDao) Edit(ctx context.Context, u User) error {
	//存更新时间
	now := time.Now().UnixMilli()
	//u.CreateTime = now
	u.UpdateTime = now

	// 构建需要更新的字段映射
	updateFields := map[string]interface{}{
		"UpdateTime":   now,
		"NickName":     u.NickName,
		"Birthday":     u.Birthday,
		"Introduction": u.Introduction,
	}

	// 构建更新条件
	updateCondition := "id = ?"
	updateParams := []interface{}{u.Id}

	// 执行更新操作
	err := dao.db.WithContext(ctx).Model(&User{}).Where(updateCondition, updateParams...).Updates(updateFields).Error
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			const uniqueConflictsErrNo uint16 = 1062
			if mysqlErr.Number == uniqueConflictsErrNo {
				// 邮箱冲突
				return ErrUseDuplicate
			}
		}
	}
	return err
}

// User 对标数据库
// 有人叫model， 也有叫 PO(persistent object)
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	//唯一的值
	Email sql.NullString `gorm:"unique"`
	// sql.NullString 允许有多个空值，但是不允许 空字符串“”
	// 不要用 指针 Phone *string  这种问题是要解引用，要判空
	Phone        sql.NullString `gorm:"unique"`
	Password     string
	NickName     string
	Birthday     string
	Introduction string
	//创建时间 -毫秒数
	CreateTime int64
	UpdateTime int64
}
