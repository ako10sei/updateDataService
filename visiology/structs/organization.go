package structs

type Column struct {
	CoollegeId         int    `json:"id_учреждения"`
	CoollegeName       string `json:"Наименование учреждения"`
	CoollegeShortName  int    `json:"Краткое наименование"`
	CoollegeOGRN       int    `json:"ОГРН"`
	CoollegeJurAddres  int    `json:"Адрес юридический"`
	CoollegeFactAddres int    `json:"Адрес фактический"`
	Director           int    `json:"Руководитель"`
	Area               int    `json:"Район"`
	Fillials           int    `json:"Филиалы"`
	Site               int    `json:"Сайт"`
	MaxOccupancy       int    `json:"Проектная мощность"`
	StudentsCount      int    `json:"Количество студентов общее"`
	MastersCount       int    `json:"Количество мастеров обучения"`
}
