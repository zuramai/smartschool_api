package models

type Kecamatan struct {
	ID            int    `json:"id" gorm:"primary_key"`
	KotaID        int    `json:"kota_id"`
	NamaKecamatan string `json:"nama_kecamatan"`
}

func (Kecamatan) TableName() string {
	return "kecamatan"
}
