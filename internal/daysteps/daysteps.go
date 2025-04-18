package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

const (
	StepLength = 0.65 // длина шага в метрах.
	mInKm      = 1000 // количество метров в километре.
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",") //Разделить строку на слайс строк.
	if len(parts) != 2 {              //чтобы длина слайса была равна 2.
		return 0, 0, errors.New("invalid string format")
	}
	//Преобразовать первый элемент слайса (количество шагов) в тип int
	//Обработать возможные ошибки.
	//При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, 0, errors.New("invalid number of steps")
	}

	duration, err := time.ParseDuration(parts[1]) //Преобразовать второй элемент слайса в time.Duration
	if err != nil {
		return 0, 0, errors.New("invalid duration format")
	}
	return steps, duration, nil //верните количество шагов, продолжительность и nil (для ошибки).
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data) //Получить данные о количестве шагов и продолжительности прогулки с помощью функции parsePackage()
	if err != nil {                            //В случае возникновения ошибки вывести её на экран и вернуть пустую строку.
		fmt.Println("parsing error:", err)
		return ""
	}
	if steps <= 0 { //Проверить, чтобы количество шагов было больше 0. В противном случае вернуть пустую строку.
		fmt.Println("invalid number of steps")
		return ""
	}
	distanceInMeters := float64(steps) * StepLength                                 //Вычислить дистанцию в метрах.
	distanceInKilometers := distanceInMeters / 1000                                 //Перевести дистанцию в километры, разделив её на 1000.
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration) //Вычислить количество калорий.
	result := fmt.Sprintf("Количество шагов: %d. \n Дистанция составила %.2f км.\n Вы сожгли %.2f ккал.", steps, distanceInKilometers, calories)
	return result //Сформировать строку, которую будете возвращать, по примеру
}
