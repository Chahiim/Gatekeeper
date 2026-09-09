package data

import "database/sql"

type Models struct {
	Consumers ConsumerModel
	Reports   ReportModel
	Images    ImageModel
	ImageJobs ImageJobModel
	Variants  VariantModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Consumers: ConsumerModel{DB: db},
		Reports:   ReportModel{DB: db},
		Images:    ImageModel{DB: db},
		ImageJobs: ImageJobModel{DB: db},
		Variants:  VariantModel{DB: db},
	}
}
