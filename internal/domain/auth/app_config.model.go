package auth

// AppConfig merepresentasikan tabel app_config di database
type AppConfig struct {
	Id            string  `db:"id" json:"id"`
	NamaSistem    string  `db:"nama_sistem" json:"namaSistem"`
	Tagline       *string `db:"tagline" json:"tagline"`
	Instansi      string  `db:"instansi" json:"instansi"`
	Status        *string `db:"status" json:"status"`
	Favicon       *string `db:"favicon" json:"favicon"`
	Logo          *string `db:"logo" json:"logo"`
	ChildLogo     *string `db:"child_logo" json:"childLogo"`
	EmailInstansi *string `db:"email_instansi" json:"emailInstansi"`
	UrlRoot       string  `db:"url_root" json:"urlRoot"`
	Jalan         *string `db:"jalan" json:"jalan"`
	Kelurahan     *string `db:"kelurahan" json:"kelurahan"`
	Kecamatan     *string `db:"kecamatan" json:"kecamatan"`
	Kabupaten     *string `db:"kabupaten" json:"kabupaten"`
	Provinsi      *string `db:"provinsi" json:"provinsi"`
	KodePos       *string `db:"kode_pos" json:"kodePos"`
	Telp          *string `db:"telp" json:"telp"`
	Fax           *string `db:"fax" json:"fax"`
}

// AppConfigDTOByID untuk konsumsi Public (Halaman Login, dsb)
type AppConfigDTOByID struct {
	Id            string  `db:"id" json:"id"`
	NamaSistem    string  `db:"nama_sistem" json:"namaSistem"`
	Tagline       *string `db:"tagline" json:"tagline"`
	Instansi      string  `db:"instansi" json:"instansi"`
	Favicon       *string `db:"favicon" json:"favicon"`
	Logo          *string `db:"logo" json:"logo"`
	ChildLogo     *string `db:"child_logo" json:"childLogo"`
	EmailInstansi *string `db:"email_instansi" json:"emailInstansi"`
	Jalan         *string `db:"jalan" json:"jalan"`
	Kelurahan     *string `db:"kelurahan" json:"kelurahan"`
	Kecamatan     *string `db:"kecamatan" json:"kecamatan"`
	Kabupaten     *string `db:"kabupaten" json:"kabupaten"`
	Provinsi      *string `db:"provinsi" json:"provinsi"`
	KodePos       *string `db:"kode_pos" json:"kodePos"`
	Telp          *string `db:"telp" json:"telp"`
	Fax           *string `db:"fax" json:"fax"`
}

// RequestAppConfigFormat adalah format JSON saat Update
type RequestAppConfigFormat struct {
	NamaSistem    string  `json:"namaSistem"`
	Tagline       *string `json:"tagline"`
	Instansi      string  `json:"instansi"`
	Status        *string `json:"status"`
	Favicon       *string `json:"favicon"`
	Logo          *string `json:"logo"`
	ChildLogo     *string `json:"childLogo"`
	EmailInstansi *string `json:"emailInstansi"`
	UrlRoot       string  `json:"urlRoot"`
	Jalan         *string `json:"jalan"`
	Kelurahan     *string `json:"kelurahan"`
	Kecamatan     *string `json:"kecamatan"`
	Kabupaten     *string `json:"kabupaten"`
	Provinsi      *string `json:"provinsi"`
	KodePos       *string `json:"kodePos"`
	Telp          *string `json:"telp"`
	Fax           *string `json:"fax"`
}

func (m *AppConfig) AppConfigFormat(req RequestAppConfigFormat) {
	m.NamaSistem = req.NamaSistem
	m.Tagline = req.Tagline
	m.Instansi = req.Instansi
	m.Status = req.Status
	m.Favicon = req.Favicon
	m.Logo = req.Logo
	m.ChildLogo = req.ChildLogo
	m.EmailInstansi = req.EmailInstansi
	m.UrlRoot = req.UrlRoot
	m.Jalan = req.Jalan
	m.Kelurahan = req.Kelurahan
	m.Kecamatan = req.Kecamatan
	m.Kabupaten = req.Kabupaten
	m.Provinsi = req.Provinsi
	m.KodePos = req.KodePos
	m.Telp = req.Telp
	m.Fax = req.Fax
}