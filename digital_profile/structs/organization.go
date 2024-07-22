package structs

import "fmt"

type Organization struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ShortName     string `json:"short_name"`
	Ogrn          any    `json:"ogrn"`
	Director      any    `json:"director,omitempty"`
	Telephone     any    `json:"telephone,omitempty"`
	Fax           any    `json:"fax"`
	Email         any    `json:"email,omitempty"`
	WebSite       any    `json:"web_site"`
	FAddressFull  any    `json:"f_address_full"`
	UAddressFull  any    `json:"u_address_full"`
	TerritoryName string `json:"territory_name"`
	Parent        any    `json:"parent"`
	MaxOccupancy  int    `json:"maximum_occupancy"`
	StudentsCount int    `json:"students_count"`
	MastersCount  int    `json:"masters_count"`
}

// GetColumnByField возвращает карту, содержащую определенные поля структуры Organization
// в форматированном виде для дальнейшей обработки. Ключи карты представляют имена столбцов в таблице Visiology,
// а значения соответствуют полям структуры Organization.
func (o *Organization) GetColumnByField() map[string]any {
	return map[string]any{
		"ОГРН":               o.Ogrn,
		"Адрес юридический":  o.UAddressFull,
		"Адрес фактический":  o.FAddressFull,
		"Руководитель":       fmt.Sprintf("%s, %s, %s", o.Director, o.Email, o.Telephone),
		"Район":              GetAreaIdByName(o.TerritoryName),
		"Филиалы":            "",
		"Сайт":               o.WebSite,
		"Проектная мощность": o.MaxOccupancy,
		"Количество студентов общее":   0,
		"Количество мастеров обучения": o.MastersCount,
	}
}

// GetAreaIdByName возвращает идентификатор района на основе его названия.
// Функция использует словарь для хранения названий районов и их соответствующих идентификаторов.
// Если указанное название района не найдено в словаре, функция возвращает -1.
func GetAreaIdByName(areaName string) int {
	areaMap := map[string]int{
		"Александровский район":  1,
		"Вязниковский район":     2,
		"Гороховецкий район":     3,
		"Гусь-Хрустальный район": 4,
		"Камешковский район":     5,
		"Киржачский район":       6,
		"Ковровский район":       7,
		"Кольчугинский район":    8,
		"Меленковский район":     9,
		"Муром":                  10,
		"Петушинский район":      11,
		"Селивановский район":    12,
		"Собинский район":        13,
		"Судогодский район":      14,
		"Суздальский район":      15,
		"Юрьев-Польский район":   16,
		"Владимир":               17,
		"Гусь-Хрустальный":       18,
		"Ковров":                 19,
		"Радужный":               20,
		"Муромский район":        21,
	}

	id, ok := areaMap[areaName]
	if !ok {
		return -1
	}

	return id
}
