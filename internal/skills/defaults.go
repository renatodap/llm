package skills

// GetDefaultSkills returns commonly useful skills
func GetDefaultSkills() []*Skill {
	return []*Skill{
		NewResearchAssistant(),
		NewSummarizer(),
		NewTranslator(),
		NewTutor(),
	}
}

// GetAllBuiltInSkills returns all available built-in skills
func GetAllBuiltInSkills() map[string]*Skill {
	skills := []*Skill{
		NewResearchAssistant(),
		NewContentCreator(),
		NewCodeReviewer(),
		NewDataAnalyst(),
		NewTutor(),
		NewTranslator(),
		NewSummarizer(),
		NewCoder(),
	}

	skillMap := make(map[string]*Skill)
	for _, skill := range skills {
		skillMap[skill.Name] = skill
	}
	return skillMap
}

// AutoRegisterDefaultSkills registers default skills to a registry
func AutoRegisterDefaultSkills(registry *SkillRegistry) {
	for _, skill := range GetDefaultSkills() {
		registry.Register(skill)
	}
}

// AutoRegisterAllSkills registers all built-in skills to a registry
func AutoRegisterAllSkills(registry *SkillRegistry) {
	for _, skill := range GetAllBuiltInSkills() {
		registry.Register(skill)
	}
}
