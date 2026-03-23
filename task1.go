//Условие задачи
//Дан массив целых чисел nums и целое число k. Нужно написать функцию,
//которая вынимает из массива nums k наиболее часто встречающихся элементов.

//Пример
//# ввод
//nums = [1,1,1,2,2,3]
//k = 2
//# вывод (в любом порядке)
//[1, 2]


func topKFrequentElements(nums []int, k int) []int {
  uniq := make(map[int]int)
	for _, num := range nums {
		uniq[num]++
	}
	amountUniq := len(uniq)

	sortArray := make([][]int, amountUniq)
	i := 0
	for key, val := range uniq {
		sortArray[i] = []int{key, val}
		i++
	}

	for i := 0; i < amountUniq-1; i++ {
		for j := i + 1; j < amountUniq; j++ {
			if sortArray[i][1] < sortArray[j][1] {
				sortArray[i], sortArray[j] = sortArray[j], sortArray[i]
			}
		}
	}
	var result []int
	for i = 0; i < k; i++ {
		result = append(result, sortArray[i][0])
	}
	return result
}
