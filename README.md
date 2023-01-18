# golang-promise

This is a small project detailing how to create a promise like structure in golang with sync.WaitGroups. I created it
to streamline fetching several results from multiple sources async. 

    result := ""
    myErr := new(error)
    wg := new(sync.WaitGroup)
    wg.Add(1)

    NewPromise[string](func() (string, error) {
      str, err := doSomething()
      if err != nil {
        return nil, err
      }

      return str, nil
    }).Then(func(str string) {
      result = str
    }, func(e error) {
      myErr = e
    }).Finally(func() {
      wg.Done()
    })

    if myErr != nil {
      handleError()
    }

    wg.Wait()