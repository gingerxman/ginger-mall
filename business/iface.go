package business

type IModel interface {
	GetId() int
}

type IUser interface {
	GetId() int
}

type ICorp interface {
	GetId() int
	GetPlatformId() int
	IsPlatform() bool
	IsValid() bool
}
