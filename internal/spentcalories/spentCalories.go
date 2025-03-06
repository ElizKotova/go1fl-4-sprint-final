package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	//Разделить строку на слайс строк.
	parts := strings.Split(data, ",")
	//Проверить, чтобы длина слайса была равна 3.
	if len(parts) != 3 {
		return 0, "", 0, errors.New("Неверный формат строки")
	}
	//Преобразовать первый элемент слайса (количество шагов) в тип int.
	steps, err := strconv.Atoi(parts[0])
	//Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("Неверное количество шагов")
	}
	activity := parts[1]
	//Преобразовать третий элемент слайса в time.Duration.
	duration, err := time.ParseDuration(parts[2])
	//Обработать возможные ошибки.
	if err != nil {
		return 0, "", 0, errors.New("Неверный формат продолжительности")
	}
	//верните количество шагов, вид активности, продолжительность и nil (для ошибки)
	return steps, activity, duration, nil
}

func distance(steps int) float64 {
	//Для вычисления дистанции умножьте шаги на длину шага lenStep и разделите на mInKm.
	return float64(steps) * lenStep / float64(mInKm)
}

func meanSpeed(steps int, duration time.Duration) float64 {
	//Проверить, что продолжительность duration больше 0.
	if duration <= 0 {
		//Если это не так, вернуть 0.
		return 0
	}
	//Вычислить дистанцию с помощью distance().
	dist := distance(steps)
	//Вычислить и вернуть среднюю скорость.
	speed := dist / duration.Hours()
	return speed
}

func TrainingInfo(data string, weight, height float64) string {
	//Получить значения из строки данных с помощью функции parseTraining(), обработать возможные ошибки.
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Ошибка парсинга: %v", err)
	}
	//какой вид тренировки был передан в строке, которую парсили (лучше использовать switch).
	switch activityType {
	case "Бег":
		//рассчитать дистанцию, среднюю скорость и калории.
		distance := distance(steps)
		speed := meanSpeed(steps, duration)
		calories := RunningSpentCalories(steps, weight, duration)
		//сформируйте и верните строку, образец которой был представлен выше.
		return fmt.Sprintf("Тип тренировки: %s\n Длительность: %.2f ч.\n Дистанция: %.2f км.\n Скорость: %.2f км/ч\n Сожгли калорий: %.2f", activityType, duration.Hours(), distance, speed, calories)
	case "Ходьба":
		distance := distance(steps)
		speed := meanSpeed(steps, duration)
		calories := WalkingSpentCalories(steps, weight, height, duration)
		return fmt.Sprintf("Тип тренировки: %s\n Длительность: %.2f ч.\n Дистанция: %.2f км.\n Скорость: %.2f км/ч\n Сожгли калорий: %.2f", activityType, duration.Hours(), distance, speed, calories)
	//Если был передан неизвестный тип тренировки, верните "неизвестный тип тренировки".
	default:
		return "Неизвестный тип тренировки"
	}
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	//Рассчитать среднюю скорость с помощью meanSpeed().
	speed := meanSpeed(steps, duration)
	//Рассчитать и вернуть количество калорий.
	calories := ((runningCaloriesMeanSpeedMultiplier * speed) - runningCaloriesMeanSpeedShift) * weight
	return calories
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	//Рассчитать среднюю скорость с помощью meanSpeed().
	speed := meanSpeed(steps, duration)
	//Продолжительность duration нужно перевести в часы.
	durationInHours := duration.Hours()
	//Рассчитать и вернуть количество калорий.
	calories := ((walkingCaloriesWeightMultiplier * weight) + ((speed * speed / height) * walkingSpeedHeightMultiplier)) * durationInHours * float64(minInH)
	return calories
}
