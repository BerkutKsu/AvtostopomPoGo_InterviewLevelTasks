// Нужно реализовать функцию, которая выполняет
// поиск query во всех переданных SearchFunc
// Когда получаем первый успешный результат -
// отдаем его сразу. Если все SearchFunc отработали
// с ошибкой - отдаем последнюю полученную ошибку
type Result struct{
  Data string
}

type SearchFunc func(ctx context.Context, query string) (Result, error)

func MultiSearch(ctx context.Context, query string, sfs []SearchFunc) (Result, error) {
	if len(sfs) == 0 {
		return Result{}, errors.New("no search functions provided")
	}

	resultCh := make(chan Result, 1)
	errCh := make(chan error, len(sfs))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	for _, sf := range sfs {
		wg.Add(1)
		go func(fn SearchFunc) {
			defer wg.Done()
			result, err := fn(ctx, query)

			select {
			case <-ctx.Done():
				return
			default:
			}

			if err == nil {
				select {
				case resultCh <- result:
					cancel()
				default:
				}
			} else {
				select {
				case errCh <- err:
				default:
				}
			}
		}(sf)
	}

	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	for {
		select {
		case result, ok := <-resultCh:
			if ok {
				return result, nil
			}
			var lastErr error
			for err := range errCh {
				lastErr = err
			}
			if lastErr != nil {
				return Result{}, lastErr
			}
			return Result{}, errors.New("all searches failed")
		case <-ctx.Done():
			var lastErr error
			for err := range errCh {
				lastErr = err
			}
			if lastErr != nil {
				return Result{}, lastErr
			}
			return Result{}, context.Canceled
		}
	}
}
