package database

// TODO: Why?
/*
func AdaptFunc(
	monitor mntr.Monitor,
	dbClient db.Client,
) (
	operator.QueryFunc,
	error,
) {

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {

		dbHost, dbPort, dbConnectionURL, err := dbClient.GetConnectionInfo(monitor, k8sClient)
		if err != nil {
			return nil, err
		}

		users, err := dbClient.ListUsers(monitor, k8sClient)
		if err != nil {
			return nil, err
		}

		curr := &Current{
			Host:          dbHost,
			Port:          dbPort,
			ConnectionURL: dbConnectionURL,
			Users:         users,
		}

		SetDatabaseInQueried(queried, curr)

		return func(k8sClient kubernetes.ClientInt) error {
			return nil
		}, nil
	}, nil
}
*/
