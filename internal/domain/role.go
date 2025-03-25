package domain

type Role struct {
	Id        int
	Name      string
	ProjectId string
}

func NewRole(name, projectId string) *Role {
	return &Role{
		Name:      name,
		ProjectId: projectId,
	}
}

func (r *Role) IsSystemRole() bool {
	return r.ProjectId == ""
}
