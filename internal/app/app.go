package app

type Application interface {
    Run() error
}

func Run(application Application) error {
    return application.Run()
}