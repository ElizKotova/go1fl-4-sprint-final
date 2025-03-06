package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	//Разделить строку на слайс строк.
	//чтобы длина слайса была равна 2.
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("Неверный формат строки")
	}
	//Преобразовать первый элемент слайса (количество шагов) в тип int
	//Обработать возможные ошибки.
	//При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, 0, errors.New("Неверное количество шагов")
	}
	//Преобразовать второй элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, errors.New("Неверный формат продолжительности")
	}
	//верните количество шагов, продолжительность и nil (для ошибки).
	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	//Получить данные о количестве шагов и продолжительности прогулки с помощью функции parsePackage()
	steps, duration, err := parsePackage(data)
	//В случае возникновения ошибки вывести её на экран и вернуть пустую строку.
	if err != nil {
		fmt.Println("Ошибка парсинга:", err)
		return ""
	}
	//Проверить, чтобы количество шагов было больше 0. В противном случае вернуть пустую строку.
	if steps <= 0 {
		return ""
	}
	//Вычислить дистанцию в метрах.
	distanceInMeters := float64(steps) * StepLength
	//Перевести дистанцию в километры, разделив её на 1000.
	distanceInKilometers := distanceInMeters / 1000
	//Вычислить количество калорий.
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	//Сформировать строку, которую будете возвращать, по примеру
	result := fmt.Sprintf("Количество шагов: %d. \n Дистанция составила %.2f км.\n Вы сожгли %.2f ккал.", steps, distanceInKilometers, calories)
	return result
}
