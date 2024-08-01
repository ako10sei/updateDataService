package structs

type Column struct {
	CoollegeShortName  int `json:"Краткое наименование"`
	CoollegeOGRN       int `json:"ОГРН"`
	CoollegeJurAddres  int `json:"Адрес юридический"`
	CoollegeFactAddres int `json:"Адрес фактический"`
	Director           int `json:"Руководитель"`
	Area               int `json:"Район"`
	Fillials           int `json:"Филиалы"`
	Site               int `json:"Сайт"`
	MaxOccupancy       int `json:"Проектная мощность"`
	StudentsCount      int `json:"Количество студентов общее"`
	MastersCount       int `json:"Количество мастеров обучения"`
}

// GetAllFields возвращает срез строк, содержащий имена всех полей в структуре Column.
//
// Функция не принимает никаких параметров и возвращает срез строк, где каждая строка
// представляет имя поля в структуре Column. Порядок строк в срезе соответствует порядку
// полей в структуре.
//
// Функция не возвращает никакого дополнительного кода за пределами непосредственной области действия.
func (c *Column) GetAllFields() []string {
	return []string{
		"ОГРН",
		"Адрес юридический",
		"Адрес фактический",
		"Руководитель",
		"Район",
		"Филиалы",
		"Проектная мощность",
	}
}
