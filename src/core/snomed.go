package core

import "strconv"

type SnomedDescription struct {
	Id                 int    `redis:"id"`
	EffectiveTime      int    `redis:"effective_time"`
	Active             bool   `redis:"active"`
	ModuleId           uint64 `redis:"module_id"`
	ConceptId          int    `redis:"concept_id"`
	LanguageCode       string `redis:"language_code"`
	TypeId             uint64 `redis:"type_id"`
	Term               string `redis:"term"`
	CaseSignificanceId uint64 `redis:"case_significance_id"`
}

func (sd SnomedDescription) GetMap() map[string]string {
	active := "0"
	if sd.Active {
		active = "1"
	}
	return map[string]string{
		"Id":                 strconv.Itoa(sd.Id),
		"EffectiveTime":      strconv.Itoa(sd.EffectiveTime),
		"Active":             active,
		"ModuleId":           strconv.FormatUint(sd.ModuleId, 10),
		"ConceptId":          strconv.Itoa(sd.ConceptId),
		"LanguageCode":       sd.LanguageCode,
		"TypeId":             strconv.FormatUint(sd.TypeId, 10),
		"Term":               sd.Term,
		"CaseSignificanceId": strconv.FormatUint(sd.CaseSignificanceId, 10),
	}
}
