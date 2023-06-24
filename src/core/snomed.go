package core

type SnomedDescription struct {
	Id                 int
	EffectiveTime      int
	Active             bool
	ModuleId           uint64
	ConceptId          int
	LanguageCode       string
	TypeId             uint64
	Term               string
	CaseSignificanceId uint64
}
