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
	TerritoryName any    `json:"territory_name"`
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
		"id_ЦП":         o.ID,
		"id_учреждения": o.ID,
		"Наименование учреждения": o.Name,
		"Краткое наименование":    o.ShortName,
		"ОГРН":               o.Ogrn,
		"Адрес юридический":  o.UAddressFull,
		"Адрес фактический":  o.FAddressFull,
		"Руководитель":       fmt.Sprintf("%s, %s, %s", o.Director, o.Email, o.Telephone),
		"Район":              o.TerritoryName,
		"Филиалы":            "",
		"Сайт":               o.WebSite,
		"Проектная мощность": o.MaxOccupancy,
		"Количество студентов общее":   0,
		"Количество мастеров обучения": o.MastersCount,
	}
}
