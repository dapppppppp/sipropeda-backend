package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"sipropeda-backend/infras"
)

type AppConfigRepository interface {
	ResolveByID(id string) (data AppConfig, err error)
	ResolveDTOByID(id string) (data AppConfigDTOByID, err error)
	Update(req AppConfig) error
}

type appConfigRepositoryPostgreSQL struct {
	db *infras.PostgresqlConn
}

func ProvideAppConfigRepositoryPostgreSQL(db *infras.PostgresqlConn) AppConfigRepository {
	return &appConfigRepositoryPostgreSQL{db: db}
}

const selectAppConfigQuery = `
	SELECT id, nama_sistem, tagline, instansi, status, favicon, logo, child_logo, 
	email_instansi, url_root, jalan, kelurahan, kecamatan, kabupaten, provinsi, kode_pos, 
	telp, fax
	FROM app_config
`

func (r *appConfigRepositoryPostgreSQL) ResolveByID(id string) (data AppConfig, err error) {
	err = r.db.Read.Get(&data, selectAppConfigQuery+" WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New(fmt.Sprintf("Config dengan ID [%s] tidak ditemukan!", id))
		}
		return
	}
	return
}

const selectDTOByIDAppConfigQuery = `
	SELECT id, nama_sistem, tagline, instansi, favicon, logo, child_logo, 
	email_instansi, jalan, kelurahan, kecamatan, kabupaten, provinsi, kode_pos, 
	telp, fax 
	FROM app_config 
`

func (r *appConfigRepositoryPostgreSQL) ResolveDTOByID(id string) (data AppConfigDTOByID, err error) {
	err = r.db.Read.Get(&data, selectDTOByIDAppConfigQuery+" WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New(fmt.Sprintf("Config dengan ID [%s] tidak ditemukan!", id))
		}
		return
	}
	return
}

const updateAppConfigQuery = `
	UPDATE app_config SET 
	nama_sistem=$1, tagline=$2, instansi=$3, status=$4, favicon=$5, logo=$6, child_logo=$7, 
	email_instansi=$8, url_root=$9, jalan=$10, kelurahan=$11, kecamatan=$12, kabupaten=$13, 
	provinsi=$14, kode_pos=$15, telp=$16, fax=$17 
	WHERE id=$18
`

func (r *appConfigRepositoryPostgreSQL) Update(req AppConfig) error {
	_, err := r.db.Write.Exec(updateAppConfigQuery,
		req.NamaSistem, req.Tagline, req.Instansi, req.Status, req.Favicon, req.Logo, req.ChildLogo,
		req.EmailInstansi, req.UrlRoot, req.Jalan, req.Kelurahan, req.Kecamatan, req.Kabupaten,
		req.Provinsi, req.KodePos, req.Telp, req.Fax, req.Id,
	)
	return err
}