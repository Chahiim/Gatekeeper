package main

import (
	"database/sql"
	"sync"

	"github.com/Chahiim/Gatekeeper/internal/data"
)

type Models struct {
	Consumers data.ConsumerModel
	Reports   data.ReportModel
	Images    data.ImageModel
	ImageJobs data.ImageJobModel
	Variants  data.VariantModel
	wg        sync.WaitGroup
}

func newModels(db *sql.DB) Models {
	return Models{
		Consumers: data.ConsumerModel{DB: db},
		Reports:   data.ReportModel{DB: db},
		Images:    data.ImageModel{DB: db},
		ImageJobs: data.ImageJobModel{DB: db},
		Variants:  data.VariantModel{DB: db},
	}
}
