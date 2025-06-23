package imple

import (
	"context"
	"database/sql"
	"time"

	db "github.com/chise0904/golang_template/pkg/database/gorm"
	"github.com/chise0904/golang_template/pkg/errors"
	"github.com/chise0904/golang_template/pkg/pagination"
	libtime "github.com/chise0904/golang_template/pkg/time"
	"github.com/chise0904/golang_template/repository"
	"github.com/rs/zerolog/log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repo struct {
	readDB          *gorm.DB
	writeDB         *gorm.DB
	selectForUpdate bool
}

// TODO: impl repository interface
func NewIdentityRepo(db db.Connection) repository.IdentityRepository {

	return &repo{
		readDB:          db.ReadDB,
		writeDB:         db.WriteDB,
		selectForUpdate: false,
	}
}

func (r *repo) Begin(ctx context.Context) repository.IdentityRepository {
	db := r.writeDB.WithContext(ctx).Begin(&sql.TxOptions{})
	return &repo{
		writeDB:         db,
		readDB:          db,
		selectForUpdate: true,
	}
}

func (r *repo) Commit() error {
	return r.writeDB.Commit().Error
}
func (r *repo) Rollback() error {
	return r.writeDB.Rollback().Error
}

// UpsertAccount implements repository.SessionRepository
func (r *repo) UpsertAccount(ctx context.Context, in *repository.AccountInfo) (err error) {
	nowMS := time.Now().UnixMilli()
	in.UpdatedAt = nowMS

	read_body := &repository.AccountInfo{}
	dbr := r.readDB.Table((&repository.AccountInfo{}).GetTableName())
	dbr.Statement.RaiseErrorOnNotFound = true
	err = (dbr.Where("id = ?", in.ID).Find(&read_body)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// in.CreatedAt = nowMS

			err = (r.writeDB.Table(in.GetTableName()).Create(in)).Error
			if err != nil {
				return errors.NewErrorf(errors.ErrorInternalError, "upsert account info by account id failed: %v", err.Error())
			}
		} else {
			return errors.NewErrorf(errors.ErrorInternalError, "upsert account info by account id failed: %v", err.Error())
		}
	} else {
		dbw := r.writeDB.Table(in.GetTableName()).Where(repository.AccountInfo{
			ID: in.ID,
		}).Save(in)
		err = dbw.Error
		if err != nil {
			return errors.NewErrorf(errors.ErrorInternalError, "upsert account info by account id failed: %v", err.Error())
		}
	}
	return
}

// UpdateAccount implements repository.IdentityRepository
func (r *repo) UpdateAccount(ctx context.Context, in *repository.AccountInfo) error {
	nowMS := time.Now().UnixMilli()
	in.UpdatedAt = nowMS

	dbw := r.writeDB.Table(in.GetTableName()).Where(repository.AccountInfo{
		ID: in.ID,
	}).Save(in)
	err := dbw.Error
	if err != nil {
		return errors.NewErrorf(errors.ErrorInternalError, "update account info failed: %v", err.Error())
	}
	return nil
}

// DeleteAccount implements repository.IdentityRepository
func (r *repo) DeleteAccount(ctx context.Context, id string) (err error) {

	code := &repository.VerificationCode{}
	db_code := r.writeDB.Table(code.GetTableName()).Where("account_id = ?", id).Delete(&repository.VerificationCode{})
	err = db_code.Error
	if err != nil {
		return errors.NewErrorf(errors.ErrorInternalError, "delete account verification code failed: %v", err.Error())
	}

	pf := &repository.AccountProfile{}
	db_pf := r.writeDB.Table(pf.GetTableName()).Where("account_id = ?", id).Delete(&repository.AccountProfile{})
	err = db_pf.Error
	if err != nil {
		return errors.NewErrorf(errors.ErrorInternalError, "delete account profile failed: %v", err.Error())
	}

	info := &repository.AccountInfo{}
	db_info := r.writeDB.Table(info.GetTableName()).Where("id = ?", id).Delete(&repository.AccountInfo{})
	err = db_info.Error
	if err != nil {
		return errors.NewErrorf(errors.ErrorInternalError, "delete account info failed: %v", err.Error())
	}

	return
}

// GetByEmailForNormal implements repository.IdentityRepository
func (r *repo) GetByEmailForNormal(ctx context.Context, email string) (out *repository.AccountInfo, record_not_found bool, err error) {
	dbr := r.readDB.Table((&repository.AccountInfo{}).GetTableName())
	if r.selectForUpdate {
		dbr = dbr.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	dbr.Statement.RaiseErrorOnNotFound = true

	err = (dbr.Where("regis_mode = ? AND email = ?", "NORMAL", email).Find(&out)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return out, true, errors.NewErrorf(errors.ErrorResourceNotFound, "get account info by account email (%v) failed: %v", email, err)
		} else {
			return out, false, errors.NewErrorf(errors.ErrorInternalError, "get Account Info by email failed: %v", err.Error())
		}
	}
	return
}

// GetByEmailForThirdParty implements repository.IdentityRepository
func (r *repo) GetByEmailForThirdParty(ctx context.Context, email string, regis_mode string) (out *repository.AccountInfo, record_not_found bool, err error) {
	dbr := r.readDB.Table((&repository.AccountInfo{}).GetTableName())
	if r.selectForUpdate {
		dbr = dbr.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	dbr.Statement.RaiseErrorOnNotFound = true

	err = (dbr.Where("regis_mode = ? AND email = ?", regis_mode, email).Find(&out)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return out, true, errors.NewErrorf(errors.ErrorResourceNotFound, "get account info by account email (%v) failed: %v", email, err)
		} else {
			return out, false, errors.NewErrorf(errors.ErrorInternalError, "get Account Info by email failed: %v", err.Error())
		}
	}
	return
}

// GetByID implements repository.IdentityRepository
func (r *repo) GetByID(ctx context.Context, id string) (out *repository.AccountInfo, record_not_found bool, err error) {
	dbr := r.readDB.Table((&repository.AccountInfo{}).GetTableName())
	if r.selectForUpdate {
		dbr = dbr.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	dbr.Statement.RaiseErrorOnNotFound = true
	err = (dbr.Where("id = ?", id).Find(&out)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return out, true, errors.NewErrorf(errors.ErrorResourceNotFound, "get account info by account id (%v) failed: %v", id, err)
		} else {
			return out, false, errors.NewErrorf(errors.ErrorInternalError, "get account info by account id failed: %v", err.Error())
		}
	}
	return
}

// GetAccountList implements repository.IdentityRepository
func (r *repo) GetAccountList(ctx context.Context, opts *repository.GetAccountListOpts) (out []repository.AccountInfo, page *pagination.Pagination, err error) {
	dbr := r.readDB.Table((&repository.AccountInfo{}).GetTableName())
	if r.selectForUpdate {
		dbr = dbr.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	dbr = dbr.Where("deleted_at = ?", 0)

	if len(opts.AccountTypes) > 0 {
		dbr = dbr.Where("account_type in (?)", opts.AccountTypes)
	}

	pn := &pagination.Pagination{Page: opts.Page, PerPage: opts.Perpage}

	var sortType, filter string
	var s bool
	var total int64

	switch opts.By {
	case "id":
		sortType = opts.By
		if len(opts.Filter) > 0 {
			filter = "%" + opts.Filter + "%"
			dbr = dbr.Where("id like ?", filter)
		}
	case "email":
		sortType = opts.By
		if len(opts.Filter) > 0 {
			filter = "%" + opts.Filter + "%"
			dbr = dbr.Where("email like ?", filter)
		}
	default:
		sortType = "created_at"
	}

	if opts.Sort == "asc" {
		s = false
	} else if opts.Sort == "desc" {
		s = true
	}
	dbr = dbr.Order(clause.OrderByColumn{
		Column: clause.Column{
			Name: sortType,
		},
		Desc: s,
	})

	if len(opts.IDs) > 0 {
		dbr = dbr.Where("id in (?)", opts.IDs)
	}

	err = dbr.Count(&total).Error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, errors.NewErrorf(errors.ErrorResourceNotFound, "get accounts failed: %v", err.Error())
		}
		return nil, nil, errors.NewErrorf(errors.ErrorInternalError, "list Accounts fail detail: %+v db error: %s", *opts, err.Error())
	}

	limit, offset := pn.LimitAndOffset()
	if limit > 0 {
		dbr = dbr.Limit(limit).Offset(offset)
	}

	err = dbr.Find(&out).Error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, errors.NewErrorf(errors.ErrorResourceNotFound, "get accounts failed: %v", err.Error())
		}
		return nil, nil, errors.NewErrorf(errors.ErrorInternalError, "list Accounts fail detail: %+v db error: %s", *opts, err.Error())
	}

	pn.SetTotalCountAndPage(uint32(total))
	return out, pn, nil

}

// GetSimpAccountList implements repository.IdentityRepository
func (r *repo) GetSimpAccountList(ctx context.Context, opts *repository.GetAccountListOpts) (out []repository.SimpAccountInfo, page *pagination.Pagination, err error) {
	dbr := r.readDB.Table((&repository.AccountInfo{}).GetTableName())
	if r.selectForUpdate {
		dbr = dbr.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	dbr = dbr.Where("deleted_at = ?", 0)

	pn := &pagination.Pagination{Page: opts.Page, PerPage: opts.Perpage}

	var sortType, filter string
	var s bool
	var total int64

	switch opts.By {
	case "id":
		sortType = opts.By
		if len(opts.Filter) > 0 {
			filter = "%" + opts.Filter + "%"
			dbr = dbr.Where("id like ?", filter)
		}
	case "email":
		sortType = opts.By
		if len(opts.Filter) > 0 {
			filter = "%" + opts.Filter + "%"
			dbr = dbr.Where("email like ?", filter)
		}
	default:
		sortType = "created_at"
	}

	if opts.Sort == "asc" {
		s = false
	} else if opts.Sort == "desc" {
		s = true
	}
	dbr = dbr.Order(clause.OrderByColumn{
		Column: clause.Column{
			Name: sortType,
		},
		Desc: s,
	})

	err = dbr.Count(&total).Error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, errors.NewErrorf(errors.ErrorResourceNotFound, "get accounts failed: %v", err.Error())
		}
		return nil, nil, errors.NewErrorf(errors.ErrorInternalError, "list Accounts fail detail: %+v db error: %s", *opts, err.Error())
	}

	limit, offset := pn.LimitAndOffset()
	if limit > 0 {
		dbr = dbr.Limit(limit).Offset(offset)
	}

	err = dbr.Find(&out).Error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, errors.NewErrorf(errors.ErrorResourceNotFound, "get accounts failed: %v", err.Error())
		}
		return nil, nil, errors.NewErrorf(errors.ErrorInternalError, "list Accounts fail detail: %+v db error: %s", *opts, err.Error())
	}

	pn.SetTotalCountAndPage(uint32(total))
	return out, pn, nil
}

// UpsertProfile implements repository.IdentityRepository
func (r *repo) UpsertProfile(ctx context.Context, in *repository.AccountProfile) (err error) {
	nowMS := time.Now().UnixMilli()
	in.UpdatedAt = nowMS

	read_body := &repository.AccountProfile{}
	dbr := r.readDB.Table((&repository.AccountProfile{}).GetTableName())
	dbr.Statement.RaiseErrorOnNotFound = true
	err = (dbr.Where("account_id = ?", in.AccountID).Find(&read_body)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			in.CreatedAt = nowMS

			err = (r.writeDB.Table(in.GetTableName()).Create(in)).Error
			if err != nil {
				return errors.NewErrorf(errors.ErrorInternalError, "upsert account profile by account id failed: %v", err.Error())
			}
		} else {
			return errors.NewErrorf(errors.ErrorInternalError, "upsert account profile by account id failed: %v", err.Error())
		}
	} else {
		dbw := r.writeDB.Table(in.GetTableName()).Where(repository.AccountProfile{
			AccountID: in.AccountID,
		}).Save(in)
		err = dbw.Error
		if err != nil {
			return errors.NewErrorf(errors.ErrorInternalError, "upsert account profile by account id failed: %v", err.Error())
		}
	}

	return
}

// UpdateProfile implements repository.IdentityRepository
func (r *repo) UpdateProfile(ctx context.Context, in *repository.DBProfileOptions) error {
	nowMS := time.Now().UnixMilli()
	dbw := r.writeDB.Table((&repository.AccountProfile{}).GetTableName())

	values := map[string]interface{}{
		"updated_at": nowMS,
	}
	if in.UserName != "" {
		values["user_name"] = in.UserName
	}
	if in.Icon != nil {
		values["icon"] = in.Icon
	}
	if in.Description != "" {
		values["description"] = in.Description
	}
	if in.Gender != "" {
		values["gender"] = in.Gender
	}
	if in.Birthday != nil {
		values["birthday"] = in.Birthday
	}
	if in.Job != "" {
		values["job"] = in.Job
	}
	if in.Country != "" {
		values["country"] = in.Country
	}
	if in.City != "" {
		values["city"] = in.City
	}
	if in.District != "" {
		values["district"] = in.District
	}
	if in.ZipCode != "" {
		values["zip_code"] = in.ZipCode
	}
	if in.Address != "" {
		values["address"] = in.Address
	}
	if in.ShippingAddress != nil {
		values["shipping_address"] = in.ShippingAddress
	}
	if in.Language != "" {
		values["language"] = in.Language
	}
	if in.Phone_noti != nil {
		values["phone_noti"] = in.Phone_noti
	}
	if in.Email_noti != nil {
		values["email_noti"] = in.Email_noti
	}

	err := dbw.Where("account_id = ?", in.AccountID).Updates(values).Error
	if err != nil {
		return errors.NewErrorf(errors.ErrorInternalError, "update profile failed: %v", err.Error())
	}
	return nil
}

// GetProfileByID implements repository.IdentityRepository
func (r *repo) GetProfileByID(ctx context.Context, id string) (out *repository.AccountProfile, err error) {
	dbr := r.readDB.Table((&repository.AccountProfile{}).GetTableName())
	if r.selectForUpdate {
		dbr = dbr.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	dbr.Statement.RaiseErrorOnNotFound = true
	err = (dbr.Where("account_id = ?", id).Find(&out)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return out, errors.NewErrorf(errors.ErrorResourceNotFound, "get account profile by account id (%v) failed: %v", id, err)
		} else {
			return out, errors.NewErrorf(errors.ErrorInternalError, "get account profile by account id failed: %v", err.Error())
		}
	}
	return
}

// GetVeriCodeByActionAndID implements repository.IdentityRepository
func (r *repo) GetVeriCodeByActionAndID(ctx context.Context, id string, action string) (out *repository.VerificationCode, record_not_found bool, err error) {
	dbr := r.readDB.Table((&repository.VerificationCode{}).GetTableName())
	if r.selectForUpdate {
		dbr = dbr.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	dbr.Statement.RaiseErrorOnNotFound = true
	err = (dbr.Where("account_id = ? AND action = ?", id, action).Find(&out)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true, errors.NewErrorf(errors.ErrorResourceNotFound, "get verification code by account id and action (%v - %v) failed: %v", id, action, err)
		} else {
			return nil, false, errors.NewErrorf(errors.ErrorInternalError, "get verification code by account id and action: %v", err.Error())
		}
	}
	return
}

// UpsertVeriCodeByActionAndID implements repository.IdentityRepository
func (r *repo) UpsertVeriCodeByActionAndID(ctx context.Context, in *repository.VerificationCode) (err error) {
	nowMS := time.Now().UnixMilli()

	in.UpdatedAt = nowMS

	out := &repository.VerificationCode{}
	dbr := r.readDB.Table((&repository.VerificationCode{}).GetTableName())
	dbr.Statement.RaiseErrorOnNotFound = true
	err = (dbr.Where("account_id = ? AND action = ?", in.AccountID, in.Action).Find(&out)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			in.CreatedAt = nowMS
			err = (r.writeDB.Table(in.GetTableName()).Create(in)).Error
			if err != nil {
				return errors.NewErrorf(errors.ErrorInternalError, "upsert verification code by account id and action: %v", err.Error())
			}
		} else {
			return errors.NewErrorf(errors.ErrorInternalError, "upsert verification code by account id and action: %v", err.Error())
		}
	} else {
		dbw := r.writeDB.Table(in.GetTableName()).Where(repository.VerificationCode{
			AccountID: in.AccountID,
			Action:    in.Action,
		}).Save(in)
		err = dbw.Error
		if err != nil {
			return errors.NewErrorf(errors.ErrorInternalError, "upsert verification code by account id and action: %v", err.Error())
		}
	}
	return
}

// DeleteVeriCodeByActionAndID implements repository.IdentityRepository
func (r *repo) DeleteVeriCodeByActionAndID(ctx context.Context, id string, action string) (err error) {
	code := &repository.VerificationCode{}
	db_code := r.writeDB.Table(code.GetTableName()).Where("account_id = ? AND action = ?", id, action).Delete(&repository.VerificationCode{})
	err = db_code.Error
	if err != nil {
		return errors.NewErrorf(errors.ErrorInternalError, "delete verification code by account id and action failed: %v", err.Error())
	}
	return
}

// PingDB implements repository.IdentityRepository
func (r *repo) PingDB() (respTimeMilliSec float32, err error) {
	tm := time.Now()

	sqlDB, err := r.readDB.DB()
	if err != nil {
		log.Error().Msgf("err: %s", err.Error())
		return 0, err
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Error().Msgf("err: %s", err.Error())
		return 0, err
	}
	etm := float32(time.Since(tm).Nanoseconds()) / 1000000

	return etm, err
}

func (r *repo) ListProfiles(ctx context.Context, filter *repository.ListProfilesFilter) (result []*repository.AccountProfile, pg *pagination.Pagination, err error) {

	dbr := r.readDB.Table((&repository.AccountProfile{}).GetTableName())
	if r.selectForUpdate {
		dbr = dbr.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if len(filter.AccountIDs) > 0 {
		dbr = dbr.Where("account_id in (?)", filter.AccountIDs)
	}
	var total int64

	err = dbr.Count(&total).Error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, errors.NewErrorf(errors.ErrorResourceNotFound, "get accounts failed: %v", err.Error())
		}
		return nil, nil, errors.NewErrorf(errors.ErrorInternalError, "list Accounts fail detail: %+v db error: %s", *filter, err.Error())
	}
	pn := &pagination.Pagination{Page: filter.Page, PerPage: filter.Perpage}

	limit, offset := pn.LimitAndOffset()
	if limit > 0 {
		dbr = dbr.Limit(limit).Offset(offset)
	}

	var out []*repository.AccountProfile

	err = dbr.Find(&out).Error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, errors.NewErrorf(errors.ErrorResourceNotFound, "get accounts failed: %v", err.Error())
		}
		return nil, nil, errors.NewErrorf(errors.ErrorInternalError, "list Accounts fail detail: %+v db error: %s", *filter, err.Error())
	}

	pn.SetTotalCountAndPage(uint32(total))
	return out, pn, nil
}

func (r *repo) UpdateAccountStatus(ctx context.Context, id string, accountStatus int32) error {
	now := libtime.NowMS()

	values := map[string]interface{}{
		"updated_at": now,
		"status":     accountStatus,
	}

	err := r.writeDB.Table((&repository.AccountInfo{}).GetTableName()).Where("id = ? ", id).Updates(values).Error
	if err != nil {
		return errors.NewError(errors.ErrorInternalError, err.Error())
	}

	return nil
}
