package presenter

type ViewModel interface{}

type Presenter[T ViewModel] interface {
	BuildViewModel() T
}
