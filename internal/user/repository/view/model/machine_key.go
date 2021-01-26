package model

//
//const (
//	MachineKeyKeyID      = "id"
//	MachineKeyObjectID   = "object_id"
//	MachineKeyObjectType = "object_type"
//)
//
//type MachineKeyView struct {
//	ID             string    `json:"keyId" gorm:"column:id;primary_key"`
//	ObjectID       string    `json:"-" gorm:"column:object_id;primary_key"`
//	ObjectType     int32     `json:"-" gorm:"column:object_type;primary_key"`
//	Type           int32     `json:"type" gorm:"column:machine_type"`
//	ExpirationDate time.Time `json:"expirationDate" gorm:"column:expiration_date"`
//	Sequence       uint64    `json:"-" gorm:"column:sequence"`
//
//	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
//
//	PublicKey []byte `json:"publicKey" gorm:"column:public_key"`
//}
//
//func MachineKeyViewFromModel(key *model.MachineKeyView) *MachineKeyView {
//	return &MachineKeyView{
//		ID:             key.ID,
//		ObjectID:       key.ObjectID,
//		ObjectType:     int32(key.ObjectType),
//		Type:           int32(key.Type),
//		ExpirationDate: key.ExpirationDate,
//		Sequence:       key.Sequence,
//		CreationDate:   key.CreationDate,
//	}
//}
//
//func MachineKeyToModel(key *MachineKeyView) *model.MachineKeyView {
//	return &model.MachineKeyView{
//		ID:             key.ID,
//		ObjectID:       key.ObjectID,
//		ObjectType:     model.ObjectType(key.ObjectType),
//		Type:           model.MachineKeyType(key.Type),
//		ExpirationDate: key.ExpirationDate,
//		Sequence:       key.Sequence,
//		CreationDate:   key.CreationDate,
//		PublicKey:      key.PublicKey,
//	}
//}
//
//func MachineKeysToModel(keys []*MachineKeyView) []*model.MachineKeyView {
//	result := make([]*model.MachineKeyView, len(keys))
//	for i, key := range keys {
//		result[i] = MachineKeyToModel(key)
//	}
//	return result
//}
//
//func (k *MachineKeyView) AppendEvent(event *models.Event) (err error) {
//	k.Sequence = event.Sequence
//	switch event.Type {
//	case es_model.MachineKeyAdded:
//		k.setRootData(event)
//		k.CreationDate = event.CreationDate
//		err = k.SetData(event)
//	}
//	return err
//}
//
//func (k *MachineKeyView) setRootData(event *models.Event) {
//	k.ObjectID = event.AggregateID
//}
//
//func (r *MachineKeyView) SetData(event *models.Event) error {
//	if err := json.Unmarshal(event.Data, r); err != nil {
//		logging.Log("EVEN-Sj90d").WithError(err).Error("could not unmarshal event data")
//		return caos_errs.ThrowInternal(err, "MODEL-lub6s", "Could not unmarshal data")
//	}
//	return nil
//}
