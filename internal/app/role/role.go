package role

type Role int

const (
	User Role = iota
	Moderator
)

func (r Role) String() string {
	switch r {
	case User:
		return "user"
	case Moderator:
		return "moderator"
	default:
		return "guest"
	}
}

func FromString(roleStr string) Role {
	switch roleStr {
	case "user":
		return User
	case "moderator":
		return Moderator
	default:
		return User
	}
}

func (r Role) GetScopes() []string {
	switch r {
	case User:
		return []string{"read:stages", "read:requests", "create:requests", "update:requests"}
	case Moderator:
		return []string{"read:stages", "read:requests", "create:requests", "update:requests", "resolve:requests", "reject:requests", "create:stages", "update:stages", "delete:stages", "manage:users"}
	default:
		return []string{"read:stages", "read:requests"}
	}
}
