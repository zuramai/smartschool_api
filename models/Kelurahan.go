package models

type Kelurahan struct {
	ID            int    `json:"id"`
	KecamatanID   int    `json:"kecamatan_id"`
	Kecamatan Kecamatan `json:"kecamatan" gorm:"foreignkey:KecamatanID"`
	NamaKelurahan string `json:"nama_kelurahan"`
}

func (Kelurahan) TableName() string {
	return "kelurahan"
}
