//1. Merge n channels
//2. Если один из входных каналов закрывается,
//то нужно закрыть все остальные каналы

func case3(channels ...chan int) chan int {
 out := make(chan int)
	done := make(chan struct{})

	var wg sync.WaitGroup
	var closeOnce sync.Once
	var closeDoneOnce sync.Once

	forward := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			select {
			case out <- v:
			case <-done:
				return
			}
		}
		closeDoneOnce.Do(func() {
			close(done)
		})
	}

	go func() {
		for _, ch := range channels {
			wg.Add(1)
			go forward(ch)
		}

		wg.Wait()
		close(out)
	}()

	go func() {
		<-done
		closeOnce.Do(func() {
			for _, ch := range channels {
				if ch != nil {
					func() {
						defer func() { recover() }()
						close(ch)
					}()
				}
			}
		})
	}()

	return out 
}
