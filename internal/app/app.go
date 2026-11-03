package app

type Application interface {
	Run() error
}

func Run(app Application) error {

	return app.Run()
}
