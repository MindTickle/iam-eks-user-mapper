package main

type MapUserConfig struct {
	UserArn  string   `yaml:"userarn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups"`
}

type UserRoles struct {
	K8sRoles    []string
	IAMArn      string
	IAMUsername string
}

func (ur UserRoles) SetK8sRoles(userK8sRoles []string) UserRoles {
	ur.K8sRoles = userK8sRoles
	return ur
}

func (ur UserRoles) UniqueK8sRoles() UserRoles {
	ur.K8sRoles = unique(ur.K8sRoles)
	return ur
}
