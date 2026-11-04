package app

type Application interface {
    Run() error
}

func Run(application Application) error {
    return application.Run()
}

type mockApplication struct {
    err error
}

func (m *mockApplication) Run() error {
    return m.err
}