package view

//
//func MachineKeyByIDs(db *gorm.DB, table, userID, keyID string) (*model.MachineKeyView, error) {
//	key := new(model.MachineKeyView)
//	query := repository.PrepareGetByQuery(table,
//		model.MachineKeySearchQuery{Key: usr_model.MachineKeyObjectID, Method: global_model.SearchMethodEquals, Value: userID},
//		model.MachineKeySearchQuery{Key: usr_model.MachineKeyKeyID, Method: global_model.SearchMethodEquals, Value: keyID},
//	)
//	err := query(db, key)
//	if caos_errs.IsNotFound(err) {
//		return nil, caos_errs.ThrowNotFound(nil, "VIEW-3Dk9s", "Errors.User.KeyNotFound")
//	}
//	return key, err
//}
//
//func SearchMachineKeys(db *gorm.DB, table string, req *usr_model.MachineKeySearchRequest) ([]*model.MachineKeyView, uint64, error) {
//	members := make([]*model.MachineKeyView, 0)
//	query := repository.PrepareSearchQuery(table, model.MachineKeySearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
//	count, err := query(db, &members)
//	if err != nil {
//		return nil, 0, err
//	}
//	return members, count, nil
//}
//
//func MachineKeysByUserID(db *gorm.DB, table string, userID string) ([]*model.MachineKeyView, error) {
//	keys := make([]*model.MachineKeyView, 0)
//	queries := []*usr_model.MachineKeySearchQuery{
//		{
//			Key:    usr_model.MachineKeyObjectID,
//			Value:  userID,
//			Method: global_model.SearchMethodEquals,
//		},
//	}
//	query := repository.PrepareSearchQuery(table, model.MachineKeySearchRequest{Queries: queries})
//	_, err := query(db, &keys)
//	if err != nil {
//		return nil, err
//	}
//	return keys, nil
//}
//
//func MachineKeyByID(db *gorm.DB, table string, keyID string) (*model.MachineKeyView, error) {
//	key := new(model.MachineKeyView)
//	query := repository.PrepareGetByQuery(table,
//		model.MachineKeySearchQuery{Key: usr_model.MachineKeyKeyID, Method: global_model.SearchMethodEquals, Value: keyID},
//	)
//	err := query(db, key)
//	if caos_errs.IsNotFound(err) {
//		return nil, caos_errs.ThrowNotFound(nil, "VIEW-BjN6x", "Errors.User.KeyNotFound")
//	}
//	return key, err
//}
//
//func PutMachineKey(db *gorm.DB, table string, role *model.MachineKeyView) error {
//	save := repository.PrepareSave(table)
//	return save(db, role)
//}
//
//func DeleteMachineKey(db *gorm.DB, table, keyID string) error {
//	delete := repository.PrepareDeleteByKey(table, model.MachineKeySearchKey(usr_model.MachineKeyKeyID), keyID)
//	return delete(db)
//}
//
//func DeleteMachineKeysByUserID(db *gorm.DB, table, userID string) error {
//	delete := repository.PrepareDeleteByKey(table, model.MachineKeySearchKey(usr_model.MachineKeyObjectID), userID)
//	return delete(db)
//}
