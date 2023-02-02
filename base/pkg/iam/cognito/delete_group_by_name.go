package cognito

func (i *cognitoIAM) DeleteGroupByName(name string) error {
	// id and name are the same with cognito
	return i.DeleteGroup(name)
}
