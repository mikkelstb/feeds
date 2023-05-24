package iptctop

type iptctop int

const (
	Arts iptctop = iota
	Crime
	Disaster
	Economy
	Education
	Environment
	Health
	Human_Interest
	Labour
	Lifestyle
	Politics
	Religion
	Science
	Society
	Sport
	Conflict
	Weather
)

func (c iptctop) GetName() string {
	return topic_names[c]
}

var topic_names = [17]string{
	"arts, culture, entertainment and media",
	"crime, law and justice",
	"disaster, accident and emergency incident",
	"economy, business and finance",
	"education",
	"environment",
	"health",
	"human interest",
	"labour",
	"lifestyle and leisure",
	"politics",
	"religion",
	"science and technology",
	"society",
	"sport",
	"conflict, war and peace",
	"weather",
}

func (c iptctop) GetDescription() string {
	return topic_description[c]
}

var topic_description = [17]string{
	"All forms of arts, entertainment, cultural heritage and media",
	"The establishment and/or statement of the rules of behaviour in society, the enforcement of these rules, breaches of the rules, the punishment of offenders and the organisations and bodies involved in these activities",
	"Man made or natural event resulting in loss of life or injury to living creatures and/or damage to inanimate objects or property",
	"All matters concerning the planning, production and exchange of wealth",
	"All aspects of furthering knowledge, formally or informally",
	"All aspects of protection, damage, and condition of the ecosystem of the planet earth and its surroundings",
	"All aspects of physical and mental well-being",
	"Item that discusses individuals, groups, animals, plants or other objects in an emotional way",
	"Social aspects, organisations, rules and conditions affecting the employment of human effort for the generation of wealth or provision of services and the economic support of the unemployed",
	"Activities undertaken for pleasure, relaxation or recreation outside paid employment, including eating and travel",
	"Local, regional, national and international exercise of power, or struggle for power, and the relationships between governing bodies and states",
	"Belief systems, institutions and people who provide moral guidance to followers",
	"All aspects pertaining to human understanding of, as well as methodical study and research of natural, formal and social sciences, such as astronomy, linguistics or economics",
	"The concerns, issues, affairs and institutions relevant to human social interactions, problems and welfare, such as poverty, human rights and family planning",
	"Competitive activity or skill that involves physical and/or mental effort and organisations and bodies involved in these activities",
	"Acts of socially or politically motivated protest or violence, military activities, geopolitical conflicts, as well as resolution efforts",
	"The study, prediction and reporting of meteorological phenomena",
}
