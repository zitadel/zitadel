package v1

//
//type Config struct {
//	Repository  z_sql.Config
//	ServiceName string
//	Cache       *config.CacheConfig `mapstructure:"cache"`
//}
//
//func (c *Config) UnmarshalMap(value interface{}) error {
//	var (
//		cfg map[string]interface{}
//		ok  bool
//	)
//	if cfg, ok = value.(map[string]interface{}); !ok {
//		return errors.ThrowInternal(nil, "", "invalid UDP port range")
//	}
//	cache, _ := cfg["cache"].(map[string]interface{})
//
//	fmt.Println(cache)
//
//	return nil
//}
//
//func Start(conf Config) (Eventstore, error) {
//	repo, _, err := z_sql.Start(conf.Repository)
//	if err != nil {
//		return nil, err
//	}
//
//	return &eventstore{
//		repo: repo,
//	}, nil
//}
