package parse

type StatHandler func(obj *Base, raw, key string)

var propsHandlers = map[string]StatHandler{
	"Min/Max Defense":                    baseHandler,
	"Min Strength":                       baseHandler,
	"Min Dexterity":                      baseHandler,
	"Durability":                         baseHandler,
	"Sockets":                            baseHandler,
	"Quality Level":                      baseHandler,
	"Type":                               baseHandler,
	"% Block":                            baseHandler,
	"Paladin Smite Damage":               baseHandler,
	"Assassin Kick Damage":               baseHandler,
	"#Boxes":                             baseHandler,
	"Min/Max 1h Damage":                  baseHandler,
	"Rangeadder":                         baseHandler,
	"Speed by Class":                     baseHandler,
	"Two-Hand Damage":                    baseHandler,
	"Min/Max Dagger Damage":              baseHandler,
	"Min/Max 2h Damage":                  baseHandler,
	"Two Hand Damage":                    baseHandler,
	"Class":                              baseHandler,
	"Throw Damage":                       baseHandler,
	"Melee Rangeadder":                   baseHandler,
	"Max Stack":                          baseHandler,
	"Speed for Paladin":                  baseHandler,
	"Magic Level":                        baseHandler,
	"Speed for Sorceress":                baseHandler,
	"Speed for Necromancer":              baseHandler,
	"Level Required":                     baseHandler,
	"Level Requirement":                  baseHandler,
	"Speed for Amazon":                   baseHandler,
	"Speed for Assassin":                 baseHandler,
	"+Skills":                            baseHandler,
	"Min/Max 1h Damage(Barbarian Only)":  baseHandler,
	"Min/Max 1h Damage (Barbarian Only)": baseHandler,
	"Speed By Class":                     baseHandler,
	"Speeds by Class":                    baseHandler,
	"Required Level":                     baseHandler,
}

var baseHandler = func(obj *Base, raw, key string) {
	obj.Stats = append(obj.Stats, Stat{
		Key: key,
		Val: raw,
	})
}

func applyStat(obj *Base, raw, key string) {
	if h, ok := propsHandlers[key]; ok {
		h(obj, raw, key)
	} else {
		obj.Stats = append(obj.Stats, Stat{Key: key, Val: raw})
	}
}
