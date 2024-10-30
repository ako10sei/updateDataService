package structs

type Column struct {
	CoollegeOGRN       int `json:"ОГРН"`
	CoollegeJurAddres  int `json:"Адрес юридический"`
	CoollegeFactAddres int `json:"Адрес фактический"`
	Director           int `json:"Руководитель"`
	Area               int `json:"Район"`
	MaxOccupancy       int `json:"Проектная мощность"`
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
		"Проектная мощность",
	}
}