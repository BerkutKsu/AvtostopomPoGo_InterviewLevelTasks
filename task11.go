//Написать код функции, которая делает merge N каналов.
//Весь входной поток перенаправляется в один канал.

func joinChannels(chs ...<-chan int) <-chan int {
  var wg sync.WaitGroup
  out := make(chan int)
	
	forward := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			out <- v
		}
	}
	wg.Add(len(chs))

	for _, ch := range chs {
		go forward(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
