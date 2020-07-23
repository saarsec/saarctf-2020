package stahl4

import "sync"

type State struct {
	Production *Production
	Logistic   *Logistic
	Mux        sync.Mutex
}

type Logistic struct {
	Busy                      bool
	TotalMaterialCapacity     uint16
	MaterialsStored           uint16
	TotalProductCapacity      uint16
	ProductsStored            uint16
	NumberOfNoMaterialAnswers uint16
}

func NewLogistic() *Logistic {
	return &Logistic{Busy: false, TotalMaterialCapacity: 150, MaterialsStored: 0, TotalProductCapacity: 150, ProductsStored: 0}
}

type Production struct {
	ProductCapacity    uint16
	MaterialCapacity   uint16
	ProductAmount      uint16
	MaterialAmount     uint16
	AskedForMaterials  bool
	AskedToGetProducts bool
}

func NewProduction() *Production {
	return &Production{ProductCapacity: 300, MaterialCapacity: 300, ProductAmount: 0, MaterialAmount: 0, AskedForMaterials: false, AskedToGetProducts: false}
}
